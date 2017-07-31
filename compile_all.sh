GOOS=darwin go build -o marcli_mac
GOOS=linux go build -o marcli_linux
GOOS=windows GOARCH=386 go build -o marcli.exe
