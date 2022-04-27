start-test-env:
	test -f docker-compose.yml || wget -O  docker-compose.yml https://raw.githubusercontent.com/posteris/ci-database-starter/main/docker-compose.yml
	docker-compose up

bench:
	go test -bench 'Benchmark' ./...

test:
	go test ./...

