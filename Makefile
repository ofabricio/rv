
test:
	GOEXPERIMENT=jsonv2 go test -count=1 ./...

test-upd:
	GOEXPERIMENT=jsonv2 go run ./cmd/script/main.go
