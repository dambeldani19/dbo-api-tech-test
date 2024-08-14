package service

import (
	"database/sql"
	"dbo-api/dto"
	"dbo-api/errorhandler"
	"dbo-api/repository"
	"fmt"
	"time"
)

type CustomerService interface {
	GetDetail(id int) (*dto.CustomerResponse, error)
	GetList(param dto.ParamRequest) ([]dto.CustomerResponse, dto.ResponseParam, error)
	Update(req dto.CustomerRequest, id int) error
	Delete(id int) error
}

type customerService struct {
	repository repository.CustomerRepository
	authrepo   repository.AuthRepository
}

func NewCustomerService(r repository.CustomerRepository, a repository.AuthRepository) *customerService {
	return &customerService{
		repository: r,
		authrepo:   a,
	}
}

func (s customerService) GetDetail(id int) (*dto.CustomerResponse, error) {

	cust, err := s.repository.GetDetail(id)
	if err != nil {
		return nil, &errorhandler.BadRequestError{Message: "data customer tidak di temukan"}
	}

	res := dto.CustomerResponse{
		ID:      cust.ID,
		Name:    cust.Name,
		Email:   cust.Email,
		Phone:   cust.Phone,
		Address: cust.Address,
	}

	return &res, nil
}

func (s *customerService) GetList(param dto.ParamRequest) ([]dto.CustomerResponse, dto.ResponseParam, error) {

	var arrData []dto.CustomerResponse
	var pagination dto.ResponseParam

	list, pagination, err := s.repository.GetList(&param)
	if err != nil {
		return nil, pagination, &errorhandler.BadRequestError{Message: "data kosong"}
	}

	for _, cust := range list {
		res := dto.CustomerResponse{
			ID:      cust.ID,
			Name:    cust.Name,
			Email:   cust.Email,
			Phone:   cust.Phone,
			Address: cust.Address,
		}
		arrData = append(arrData, res)
	}

	return arrData, pagination, err
}

func (s *customerService) Update(req dto.CustomerRequest, id int) error {

	fmt.Println(id)
	cust, err := s.repository.GetDetail(id)
	if err != nil {
		return &errorhandler.BadRequestError{Message: "data customer tidak di temukan"}
	}

	if req.Name != "" {
		cust.Name = req.Name
	}

	if req.Email != "" {
		user, err := s.authrepo.GetUserByEmail(cust.Email)
		if err != nil {
			return &errorhandler.InternalServerError{Message: err.Error()}
		}
		user.Email = req.Email

		err = s.authrepo.Update(user)
		if err != nil {
			return &errorhandler.InternalServerError{Message: err.Error()}
		}
		cust.Email = req.Email
	}

	if req.Phone != "" {
		cust.Phone = req.Phone
	}

	if req.Address != "" {
		cust.Address = req.Address
	}

	err = s.repository.Update(cust)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return err
}

func (s *customerService) Delete(id int) error {

	cust, err := s.repository.GetDetail(id)
	if err != nil {
		return &errorhandler.BadRequestError{Message: "data customer tidak di temukan"}
	}
	cust.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// update deleted at customer
	err = s.repository.Update(cust)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user, err := s.authrepo.GetUserByEmail(cust.Email)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}
	// user.DeletedAt.Time = time.Now()
	user.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// update deleted at user
	err = s.authrepo.Update(user)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return err
}
