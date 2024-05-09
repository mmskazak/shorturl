lint:
	golangci-lint run -c .golangci.yml > golangci-lint/report-unformatted.json # | jq > ./golangci-lint/report.json

lint-clean:
	sudo rm -rf ./golangci-lint
