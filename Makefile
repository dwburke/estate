
default:
	go build

test:
	go test ./...

cert:
	./lode generatecert

run:
	./lode api --https

