module klimaguessr

go 1.24.1

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/zishang520/engine.io/v2 v2.5.0
	github.com/zishang520/socket.io/v2 v2.5.0
	golang.org/x/crypto v0.36.0
	gorm.io/driver/sqlite v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.25.10
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.53.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	github.com/zishang520/engine.io-go-parser v1.3.2 // indirect
	github.com/zishang520/socket.io-go-parser/v2 v2.5.0 // indirect
	github.com/zishang520/webtransport-go v0.9.1 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
)

replace golang.org/x/crypto => github.com/golang/crypto v0.24.0

replace gorm.io/gorm => github.com/go-gorm/gorm v1.25.10

replace gopkg.in/yaml.v2 => github.com/go-yaml/yaml/v2 v2.4.0

replace gopkg.in/check.v1 => github.com/go-check/check v0.0.0-20180628173108-788fd7840127

replace golang.org/x/sys => github.com/golang/sys v0.21.0

replace golang.org/x/net => github.com/golang/net v0.26.0

replace golang.org/x/text => github.com/golang/text v0.16.0

replace golang.org/x/sync => github.com/golang/sync v0.7.0

replace gorm.io/driver/sqlite => github.com/go-gorm/sqlite v1.5.6
