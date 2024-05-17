package fuzzing

import (
	"reflect"
	"testing"
)

func Add[T any](f *testing.F, t T) {
	tValue := reflect.ValueOf(t)
	tType := tValue.Type()
	if tType.Kind() != reflect.Struct {
		panic("t must be a struct")
	}
	// TODO
	fieldValues := []any{}
	fields := structToFields[T]()
	for _, field := range fields {
		fieldValue := tValue.FieldByIndex(field.Path)
		fieldValues = append(fieldValues, fieldValue.Interface())
	}
	f.Add(fieldValues...)
}

func FuzzStruct[T any](f *testing.F, fn func(*testing.T, T)) {
	in := []reflect.Type{
		reflect.TypeFor[*testing.T](),
	}
	fields := structToFields[T]()
	for _, field := range fields {
		in = append(in, field.Type)
	}
	out := []reflect.Type{}
	fuzzTargetType := reflect.FuncOf(in, out, false)

	fuzzTargetValue := reflect.MakeFunc(fuzzTargetType, func(args []reflect.Value) (results []reflect.Value) {
		testingT := args[0].Interface().(*testing.T)
		t := valuesToStruct[T](fields, args)
		fn(testingT, t)
		return []reflect.Value{}
	})
	f.Fuzz(fuzzTargetValue.Interface())
}

type field struct {
	Type reflect.Type
	Path []int
}

func structToFields[T any]() []field {
	tType := reflect.TypeFor[T]()
	if tType.Kind() != reflect.Struct {
		panic("T must be a struct")
	}
	fields := []field{}
	numFields := tType.NumField()
	for i := 0; i < numFields; i++ {
		elem := tType.Field(i)
		elemType := elem.Type
		elemKind := elemType.Kind()
		switch elemKind {
		// TODO more types
		default:
			fields = append(fields, field{
				Type: elemType,
				Path: []int{i},
			})
			break
		}
	}
	return fields
}

func valuesToStruct[T any](fields []field, values []reflect.Value) T {
	if len(fields)+1 != len(values) {
		panic("values must have the same number of elements as fields, plus an additional element for *testing.T")
	}
	tType := reflect.TypeFor[T]()
	tPtrValue := reflect.New(tType)
	tValue := tPtrValue.Elem()
	for i := 1; i < len(values); i++ {
		tValue.FieldByIndex(fields[i-1].Path).Set(values[i])
	}
	myT := tValue.Interface().(T)

	return myT
}
