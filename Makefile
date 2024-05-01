run:
	@templ generate
	@go build -o ./tmp/dobby-server.exe cmd/dobby-server/dobby-server.go
	@air

gobs:
	@go build -o ./tmp/gobs/gobs.exe cmd/gobs-client/gobs-client.go