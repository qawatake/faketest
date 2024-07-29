# faketest

[![Go Reference](https://pkg.go.dev/badge/github.com/qawatake/faketest.svg)](https://pkg.go.dev/github.com/qawatake/faketest)
[![test](https://github.com/qawatake/faketest/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/faketest/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/qawatake/faketest/graph/badge.svg)](https://codecov.io/gh/qawatake/faketest)

```go
type Book struct {
	ID   int64
	Name string
}

func Good() *Book {
	return &Book{
		ID:   rand.Int64(),
		Name: []string{"Macbeth", "Hamlet", "Othello"}[rand.Intn(3)],
	}
}

func Bad() *Book {
	return &Book{
		ID:   rand.Int64(),
		Name: "Fixed",
	}
}
```
