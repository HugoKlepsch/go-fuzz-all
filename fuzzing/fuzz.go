package fuzzing

import (
	"fmt"
	"reflect"
	"testing"
)

func Fuzz[T any](f TestingF, fn func(*testing.T, T)) {
	tType := reflect.TypeFor[T]()
	in := []reflect.Type{
		reflect.TypeFor[*testing.T](),
	}
	fieldsTraverser := anyToFieldsTraverser{}
	fieldsTraverser.traverseType(tType)
	in = append(in, fieldsTraverser.fieldsTypes...)

	out := []reflect.Type{}

	fuzzTargetType := reflect.FuncOf(in, out, false)

	fuzzTargetValue := reflect.MakeFunc(fuzzTargetType, func(args []reflect.Value) (results []reflect.Value) {
		testingT := args[0].Interface().(*testing.T)
		builder := buildAnyTraverser{
			fields: args[1:],
		}
		fn(testingT, builder.traverseType(tType).Interface().(T))
		return nil
	})
	f.Fuzz(fuzzTargetValue.Interface())
}

type buildAnyTraverser struct {
	fields []reflect.Value
	value  reflect.Value
}

func (a *buildAnyTraverser) popValue() reflect.Value {
	value := a.fields[0]
	a.fields = a.fields[1:]
	return value
}

func (a *buildAnyTraverser) traverseType(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return a.popValue()
	case reflect.Array:
		// TODO
		fallthrough
	case reflect.Chan:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		fallthrough
	case reflect.Func:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		fallthrough
	case reflect.Interface:
		// TODO Can't know what reflect.Type to use
		// true/false whether or not it is nil
		fallthrough
	case reflect.Map:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		return reflect.New(t).Elem()
	case reflect.Pointer:
		// Pointer is encoded like this:
		// First value is bool - whether or not the pointer is set.
		// subsequent value(s) are the fields from what the pointer points at.
		isSet := a.popValue().Bool()
		// TODO check this logic
		pointedToType := t.Elem()
		valueToSet := a.traverseType(pointedToType)
		if isSet {
			valueToSetPtrValue := reflect.New(pointedToType)
			valueToSetPtrValue.Elem().Set(valueToSet)
			return valueToSetPtrValue
		} else {
			tPointer := reflect.New(t)
			tPointer.Elem().Set(reflect.Zero(t))
			return tPointer.Elem()
		}
	case reflect.Slice:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		return reflect.New(t).Elem()
	case reflect.String:
		return a.popValue()
	case reflect.Struct:
		// TODO: avoid infinite recursion for types that self-reference
		structValue := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			if !structField.IsExported() {
				continue
			}
			fieldValue := a.traverseType(structField.Type)
			structValue.Field(i).Set(fieldValue)
		}
		return structValue
	case reflect.UnsafePointer:
		// TODO Can we even do anything?
		return reflect.New(t).Elem()
	default:
		panic(fmt.Errorf("unknown kind %v", t.Kind()))
	}
}
