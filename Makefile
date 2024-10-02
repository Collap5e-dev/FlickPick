build_windows:
	go build -o bin/server.exe cmd/server/main.go
	go build -o bin/migrator.exe cmd/migrator/migrator.go
	go build -o bin/apply_dump.exe cmd/apply_dump/apply_dump.go

build_macos:
	go build -o bin/server cmd/server/main.go
	go build -o bin/migrator cmd/migrator/migrator.go
	go build -o bin/apply_dump cmd/apply_dump/apply_dump.go

