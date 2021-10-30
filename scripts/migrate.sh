export POSTGRESQL_URL=postgres://learn-go-database:password@localhost:5555/learn-go-database?sslmode=disable

migrate -database ${POSTGRESQL_URL} -path db/migrations $1 $2
