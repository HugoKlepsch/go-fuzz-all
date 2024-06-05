package fuzzing

import (
	"reflect"
)

type path struct {
	Indexes []int
	Names   []string
}

func (p path) visitChild(index int, name string) path {
	newPath := path{Indexes: p.Indexes[:], Names: p.Names[:]}
	newPath.Indexes = append(newPath.Indexes, index)
	newPath.Names = append(newPath.Names, name)
	return newPath
}

type anyVisitor interface {
	VisitBool(interface{}, path)
	VisitInt(interface{}, path)
	VisitUInt(interface{}, path)
	VisitFloat(interface{}, path)
	VisitComplex(interface{}, path)
	VisitArray(interface{}, path)
	VisitChan(interface{}, path)
	VisitFunc(interface{}, path)
	VisitInterface(interface{}, path)
	VisitMap(interface{}, path)
	VisitPointer(interface{}, path)
	VisitSlice(interface{}, path)
	VisitString(interface{}, path)
	VisitStruct(interface{}, path)
	VisitUnsafePointer(interface{}, path)
}

type typeVisitor interface {
	VisitBool(reflect.Type, path)
	VisitInt(reflect.Type, path)
	VisitUInt(reflect.Type, path)
	VisitFloat(reflect.Type, path)
	VisitComplex(reflect.Type, path)
	VisitArray(reflect.Type, path)
	VisitChan(reflect.Type, path)
	VisitFunc(reflect.Type, path)
	VisitInterface(reflect.Type, path)
	VisitMap(reflect.Type, path)
	VisitPointer(reflect.Type, path)
	VisitSlice(reflect.Type, path)
	VisitString(reflect.Type, path)
	VisitStruct(reflect.Type, path)
	VisitUnsafePointer(reflect.Type, path)
}

type valueVisitor interface {
	VisitBool(reflect.Value, path)
	VisitInt(reflect.Value, path)
	VisitUInt(reflect.Value, path)
	VisitFloat(reflect.Value, path)
	VisitComplex(reflect.Value, path)
	VisitArray(reflect.Value, path)
	VisitChan(reflect.Value, path)
	VisitFunc(reflect.Value, path)
	VisitInterface(reflect.Value, path)
	VisitMap(reflect.Value, path)
	VisitPointer(reflect.Value, path)
	VisitSlice(reflect.Value, path)
	VisitString(reflect.Value, path)
	VisitStruct(reflect.Value, path)
	VisitUnsafePointer(reflect.Value, path)
}

type RealValueVisitor interface {
	VisitBool(bool, path)
	VisitInt(int64, path)
	VisitUInt(uint64, path)
	VisitFloat(float64, path)
	VisitComplex(complex128, path)
	VisitArray(interface{}, path)
	VisitChan(interface{}, path)
	VisitFunc(interface{}, path)
	VisitInterface(interface{}, path)
	VisitMap(interface{}, path)
	VisitPointer(interface{}, path)
	VisitSlice(interface{}, path)
	VisitString(string, path)
	VisitStruct(interface{}, path)
	VisitUnsafePointer(interface{}, path)
}

func TraverseRealValue(v interface{}, visitor RealValueVisitor) {
	traverseValue(v, &anyToRealValueVisitor{visitor: visitor})
}

func TraverseValue(v interface{}, visitor valueVisitor) {
	traverseValue(v, &anyToValueVisitor{visitor: visitor})
}

func TraverseType(v interface{}, visitor typeVisitor) {
	traverseValue(v, &anyToTypeVisitor{visitor: visitor})
}

// traverseValue will recursively traverse any Go value. Does not traverse inside collection types like Array, Map,
// Slice, Chan. It will simply visit those nodes. The visitor may implement their own traversal.
func traverseValue(v interface{}, visitor anyVisitor) {
	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	traverseValueRecursive(vType, vValue, visitor, path{})
}

