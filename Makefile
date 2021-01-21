
appName=gogurt
rootOut = build
linuxOut = $(rootOut)/$(appName)
winOut = $(rootOut)/$(appName).exe
releaseFlags = -s -w
buildCmd = go build

compile:
	$(buildCmd) -o $(linuxOut)

listen: compile
	./$(linuxOut) listen

listen-music: compile
	./$(linuxOut) -a "Music" listen 

wsl: compile
	./$(linuxOut) -f "/mnt/c/Users/cueva/Sync/" listen 

help: compile
	./$(linuxOut) -h

dial: compile
	./$(linuxOut) dial

all: 
	GOOS=windows $(buildCmd) -o $(winOut)
	GOOS=linux $(buildCmd) -o $(linuxOut)

release:
	GOOS=windows $(buildCmd) -o $(winOut) -tags release -ldflags="$(releaseFlags)"
	GOOS=linux $(buildCmd) -o $(linuxOut) -tags release -ldflags="$(releaseFlags)"

xgo:
	mkdir -p $(rootOut)
	xgo -v -x -tags='release' -ldflags='-s -w' -dest ./$(rootOut) --targets=windows/*,linux/amd64,linux/386 github.com/JayCuevas/$(appName)

clean:
	rm -rvf $(rootOut)

install:
	go install .

uninstall:
	rm -rvf "${GOPATH}/bin/$(appName)"