package inscribe

import (
	"goBTC/db"

	"gorm.io/gorm"
)

type InscribeInfo struct {
	gorm.Model
	InscribeID   string    `gorm:"not null;uniqueIndex;column:inscribe_id"`       // 铭文ID，唯一索引
	InscribeNum  string    `gorm:"not null;column:inscribe_number"`               // 铭文编号
	TokenType    string    `gorm:"not null;column:token_type"`                    // 铭文类型：BRC20、NFT
	Action       OrdAction `gorm:"not null;index:idx_token_action;column:action"` // 铭文操作类型：deploy、mint、inscribeTransfer、transfer，联合索引
	LogoURL      string    `gorm:"not null;column:logoUrl"`                       // 铭文链接
	OwnerAddress string    `gorm:"not null;column:owner_address"`                 // 铭文所有者地址
	CreateAddr   string    `gorm:"column:create_addr"`                            // 创建地址
	Vout         string    `gorm:"column:vout"`                                   // 流转位置
	TxID         string    `gorm:"column:txId"`                                   // 最新流转HASH
	Content      string    `gorm:"not null;column:content"`                       // 铭文状态：0-不可用、1-可用
	ContentType  string    `gorm:"column:content_type"`                           // 存储信息类型
	ContentSize  int64     `gorm:"column:content_size"`                           // 存储信息大小
	BlockHeight  string    `gorm:"column:block_height"`                           // 最新区块高度
	BlockTime    int64     `gorm:"column:block_time"`                             // 交易区块时间
}

func (this *InscribeInfo) TableName() string {
	return InscribeInfoName
}

func (this *InscribeInfo) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *InscribeInfo) Create() error {
	return this.getDB().Create(this).Error
}
