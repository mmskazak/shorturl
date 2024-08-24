lint:
	# | jq > ./golangci-lint/report.json
	golangci-lint run --fix -c .golangci.yml > golangci-lint/report-unformatted.json
	goimports -local mmskazak -w .

lint-clean:
	sudo rm -rf ./golangci-lint

test:
	go test ./...

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
