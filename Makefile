# build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o app cmd/api/main.go

# run tests
test:
	go test ./providers/mysql
	go test ./service

docker-build:
	docker build -t fb-tinder-app .

# dev commands
dev-build:
	go build -o app cmd/api/main.go

compose-up: dev-build
	docker-compose up -d

compose-restart-web: dev-build
	docker-compose restart web
