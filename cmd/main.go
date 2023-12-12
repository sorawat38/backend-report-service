package main

import (
	"net/http"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/config"
	menuservice "github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway/menuService"
	paymentservice "github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway/paymentService"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/handler"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/handler/reporthdl"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/service/reportsrv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger.InitLog(cfg.Log)
	defer logger.CloseLogger()
	logger.Info("", zap.Any("config", cfg))
	sess, err := initAWSSession(cfg.AWSSession)
	if err != nil {
		panic(err)
	}

	// Init gateway
	paymentServiceGw := paymentservice.New(cfg.Gateway.PaymentService)
	menuServeGw := menuservice.New(cfg.Gateway.MenuService)

	reportSrv := reportsrv.New(paymentServiceGw, menuServeGw, sess)
	reportHandler := reporthdl.NewHTTPHandler(reportSrv)

	// Starting server
	e := echo.New()
	handler.InitRoute(e, reportHandler)

	logger.Infof("Starting server on port %v...\n", cfg.App.Port)
	if err := e.Start(":" + cfg.App.Port); err != http.ErrServerClosed {
		logger.Fatal(err.Error())
	}
}

func initAWSSession(config config.AWSSession) (*session.Session, error) {

	// Initialize AWS session using your credentials or IAM role.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		// Credentials: credentials.NewStaticCredentials(config.Id, config.Secret, ""),
	})
	if err != nil {
		logger.Error("Error creating AWS session", zap.Error(err))
		return nil, err
	}

	s3Svc := s3.New(sess)

	// List buckets to test the session.
	_, err = s3Svc.ListBuckets(nil)
	if err != nil {
		logger.Error("Error listing AWS bucket", zap.Error(err))
		return nil, err
	}

	logger.Info("Init AWS session successfully")

	return sess, nil
}
