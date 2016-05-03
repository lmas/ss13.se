
build: generate
	go build -tags "embed"

generate:
	cd src/assettemplates/  && yaber --pkg "assettemplates" -strip "../../"  ../../templates/
	cd src/assetstatic && yaber --pkg "assetstatic" -strip "../../" ../../static/

test_run:
	go run ss13.go -verbose run

test_run_embed:
	go run -tags "embed" ss13.go -verbose run
