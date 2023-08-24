package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

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
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at" format:"2006-01-02 15:04:05"`
	UpdatedAt   time.Time      `json:"updated_at" format:"2006-01-02 15:04:05"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Uuid        string         `gorm:"not null;type:varchar(36);comment:uuid" json:"uuid"` // uuid
	Header      Dict           `gorm:"type:json;comment:请求头" json:"header"`
	Raw         string         `gorm:"comment:请求体" json:"raw"`
	QueryParams Dict           `gorm:"type:json;comment:请求参数" json:"query_params"`
	FormData    Dict           `gorm:"type:json;comment:表单数据" json:"form_data"`
	Host        string         `gorm:"type:varchar(256);comment:请求地址" json:"host"`
	Method      string         `gorm:"type:varchar(16);comment:请求方法" json:"method"`
	UserID      int            `gorm:"comment:用户id" json:"user_id"`
}

// 自定义下序列化的格式
func (r ResourceCalledHistory) MarshalJSON() ([]byte, error) {
	// 定义一个该结构体的别名
	type R ResourceCalledHistory
	// 定义一个新的结构体
	temp := struct {
		R
		UpdatedAt string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
	}{
		R:         (R)(r),
		UpdatedAt: r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(temp)
}

func (ResourceCalledHistory) TableName() string {
	return "alinlab_webhook_resource_called_history"
}

func CreateNewResourceCalledHistory(db *gorm.DB, history *ResourceCalledHistory) error {
	return db.Create(history).Error
}

func QueryAllHistoryByUUid(uuid string, db *gorm.DB) ([]*ResourceCalledHistory, error) {
	var history []*ResourceCalledHistory
	err := db.Where("uuid = ?", uuid).Find(&history).Error
	return history, err
}
