package handler

import (
	"dbo-api/dto"
	"dbo-api/errorhandler"
	"dbo-api/helper"
	"dbo-api/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *customerHandler {
	return &customerHandler{
		service: service,
	}
}

func (h *customerHandler) GetList(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	perpage, _ := strconv.Atoi(c.Query("per_page"))

	if perpage <= 0 {
		perpage = 10
	}

	params := dto.ParamRequest{
		Search: c.Query("search"),
		Pagination: dto.Pagination{
			Page:    page,
			PerPage: perpage,
		},
	}

	data, param, err := h.service.GetList(params)
	if err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
		Pagination: param.Pagination,
	})

	c.JSON((http.StatusOK), res)

}

func (h *customerHandler) GetDetail(c *gin.Context) {

	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	data, err := h.service.GetDetail(intID)
	if err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data:       data,
	})

	c.JSON((http.StatusOK), res)

}

func (h *customerHandler) Update(c *gin.Context) {

	var req dto.CustomerRequest
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Update(req, intID); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success",
	})

	c.JSON((http.StatusOK), res)

}

func (h *customerHandler) Delete(c *gin.Context) {

	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	if err := h.service.Delete(intID); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success delete",
	})

	c.JSON((http.StatusOK), res)

}
