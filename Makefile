GOCMD=go
GORUN=$(GOCMD) run
GCC= gcc
BINARY= test.out

build:
	$(GORUN) assembler.go dsutil.go interpret.go lang.go parser.go token.go
	$(GCC) test.asm -o test.out -e execute