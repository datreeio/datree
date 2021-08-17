run: 
	go run -tags staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go test ./internal/fixtures/**/*.yaml

test:
	go test ./...

build:
	go build -tags staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1"

create-bin:
	goreleaser --snapshot --skip-publish --rm-dist

print-version:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go version

set-token:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go config set token testtoken

publish:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go publish --help #./internal/fixtures/kube/fail-30.yaml