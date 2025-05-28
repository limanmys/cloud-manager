package entities

type Cloud struct {
	Base
	Name     string    `json:"name" gorm:"uniqueIndex" validate:"required"`
	Type     string    `json:"type" validate:"required"`
	Machines []Machine `json:"machines" gorm:"many2many:cloud_machines"`
}
