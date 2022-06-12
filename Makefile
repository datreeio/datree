run: 
	go run -tags $(or $(datree_build_env),staging) -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0" main.go test ./internal/fixtures/**/*.yaml
run-production:
	make datree_build_env=main run
run-staging:
	make datree_build_env=staging run
run-dev:
	make datree_build_env=dev run

build:
	go build -tags $(or $(datree_build_env),staging) -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0"
build-production:
	make datree_build_env=main build
build-staging:
	make datree_build_env=staging build
build-dev:
	make datree_build_env=dev build

build-windows-amd:
	GOOS=windows GOARCH=amd64 go build -tags $(or $(datree_build_env),staging) -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0"

test:
	go test ./...

create-bin:
	goreleaser --snapshot --skip-publish --rm-dist

print-version:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0" main.go version

set-token:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0" main.go config set token testtoken

publish:
	go run -tags=staging -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=1.0.0" main.go publish ./internal/fixtures/policyAsCode/valid-schema.yaml

datree_test_to_JUnit_report: # install junit2html first: https://github.com/inorton/junit2html
	make build && ./datree test ./internal/fixtures/kube/skipRule/k8s-demo-skip-two.yaml -o JUnit > ./internal/fixtures/junit_support/datree_test_junit.xml || true && junit2html ./internal/fixtures/junit_support/datree_test_junit.xml ./internal/fixtures/junit_support/datree_test_junit_report.html
