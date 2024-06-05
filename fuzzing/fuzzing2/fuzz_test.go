package fuzzing2

import (
	"github.com/hugoklepsch/go-fuzz-all/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestFuzz2Struct_NoNesting_NoPtrs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Foo struct {
		S string
		B bool
		I int
		F float64
	}
	mockF.EXPECT().Fuzz(gomock.AssignableToTypeOf(func(t *testing.T, s string, b bool, i int, f float64) {}))
	Fuzz2(mockF, func(t *testing.T, foo Foo) {})
}

func TestFuzz2Struct_NoNesting_Ptrs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	mockF.EXPECT().Fuzz(gomock.AssignableToTypeOf(
		func(t *testing.T, b1 bool, s1 string, b2, b3, b4 bool, i1 int, b5 bool, f1 float64) {}),
	)
	Fuzz2(mockF, func(t *testing.T, foo Foo) {})
}

func TestFuzz2Struct_NoNesting_PtrsWithNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	f1 := Foo{nil, nil, nil, nil}

	mockF.EXPECT().Add(false, "", false, false, false, 0, false, 0.0)

	Add2(mockF, f1)
}

func TestFuzz2Struct_Nesting_PtrsWithNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Nested struct {
		F *float32
		C *complex64
	}

	type Foo struct {
		S *string
		B *bool
		i *int // n.b. i is not exported, so we should not add it.
		F *float64
		N *Nested
	}
	f1 := Foo{nil, nil, nil, nil, ptr(Nested{ptr(float32(3.14)), ptr(complex64(12))})}

	mockF.EXPECT().Add(
		false, "",
		false, false,
		false, 0.0,
		true, true, float32(3.14), true, complex64(12))

	Add2(mockF, f1)
}

func TestBuildAnyTraverser_NoNesting_NoPts(t *testing.T) {
	type Foo struct {
		S string
		B bool
		I int
		F float64
	}
	fieldsAny := []any{"foo", true, 42, 3.14}
	fields := make([]reflect.Value, 0, len(fieldsAny))
	for _, fieldAny := range fieldsAny {
		fields = append(fields, reflect.ValueOf(fieldAny))
	}
	builder := buildAnyTraverser{
		fields: fields,
	}
	builtFoo := builder.traverseType(reflect.TypeFor[Foo]()).Interface().(Foo)
	assert.Equal(t, fieldsAny[0].(string), builtFoo.S)
	assert.Equal(t, fieldsAny[1].(bool), builtFoo.B)
	assert.Equal(t, fieldsAny[2].(int), builtFoo.I)
	assert.Equal(t, fieldsAny[3].(float64), builtFoo.F)
}

func TestBuildAnyTraverser_NoNesting_Pts(t *testing.T) {
	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	fieldsAny := []any{true, "foo", true, true, true, 42, true, 3.14}
	fields := make([]reflect.Value, 0, len(fieldsAny))
	for _, fieldAny := range fieldsAny {
		fields = append(fields, reflect.ValueOf(fieldAny))
	}
	builder := buildAnyTraverser{
		fields: fields,
	}
	builtFoo := builder.traverseType(reflect.TypeFor[Foo]()).Interface().(Foo)
	assert.Equal(t, fieldsAny[1].(string), *builtFoo.S)
	assert.Equal(t, fieldsAny[3].(bool), *builtFoo.B)
	assert.Equal(t, fieldsAny[5].(int), *builtFoo.I)
	assert.Equal(t, fieldsAny[7].(float64), *builtFoo.F)
}

func TestBuildAnyTraverser_NoNesting_PtsWithNil(t *testing.T) {
	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	// Although I set a value for the fields, the `false` should tell the builder to
	// not set the pointers.
	fieldsAny := []any{false, "foo", false, true, false, 42, false, 3.14}
	fields := make([]reflect.Value, 0, len(fieldsAny))
	for _, fieldAny := range fieldsAny {
		fields = append(fields, reflect.ValueOf(fieldAny))
	}
	builder := buildAnyTraverser{
		fields: fields,
	}
	builtFoo := builder.traverseType(reflect.TypeFor[Foo]()).Interface().(Foo)
	expected := Foo{}
	assert.Equal(t, expected, builtFoo)
}

func TestBuildAnyTraverser_Nesting(t *testing.T) {
	type Nested struct {
		F *float32
		C *complex64
	}

	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
		N *Nested
	}

	// This set of fields has both set and unset pointers
	fieldsAny := []any{
		false, "",
		false, false,
		false, 0,
		false, 0.0,
		true, true, float32(3.14), true, complex64(12),
	}
	fields := make([]reflect.Value, 0, len(fieldsAny))
	for _, fieldAny := range fieldsAny {
		fields = append(fields, reflect.ValueOf(fieldAny))
	}
	builder := buildAnyTraverser{
		fields: fields,
	}
	builtFoo := builder.traverseType(reflect.TypeFor[Foo]()).Interface().(Foo)
	expected := Foo{N: &Nested{
		F: ptr(fieldsAny[10].(float32)),
		C: ptr(fieldsAny[12].(complex64)),
	}}
	assert.Equal(t, expected, builtFoo)
}

func TestBuildAnyTraverser_Unexported(t *testing.T) {
	// Foo.n is unexported, but Nested has exported fields. We must not assign n.
	type Nested struct {
		F *float32
		C *complex64
	}
	type Foo struct {
		S string
		// Must not set i or f
		i int
		f float64
		B bool
		n *Nested
	}
	// Can only set Foo.S, Foo.B
	fieldsAny := []any{"foo", true}
	fields := make([]reflect.Value, 0, len(fieldsAny))
	for _, fieldAny := range fieldsAny {
		fields = append(fields, reflect.ValueOf(fieldAny))
	}
	builder := buildAnyTraverser{
		fields: fields,
	}
	builtFoo := builder.traverseType(reflect.TypeFor[Foo]()).Interface().(Foo)
	expected := Foo{
		S: fieldsAny[0].(string),
		B: fieldsAny[1].(bool),
	}
	assert.Equal(t, expected, builtFoo)
}
