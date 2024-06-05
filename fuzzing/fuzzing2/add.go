package fuzzing2

import (
	"fmt"
	"github.com/hugoklepsch/go-fuzz-all/fuzzing"
	"reflect"
)

func Add2[T any](f fuzzing.TestingF, t T) {
	tValue := reflect.ValueOf(t)

	fieldsTraverser := anyToFieldsTraverser{}
	fieldsTraverser.traverseValue(tValue)

	f.Add(fieldsTraverser.fields...)
}

type anyToFieldsTraverser struct {
	fields      []any
	fieldsTypes []reflect.Type
}

func (a *anyToFieldsTraverser) addValue(i any) {
	a.fields = append(a.fields, i)
	a.fieldsTypes = append(a.fieldsTypes, reflect.TypeOf(i))
}

func (a *anyToFieldsTraverser) addZeroValue(t reflect.Type) {
	a.fields = append(a.fields, reflect.Zero(t).Interface())
	a.fieldsTypes = append(a.fieldsTypes, t)
}

func (a *anyToFieldsTraverser) traverseValue(value reflect.Value) {
	switch value.Kind() {
	case reflect.Bool:
		a.addValue(value.Bool())
		break
	case reflect.Int:
		a.addValue(int(value.Int()))
		break
	case reflect.Int8:
		a.addValue(int8(value.Int()))
		break
	case reflect.Int16:
		a.addValue(int16(value.Int()))
		break
	case reflect.Int32:
		a.addValue(int32(value.Int()))
		break
	case reflect.Int64:
		a.addValue(value.Int())
		break
	case reflect.Uint:
		a.addValue(uint(value.Uint()))
		break
	case reflect.Uint8:
		a.addValue(uint8(value.Uint()))
		break
	case reflect.Uint16:
		a.addValue(uint16(value.Uint()))
		break
	case reflect.Uint32:
		a.addValue(uint32(value.Uint()))
		break
	case reflect.Uint64:
		a.addValue(value.Uint())
		break
	case reflect.Uintptr:
		a.addValue(uintptr(value.Uint()))
		break
	case reflect.Float32:
		a.addValue(float32(value.Float()))
		break
	case reflect.Float64:
		a.addValue(value.Float())
		break
	case reflect.Complex64:
		a.addValue(complex64(value.Complex()))
		break
	case reflect.Complex128:
		a.addValue(value.Complex())
		break
	case reflect.Array:
		// TODO
		break
	case reflect.Chan:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Func:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Interface:
		// TODO Can't know what reflect.Type to use
		// true/false whether or not it is nil
		break
	case reflect.Map:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Pointer:
		// Pointer is encoded like this:
		// First value is bool - whether or not the pointer is set.
		// subsequent value(s) are the fields from what the pointer points at.
		isSet := !value.IsNil()
		a.addValue(isSet)
		if isSet {
			a.traverseValue(value.Elem())
		} else {
			a.traverseType(value.Type().Elem())
		}
		break
	case reflect.Slice:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.String:
		a.addValue(value.String())
		break
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			iValue := value.Field(i)
			// Do not traverse into unexported struct fields
			if !value.Type().Field(i).IsExported() {
				continue
			}
			a.traverseValue(iValue)
		}
		break
	case reflect.UnsafePointer:
		// TODO Can we even do anything?
		break
	default:
		panic(fmt.Errorf("unknown kind %v", value.Kind()))
	}
}

func (a *anyToFieldsTraverser) traverseType(t reflect.Type) {
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
		a.addZeroValue(t)
		break
	case reflect.Array:
		// TODO
		break
	case reflect.Chan:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Func:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Interface:
		// TODO Can't know what reflect.Type to use
		// true/false whether or not it is nil
		break
	case reflect.Map:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.Pointer:
		// Pointer is encoded like this:
		// First value is bool - whether or not the pointer is set.
		// subsequent value(s) are the fields from what the pointer points at.
		isSet := false
		a.addValue(isSet)
		a.traverseType(t.Elem())
		break
	case reflect.Slice:
		// TODO Can we even do anything?
		// true/false whether or not it is nil
		break
	case reflect.String:
		a.addZeroValue(t)
		break
	case reflect.Struct:
		// TODO: avoid infinite recursion for types that self-reference
		for i := 0; i < t.NumField(); i++ {
			iStructField := t.Field(i)
			if !iStructField.IsExported() {
				continue
			}
			a.traverseType(iStructField.Type)
		}
		break
	case reflect.UnsafePointer:
		// TODO Can we even do anything?
		break
	default:
		panic(fmt.Errorf("unknown kind %v", t.Kind()))
	}
}
