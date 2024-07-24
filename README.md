# DaftWatch

### Summary
Create alerts in daft.ie matching a specific search criteria and receive live updates when a new property ad is added or deleted.

## Info about Go structure
* https://github.com/golang-standards/project-layout
* https://www.youtube.com/watch?v=1ZbQS6pOlSQ

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
