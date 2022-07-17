
build:
	go build -o bin/splinter main.go

compile-dev:
	echo "Dev-Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/splinter-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/splinter-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/splinter-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/splinter-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/splinter-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -o bin/splinter-windows-arm64.exe main.go

compile-release:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build -o bin/splinter-${version}-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/splinter-${version}-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/splinter-${version}-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/splinter-${version}-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/splinter-${version}-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -o bin/splinter-${version}-windows-arm64.exe main.go

clean:
	echo "Cleaning Binaries"
	rm -rf bin/*