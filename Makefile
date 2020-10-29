.PHONY: build

all: telnet-chat-go telnet-chat-c

check-dir:
	@if [ ! -d "build" ]; then mkdir build; fi

telnet-chat-go: check-dir
	@go build -o build/telnet-chat-go go/*.go
	@echo "telnet-chat-go build success."

telnet-chat-c: check-dir
	@cc c/main.c -o build/telnet-chat-c
	@echo "telnet-chat-c build success."
