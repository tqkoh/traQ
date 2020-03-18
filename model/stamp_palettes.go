package model

import (
	"database/sql/driver"
	"errors"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

type UUIDs []uuid.UUID

func (arr UUIDs) Value() (driver.Value, error) {
	idStr := []string{}
	for _, id := range arr {
		idStr = append(idStr, id.String())
	}
	return strings.Join(idStr, ","), nil
}

func (arr *UUIDs) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*arr = UUIDs{}
	case string:
		idSlice := strings.Split(s, ",")
		for _, id := range idSlice {
			stampID, _ := uuid.FromString(id)
			*arr = append(*arr, stampID)
		}
	case []byte:
		str := string(s)
		idSlice := strings.Split(str, ",")
		for _, id := range idSlice {
			stampID, _ := uuid.FromString(id)
			*arr = append(*arr, stampID)
		}
	default:
		return errors.New("failed to scan Stamps")
	}
	return nil
}

type StampPalette struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	Name        string    `gorm:"type:varchar(30);not null"`
	Description string    `gorm:"type:text(1000);not null"`
	Stamps      UUIDs     `gorm:"type:text;not null"`
	CreatorID   uuid.UUID `gorm:"type:char(36);not null"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

// TableName StampPalettes構造体のテーブル名
func (*StampPalette) TableName() string {
	return "stamp_palettes"
}
