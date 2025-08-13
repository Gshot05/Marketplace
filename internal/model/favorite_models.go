package model

type (
	FavoriteReq struct {
		ID         uint `json:"id"`
		CustomerID uint `json:"customer_id"`
		ServiceID  uint `json:"service_id"`
	}

	FavoriteInfoReq struct {
		ID                 uint   `json:"id"`
		CustomerName       string `json:"customer_name"`
		ServiceTitle       string `json:"service_title"`
		ServiceDescription string `json:"service_description"`
	}

	FavoriteDeleteReq struct {
		ServiceID uint `json:"serviceID"`
	}
)
