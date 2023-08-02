package elastic

type ElasticConfig struct {
	Host     string `mapstructure:"host"     json:"host"     yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type InscribeInfo struct {
	InscribeId      string      `json:"inscribe_id"`
	InscribeNumber  int64       `json:"inscribe_number"`
	InscribeContent string      `json:"inscribe_content"`
	InscribeType    string      `json:"inscribe_type"`
	TxHash          string      `json:"tx_hash"`
	ContentType     string      `json:"content_type"`
	OwnerAddress    string      `json:"owner_address"`
	OwnerOutput     string      `json:"owner_output"`
	GenesisAddress  string      `json:"genesis_address"`
	GenesisOutput   string      `json:"genesis_output"`
	BlockHeight     int64       `json:"block_height"`
	State           string      `json:"state"`
	SyncState       string      `json:"sync_state"`
	SatsName        string      `json:"sats_name,omitempty"`
	Brc20Info       interface{} `json:"brc20_info,omitempty"`
}

type OrdToken struct {
	InscribeId   string      `json:"inscribe_id"`
	TxHash       string      `json:"tx_hash"`
	OwnerAddress string      `json:"owner_address"`
	OwnerOutput  string      `json:"owner_output"`
	Value        string      `json:"value"`
	Action       string      `json:"action"`
	Brc20Info    interface{} `json:"brc20_info"`
	BlockHeight  int64       `json:"block_height"`
	BlockTime    int64       `json:"block_time"`
	State        string      `json:"state"`
	SyncState    string      `json:"sync_state"`
}

type ActivityInfo struct {
	InscribeId   string `json:"inscribe_id"`
	InscribeType string `json:"inscribe_type"`
	TxHash       string `json:"tx_hash"`
	ActivityType string `json:"activity_type"`
	From         string `json:"from"`
	To           string `json:"to"`
	BlockTime    int64  `json:"block_time"`
}

type UserUTXO struct {
	InscribeId string `json:"inscribe_id"`
	Tick       string `json:"tick"`
	Address    string `json:"address"`
	vout       string `json:"vout"`
	State      string `json:"state"`
}

type DeployToken struct {
	InscribeId string `json:"inscribe_id"`
	TxHash     string `json:"tx_hash"`
	Tick       string `json:"tick"`
	Lim        string `json:"lim"`
	Max        string `json:"max"`
	DeployAddr string `json:"deploy_addr"`
	DeployTime int64  `json:"deploy_time"`
	Minted     string `json:"minted"`
	State      string `json:"state"`
}
