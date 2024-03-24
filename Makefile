BUILD_FLAG = -ldflags "-s -w"
BUILD_DIR = bin

default: build

build:
	env CGO_ENABLED=0  GOOS=windows GOARCH=amd64 go build $(BUILD_FLAG) -o $(BUILD_DIR)/Serverless_PortScan.exe main/main.go
	env CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build $(BUILD_FLAG) -o $(BUILD_DIR)/Serverless_PortScan main/main.go

clean:
	rm -rf ./$(BUILD_FLAG)