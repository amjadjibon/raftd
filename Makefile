.PHONY: gen
gen:
	@echo "Starting..."
	@echo "command: buf generate"
	@buf generate
	@echo "Done!"

run:
	go run server/main.go
