build: static
	go build gallery.go

deps:
	go get github.com/rakyll/statik

static:
	statik -src=assets

clear:
	rm -rf statik
