migrate reference data - 
go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
migrate create -ext mongodb -seq -digits 5 product_collection
migrate -path ./migration/ -database mongodb://localhost:27017/gymondodb up 1

- generate mock for testing
go generate ./...