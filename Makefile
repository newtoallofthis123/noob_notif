build:
	@go build -o ./bin/notif
run: build
	@./bin/notif
