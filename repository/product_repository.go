package repository

import (
	"dbo-api/dto"
	"dbo-api/entity"
	"math"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Insert(product *entity.Product) error
	GetDetail(id int) (*entity.Product, error)
	GetList(param *dto.ParamRequest) ([]entity.Product, dto.ResponseParam, error)
	Update(cust *entity.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Insert(product *entity.Product) error {

	err := r.db.Create(&product).Error
	return err
}

func (r *productRepository) GetDetail(id int) (*entity.Product, error) {
	var prod *entity.Product

	err := r.db.Where("deleted_at IS NULL").First(&prod, "id = ?", id).Error

	return prod, err
}

func (r *productRepository) GetList(param *dto.ParamRequest) ([]entity.Product, dto.ResponseParam, error) {
	var arrProduct []entity.Product
	var paginationInfo dto.ResponseParam
	var total int64

	// Hitung total data
	err := r.db.Model(&entity.Product{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, paginationInfo, err
	}

	// Pagination
	offset := (param.Pagination.Page - 1) * param.Pagination.PerPage

	query := r.db.Where("deleted_at IS NULL")

	if param.Search != "" {
		query = query.Where("name LIKE ?", param.Search)
	}

	err = query.Offset(offset).
		Limit(param.Pagination.PerPage).
		Find(&arrProduct).Error
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

	return arrProduct, paginationInfo, err
}

func (r *productRepository) Update(product *entity.Product) error {
	err := r.db.Model(&product).Updates(&product).Error

	return err
}
