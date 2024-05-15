package examples

import (
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
	s string
	b bool
	i int
	f float64
}

func FunctionToTestStruct(m MyStruct) bool {
	return (utf8.ValidString(m.s) && m.b) || (m.i > 0 && m.f < 0)
}

func FuzzFunctionToTestStruct_NotWorking(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	f.Fuzz(func(t *testing.T, m MyStruct) {
		FunctionToTestStruct(m)
	})
}

func FuzzFunctionToTestStruct_GeneratedFuzzTarget(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	/* TODO:
	ex1 := MyStruct{i: 42}
	ex2 := MyStruct{f: 42.0}

	addArguments := GenerateCorpus(ex1)
	f.Add(addArguments...)
	addArguments = GenerateCorpus(ex2)
	f.Add(addArguments...)

	fuzz_target := GenerateFuzzTarget[MyStruct](func(t *testing.T, m MyStruct){})
	f.Fuzz(fuzz_target)
	*/
}
