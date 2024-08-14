package service

import (
	"database/sql"
	"dbo-api/dto"
	"dbo-api/entity"
	"dbo-api/errorhandler"
	"dbo-api/repository"
	"time"
)

type ProductService interface {
	Insert(req dto.ProductRequest) error
	GetDetail(id int) (*dto.ProductResponse, error)
	GetList(param dto.ParamRequest) ([]dto.ProductResponse, dto.ResponseParam, error)
	Update(req dto.ProductRequest, id int) error
	Delete(id int) error
}

type productService struct {
	repository repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) *productService {
	return &productService{
		repository: r,
	}
}

func (s *productService) Insert(req dto.ProductRequest) error {

	prod := entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
	}
	err := s.repository.Insert(&prod)

	return err
}

func (s *productService) GetDetail(id int) (*dto.ProductResponse, error) {

	p, err := s.repository.GetDetail(id)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "data product tidak di temukan"}
	}

	res := dto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}

	return &res, nil
}

func (s *productService) GetList(param dto.ParamRequest) ([]dto.ProductResponse, dto.ResponseParam, error) {

	var arrData []dto.ProductResponse
	var pagination dto.ResponseParam

	list, pagination, err := s.repository.GetList(&param)
	if err != nil {
		return nil, pagination, &errorhandler.BadRequestError{Message: "data kosong"}
	}

	for _, p := range list {
		res := dto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		}
		arrData = append(arrData, res)
	}

	return arrData, pagination, err
}

func (s *productService) Update(req dto.ProductRequest, id int) error {

	p, err := s.repository.GetDetail(id)
	if err != nil {
		return &errorhandler.BadRequestError{Message: "data product tidak di temukan"}
	}

	if req.Name != "" {
		p.Name = req.Name
	}

	if req.Description != "" {
		p.Description = req.Description

	}

	if req.Price > 0 {
		p.Price = req.Price
	}

	if req.Stock > 0 {
		p.Stock = req.Stock
	}

	err = s.repository.Update(p)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return err
}

func (s *productService) Delete(id int) error {

	p, err := s.repository.GetDetail(id)
	if err != nil {
		return &errorhandler.BadRequestError{Message: "data product tidak di temukan"}
	}
	p.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// update deleted at product
	err = s.repository.Update(p)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return err
}
