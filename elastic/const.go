package elastic

const (
	// http method
	HttpPost   = "POST"
	HttpPut    = "PUT"
	HttpDelete = "DELETE"

	// Index分类
	InscribeInfoType = "inscription"
	OrdTokenType     = "ordtoken"
	ActivityType     = "activity"
	DeployType       = "deploy"

	// 铭文类型
	InscribeTypeNFT   = "NFT"
	InscribeTypeBRC20 = "BRC20"

	// 铭文信息同步状态
	StateSyncIsFalse = "0"
	StateSyncIsTrue  = "1"

	// 铭文状态
	InscriptionStateSuccess = "success"
	InscriptionStateInvalid = "invalid"
	InscriptionStatePending = "pending"

	// Activity type
	ActivityTypeInscribed = "Inscribed"
	ActivityTypeTransfer  = "Transfer"

	// Activity action
	ActivityActionNew      = "New inscription"
	ActivityActionReceived = "Received"
	ActivityActionSent     = "Sent"
)
