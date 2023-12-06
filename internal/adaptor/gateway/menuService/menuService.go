package menuservice

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/adaptor/gateway/config"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/helper/logger"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
	"go.uber.org/zap"
)

type menuGateway struct {
	cfg config.MenuServiceCfg
}

func New(cfg config.MenuServiceCfg) gateway.MenuService {
	return menuGateway{cfg: cfg}
}

func (gw menuGateway) GetMenuById(id string) (models.MenuGetByIdResponse, error) {

	req, err := http.NewRequest(http.MethodGet, gw.cfg.HostURL+"/menu/"+id, nil)
	if err != nil {
		logger.Error("client: could not create request of GetCartByDateMonth", zap.String("id", id), zap.Error(err))
		return models.MenuGetByIdResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("client: error making http request of GetCartByDateMonth", zap.String("id", id), zap.Error(err))
		return models.MenuGetByIdResponse{}, err
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		return models.MenuGetByIdResponse{}, err
	}

	logger.Info("client: got response!")

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("client: could not read response body of GetCartByDateMonth", zap.String("id", id), zap.Error(err))
		return models.MenuGetByIdResponse{}, err
	}

	logger.Info("client: response body", zap.String("response_body", string(resBody)))

	var response models.MenuGetByIdResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		logger.Error("client: can't umarshal response body of GetCartByDateMonth", zap.String("id", id), zap.Error(err))
		return models.MenuGetByIdResponse{}, err
	}

	return response, nil
}
