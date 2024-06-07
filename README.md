# `go-fuzz-any`

Use this package to fuzz any Go type.

# Non-AGPL license is available with the purchase of a support contract. Contact [hugo.klepsch@gmail.com][2] for details

## Problem and Solution example

You may have tried to fuzz a function that takes a struct as a parameter. The direct way of doing this does not work, 
because the Go fuzzing framework is limited to `string`, `[]byte`, `int`, `int8`, `int16`, `int32`/`rune`, `int64`, 
`uint`, `uint8`/`byte`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, or `bool`.

You will encounter errors like:

```
panic: testing: unsupported type for fuzzing examples.MyStruct
```

Below you will find two examples using the broken direct method, and one working example using `go-fuzz-all` to 
fuzz a function that takes a struct.

```go
package examples

import (
	"github.com/hugoklepsch/go-fuzz-all/fuzzing"
	"testing"
	"unicode/utf8"
)

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

func FuzzFunctionToTestWithPanicBug_NotWorking(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing examples.MyStruct
	f.Fuzz(func(t *testing.T, m MyStruct) {
		FunctionToTestWithPanicBug(m)
	})
}

func FuzzFunctionToTestWithPanicBug_NotWorking2(f *testing.F) {
	// Does not work:
	// panic: testing: unsupported type for fuzzing []interface {}
	f.Add("foo", false, 42, 42.0)
	f.Fuzz(func(t *testing.T, args ...any) {
		FunctionToTestWithPanicBug(MyStruct{S: args[0].(string), B: args[1].(bool), I: args[2].(int), F: args[3].(float64)})
	})
}

// Solution: Use go-fuzz-any to fuzz Go structs directly.
func FuzzFunctionToTestWithPanicBug_Working(f *testing.F) {
	ex1 := MyStruct{I: 42}
	ex2 := MyStruct{F: 42.0}
	// Works!
	// Add examples to the fuzz corpus
	fuzzing.Add(f, ex1)
	fuzzing.Add(f, ex2)
	// Begin fuzzing
	fuzzing.Fuzz(f, func(t *testing.T, m MyStruct) {
		t.Logf("%v", m)
		FunctionToTestWithPanicBug(m)
	})
}
```

## Anatomy of a Fuzz test

```

         ╔              func FuzzMyFunc(f *testing.F) {
         ║                  ex1 := MyStruct{I: 42}
         ║                  ex2 := MyStruct{F: 42.0}
         ║                  // Add examples to the fuzz corpus
         ║                  fuzzing.Add(f, ex1) <<----- Seed corpus
         ║                  fuzzing.Add(f, ex2) <<----- Same type as fuzzing argument
fuzz test╣                  // Begin fuzzing
         ║              ╔   fuzzing.Fuzz(f, func(t *testing.T, m MyStruct) {
         ║              ║       t.Logf("%v", m)                ^^^^^^^^^^----fuzzing argument
         ║  fuzz target ╣       err := MyFunc(m)
         ║              ║       if err != nil { t.Fail() }
         ║              ╚   })
         ╚              }

```

## Add to corpus

### `fuzzing.Add[T any](f *testing.F, t T)`

`fuzzing.Add` will add the given `t` to the fuzz corpus, giving the fuzzer examples to start from.

If your function takes multiple arguments, create a single struct that contains each argument as a struct field, and
fuzz using that type.

## Fuzzing

### `fuzzing.Fuzz[T any](f *testing.F, fuzzTarget func(t *testing.T, myT T))`

`fuzzing.Fuzz` is called to set up fuzzer. Provide a function `fuzzTarget` that is called for each iteration of the 
fuzz test. It should be safe to call from multiple threads and fast.

## Running fuzz tests

```sh
go test -fuzz={FuzzTestName}
```

See [documentation][1].

[1]: https://pkg.go.dev/cmd/go
[2]: mailto:hugo.klepsch@gmail.com
