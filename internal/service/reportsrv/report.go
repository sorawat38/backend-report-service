package reportsrv

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jung-kurt/gofpdf"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type service struct {
	paymentGw  gateway.PaymentService
	menuGw     gateway.MenuService
	awsSession *session.Session
}

func New(paymentGw gateway.PaymentService, menuGw gateway.MenuService, awsSession *session.Session) service {
	return service{paymentGw: paymentGw, menuGw: menuGw, awsSession: awsSession}
}

func (srv service) GenerateMontlyReport(date time.Time) error {

	var (
		totalOrder    int
		subTotal      decimal.Decimal
		totalDiscount decimal.Decimal
		total         decimal.Decimal
	)

	// Get orders from payment service
	ordersResp, err := srv.paymentGw.GetOrderByDateMonth(date)
	if err != nil {
		logger.Error("can't get orders by date month", zap.Error(err))
		return err
	}

	// Check lenght of response data
	if len(ordersResp.Data) == 0 {
		logger.Warnf("orders from getting by `%v` is empty", date.Format(time.DateOnly))
	}

	mReport := make(map[string]int, 0)
	menuCache := make(map[string]models.MenuGetByIdResponseBody, 0)
	discountCache := make(map[string]models.GetDiscountByCodeResponseBody, 0)

	// Get cart detail from `cart_id`
	// WARNING: THE MEMMORY USAGE CONCERN HERE
	for _, eachOrder := range ordersResp.Data {

		if eachOrder.Status == "Cancel" {
			continue
		}

		cartsResp, err := srv.paymentGw.GetCartById(eachOrder.CartId)
		if err != nil {
			logger.Error("can't get carts by id", zap.String("cart_id", eachOrder.CartId), zap.Error(err))
			return err
		}

		// Check lenght of response data
		if len(cartsResp.Data) == 0 {
			newErr := errors.New("carts from getting by `id` is empty") // this must not empty
			logger.Error(newErr.Error())
			return newErr
		}
		totalOrder += len(cartsResp.Data)
		for _, eachCart := range cartsResp.Data {

			// caching here
			_, ok := menuCache[eachCart.MenuId]
			if !ok {
				menuResp, err := srv.menuGw.GetMenuById(eachCart.MenuId)
				if err != nil {
					logger.Error("can't get menu by id", zap.String("menu_id", eachCart.MenuId), zap.Error(err))
					return err
				}
				menuCache[eachCart.MenuId] = menuResp.Data
			}

			// var discountCode string
			if eachOrder.DiscountCode != "" {
				// caching here
				_, ok := discountCache[eachOrder.DiscountCode]
				if !ok {
					discountResp, err := srv.paymentGw.GetDiscountByCode(eachOrder.DiscountCode)
					if err != nil {
						logger.Error("can't get discount by code", zap.String("code", eachOrder.DiscountCode), zap.Error(err))
						return err
					}
					discountCache[eachOrder.DiscountCode] = discountResp.Data
				}
			}

			mReportKey := generateMapReportKey(eachCart.MenuId, eachOrder.DiscountCode)
			qty, ok := mReport[mReportKey]
			if !ok {
				mReport[mReportKey] = 1
			} else {
				mReport[mReportKey] = qty + 1
			}
		}

		subTotal = subTotal.Add(eachOrder.SubTotal).Round(2)
		total = total.Add(eachOrder.TotalAmount).Round(2)
		totalDiscount = subTotal.Sub(total).Round(2)
	}

	items := generateItemData(mReport, menuCache, discountCache)

	// reset
	mReport = nil
	menuCache = nil
	discountCache = nil

	pdfBytes, err := generatePDF(date, totalOrder, items, subTotal.StringFixed(2), totalDiscount.StringFixed(2), total.StringFixed(2))
	if err != nil {
		return err
	}

	// Send to SNS to send email
	err = srv.notifyEmail(pdfBytes)
	if err != nil {
		return err
	}

	return nil
}

func generateMapReportKey(menuName string, discountCode string) string {

	if discountCode == "" {
		return menuName
	} else {
		return menuName + "_" + discountCode
	}
}

func generateItemData(mReport map[string]int, menuCache map[string]models.MenuGetByIdResponseBody, discountCache map[string]models.GetDiscountByCodeResponseBody) [][]string {

	var result [][]string

	for key, qty := range mReport {

		splitedKey := strings.Split(key, "_")

		switch len(splitedKey) {
		case 1: // no discount
			val := menuCache[key]
			lineTotal := decimal.NewFromFloat(val.Price * float64(qty))
			result = append(result, []string{strconv.Itoa(qty), val.FNname, fmt.Sprintf("%.2f", val.Price), "0%", lineTotal.Round(2).StringFixed(2)})
		case 2: // with discount
			val := menuCache[splitedKey[0]]
			discount := discountCache[splitedKey[1]]
			lineSubTotal := decimal.NewFromFloat(val.Price * float64(qty))
			lineTotal := lineSubTotal.Sub(lineSubTotal.Mul(discount.Value.Div(decimal.NewFromFloat(100.00))))
			result = append(result, []string{strconv.Itoa(qty), val.FNname, fmt.Sprintf("%.2f", val.Price), discount.Value.String() + "%", lineTotal.Round(2).StringFixed(2)})
		}
	}

	return result
}

