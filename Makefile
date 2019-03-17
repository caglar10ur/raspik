DIRS = sub pub graph
all: format
	@for dir in $(DIRS); do (cd $$dir; make || exit 1) || exit 1; done
clean:
	@for dir in $(DIRS); do (cd $$dir; make clean || exit 1) || exit 1; done
format:
	gofmt -s -w *.go */*.go
