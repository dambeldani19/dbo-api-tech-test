package helper

import (
	"dbo-api/dto"
)

type ResponseWithData struct {
	Code       int             `json:"code"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Pagination *dto.Pagination `json:"pagination,omitempty"`
	Data       any             `json:"data"`
}

type ResponseWithOutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(params dto.ResponseParam) any {
	var response any
	var status string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	} else {
		status = "failed"
	}

	if params.Data != nil {
		response = &ResponseWithData{
			Code:       params.StatusCode,
			Status:     status,
			Message:    params.Message,
			Pagination: params.Pagination,
			Data:       params.Data,
		}
	} else {
		response = &ResponseWithOutData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
		}
	}

	return response
}
