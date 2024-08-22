build_linux:
	go build -o server cmd/server/main.go
	go build -o migrator cmd/migrator/migrator.go
	go build -o apply_dump cmd/apply_dump/apply_dump.go