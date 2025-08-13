package model

type (
	Service struct {
		ID          uint    `json:"id"`
		PerformerID uint    `json:"performer_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	ServiceCreateReq struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	ServiceUpdateReq struct {
		ServiceID   uint    `json:"serviceID"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
)
