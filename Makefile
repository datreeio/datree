run: 
	go run -tags staging main.go test "./internal/fixtures/**/*.yaml"

test:
	go test ./...

create-bin:
	goreleaser --snapshot --skip-publish --rm-dist
