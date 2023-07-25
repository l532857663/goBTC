module goBTC

go 1.19

require (
	github.com/btcsuite/btcd v0.23.4
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/minchenzz/brc20tool v0.0.0-20230430042556-a383f1b5d34a
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.3.1
	go.uber.org/zap v1.24.0
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.1
)

require (
	github.com/aead/siphash v1.0.1 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/kkdai/bstream v1.0.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)

replace github.com/btcsuite/btcd => github.com/qhxcWallet/btcd v0.0.0-20221101160504-9f81f6fbe13b
