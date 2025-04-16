.PHONY: test

test:
	$(eval COMMAND := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS)))
	@if [ "$(COMMAND)" ]; then \
		godotenv -f ./.testing.env go test ./$(COMMAND); \
	else \
		godotenv -f ./.testing.env go test ./... -v; \
	fi

