DIRS = sub pub graph
all:
	for dir in $(DIRS); do (cd $$dir; make $1 || exit 1) || exit 1; done
format:
	gofmt -s -w *.go */*.go
