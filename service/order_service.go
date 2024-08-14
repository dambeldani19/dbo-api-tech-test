package service

import (
	"dbo-api/dto"
	"dbo-api/entity"
	"dbo-api/errorhandler"
	"dbo-api/helper"
	"dbo-api/repository"
	"fmt"
	"time"
)

type OrderService interface {
	CreateOrder(req *dto.OrderRequest, custID int) error
	GetList(param dto.ParamRequest) ([]dto.OrderResponse, dto.ResponseParam, error)
	GetDetail(id, customerID int) (*dto.OrderResponse, error)

	//order detail
	UpdateOrder(req dto.OrderRequest, id, customerID int) error
	DeleteProduct(req dto.OrderRequest, id, customerID int) error
}

type orderService struct {
	repository repository.OrderRepository
	prodRepo   repository.ProductRepository
}

func NewOrderService(r repository.OrderRepository, p repository.ProductRepository) *orderService {
	return &orderService{
		repository: r,
		prodRepo:   p,
	}
}

const StatusProcess = "process"
const StatusCompleted = "completed"
const StatusCancelled = "cancelled"

func (s *orderService) CreateOrder(req *dto.OrderRequest, custID int) error {
	var ordDetail []entity.OrderDetail
	var totalAmount int

	for _, od := range req.OrderReq {

		prod, err := s.prodRepo.GetDetail(od.ProductID)
		if err != nil {
			return &errorhandler.BadRequestError{Message: fmt.Sprintf("product id %d tidak ada/tidak aktif", od.ProductID)}
		}
		ordDetail = append(ordDetail, entity.OrderDetail{
			ProductID:  od.ProductID,
			Quantity:   od.Quantity,
			Price:      prod.Price,
			TotalPrice: prod.Price * od.Quantity,
		})

		totalAmount += prod.Price * od.Quantity
	}

	code, _ := helper.Generate(`ord-[a-z0-9]{16}`)

	ord := entity.Order{
		OrderCode:   code,
		CustomerID:  custID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Status:      StatusProcess,
		OrderDetail: ordDetail,
	}
	err := s.repository.InsertOrder(&ord)

	return err
}

func (s *orderService) GetList(param dto.ParamRequest) ([]dto.OrderResponse, dto.ResponseParam, error) {

	var arrData []dto.OrderResponse
	var pagination dto.ResponseParam

	list, pagination, err := s.repository.GetList(&param)
	if err != nil {
		return nil, pagination, &errorhandler.BadRequestError{Message: "data kosong"}
	}

	for _, o := range list {
		var ordRes []dto.OrderDetailResponse

		for _, od := range o.OrderDetail {
			ordRes = append(ordRes, dto.OrderDetailResponse{
				Quantity:   od.Quantity,
				Price:      od.Price,
				TotalPrice: od.TotalPrice,
				Product: dto.ProductResponse{
					ID:          od.ProductID,
					Name:        od.Product.Name,
					Description: od.Product.Description,
					Price:       od.Product.Price,
					Stock:       od.Product.Stock,
				},
			})
		}
		res := dto.OrderResponse{
			ID:          o.ID,
			OrderCode:   o.OrderCode,
			TotalAmount: o.TotalAmount,
			CreatedAt:   o.CreatedAt,
			OrderDetail: ordRes,
		}
		arrData = append(arrData, res)
	}

	return arrData, pagination, err
}

func (s *orderService) GetDetail(id, customerID int) (*dto.OrderResponse, error) {

	o, err := s.repository.GetDetail(id, customerID)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "data order tidak di temukan"}
	}

	var ordRes []dto.OrderDetailResponse

	for _, od := range o.OrderDetail {
		ordRes = append(ordRes, dto.OrderDetailResponse{
			Quantity:   od.Quantity,
			Price:      od.Price,
			TotalPrice: od.TotalPrice,
			Product: dto.ProductResponse{
				ID:          od.ProductID,
				Name:        od.Product.Name,
				Description: od.Product.Description,
				Price:       od.Product.Price,
				Stock:       od.Product.Stock,
			},
		})
	}

	res := dto.OrderResponse{
		ID:          o.ID,
		OrderCode:   o.OrderCode,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
		OrderDetail: ordRes,
	}

	return &res, nil
}

func (s *orderService) UpdateOrder(req dto.OrderRequest, id, customerID int) error {
	var totalAmount int

	// check product
	for _, od := range req.OrderReq {

		prod, err := s.prodRepo.GetDetail(od.ProductID)
		if err != nil {
			return &errorhandler.BadRequestError{Message: fmt.Sprintf("product id %d tidak ada/tidak aktif", od.ProductID)}
		}

		//update qty
		err = s.repository.UpdateQty(id, od.ProductID, od.Quantity, prod.Price)
		if err != nil {
			return &errorhandler.InternalServerError{Message: err.Error()}
		}

	}

	// check order
	o, err := s.repository.GetDetail(id, customerID)
	if err != nil {
		return &errorhandler.BadRequestError{Message: "data order tidak di temukan"}
	}

	if o.Status == StatusCompleted {
		return &errorhandler.BadRequestError{Message: "status order completed, tidak bisa di update"}
	}

	if req.Status != "" {
		if !s.validateOrderStatus(req.Status) {
			return &errorhandler.BadRequestError{
				Message: "Status order hanya bisa 'process' atau 'completed'",
			}
		}

		o.Status = req.Status
	}

	for _, ord := range o.OrderDetail {
		totalAmount += ord.TotalPrice
	}
	o.TotalAmount = totalAmount

	//update total amount dan status
	err = s.repository.Update(o)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return err
}

func (s *orderService) DeleteProduct(req dto.OrderRequest, id, customerID int) error {

	for _, v := range req.OrderReq {
		_, err := s.repository.GetDetail(id, customerID)
		if err != nil {
			return &errorhandler.BadRequestError{Message: "data order tidak di temukan"}
		}

		if err := s.repository.DeleteProduct(id, v.ProductID); err != nil {
			return &errorhandler.InternalServerError{Message: err.Error()}
		}
	}

	return nil
}

func (s *orderService) validateOrderStatus(status string) bool {
	switch status {
	case StatusProcess, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}
