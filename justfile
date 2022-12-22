build:
	go build -v -o madge

clean:
	go clean
	rm *.exe
	rm go_build*
