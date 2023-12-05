package apiresponses

import (
	"errors"

	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/constant"
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/internal/models"
)

func SuccessResponse() models.CommonResponse {
	return models.CommonResponse{
		Code:    constant.SuccessCode,
		Message: constant.SuccessMessage,
	}
}

func InvalidInputError(err error) models.CommonResponse {
	res := models.CommonResponse{
		Code:    constant.InvalidInputCode,
		Message: constant.InvalidInputMsg,
	}
	if err != nil {
		res.Error = err.Error()
	}

	if errors.Unwrap(err) != nil {
		res.ErrorDetail = errors.Unwrap(err).Error()
	}
	return res
}

func NotFoundError(err error) models.CommonResponse {
	res := models.CommonResponse{
		Code:    constant.NotFoundCode,
		Message: constant.NotFoundMsg,
	}
	if err != nil {
		res.Error = err.Error()
	}

	if errors.Unwrap(err) != nil {
		res.ErrorDetail = errors.Unwrap(err).Error()
	}

	return res
}

func InternalError(err error) models.CommonResponse {

	res := models.CommonResponse{
		Code:    constant.InternalErrorCode,
		Message: constant.InternalError,
	}
	if err != nil {
		res.Error = err.Error()
	}

	if errors.Unwrap(err) != nil {
		res.ErrorDetail = errors.Unwrap(err).Error()
	}

	return res
}
