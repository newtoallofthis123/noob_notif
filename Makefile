build:
	@go build -o ./bin/notif
run: build
	@./bin/notif
css:
	bun tailwindcss -i ./static/input.css -o ./static/output.css --watch
