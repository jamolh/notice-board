package models

import "time"

type Notice struct {
	ID          string    `json:"id" example:"5dcca8d2-a5d6-11eb-bcbc-0242ac130002" format:"uuid"`
	Title       string    `json:"title" example:"If you find a cat, let me know!"`
	Description string    `json:"description" example:"Please help to find missing cat Salem"`
	Price       float64   `json:"price" example:"10000.0" format:"number"`
	CreatedAt   time.Time `json:"created_at" example:"2021-04-24T16:34:19.61868Z"`
	Image       []string  `json:"image" example:"salem_image.jpeg"`
}

type GetNoticesRequestDto struct {
	Field string
	Order string
}

type GetNoticeByIDRequestDto struct {
	ID           string `json:"id" example:"5dcca8d2-a5d6-11eb-bcbc-0242ac130002" format:"uuid"`
	GetAllImages bool   `json:"-"`
}
