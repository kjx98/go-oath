#
#	Makefile for oath
#
# switches:
#	define the ones you want in the CFLAGS definition...
#
#	TRACE		- turn on tracing/debugging code
#
#
#
#

# Version for distribution
VER=1_0r1

MAKEFILE=GNUmakefile

# We Use Compact Memory Model

all: bin/authK
	@[ -d bin ] || exit

win64: bin/authK64.exe
	@[ -d bin ] || exit

bin/authK:	cmd/authK/readAcct.go cmd/authK/main.go
	@[ -d bin ] || mkdir bin
	@go build -o $@ ./cmd/authK
	@strip $@ || echo "authK OK"

bin/authK64.exe:	cmd/authK/readAcct.go cmd/authK/main.go
	@[ -d bin ] || mkdir bin
	#(. ./mingw64-env.sh; go build -o $@ ./cmd/authK)
	GOOS=windows GOARCH=amd64 go build -o $@ ./cmd/authK
	@echo "authK64.exe OK"

test:
	@go test -v

clean:
	@rm -f bin/*

distclean: clean
	@rm -rf bin
