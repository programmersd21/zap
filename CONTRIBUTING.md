# contributing to zap

thanks for considering contributing to zap! this document outlines the process and guidelines.

## code of conduct

be respectful and constructive. we're all here to build better tools.

## how to contribute

### reporting bugs

before creating a bug report:
- check if the issue already exists
- use the latest version of zap
- include reproduction steps and system info (os, go version)

### suggesting features

feature requests are welcome! please:
- check existing issues first
- clearly describe the use case
- explain why it fits zap's scope (see non-goals in readme)

### pull requests

1. **fork and clone** the repository
2. **create a branch** for your feature (`git checkout -b feature/amazing-feature`)
3. **make your changes** and add tests
4. **run tests** (`go test ./...`)
5. **build** to ensure it compiles (`go build ./cmd/zap`)
6. **commit** with clear messages
7. **push** to your fork and create a pr

### code style

- follow standard go conventions (`gofmt`, `go vet`)
- write tests for new functionality
- keep functions focused and documented
- avoid external dependencies unless necessary

### testing

all code changes should include tests:
- unit tests for logic (`internal/` packages)
- integration tests for operations (`internal/ops/`)
- ui tests using `teatest` for bubble tea components

run tests:
```bash
go test ./...
go test -race ./...  # check for race conditions
go test -cover ./... # check coverage
```

## ai-generated code

if you use ai assistance (copilot, chatgpt, etc.):
- mark prs with `[ai-assisted]` if significant portions are ai-generated
- ensure you understand the code you're submitting
- test thoroughly and be prepared to explain implementation choices

## development setup

```bash
# clone your fork
git clone https://github.com/YOUR_USERNAME/zap.git
cd zap

# install dependencies
go mod download

# build
go build -o zap ./cmd/zap

# run tests
go test ./...

# install locally for testing
go install ./cmd/zap
```

## project structure

```
zap/
├── cmd/zap/          # main entry point
├── internal/
│   ├── ops/          # file operations (copy, move, delete)
│   ├── ui/           # bubble tea ui and theming
│   ├── walk/         # directory traversal
│   └── errs/         # error handling
└── tests/            # integration tests
```

## commit messages

- use present tense ("add feature" not "added feature")
- keep first line under 72 characters
- reference issues when applicable (#123)

examples:
```
add support for symlink preservation in copy

fix progress bar flickering on rapid updates (#45)

improve error messages for permission denied cases
```

## questions?

open an issue or discussion thread. we're happy to help!
