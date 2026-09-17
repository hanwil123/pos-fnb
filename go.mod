module github.com/yourorg/pos-fnb-backend

go 1.22.2

replace golang.org/x/sys => github.com/golang/sys v0.20.0

replace golang.org/x/crypto => github.com/golang/crypto v0.23.0

replace golang.org/x/text => github.com/golang/text v0.15.0

replace golang.org/x/net => github.com/golang/net v0.25.0

replace golang.org/x/sync => github.com/golang/sync v0.7.0

require (
	github.com/go-chi/chi/v5 v5.0.14
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

replace gorm.io/gorm => github.com/go-gorm/gorm v1.25.10

replace gorm.io/driver/postgres => github.com/go-gorm/postgres v1.5.9

replace gopkg.in/yaml.v3 => github.com/go-yaml/yaml/v3 v3.0.1

replace gopkg.in/check.v1 => github.com/go-check/check v0.0.0-20200902074654-038fdea0a05b
