# PC Core Backend

### About
`PC Core` is a store for computer hardware and related products. This repository contains the entire source code for the backend part of the `PC Core` website

### Requirenments
- PostgreSQL
- Redis
- Go
- MinIO
- Required Go packages:
    - `gin` - `go get github.com/gin-gonic/gin`
    - `gorm` - `go get gorm.io/gorm`
    - `gorm-postgres` - `go get gorm.io/driver/postgres`
    - `gin-cors` - `go get github.com/gin-contrib/cors`
    - `gin-swagger` - `go get github.com/swaggo/gin-swagger`
    - `swagger-files` - `go get github.com/swaggo/files`
    - `go-redis` - `go get github.com/redis/go-redis/v9`
    - `godotenv` - `go get github.com/joho/godotenv`
    - `swag` - `go get github.com/swaggo/swag`
    - `jwt` - `go get github.com/golang-jwt/jwt/v5`
    - `headers` - `go get -u github.com/go-http-utils/headers`
    - `minio` - `go get github.com/minio/minio-go/v7`

### How to run
1. Create the Postgres database
2. Migrate the database (sql/migrations)
3. Setup the [environment](#env-variables)
4. Initialize Redis by starting `init_redis.sh` script
5. Start the /cmd/pccore/main.go file with [arguments](#cli-arguments)

### ENV Variables
- `PCCORE_POSTGRES_CONN` - Postgres connection string
- `PCCORE_JWT_KEY` - JWT secret key file path
- `PCCORE_REDIS_PASSWORD` - Redis password
- `MINIO_ACCESS` - MinIO login
- `MINIO_SECRET` - MinIO password

### CLI Arguments
- `--working-dir` - The directory containing the config files. The default value is './'

### Swagger
To open the Swagger page:
1. Start the server in debug mode 
2. Go to localhost:{port}/swagger/index.html

### Seeds
To setup seeds, run the `seeds.go` file. If no flags are provided, a list of available options will be displayed. The execution order follows the order in which the flags are passed.
- `go run seeds.go all` - download all

- `go run seeds.go media` - download all media

- `go run seeds.go laptop` - download all laptops