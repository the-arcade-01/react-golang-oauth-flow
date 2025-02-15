module server

go 1.23.2

require (
	github.com/go-chi/chi/v5 v5.2.0
	github.com/go-chi/cors v1.2.1
	github.com/go-sql-driver/mysql v1.8.1
	golang.org/x/oauth2 v0.25.0
)

require github.com/cespare/xxhash/v2 v2.3.0 // indirect

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-chi/httprate v0.14.1
	github.com/joho/godotenv v1.5.1
)
