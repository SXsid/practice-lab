.PHONY: build run clean test debug

CC = gcc
CFLAGS = -Wall -Wextra -Werror -pedantic -std=c11
SRC = ./server/src/main.c
BIN = ./bin/main

build:clean
	@echo "Building..."
	@mkdir -p bin
	$(CC) $(CFLAGS) $(SRC) -o $(BIN)

run: build
	@echo "Running..."
	@$(BIN)

debug:
	@echo "Building with AddressSanitizer..."
	@mkdir -p bin
	$(CC) $(CFLAGS) -g -fsanitize=address $(SRC) -o $(BIN)
	@echo "Running with sanitizer..."
	@$(BIN)

test:
	@echo "Running memory checks with Valgrind..."
	@valgrind --leak-check=full --track-origins=yes $(BIN)

clean:
	@echo "Cleaning..."
	@rm -rf bin
