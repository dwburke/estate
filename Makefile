
default:
	go build

test:
	go test ./...

cert:
	./prefs generatecert

run:
	./prefs api --https

