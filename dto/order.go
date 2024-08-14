package dto

import "time"

type OrderRequest struct {
	OrderReq []ordReq `json:"orders"`
	Status   string   `json:"status"`
}

type ordReq struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderResponse struct {
	ID          int                   `json:"id"`
	OrderCode   string                `json:"order_code"`
	TotalAmount int                   `json:"total_amount"`
	CreatedAt   time.Time             `json:"created_at"`
	OrderDetail []OrderDetailResponse `json:"order_detail"`
}

type OrderDetailResponse struct {
	Product    ProductResponse `json:"product"`
	Quantity   int             `json:"quantity"`
	Price      int             `json:"price"`
	TotalPrice int             `json:"total_price"`
}
