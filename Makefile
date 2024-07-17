BINARY_NAME=main.out
MODULE_PATH=github.com/venturarome/DaftWatch

all: build test

build:
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm -f ${BINARY_NAME}

# install:
# 	go install ${MODULE_PATH}

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

test:
	go test -v main.go

tidy:
	go mod tidy
