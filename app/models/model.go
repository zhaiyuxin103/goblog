package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time      `gorm:"column:created_at;index;comment:创建时间" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;comment:最后编辑时间" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
}
