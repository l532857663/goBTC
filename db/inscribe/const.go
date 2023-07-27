package inscribe

const (
	// 表名
	InscribeInfoName    = "inscribe_info"
	OrdTokenName        = "ord_tokens"
	Brc20AssetsName     = "brc20_assets"
	UserBrc20AssetsName = "user_brc20_assets"
)

type OrdAction string

const (
	// 铭文信息同步状态
	StateSyncIsFalse = "0"
	StateSyncIsTrue  = "1"
	// 铭文操作类型
	ActionForDeploy   OrdAction = "deploy"
	ActionForMint     OrdAction = "mint"
	ActionForTransfer OrdAction = "transfer"
	ActionForSend     OrdAction = "send"
	ActionForReceive  OrdAction = "receive"
	// 铭文类型
	InscribeTypeNFT   = "NFT"
	InscribeTypeBRC20 = "BRC20"
	// 铭文状态
	Brc20StateSuccess = "success"
	Brc20StateInvalid = "invalid"
	Brc20StatePending = "pending"
)
