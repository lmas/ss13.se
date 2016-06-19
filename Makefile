
build: generate
	go build

generate:
	yaber -out src/assets.go templates/ static/

test_run:
	go run ss13.go -verbose run

test_run_embed:
	go run ss13.go -verbose run
