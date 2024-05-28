# This is a work in progress.

---

# `go-fuzz-any`

Use this package to fuzz any Go type.

```go
package examples

import (
	"github.com/hugoklepsch/go-fuzz-all/fuzzing"
	"testing"
	"unicode/utf8"
)

type MyStruct struct {
    S  string
    B  bool
    I  int
    F  float64
}

func FunctionToTestPanics(m MyStruct) {
    if utf8.ValidString(m.S) && m.B && m.I > 0 && m.F < 0.1 && m.F > 0.099 {
        panic("uh oh")
    }
}

func FuzzFunctionToTestStruct_NotWorking(f *testing.F) {
    // Does not work:
    // panic: testing: unsupported type for fuzzing examples.MyStruct
    f.Fuzz(func(t *testing.T, m MyStruct) {
        FunctionToTestPanics(m)
    })
}

func FuzzFunctionToTestStruct_NotWorking2(f *testing.F) {
    // Does not work:
    // panic: testing: unsupported type for fuzzing []interface {}
    f.Add("foo", false, 42, 42.0)
    f.Fuzz(func(t *testing.T, args ...any) {
        FunctionToTestPanics(MyStruct{S: args[0].(string), B: args[1].(bool), I: args[2].(int), F: args[3].(float64)})
    })
}

func FuzzFunctionToTestStruct_GeneratedFuzzTarget(f *testing.F) {
    ex1 := MyStruct{I: 42}
    ex2 := MyStruct{F: 42.0}
    fuzzing.Add(f, ex1)
    fuzzing.Add(f, ex2)
    fuzzing.Fuzz[MyStruct](f, func(t *testing.T, m MyStruct) {
        t.Logf("%v", m)
        FunctionToTestPanics(m)
    })
}
```

# TODO

* I started work on visitor pattern implementation, but I realized that
  this approach would be unable to fuzz the fields of any struct
  that is referenced by a pointer that is `nil` when passed into
  `Add`.

From a source code comment:

```go
        /*
                TODO: make this work with pointers to more complicated types.
                If there is a pointer that points to a structure, we want to fuzz its fields too.
                If the user passes in an object where this pointer is set, then we will dereference the
                pointer and continue to traverse the struct, adding its fields.
                However, if the user passes in an object where this pointer is not set, then right now we
                cannot dereference the pointer in the `reflect.Value`, so we cannot traverse the pointed to
                structure, which means the fuzzer will not fuzz any fields in that struct.
                When adding, we normally traverse using `reflect.Value`s of the passed in values.
                To deal with this scenario, we need to switch from traversing the real value to traversing
                the `reflect.Type`, and adding the zero values of the types to the fuzz corpus.
        */
```

* Some further refactoring is needed to support that scenario, but I
  wanted to get this library out there. It works in some basic
  scenarios, like outlined above. 
