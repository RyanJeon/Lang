GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GCC= gcc
BINARY= test.out
UTIL= ./Util/

compile:
	$(GORUN) assembler.go dsutil.go interpret.go lang.go parser.go token.go asmutil.go
	$(GCC) test.asm -o test.out -e execute

build:
	$(GOBUILD) -o rlang assembler.go dsutil.go interpret.go lang.go parser.go token.go asmutil.go