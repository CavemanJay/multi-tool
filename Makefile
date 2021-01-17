compile:
	go build -o build/server

listen: compile
	./build/server listen

help: compile
	./build/server -h

dial: compile
	./build/server dial