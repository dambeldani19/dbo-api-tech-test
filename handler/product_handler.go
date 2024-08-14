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

type productHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *productHandler {
	return &productHandler{
		service: service,
	}
}

func (h *productHandler) AddProduct(c *gin.Context) {
	var register dto.ProductRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.HandlerError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Insert(register); err != nil {
		errorhandler.HandlerError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParam{
		StatusCode: http.StatusCreated,
		Message:    "add product success",
	})

	c.JSON(http.StatusCreated, res)

}

func (h *productHandler) GetListProduct(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	perpage, _ := strconv.Atoi(c.Query("per_page"))
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

func (h *productHandler) GetDetail(c *gin.Context) {

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

func (h *productHandler) Update(c *gin.Context) {

	var req dto.ProductRequest
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

func (h *productHandler) Delete(c *gin.Context) {

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
