
default:
	go build

test:
	go test ./...

cert:
	./estate generatecert

run:
	./estate api --https

