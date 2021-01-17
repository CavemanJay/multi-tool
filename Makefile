
rootOut = build
linuxOut = $(rootOut)/server_linux
winOut = $(rootOut)/server_win.exe

compile:
	go build -o $(linuxOut)

listen: compile
	./$(linuxOut) listen

help: compile
	./$(linuxOut) -h

dial: compile
	./$(linuxOut) dial

all:
	GOOS=windows go build -o $(winOut)
	GOOS=linux go build -o $(linuxOut)

clean:
	rm -rvf $(rootOut)