run: 
	go run -tags $(or $(datree_build_env),staging) -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go test ./internal/fixtures/**/*.yaml
run-production:
	make datree_build_env=main run
run-staging:
	make datree_build_env=staging run
run-dev:
	make datree_build_env=dev run

build:
	go build -tags $(or $(datree_build_env),staging) -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1"
build-production:
	make datree_build_env=main build
build-staging:
	make datree_build_env=staging build
build-dev:
	make datree_build_env=dev build

test:
	go test ./...

create-bin:
	goreleaser --snapshot --skip-publish --rm-dist

print-version:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go version

set-token:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go config set token testtoken

publish:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.0.1" main.go publish ./internal/fixtures/policyAsCode/valid-schema.yaml

junit_report:
	junit2html ./internal/fixtures/junit_support/junit_datree.xml ./internal/fixtures/junit_support/junit_report.html
