<!-- create migration file -->
goose -dir ./migrations create orders_table sql

<!-- run migration -->

goose -dir ./migrations sqlite3 ./orders.db up 

<!-- run app -->
go run ./main.go

