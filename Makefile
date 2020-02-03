build-api:
	go build -o build/hoodiee
	cp -rf definitions build/

build: clean
	mkdir -p build build/client
	make build-api
	cd client && npm run build
	mv -f client/build/* build/client/

clean:
	rm -rf build client/build
