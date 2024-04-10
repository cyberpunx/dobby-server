run:
	@templ generate
	@go build -o ./tmp/dobby-server.exe cmd/dobby-server/dobby-server.go
	@air