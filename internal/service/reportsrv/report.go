package reportsrv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
	"github.com/jung-kurt/gofpdf"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type service struct {
	paymentGw gateway.PaymentService
	menuGw    gateway.MenuService
}

func New(paymentGw gateway.PaymentService, menuGw gateway.MenuService) service {
	return service{paymentGw: paymentGw, menuGw: menuGw}
}

func (srv service) GenerateMontlyReport(date time.Time) error {

	var (
		totalOrder int
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
		totalOrder = 0
	} else {
		totalOrder = len(ordersResp.Data)
	}

	mReport := make(map[string]int, 0)
	menuCache := make(map[string]models.MenuGetByIdResponseBody, 0)
	discountCache := make(map[string]models.GetDiscountByCodeResponseBody, 0)

	// Get cart detail from `cart_id`
	// WARNING: THE MEMMORY USAGE CONCERN HERE
	for _, eachOrder := range ordersResp.Data {
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

			var discountVal string
			if eachOrder.DiscountCode != "" {
				// caching here
				discountRespBodyCache, ok := discountCache[eachOrder.DiscountCode]
				if !ok {
					discountResp, err := srv.paymentGw.GetDiscountByCode(eachOrder.DiscountCode)
					if err != nil {
						logger.Error("can't get discount by code", zap.String("code", eachOrder.DiscountCode), zap.Error(err))
						return err
					}
					discountCache[eachOrder.DiscountCode] = discountResp.Data
					discountVal = discountResp.Data.Value.String()
				} else {
					discountVal = discountRespBodyCache.Value.String()
				}
			}

			mReportKey := generateMapReportKey(eachCart.MenuId, discountVal)
			qty, ok := mReport[mReportKey]
			if !ok {
				mReport[mReportKey] = 0
			} else {
				mReport[mReportKey] = qty + 1
			}
		}
	}

	items := generateItemData(mReport, menuCache, discountCache)

	// reset
	mReport = nil
	menuCache = nil
	discountCache = nil

	if err := generatePDF(date, totalOrder, items, "", "", ""); err != nil {
		return err
	}

	return nil
}

func generateMapReportKey(menuName string, discountVal string) string {

	if discountVal == "" {
		return menuName
	} else {
		return menuName + "_" + discountVal
	}
}

func generateItemData(mReport map[string]int, menuCache map[string]models.MenuGetByIdResponseBody, discountCache map[string]models.GetDiscountByCodeResponseBody) [][]string {

	var result [][]string

	const defaultTax = "5%"

	for key, qty := range mReport {

		splitedKey := strings.Split(key, "_")

		switch len(splitedKey) {
		case 1: // no discount
			val := menuCache[key]
			lineTotal := decimal.NewFromFloat(val.Price * float64(qty))
			result = append(result, []string{strconv.Itoa(qty), val.FNname, fmt.Sprintf("%.2f", val.Price), "0%", defaultTax, lineTotal.String()})
		case 2: // with discount
			val := menuCache[splitedKey[0]]
			lineTotal := decimal.NewFromFloat(val.Price * float64(qty))
			result = append(result, []string{strconv.Itoa(qty), val.FNname, fmt.Sprintf("%.2f", val.Price), splitedKey[1], defaultTax, lineTotal.String()})
		}
	}

	return result
}

const (
	ReportName = "ICS_Monthly_Report"
)

func generatePDF(date time.Time, totalOrder int, cartItems [][]string, subTotal string, totalDiscount string, total string) error {

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
	pdf.SetFont(font, "B", 12)
	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(126, 126, 126)
	pdf.CellFormat(13, 10, "Qty", "B", 0, "L", false, 0, "")
	pdf.CellFormat(60, 10, "Description", "B", 0, "L", false, 0, "")
	pdf.CellFormat(25, 10, "Price", "B", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, "Discount", "B", 0, "R", false, 0, "")
	pdf.CellFormat(35, 10, "Line Total", "B", 1, "R", false, 0, "")

	// Sample invoice items (loop through your actual invoice items here)
	// items := [][]string{
	// 	{"12", "Chocolate", "$10.00", "0%", "5%", "$184.40"},
	// 	{"20", "Vanilla", "$9.80", "0%", "5%", "$107.80"},
	// 	{"5", "Chai", "$8.59", "0%", "5%", "$191.40"},
	// }

	// Table content
	pdf.SetFont(font, "", 12)
	pdf.SetDrawColor(236, 236, 236)
	for _, item := range cartItems {
		for idx, col := range item {
			switch idx {
			case 0:
				pdf.CellFormat(13, 10, col, "B", 0, "L", false, 0, "")
			case 1:
				pdf.CellFormat(60, 10, col, "B", 0, "L", false, 0, "")
			case 2:
				pdf.CellFormat(25, 10, col, "B", 0, "R", false, 0, "")
			case 3:
				pdf.CellFormat(30, 10, col, "B", 0, "R", false, 0, "")
			case 4:
				pdf.CellFormat(35, 10, col, "B", 0, "R", false, 0, "")
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

	// Generate output file
	err := pdf.OutputFileAndClose(generateReportName())
	if err != nil {
		logger.Error("can't create output PDF file", zap.Error(err))
		return err
	}

	logger.Info("Monthly report generated successfully.")

	return nil
}

func generateReportName() string {
	return ReportName + "_" + time.Now().Format("20060102") + ".pdf"
}
