package repository

import (
	"dbo-api/dto"
	"dbo-api/entity"
	"fmt"
	"math"

	"gorm.io/gorm"
)

type OrderRepository interface {
	//order
	InsertOrder(order *entity.Order) error
	Update(order *entity.Order) error
	GetList(param *dto.ParamRequest) ([]entity.Order, dto.ResponseParam, error)
	GetDetail(id, customerID int) (*entity.Order, error)

	//order detail
	AddOrderDetail(orderDetail *entity.OrderDetail) error
	DeleteProduct(ordID, prodID int) error
	UpdateQty(ordID, prodID, newQty, price int) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) InsertOrder(order *entity.Order) error {

	err := r.db.Create(&order).Error
	if err != nil {
		return err
	}

	return err
}

func (r *orderRepository) GetDetail(id, customerID int) (*entity.Order, error) {
	var ord *entity.Order

	err := r.db.Where("deleted_at IS NULL").
		Where("customer_id = ?", customerID).
		Preload("OrderDetail.Product").
		First(&ord, "id = ?", id).Error

	return ord, err
}

func (r *orderRepository) GetList(param *dto.ParamRequest) ([]entity.Order, dto.ResponseParam, error) {
	var ord []entity.Order
	var paginationInfo dto.ResponseParam
	var total int64

	// Hitung total data
	err := r.db.Model(&entity.Order{}).
		Where("deleted_at IS NULL").
		Where("customer_id = ?", param.UserID).
		Count(&total).Error
	if err != nil {
		return nil, paginationInfo, err
	}

	// Hitung offset
	offset := (param.Pagination.Page - 1) * param.Pagination.PerPage

	query := r.db.Joins("JOIN order_details ON orders.id = order_details.order_id").
		Joins("JOIN products ON order_details.product_id = products.id").
		Where("orders.deleted_at IS NULL").
		Where("customer_id = ?", param.UserID)

	if param.Search != "" {
		query = query.Where("order_code LIKE ? OR products.name LIKE ?", param.Search, param.Search)
	}

	err = query.Offset(offset).
		Limit(param.Pagination.PerPage).
		Preload("OrderDetail.Product").
		Find(&ord).Error

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

	return ord, paginationInfo, err
}

func (r *orderRepository) AddOrderDetail(orderDetail *entity.OrderDetail) error {
	err := r.db.Create(orderDetail).Error
	return err
}

func (r *orderRepository) DeleteProduct(ordID, prodID int) error {
	result := r.db.Where("order_id = ?", ordID).
		Where("product_id = ?", prodID).
		Delete(&entity.OrderDetail{})

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return result.Error
}

func (r *orderRepository) Update(order *entity.Order) error {
	err := r.db.Model(&order).Updates(&order).Error

	return err
}

func (r *orderRepository) UpdateQty(ordID, prodID, newQty, price int) error {
	err := r.db.Model(&entity.OrderDetail{}).
		Where("order_id = ? AND product_id = ?", ordID, prodID).
		Updates(map[string]interface{}{
			"quantity":    newQty,
			"total_price": newQty * price,
		}).Error

	return err
}
