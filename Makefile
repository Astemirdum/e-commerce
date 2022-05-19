
run:
	docker-compose up --build

gen-proto:
	buf build && buf generate proto

test:
	 go test -v ./...

.PHONY: run, gen-proto, test