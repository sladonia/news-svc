VERSION := $(shell git describe --tags)

test:
	go test -p 1 -count 1 ./...

build:
ifeq ($(VERSION),)
	$(eval VERSION := "dev")
endif
	docker build -t $(VERSION) .

up:
	docker-compose up -d

.PHONY: test build up
