
appName=gogurt
rootOut = build
linuxOut = $(rootOut)/$(appName)
winOut = $(rootOut)/$(appName).exe
releaseFlags = -s -w
buildCmd = go build
version = 0.1.0
buildFlags = -X main.version=$(version)
releasePath = "release/$(version)"

compile:
	$(buildCmd) -ldflags="$(buildFlags)" -o $(linuxOut)

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

xgo:
	mkdir -p $(releasePath)
	xgo -v -x -tags='release' -ldflags='$(releaseFlags) $(buildFlags)' -dest ./$(releasePath) --targets=windows/*,linux/amd64,linux/386 github.com/CavemanJay/$(appName)

clean:
	rm -rvf $(rootOut) build data $(releasePath)/*

install:
	go install .

uninstall:
	rm -rvf "${GOPATH}/bin/$(appName)"

release: clean xgo
	find ./ -name "*.go" -o -name "go.*" -o -name "*.yml" | tar -cvf $(releasePath)/$(appName)-$(version).tar.gz -T -
