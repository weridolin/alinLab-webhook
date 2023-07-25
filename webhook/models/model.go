package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Dict map[string]string

func (j *Dict) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}

	return json.Unmarshal(b, &j)
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j Dict) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	return json.Marshal(j)
}

type ResourceCalledHistory struct {
	gorm.Model
	Uuid        string `gorm:"not null;type:varchar(36);comment:uuid" json:"uuid"` // uuid
	Header      Dict   `gorm:"type:json;comment:请求头" json:"header"`
	Raw         string `gorm:"comment:请求体" json:"raw"`
	QueryParams Dict   `gorm:"type:json;comment:请求参数" json:"query_params"`
	FormData    Dict   `gorm:"type:json;comment:表单数据" json:"form_data"`
	Host        string `gorm:"type:varchar(256);comment:请求地址" json:"host"`
	Method      string `gorm:"type:varchar(16);comment:请求方法" json:"method"`
	UserID      int    `gorm:"comment:用户id" json:"user_id"`
}

func (ResourceCalledHistory) TableName() string {
	return "alinlab_webhook_resource_called_history"
}

func CreateNewResourceCalledHistory(db *gorm.DB, history *ResourceCalledHistory) error {
	return db.Create(history).Error
}
