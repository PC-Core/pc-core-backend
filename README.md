# CM-BACKEND

### Requirenments
- PostgreSQL
- Go
- gin (go get github.com/gin-gonic/gin)
- Postgres driver for Go (go get github.com/lib/pq)
- gin-cors (go get github.com/gin-contrib/cors)
- gin-swagger (go get github.com/swaggo/gin-swagger)
- swagger-files (go get github.com/swaggo/files)

### How to run

1. Create the Postgres database
2. Migrate the database (sql/migrations)
3. Create an ENV variable named 'POSTGRES_IBYTE_CONN' with Postgres Connection String
4. Generate a secret key for JWT in file
5. Create an ENV variable named 'PCCORE_JWT_KEY' and put the path to the secret key file in it
6. Start the /cmd/pccore/main.go file

### Swagger
To open the Swagger page, start the server in debug mode and go to localhost:{port}/swagger/index.html