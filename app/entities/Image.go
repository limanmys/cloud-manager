package entities

import gormjsonb "github.com/dariubs/gorm-jsonb"

type Image struct {
	Base
	Checksum        string `json:"checksum"`
	ContainerFormat string `json:"container_format"`
	//CreatedAt       time.Time `json:"created_at"`
	DiskFormat string `json:"disk_format"`
	File       string `json:"file"`
	//ID         string          `json:"id"`
	MinDisk    int             `json:"min_disk"`
	MinRAM     int             `json:"min_ram"`
	Name       string          `json:"name"`
	Owner      string          `json:"owner"`
	Properties gormjsonb.JSONB `json:"properties" gorm:"type:jsonb"`
	Protected  bool            `json:"protected"`
	Schema     string          `json:"schema"`
	Size       int             `json:"size"`
	Status     string          `json:"status"`
	Tags       []interface{}   `json:"tags"`
	//UpdatedAt   time.Time     `json:"updated_at"`
	VirtualSize int64  `json:"virtual_size"`
	Visibility  string `json:"visibility"`
	Cloud       *Cloud `json:"cloud" gorm:"many2many:cloud_images"`
}
