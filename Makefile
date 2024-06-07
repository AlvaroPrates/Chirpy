build:
	go build -o ./bin/chirpy 

run: build
	./bin/chirpy --debug