BIN_DIR := ./bin

build:
	go build -o $(BIN_DIR)/awsctl

clean:
	rm -rf $(BIN_DIR)

.PHONY: all clean 