.PHONY: e2e

e2e:
	@set -e; \
	docker compose -f compose.e2e.yaml up -d --build; \
	trap 'docker compose -f compose.e2e.yaml down -v' EXIT; \
	go test ./tests/e2e -v -count=1
