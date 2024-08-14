package repository

import (
	"dbo-api/dto"
	"dbo-api/entity"
	"math"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	Insert(cust *entity.Customer) error
	GetDetail(id int) (*entity.Customer, error)
	GetList(param *dto.ParamRequest) ([]entity.Customer, dto.ResponseParam, error)
	Update(cust *entity.Customer) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *customerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) Insert(cust *entity.Customer) error {

	err := r.db.Create(&cust).Error
	return err
}

func (r *customerRepository) GetDetail(id int) (*entity.Customer, error) {
	var cust *entity.Customer

	err := r.db.Where("deleted_at IS NULL").First(&cust, "id = ?", id).Error

	return cust, err
}

func (r *customerRepository) GetList(param *dto.ParamRequest) ([]entity.Customer, dto.ResponseParam, error) {
	var arrCust []entity.Customer
	var paginationInfo dto.ResponseParam
	var total int64

	// Hitung total data
	err := r.db.Model(&entity.Customer{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, paginationInfo, err
	}

	// Pagination
	offset := (param.Pagination.Page - 1) * param.Pagination.PerPage

	query := r.db.Where("deleted_at IS NULL")

	if param.Search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR address LIKE ? OR phone LIKE ?", param.Search, param.Search, param.Search, param.Search)
	}

	err = query.Offset(offset).
		Limit(param.Pagination.PerPage).
		Find(&arrCust).Error
	if err != nil {
		return nil, paginationInfo, err
	}

	// Hitung total halaman
	totalPage := int(math.Ceil(float64(total) / float64(param.Pagination.PerPage)))

	paginationInfo = dto.ResponseParam{
		Pagination: &dto.Pagination{
			Page:      param.Pagination.Page,
			PerPage:   param.Pagination.Page,
			Total:     int(total),
			TotalPage: totalPage,
		},
	}

	return arrCust, paginationInfo, err
}

func (r *customerRepository) Update(cust *entity.Customer) error {
	err := r.db.Model(&cust).Updates(&cust).Error

	return err
}
