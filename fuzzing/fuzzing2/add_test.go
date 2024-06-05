package fuzzing2

import (
	"github.com/hugoklepsch/go-fuzz-all/internal/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func ptr[T any](t T) *T {
	return &t
}

func TestAdd2Struct_NoNesting_NoPtrs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Foo struct {
		S string
		B bool
		I int
		F float64
	}
	f1 := Foo{"foo", true, 42, 3.14}

	mockF.EXPECT().Add("foo", true, 42, 3.14)

	Add2(mockF, f1)
}

func TestAdd2Struct_NoNesting_Ptrs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockF := mocks.NewMockTestingF(mockCtrl)

	type Foo struct {
		S *string
		B *bool
		I *int
		F *float64
	}
	f1 := Foo{ptr("foo"), ptr(true), ptr(42), ptr(3.14)}

	mockF.EXPECT().Add(true, "foo", true, true, true, 42, true, 3.14)

	Add2(mockF, f1)
}

func TestAdd2Struct_NoNesting_PtrsWithNil(t *testing.T) {
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

func TestAdd2Struct_Nesting_PtrsWithNil(t *testing.T) {
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
