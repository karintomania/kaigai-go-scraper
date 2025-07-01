.PHONY: test fmt

test:
	$(eval COMMAND := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS)))
	@if [ "$(COMMAND)" ]; then \
		godotenv -f ./.testing.env go test ./$(COMMAND); \
	else \
		godotenv -f ./.testing.env go test ./... -v; \
	fi

fmt:
	go fmt ./...
	golangci-lint run

server:
	./kaigai-go-scraper -mode server

scrape:
	./kaigai-go-scraper -run-all
