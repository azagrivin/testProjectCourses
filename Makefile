easyjson:
	rm -rf internal/models/*easyjson.go
	easyjson --gen_build_flags="-mod=mod" internal/models/*

migrations:
	docker-compose up migrations

run: migrations
	go run cmd/app/main.go