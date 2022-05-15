
run:
	docker-compose up --build

gen-proto:
	buf build && buf generate proto

.PHONY: run, gen-proto