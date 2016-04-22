
build: generate
	go build -tags "embed"

generate:
	yaber --pkg "assettemplates" templates/
	mv asset_* src/assettemplates/
	yaber --pkg "assetstatic" static/
	mv asset_* src/assetstatic/

test_run: generate
	go run ss13.go -verbose run

test_run_embed: generate
	go run -tags "embed" ss13.go -verbose run
