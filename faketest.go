package faketest

import (
	"reflect"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertEachFieldIsRandom[T any](tb testing.TB, f func() *T, ignoreFields ...string) {
	tb.Helper()
	const count = 122 // from uuid
	value := reflect.ValueOf(new(T))
	typ := value.Type().Elem()
	if typ.Kind() != reflect.Struct {
		tb.Fatalf("expected a pointer to a struct, got %s", typ.Kind())
	}
	ignored := make([]bool, typ.NumField())
	for i := range typ.NumField() {
		if slices.Contains(ignoreFields, typ.Field(i).Name) {
			ignored[i] = true
		}
		if !typ.Field(i).IsExported() {
			ignored[i] = true
		}
	}
	initials := make([]any, typ.NumField())
	different := make([]bool, typ.NumField())
	for i := range typ.NumField() {
		if ignored[i] {
			different[i] = true
		}
	}
	for i := range count {
		v := reflect.ValueOf(f()).Elem()
		if i == 0 {
			for j := range typ.NumField() {
				if ignored[j] {
					continue
				}
				initials[j] = v.Field(j).Interface()
			}
			continue
		}
		for j := range typ.NumField() {
			if ignored[j] {
				continue
			}
			if cmp.Diff(initials[j], v.Field(j).Interface()) != "" {
				different[j] = true
				ok := true
				for _, d := range different {
					ok = ok && d
				}
				if ok {
					return
				}
			}
		}
	}
	for i, d := range different {
		if !d {
			tb.Errorf("%s.%s is not random", typ.Name(), typ.Field(i).Name)
		}
	}
}
