package brc20_market

import (
	"goBTC/db"
	"time"

	"gorm.io/gorm"
)

type Rate struct {
	ID        uint `gorm:"primarykey"`
	UpdatedAt time.Time
	Pair      string `gorm:"column:pair;not null"`
	Price     string `gorm:"column:price;not null"`
}

func (this *Rate) TableName() string {
	return "rate"
}

func (this *Rate) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *Rate) SaveRate() error {
	data := Rate{}
	// 先尝试获取已存在的记录
	db := this.getDB()
	err := db.Where("pair = ?", this.Pair).Last(&data).Error
	if err != gorm.ErrRecordNotFound {
		// 如果找到记录，则更新它
		db.Model(&data).Updates(this)
		return nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		// 如果不是因为找不到记录而产生的错误，则返回该错误
		return err
	}

	// 如果没有找到匹配的记录，则创建新的记录
	return db.Create(this).Error
}
