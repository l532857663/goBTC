package inscribe

import (
	"goBTC/db"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Brc20Assets struct {
	gorm.Model
	InscribeID    string          `gorm:"not null;uniqueIndex;column:inscribe_id"` // 铭文ID，唯一索引
	TokenType     string          `gorm:"not null;column:token_type"`              // 铭文类型：BRC20、NFT
	Tick          string          `gorm:"not null;column:tick"`                    // 铭文名称
	Supply        string          `gorm:"column:supply"`                           // 币种总发行量
	Lim           string          `gorm:"column:lim"`                              // BRC20操作数量
	Minted        decimal.Decimal `gorm:"column:minted"`                           // 已mint数量
	DeployAddr    string          `gorm:"column:deploy_addr"`                      // deploy地址
	DeployTime    string          `gorm:"column:deploy_time"`                      // deploy时间
	Holder        int             `gorm:"column:holder"`                           // 持有人数
	TransferTotal decimal.Decimal `gorm:"column:transferTotal"`                    // 交易次数
	State         string          `gorm:"column:state"`                            // 交易状态
}

func (this *Brc20Assets) TableName() string {
	return Brc20AssetsName
}

func (this *Brc20Assets) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *Brc20Assets) Create() error {
	return this.getDB().Create(this).Error
}
