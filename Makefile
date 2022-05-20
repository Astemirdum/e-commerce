
run:
	docker-compose up --build

gen-proto:
	buf build && buf generate proto

test:
	 go test -v ./...

mocks:
	#auth
	cd ./store-auth/internal/repo/mocks/; go generate;
	cd ./store-auth/internal/service/jwtoken/mocks/; go generate;
	#order
	cd ./store-order/internal/repo/mocks/; go generate;
	cd ./store-order/client/mocks/; go generate;
	#product
	cd ./store-product/internal/repo/mocks/; go generate;

.PHONY: run, gen-proto, test, mocks