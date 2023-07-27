package inscribe

import (
	"goBTC/db"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserBrc20Assets struct {
	gorm.Model
	Address      string          `gorm:"not null;column:address"`     // 铭文所有者地址
	InscribeID   string          `gorm:"not null;column:inscribe_id"` // 铭文ID，唯一索引
	TokenType    string          `gorm:"not null;column:token_type"`  // 铭文类型：BRC20、NFT
	Tick         string          `gorm:"not null;column:tick"`        // 铭文名称
	Balance      decimal.Decimal `gorm:"column:balance"`              // 用户余额
	Available    decimal.Decimal `gorm:"column:available"`            // 可用余额
	Transferable decimal.Decimal `gorm:"column:transferable"`         // 已铸造交易额
}

func (this *UserBrc20Assets) TableName() string {
	return UserBrc20AssetsName
}

func (this *UserBrc20Assets) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *UserBrc20Assets) Create() error {
	return this.getDB().Create(this).Error
}
