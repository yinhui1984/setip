APPNAME := "setip"

all:run

build:
	go build -o ./bin/$(APPNAME) main.go

clean:
	rm -rf ./bin/$(APPNAME)

run:build
	./bin/$(APPNAME)
