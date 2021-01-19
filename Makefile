
rootOut = build
linuxOut = $(rootOut)/gogurt
winOut = $(rootOut)/gogurt.exe
releaseFlags = -s -w
buildCmd = go build

compile:
	$(buildCmd) -o $(linuxOut)

listen: compile
	./$(linuxOut) listen

wsl: compile
	./$(linuxOut) listen -f "/mnt/c/Users/cueva/Sync/"

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
	xgo -v -x -tags='release' -ldflags='-s -w' -dest ./$(rootOut) --targets=windows/*,linux/amd64,linux/386 github.com/JayCuevas/gogurt

clean:
	rm -rvf $(rootOut)