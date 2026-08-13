.PHONY: build build-mac build-linux build-windows clean test

DIST := .dist

build:
	mkdir -p $(DIST)
	go build -o $(DIST)/backstage-vikunja .

build-mac:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/backstage-vikunja .

build-linux:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/backstage-vikunja-linux-amd64 .

build-windows:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/backstage-vikunja.exe .

clean:
	rm -rf $(DIST)

test:
	go test ./...
