.PHONY: test

test:
	$(eval COMMAND := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS)))
	if [ "$(COMMAND)" ]; then \
		godotenv -f ./.env.testing go test ./$(COMMAND); \
	else \
		godotenv -f ./.env.testing go test ./... -v; \
	fi

