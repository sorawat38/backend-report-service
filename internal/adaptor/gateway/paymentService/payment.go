package paymentservice

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway/config"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
	"go.uber.org/zap"
)

type paymentGateway struct {
	cfg config.PaymentServiceCfg
}

func New(cfg config.PaymentServiceCfg) gateway.PaymentService {
	return paymentGateway{cfg: cfg}
}

func (gw paymentGateway) GetOrderByDateMonth(date time.Time) (models.GetOrderByDateMonthResponse, error) {

	date_format_yyyymmdd := date.Format(time.DateOnly)

	req, err := http.NewRequest(http.MethodGet, gw.cfg.HostURL+"/order/"+date_format_yyyymmdd, nil)
	if err != nil {
		logger.Error("client: could not create request", zap.String("date", date_format_yyyymmdd), zap.Error(err))
		return models.GetOrderByDateMonthResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("client: error making http request", zap.String("date", date_format_yyyymmdd), zap.Error(err))
		return models.GetOrderByDateMonthResponse{}, err
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		return models.GetOrderByDateMonthResponse{}, err
	}

	logger.Info("client: got response!")

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("client: could not read response body", zap.String("date", date_format_yyyymmdd), zap.Error(err))
		return models.GetOrderByDateMonthResponse{}, err
	}

	logger.Info("client: response body", zap.String("response_body", string(resBody)))

	var response models.GetOrderByDateMonthResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		logger.Error("client: can't umarshal response body", zap.String("date", date_format_yyyymmdd), zap.Error(err))
		return models.GetOrderByDateMonthResponse{}, err
	}

	return response, nil
}

func (gw paymentGateway) GetCartById(cartId string) (models.GetCartByIdResponse, error) {

	req, err := http.NewRequest(http.MethodGet, gw.cfg.HostURL+"/cart/"+cartId, nil)
	if err != nil {
		logger.Error("client: could not create request of GetCartByDateMonth", zap.String("cart_id", cartId), zap.Error(err))
		return models.GetCartByIdResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("client: error making http request of GetCartByDateMonth", zap.String("cart_id", cartId), zap.Error(err))
		return models.GetCartByIdResponse{}, err
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		return models.GetCartByIdResponse{}, err
	}

	logger.Info("client: got response!")

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("client: could not read response body of GetCartByDateMonth", zap.String("cart_id", cartId), zap.Error(err))
		return models.GetCartByIdResponse{}, err
	}

	logger.Info("client: response body", zap.String("response_body", string(resBody)))

	var response models.GetCartByIdResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		logger.Error("client: can't umarshal response body of GetCartByDateMonth", zap.String("cart_id", cartId), zap.Error(err))
		return models.GetCartByIdResponse{}, err
	}

	return response, nil
}
