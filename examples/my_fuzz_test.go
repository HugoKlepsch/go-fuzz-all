package examples

import (
	"github.com/hugoklepsch/go-fuzz-all/fuzzing"
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

func FunctionToTestWithPanicBug(m MyStruct) {
	if utf8.ValidString(m.S) && m.B && m.I > 0 && m.F < 0.1 && m.F > 0.099 {
		panic("uh oh")
	}
}

func FuzzFunctionToTestStruct_NotWorking(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	f.Fuzz(func(t *testing.T, m MyStruct) {
		FunctionToTestWithPanicBug(m)
	})
}

func FuzzFunctionToTestStruct_NotWorking2(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing []interface {}
	f.Add("foo", false, 42, 42.0)
	f.Fuzz(func(t *testing.T, args ...any) {
		FunctionToTestWithPanicBug(MyStruct{S: args[0].(string), B: args[1].(bool), I: args[2].(int), F: args[3].(float64)})
	})
}

func FuzzFunctionToTestStruct_GeneratedFuzzTarget(f *testing.F) {
	ex1 := MyStruct{I: 42}
	ex2 := MyStruct{F: 42.0}
	fuzzing.Add(f, ex1)
	fuzzing.Add(f, ex2)
	fuzzing.Fuzz[MyStruct](f, func(t *testing.T, m MyStruct) {
		t.Logf("%v", m)
		FunctionToTestWithPanicBug(m)
	})
}
