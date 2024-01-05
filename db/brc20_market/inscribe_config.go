package brc20_market

import (
	"goBTC/db"
	"time"

	"gorm.io/gorm"
)

type InscribeConfig struct {
	gorm.Model
	ProtocolName  string    `gorm:"column:protocol_name;type:varchar(255);not null"`
	InscribeType  string    `gorm:"column:inscribe_type;type:varchar(255);not null"`
	InscribeID    *string   `gorm:"column:inscribe_id;type:varchar(255);default:null"`
	Tick          string    `gorm:"column:tick;type:varchar(255);not null"`
	DeployNumber  int64     `gorm:"column:deploy_number;not null;default:0"`
	ConfirmMinted int64     `gorm:"column:confirm_minted;not null;default:0"`
	LimitMint     int64     `gorm:"column:limit_mint;not null;default:0"`
	Rank          *int      `gorm:"column:rank"`
	Enable        bool      `gorm:"column:enable;not null;default:true"`       // 将tinyint映射为bool类型，1为true, 2为false
	MintFinish    bool      `gorm:"column:mint_finish;not null;default:false"` // 同上，将mint_finish字段映射为bool类型
	ImageURL      *string   `gorm:"column:image_url;type:varchar(255);default:null"`
	CreateTime    time.Time `gorm:"column:create_time;default:current_timestamp"`
	UpdateTime    time.Time `gorm:"column:update_time;default:current_timestamp on update current_timestamp"`
}

func (this *InscribeConfig) TableName() string {
	return "inscribe_config"
}

func (this *InscribeConfig) getDB() *gorm.DB {
	return db.GetDBByName(this.TableName())
}

func (this *InscribeConfig) Create() error {
	return this.getDB().Create(this).Error
}