func traverseValueRecursive(t reflect.Type, v reflect.Value, visitor anyVisitor, pather path) {
	switch t.Kind() {
	case reflect.Bool:
		visitor.VisitBool(v, pather)
		break
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		visitor.VisitInt(v, pather)
		break
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
		visitor.VisitUInt(v, pather)
		break
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		visitor.VisitFloat(v, pather)
		break
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		visitor.VisitComplex(v, pather)
		break
	case reflect.Array:
		visitor.VisitArray(v, pather)
		break
	case reflect.Chan:
		visitor.VisitChan(v, pather)
		break
	case reflect.Func:
		visitor.VisitFunc(v, pather)
		break
	case reflect.Interface:
		visitor.VisitInterface(v, pather)
		// Can't know what reflect.Type to use
		traverseValueRecursive(nil, v.Elem(), visitor, pather.visitChild(0, "Elem()"))
		break
	case reflect.Map:
		visitor.VisitMap(v, pather)
		break
	case reflect.Pointer:
		visitor.VisitPointer(v, pather)
		isPointerSet := !v.IsNil()
		if isPointerSet {
			traverseValueRecursive(t.Elem(), v.Elem(), visitor, pather.visitChild(0, "*"))
		}
		break
	case reflect.Slice:
		visitor.VisitSlice(v, pather)
		break
	case reflect.String:
		visitor.VisitString(v, pather)
		break
	case reflect.Struct:
		visitor.VisitStruct(v, pather)
		for i := 0; i < v.NumField(); i++ {
			viValue := v.Field(i)
			traverseValueRecursive(viValue.Type(), viValue, visitor, pather.visitChild(i, viValue.Type().Name()))
		}
		break
	case reflect.UnsafePointer:
		visitor.VisitUnsafePointer(v, pather)
		// I don't think you can traverse into an unsafe pointer.
		//traverseValueRecursive(v.Elem(), visitor, pather.visitChild(0, "unsafe *"))
		break
	default:
		panic(v.Kind())
	}
}

type anyToRealValueVisitor struct {
	visitor RealValueVisitor
}

func (v *anyToRealValueVisitor) VisitBool(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitBool(value.Bool(), p)
}

func (v *anyToRealValueVisitor) VisitInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInt(value.Int(), p)
}

func (v *anyToRealValueVisitor) VisitUInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUInt(value.Uint(), p)
}

func (v *anyToRealValueVisitor) VisitFloat(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFloat(value.Float(), p)
}

func (v *anyToRealValueVisitor) VisitComplex(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitComplex(value.Complex(), p)
}

func (v *anyToRealValueVisitor) VisitArray(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitArray(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitChan(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitChan(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitFunc(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFunc(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitInterface(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInterface(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitMap(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitMap(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitPointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitPointer(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitSlice(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitSlice(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitString(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitString(value.String(), p)
}

func (v *anyToRealValueVisitor) VisitStruct(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitStruct(value.Interface(), p)
}

func (v *anyToRealValueVisitor) VisitUnsafePointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUnsafePointer(value.Interface(), p)
}

type anyToTypeVisitor struct {
	visitor typeVisitor
}

func (v *anyToTypeVisitor) VisitBool(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitBool(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInt(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitUInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUInt(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitFloat(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFloat(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitComplex(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitComplex(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitArray(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitArray(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitChan(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitChan(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitFunc(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFunc(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitInterface(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInterface(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitMap(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitMap(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitPointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitPointer(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitSlice(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitSlice(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitString(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitString(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitStruct(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitStruct(value.Type(), p)
}

func (v *anyToTypeVisitor) VisitUnsafePointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUnsafePointer(value.Type(), p)
}

type anyToValueVisitor struct {
	visitor valueVisitor
}

func (v *anyToValueVisitor) VisitBool(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitBool(value, p)
}

func (v *anyToValueVisitor) VisitInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInt(value, p)
}

func (v *anyToValueVisitor) VisitUInt(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUInt(value, p)
}

func (v *anyToValueVisitor) VisitFloat(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFloat(value, p)
}

func (v *anyToValueVisitor) VisitComplex(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitComplex(value, p)
}

func (v *anyToValueVisitor) VisitArray(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitArray(value, p)
}

func (v *anyToValueVisitor) VisitChan(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitChan(value, p)
}

func (v *anyToValueVisitor) VisitFunc(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitFunc(value, p)
}

func (v *anyToValueVisitor) VisitInterface(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitInterface(value, p)
}

func (v *anyToValueVisitor) VisitMap(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitMap(value, p)
}

func (v *anyToValueVisitor) VisitPointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitPointer(value, p)
}

func (v *anyToValueVisitor) VisitSlice(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitSlice(value, p)
}

func (v *anyToValueVisitor) VisitString(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitString(value, p)
}

func (v *anyToValueVisitor) VisitStruct(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitStruct(value, p)
}

func (v *anyToValueVisitor) VisitUnsafePointer(i interface{}, p path) {
	value := i.(reflect.Value)
	v.visitor.VisitUnsafePointer(value, p)
}
