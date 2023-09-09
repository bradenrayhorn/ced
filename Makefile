.PHONY: run/server run/ui

run/server: export HTTP_PORT=8080
run/server: export DB_PATH=ced.db

run/server:
	(cd server && go run ./cmd/cedd)

run/ui: export PUBLIC_EVENT_TITLE=An Event
run/ui: export PUBLIC_EVENT_URL=https://example.com
run/ui: export PUBLIC_BASE_API_URL=

run/ui:
	(cd ui && npm run dev)

