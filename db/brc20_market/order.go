package brc20_market

import (
	"goBTC/db"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ConfigID        int     `gorm:"column:config_id;not null"`
	TxHash          string  `gorm:"column:tx_hash;type:varchar(255);not null"`
	InscribeID      *string `gorm:"column:inscribe_id;type:varchar(255);default:null"`
	InscribeNumber  *string `gorm:"column:inscribe_number;type:varchar(255);default:null"`
	InscribeContent string  `gorm:"column:inscribe_content;type:varchar(500);default:null"`
	ContentType     *string `gorm:"column:content_type;type:varchar(255);default:null"`
	BlockHeight     *int64  `gorm:"column:block_height"`
	Tick            string  `gorm:"column:tick;type:varchar(255);not null"`
	Side            int     `gorm:"column:side;not null;default:1"`
	State           int     `gorm:"column:state;not null;default:1"`
	Number          int64   `gorm:"column:number;not null"`
	Amount          int64   `gorm:"column:amount;not null"`
	TotalAmount     int64   `gorm:"column:total_amount;not null"`
	ServerFee       int64   `gorm:"column:server_fee;not null;default:0"`
	GasFee          *int64  `gorm:"column:gas_fee"`
	GasFeeTotal     *int64  `gorm:"column:gas_fee_total"`
	From            *string `gorm:"column:from;collate:utf8mb4_general_ci"`
	To              *string `gorm:"column:to;collate:utf8mb4_general_ci"`
}

func (this *Order) TableName() string {
	return "order"
}

func (this *Order) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *Order) Create() error {
	return this.getDB().Create(this).Error
}

func (this *Order) GetPendingOrder() ([]*Order, error) {
	var list []*Order
	err := this.getDB().
		Where("state in (1,4,6)").
		Find(&list).Error
	return list, err
}

func (this *Order) UpdatePendingOrderState() (int64, error) {
	result := this.getDB().
		Where("id = ?", this.ID).
		Update("state", this.State)
	return result.RowsAffected, result.Error
}
