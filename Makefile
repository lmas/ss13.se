
build: generate
	go build -tags "embed"

generate:
	yaber --prefix="tmpl_asset_" --pkg "ss13" templates/
	mv tmpl_asset_* src/

test_run: generate
	go run ss13.go -verbose run

test_run_embed: generate
	go run -tags "embed" ss13.go -verbose run
