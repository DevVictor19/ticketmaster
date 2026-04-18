.PHONY: run-dev run

run-dev:
	air

run:
	export $$(grep -v '^#' .env | xargs) && go run ./cmd/server/