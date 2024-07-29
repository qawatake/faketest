package faketest_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qawatake/faketest"
)

func TestAssertEachFieldIsRandom(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fakeFuctory  func() *Book
		ignoreFields []string
		errors       []string
	}{
		{
			name: "success when all fields are random",
			fakeFuctory: func() *Book {
				return &Book{
					ID:         rand.Int64(), //nolint:gosec
					Name:       sample("book1", "book2"),
					Tags:       shuffle([]string{"tag1", "tag2"}),
					unexported: sample("unexported1", "unexported2"),
				}
			},
			ignoreFields: nil,
			errors:       nil,
		},
		{
			name: "fail when some exported field are not random",
			fakeFuctory: func() *Book {
				return &Book{
					ID:         42,
					Name:       sample("book1", "book2"),
					Tags:       shuffle([]string{"tag1", "tag2"}),
					unexported: sample("unexported1", "unexported2"),
				}
			},
			ignoreFields: nil,
			errors:       []string{"Book.ID is not random"},
		},
		{
			name: "success when all non-random fields are unexported",
			fakeFuctory: func() *Book {
				return &Book{
					ID:         rand.Int64(), //nolint:gosec
					Name:       sample("book1", "book2"),
					Tags:       shuffle([]string{"tag1", "tag2"}),
					unexported: "const",
				}
			},
			ignoreFields: nil,
			errors:       nil,
		},
		{
			name: "success when non-random field is ignored",
			fakeFuctory: func() *Book {
				return &Book{
					ID:         42,
					Name:       sample("book1", "book2"),
					Tags:       shuffle([]string{"tag1", "tag2"}),
					unexported: sample("uneported1", "unexported2"),
				}
			},
			ignoreFields: []string{"ID"},
			errors:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reporter := &testReporter{}
			faketest.AssertEachFieldIsRandom(reporter, tt.fakeFuctory, tt.ignoreFields...)
			if diff := cmp.Diff(tt.errors, reporter.ErrorfCalls()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

type Book struct {
	ID         int64
	Name       string
	Tags       []string
	unexported string
}

type testReporter struct {
	testing.TB
	errorfCalls []string
}

func (t *testReporter) Helper() {}

func (t *testReporter) Errorf(format string, args ...any) {
	t.errorfCalls = append(t.errorfCalls, fmt.Sprintf(format, args...))
}

func (t *testReporter) ErrorfCalls() []string {
	return t.errorfCalls
}

func sample[E any](s ...E) E {
	return s[rand.IntN(len(s))] //nolint:gosec
}

func shuffle[S ~[]E, E any](s S) S {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}
