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

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repository repository.AuthRepository
	custRepo   repository.CustomerRepository
}

func NewAuthService(r repository.AuthRepository, c repository.CustomerRepository) *authService {
	return &authService{
		repository: r,
		custRepo:   c,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {

	if emailExist := s.repository.EmailExist(req.Email); emailExist {
		return &errorhandler.BadRequestError{Message: "email already registered"}
	}

	if req.Password != req.PasswordConfirmation {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user := entity.User{
		Username:  req.UserName,
		Password:  passwordHash,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	idUser, err := s.repository.Register(&user)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	cust := entity.Customer{
		UserID:    idUser,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: time.Now(),
	}

	err = s.custRepo.Insert(&cust)
	if err != nil {
		fmt.Println("insert customer")
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil

}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var data *dto.LoginResponse

	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	token, err := helper.GenerateToken(user)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	data = &dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Username,
		Token: token,
	}
	return data, nil
}
