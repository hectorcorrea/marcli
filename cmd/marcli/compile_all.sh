GOOS=darwin go build -o marcli
GOOS=linux go build -o marcli_linux
GOOS=windows GOARCH=386 go build -o marcli.exe
