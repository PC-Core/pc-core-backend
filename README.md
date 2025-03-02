# PC Core Backend

### About
`PC Core` is a store for computer hardware and related products. This repository contains the entire source code for the backend part of the `PC Core` website

### Requirenments
- PostgreSQL
- Redis
- Go
- Required Go packages:
    - `gin` - `go get github.com/gin-gonic/gin`
    - `lib/pq` - `go get github.com/lib/pq`
    - `gin-cors` - `go get github.com/gin-contrib/cors`
    - `gin-swagger` - `go get github.com/swaggo/gin-swagger`
    - `swagger-files` - `go get github.com/swaggo/files`
    - `go-redis` - `go get github.com/redis/go-redis/v9`
    - `godotenv` - `go get github.com/joho/godotenv`
    - `swag` - `go get github.com/swaggo/swag`
    - `jwt` - `go get github.com/golang-jwt/jwt/v5`

### How to run
1. Create the Postgres database
2. Migrate the database (sql/migrations)
3. Create an ENV variable named 'PCCORE_POSTGRES_CONN' with Postgres Connection String
4. Generate a secret key for JWT in file
5. Create an ENV variable named 'PCCORE_JWT_KEY' and put the path to the secret key file in it
6. Initialize Redis by starting `init_redis.sh` script
7. Create an ENV variable named 'PCCORE_REDIS_PASSWORD' and put the Redis password in it
8. Create an ENV variable named 'CFG_PATH' and put the absolute path to the `cfg.yml` file in it
9. Start the /cmd/pccore/main.go file

### Swagger
To open the Swagger page:
1. Start the server in debug mode 
2. Go to localhost:{port}/swagger/index.html

### Seeds
To setup seeds, run the `seeds.go` file. If no flags are provided, a list of available options will be displayed. The execution order follows the order in which the flags are passed.