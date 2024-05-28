package fuzzing

import (
	"fmt"
	"reflect"
	"strings"
)

func indent(level int) {
	indent := make([]string, level+1)
	fmt.Print(strings.Join(indent, "  "))
}

type addCorpusVisitor struct {
	level                int
	FieldInterfaceValues []any
}

func (a *addCorpusVisitor) debugPrint(valueType string, value reflect.Value, p path) {
	indent(len(p.Indexes))
	fmt.Printf("%s: %s, %v\n", valueType, value.Type().Name(), value.Interface())
}

func (a *addCorpusVisitor) add(value reflect.Value) {
	a.FieldInterfaceValues = append(a.FieldInterfaceValues, value.Interface())
}

func (a *addCorpusVisitor) VisitBool(value reflect.Value, p path) {
	a.debugPrint("Bool", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitInt(value reflect.Value, p path) {
	a.debugPrint("Int", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitUInt(value reflect.Value, p path) {
	a.debugPrint("UInt", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitFloat(value reflect.Value, p path) {
	a.debugPrint("Float", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitComplex(value reflect.Value, p path) {
	a.debugPrint("Complex", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitArray(value reflect.Value, p path) {
	a.debugPrint("Array", value, p)
}

func (a *addCorpusVisitor) VisitChan(value reflect.Value, p path) {
	a.debugPrint("Chan", value, p)
}

func (a *addCorpusVisitor) VisitFunc(value reflect.Value, p path) {
	a.debugPrint("Func", value, p)
}

func (a *addCorpusVisitor) VisitInterface(value reflect.Value, p path) {
	a.debugPrint("Interface", value, p)
}

func (a *addCorpusVisitor) VisitMap(value reflect.Value, p path) {
	a.debugPrint("Map", value, p)
}

func (a *addCorpusVisitor) VisitPointer(value reflect.Value, p path) {
	a.debugPrint("Pointer", value, p)
	// If it is a pointer, let's add a boolean to represent whether the pointer is nil or not. If the Go Fuzz
	// framework passes in false for the boolean, we will provide nil for the value. If it is true, we will
	// provide a value for the pointer.
	// If the provided value is nil, seed the corpus of the boolean with false.
	isPointerSet := !value.IsNil()
	isPointerSetValue := reflect.ValueOf(isPointerSet)
	a.add(isPointerSetValue)
	if !isPointerSet {
		a.add(reflect.Zero(value.Type().Elem()))
	}
	// If it is set, then we will traverse into the pointed-to field and add it there.

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
}

func (a *addCorpusVisitor) VisitSlice(value reflect.Value, p path) {
	a.debugPrint("Slice", value, p)
	if value.Type().Elem().Kind() == reflect.Uint8 {
		// it is a []byte, which is the only slice allowed
		a.add(value)
	}
}

func (a *addCorpusVisitor) VisitString(value reflect.Value, p path) {
	a.debugPrint("String", value, p)
	a.add(value)
}

func (a *addCorpusVisitor) VisitStruct(value reflect.Value, p path) {
	a.debugPrint("Struct", value, p)
}

func (a *addCorpusVisitor) VisitUnsafePointer(value reflect.Value, p path) {
	a.debugPrint("Unsafe pointer", value, p)
}
