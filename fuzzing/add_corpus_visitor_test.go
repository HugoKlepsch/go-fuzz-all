package fuzzing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoo(t *testing.T) {
	type MyStruct struct {
		S string
		B bool
		I int
		F float64
	}
	ex1 := MyStruct{"foo", true, 42, 3.14}
	TraverseValue(ex1, &addCorpusVisitor{})
}

func TestStruct_NoNesting_NoPtrs(t *testing.T) {
	type Foo struct {
		S string
		B bool
		I int
		F float64
	}
	f1 := Foo{"foo", true, 42, 3.14}

	visitor := addCorpusVisitor{}
	TraverseValue(f1, &visitor)
	assert.Equal(t, []any{"foo", true, 42, 3.14}, visitor.FieldInterfaceValues)
}

func ptr[T any](t T) *T {
	return &t
}

func TestStruct_NoNesting_Ptrs(t *testing.T) {
	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	f1 := Foo{ptr("foo"), ptr(true), ptr(42), ptr(3.14)}

	visitor := addCorpusVisitor{}
	TraverseValue(f1, &visitor)
	assert.Equal(t, []any{true, "foo", true, true, true, 42, true, 3.14}, visitor.FieldInterfaceValues)
}

func TestStruct_NoNesting_PtrsWithNil(t *testing.T) {
	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	f1 := Foo{nil, nil, nil, nil}

	visitor := addCorpusVisitor{}
	TraverseValue(f1, &visitor)
	assert.Equal(t, []any{false, "", false, false, false, 0, false, 0.0}, visitor.FieldInterfaceValues)
}
