package model

type (
	Offer struct {
		ID          uint    `json:"id"`
		CustomerID  uint    `json:"customer_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	OfferCreateReq struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	OfferUpdateReq struct {
		OfferID     uint    `json:"offerID"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	OfferDeleteReq struct {
		OfferID uint `json:"offerID"`
	}
)
