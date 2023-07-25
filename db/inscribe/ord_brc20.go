package inscribe

import (
	"goBTC/db"

	"gorm.io/gorm"
)

type OrdToken struct {
	gorm.Model
	InscribeID   string    `gorm:"not null;uniqueIndex;column:inscribe_id"` // 铭文ID，唯一索引
	InscribeNum  string    `gorm:"not null;column:inscribe_number"`         // 铭文编号
	OwnerAddress string    `gorm:"not null;column:owner_address"`           // 铭文所有者地址
	Vout         string    `gorm:"column:vout"`                             // 流转位置
	Value        string    `gorm:"column:value"`                            // UTXO金额
	Tick         string    `gorm:"column:tick"`                             // 代币名称
	TokenType    string    `gorm:"not null;column:token_type"`              // 铭文类型：BRC20 :p
	Action       OrdAction `gorm:"not null;column:action"`                  // 铭文操作类型：deploy、mint、inscribeTransfer、transfer，联合索引 : op
	Amt          string    `gorm:"column:amt"`                              // BRC20操作数量
	Lim          string    `gorm:"column:lim"`                              // BRC20操作数量
	Supply       string    `gorm:"column:Supply"`                           // BRC20代币最大数量
	TxID         string    `gorm:"column:txId"`                             // 创建HASH
	BlockHeight  string    `gorm:"column:block_height"`                     // 区块高度
	BlockTime    int64     `gorm:"column:block_time"`                       // 交易区块时间
	State        string    `gorm:"column:state"`                            // 交易状态，可选值：success, invalid, pending
	SyncState    string    `gorm:"column:sync_state"`                       // 同步状态 0-未同步，1-已同步
}

func (this *OrdToken) TableName() string {
	return OrdTokenName
}

func (this *OrdToken) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *OrdToken) Create() error {
	return this.getDB().Create(this).Error
}
