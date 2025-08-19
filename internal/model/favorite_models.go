package model

type (
	FavoriteReq struct {
		ID         uint `json:"id"`
		CustomerID uint `json:"customer_id"`
		ServiceID  uint `json:"service_id"`
	}

	FavoriteInfoReq struct {
		ID           uint   `json:"id"`
		CustomerName string `json:"customer_name"`
		ServiceTitle string `json:"service_title"`
		ServiceID    uint   `json:"serviceID"`
	}
)
