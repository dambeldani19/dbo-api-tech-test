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

type orderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *orderHandler {
	return &orderHandler{
		service: service,
	}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {

	var ord dto.OrderRequest
	custID, _ := c.Get("customerID")

	if err := c.ShouldBindJSON(&ord); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.CreateOrder(&ord, custID.(int)); err != nil {
		errorhandler.HandlerError(c, &errorhandler.InternalServerError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success",
	})

	c.JSON((http.StatusOK), res)
}

func (h *orderHandler) GetListOrder(c *gin.Context) {

	custID, _ := c.Get("customerID")

	page, _ := strconv.Atoi(c.Query("page"))
	perpage, _ := strconv.Atoi(c.Query("per_page"))
	params := dto.ParamRequest{
		Search: c.Query("search"),
		UserID: custID.(int),
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

func (h *orderHandler) GetDetailOrder(c *gin.Context) {

	custID, _ := c.Get("customerID")
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	data, err := h.service.GetDetail(intID, custID.(int))
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

func (h *orderHandler) DeleteOrderProduct(c *gin.Context) {
	var ord dto.OrderRequest
	id := c.Param("id")
	custID, _ := c.Get("customerID")
	intID, _ := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&ord); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.DeleteProduct(ord, intID, custID.(int)); err != nil {
		errorhandler.HandlerError(c, &errorhandler.InternalServerError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success delete product order",
	})

	c.JSON((http.StatusOK), res)
}

func (h *orderHandler) UpdateOrder(c *gin.Context) {
	var ord dto.OrderRequest
	custID, _ := c.Get("customerID")
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&ord); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.UpdateOrder(ord, intID, custID.(int)); err != nil {
		errorhandler.HandlerError(c, &errorhandler.InternalServerError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusOK,
		Message:    "success update order",
	})

	c.JSON((http.StatusOK), res)
}
