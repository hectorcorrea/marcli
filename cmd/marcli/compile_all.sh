GOOS=darwin GOARCH=amd64 go build -o marcli
GOOS=linux GOARCH=amd64 go build -o marcli_linux
GOOS=windows GOARCH=386 go build -o marcli.exe
