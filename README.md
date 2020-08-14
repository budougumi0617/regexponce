regexponce
===

[![PkgGoDev](https://pkg.go.dev/badge/budougumi0617/regexponce)][godoc]


[godoc]:https://godoc.org/github.com/budougumi0617/regexponce

`regexp.{Must}Compile{POSIX}` should be called at once for performance.


## Install

You can get `regexponce` by `go get` command.

```bash
$ go get -u github.com/budougumi0617/regexponce/cmd/regexponce
```

## QuickStart

`regexponce` run with `go vet` as below when Go is 1.12 and higher.

```bash
$ go vet -vettool=$(which regexponce) ./...
```

## Analyzer
`regexp.Compile` is heavy, therefore the functions should be called jsut once in the process.
Analyzer confirms that the below functions are not called multiple times.

### Target functions
- [regexp.Compile](https://golang.org/pkg/regexp/#Compile)
- [regexp.CompilePOSIX](https://golang.org/pkg/regexp/#CompilePOSIX)
- [regexp.MustCompile](https://golang.org/pkg/regexp/#MustCompile)
- [regexp.MustCompilePOSIX](https://golang.org/pkg/regexp/#MustCompilePOSIX)
### Allow condition
- Target functions are called in the package scope.
- Target functions are called in `init` function.
- Target functions are called in `main` function.
  - Except if they are called in `for` loop.
- Add [staticcheck's style comments](https://staticcheck.io/docs/#ignoring-problems)
  - `//lint:ignore regexponce REASON`
### Error condition
- Target functions are called in normal function.
- Target functions are called in for loop.

The warning sample is below.

```go
```

## Description
```bash
$ regexponce help regexponce
go run ./cmd/regexponce help regexponce
regexponce: regexp.Compile and below functions should be called at once for performance.
- regexp.MustCompile
- regexp.CompilePOSIX
- regexp.MustCompilePOSIX

Allow call in init, and main(exept for in for loop) functions because each function is called only once.
```

## Contribution
1. Fork ([https://github.com/budougumi0617/regexponce/fork](https://github.com/budougumi0617/regexponce/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

[MIT](https://github.com/budougumi0617/regexponce/blob/master/LICENSE)

## Author
[Yoichiro Shimizu(@budougumi0617)](https://github.com/budougumi0617)
