module github.com/noah-blockchain/noah-explorer-api

go 1.12

replace mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1

require (
	github.com/centrifugal/centrifuge-go v0.2.3
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jpillora/backoff v0.0.0-20180909062703-3050d21c67d7 // indirect
	github.com/noah-blockchain/noah-explorer-tools v0.1.1
	github.com/noah-blockchain/noah-go-node v0.1.2
	github.com/spf13/viper v1.4.0
	gopkg.in/go-playground/validator.v8 v8.18.2
)
