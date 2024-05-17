package examples

import (
	"github.com/hugoklepsch/go-fuzz-struct/fuzzing"
	"testing"
	"unicode/utf8"
)

func FunctionToTestBasicTypes(s string) bool {
	return utf8.ValidString(s)
}

func FuzzFunctionToTestBasicTypes(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		FunctionToTestBasicTypes(s)
	})
}

type MyStruct struct {
	S string
	B bool
	I int
	F float64
}

func FunctionToTestStruct(m MyStruct) bool {
	return (utf8.ValidString(m.S) && m.B) || (m.I > 0 && m.F < 0)
}

func FunctionToTestPanics(m MyStruct) {
	if m.F < 0.1 && m.F > 0.099 {
		panic("uh oh")
	}
}

func FuzzFunctionToTestStruct_NotWorking(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	f.Fuzz(func(t *testing.T, m MyStruct) {
		FunctionToTestStruct(m)
	})
}

func FuzzFunctionToTestStruct_NotWorking2(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	f.Add("foo", false, 42, 42.0)
	f.Fuzz(func(t *testing.T, args ...any) {
		FunctionToTestStruct(MyStruct{S: args[0].(string), B: args[1].(bool), I: args[2].(int), F: args[3].(float64)})
	})
}

func FuzzFunctionToTestStruct_GeneratedFuzzTarget(f *testing.F) {
	ex1 := MyStruct{I: 42}
	ex2 := MyStruct{F: 42.0}
	fuzzing.Add(f, ex1)
	fuzzing.Add(f, ex2)
	fuzzing.FuzzStruct(f, func(t *testing.T, m MyStruct) {
		t.Logf("%v", m)
		FunctionToTestStruct(m)
	})
}
