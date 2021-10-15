
build:
	@go run generate.go && mjml generated.mjml -o generated.html

preview: build
	open generated.html
