lint:
	# | jq > ./golangci-lint/report.json
	golangci-lint run --fix -c .golangci.yml > golangci-lint/report-unformatted.json
	goimports -local mmskazak -w .

lint-clean:
	sudo rm -rf ./golangci-lint

test:
	go test ./...

proto:
	protoc --proto_path=internal/proto \
        --go_out=internal/proto \
        --go_opt=paths=source_relative \
        --go-grpc_out=internal/proto \
        --go-grpc_opt=paths=source_relative \
        internal/proto/delete_user_urls_request.proto \
        internal/proto/delete_user_urls_response.proto \
        internal/proto/find_user_urls_request.proto \
        internal/proto/find_user_urls_response.proto \
        internal/proto/internal_stats_request.proto \
        internal/proto/internal_stats_response.proto \
        internal/proto/save_shorten_url_batch_request.proto \
        internal/proto/save_shorten_url_batch_response.proto \
        internal/proto/handle_create_short_url_request.proto \
        internal/proto/handle_create_short_url_response.proto \
        internal/proto/shorturl.proto


# Параметры контейнера и образа
CONTAINER_NAME=my_postgres
IMAGE=postgres:16.3
POSTGRES_USER=pguser
POSTGRES_PASSWORD=pgpassword
POSTGRES_DB=dbshorturl
VOLUME_NAME=postgres_data

# Команда для запуска контейнера PostgreSQL
db:
	docker run -d \
        --name $(CONTAINER_NAME) \
        -e POSTGRES_USER=$(POSTGRES_USER) \
        -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
        -e POSTGRES_DB=$(POSTGRES_DB) \
        -p 5432:5432 \
        -v $(VOLUME_NAME):/var/lib/postgresql/data \
        $(IMAGE)