const (
	ReportName = "ICS_Monthly_Report"
)

func generatePDF(date time.Time, totalOrder int, cartItems [][]string, subTotal string, totalDiscount string, total string) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	const (
		font = "Arial"
	)

	// Set font
	pdf.SetFont(font, "", 12)

	// Title
	pdf.SetFont(font, "B", 36)
	pdf.CellFormat(0, 10, "Monthly Report", "0", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Contract Information
	pdf.SetFont(font, "B", 12)
	pdf.CellFormat(0, 6, "Ice Cream Shop", "0", 1, "L", false, 0, "")
	pdf.SetFont(font, "", 12)
	pdf.CellFormat(0, 6, "(123) 456-7890", "0", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Address
	pdf.CellFormat(0, 5, "154 North West", "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, "Calgary, Alberta", "0", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, "T3P9J6", "0", 1, "", false, 0, "")
	pdf.Ln(15)

	// Report information
	pdf.SetFont(font, "B", 12)
	pdf.CellFormat(12, 5, "Date:", "0", 0, "L", false, 0, "")
	pdf.SetFont(font, "", 12)
	pdf.CellFormat(0, 5, date.Format(time.DateOnly), "0", 1, "L", false, 0, "")
	pdf.Ln(15)

	pdf.CellFormat(25, 5, "Total order: "+strconv.Itoa(totalOrder), "0", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Table header
	qtyWidth := 20.0
	descpWidth := 60.0
	priceWidth := 35.0
	discountWidth := 35.0
	lineTotalWidth := 40.0

	pdf.SetFont(font, "B", 12)
	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(126, 126, 126)
	pdf.CellFormat(qtyWidth, 10, "Qty", "B", 0, "L", false, 0, "")
	pdf.CellFormat(descpWidth, 10, "Description", "B", 0, "L", false, 0, "")
	pdf.CellFormat(priceWidth, 10, "Price", "B", 0, "R", false, 0, "")
	pdf.CellFormat(discountWidth, 10, "Discount", "B", 0, "R", false, 0, "")
	pdf.CellFormat(lineTotalWidth, 10, "Line Total", "B", 1, "R", false, 0, "")

	// Table content
	pdf.SetFont(font, "", 12)
	pdf.SetDrawColor(236, 236, 236)
	for _, item := range cartItems {
		for idx, col := range item {
			switch idx {
			case 0:
				pdf.CellFormat(qtyWidth, 10, col, "B", 0, "L", false, 0, "")
			case 1:
				pdf.CellFormat(descpWidth, 10, col, "B", 0, "L", false, 0, "")
			case 2:
				pdf.CellFormat(priceWidth, 10, col, "B", 0, "R", false, 0, "")
			case 3:
				pdf.CellFormat(discountWidth, 10, col, "B", 0, "R", false, 0, "")
			case 4:
				pdf.CellFormat(lineTotalWidth, 10, col, "B", 0, "R", false, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	// Set default
	pdf.SetLineWidth(0.2)
	pdf.SetDrawColor(0, 0, 0) // black
	pdf.Ln(10)

	// Total
	pdf.SetFont(font, "B", 12)
	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(126, 126, 126)

	pdf.CellFormat(20, 10, "Subtotal", "T", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, subTotal, "T", 1, "R", false, 0, "")

	pdf.SetDrawColor(236, 236, 236)

	pdf.CellFormat(20, 10, "Discount", "T", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, totalDiscount, "T", 1, "R", false, 0, "")

	pdf.SetDrawColor(126, 126, 126)
	pdf.CellFormat(0, 0, "", "B", 1, "R", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont(font, "B", 16)
	pdf.CellFormat(20, 10, "Total", "0", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, total, "0", 1, "R", false, 0, "")

	// Create a bytes.Buffer to hold the PDF content
	var buf bytes.Buffer

	// Output the PDF content to the buffer
	if err := pdf.Output(&buf); err != nil {
		logger.Error("Error generating PDF", zap.Error(err))
		return nil, err
	}

	// Convert the buffer content to a byte slice
	pdfBytes := buf.Bytes()

	logger.Info("Monthly report generated successfully.")

	return pdfBytes, nil
}

func generateReportName() string {
	return ReportName + "_" + time.Now().Format("20060102150405") + ".pdf"
}

func (srv service) notifyEmail(fileContent []byte) error {

	// Create new seession for S3
	s3Svc := s3.New(srv.awsSession)

	bucketName := "clcm3102-group-4-project"
	objectKey := generateReportName()

	// Upload the file to S3.
	_, err := s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileContent),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		logger.Error("Error uploading file to S3", zap.Error(err))
		return err
	}

	// Create new session for sns
	svc := sns.New(srv.awsSession)
	topicArn := "arn:aws:sns:us-east-1:698668199773:clcm3102-group-4-project"

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey)
	message := fmt.Sprintf("Message with attachment. File URL: %s", fileURL)

	params := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
	}

	// Publish the message.
	result, err := svc.Publish(params)
	if err != nil {
		logger.Error("Error publishing message", zap.Error(err))
		return err
	}

	logger.Infof("Messsage published with message ID: %s", *result.MessageId)

	return nil
}
