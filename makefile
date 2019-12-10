GO111MODULE=on

DIRECTORIES= externals
MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))

.PHONY: test test_race lint vet get mocks clean-mocks 
.ONESHELL:

test:
	go test ./...

test_race:
	go test ./... -race 

lint:
	golint $(go list ./... | grep -v mocks)

vet:
	go vet $(go list ./... | grep -v mocks)

get:
	go get ./...

clean-mocks:
	rm -rf mocks

mocks: $(MOCKS) 
	mockery -output=./mocks -dir=./ -all
	
$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all
