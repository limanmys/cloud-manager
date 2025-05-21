package entities

type Book struct {
	Base
	Name      string `json:"name" validate:"required"`
	Vendor    string `json:"vendor" validate:"required"`
	PageCount int    `json:"page_count" validate:"required"`
}
