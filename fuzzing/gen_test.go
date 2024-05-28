package fuzzing

import (
	"github.com/hugoklepsch/go-fuzz-all/internal/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAddStruct_NoNesting_NoPtrs(t *testing.T) {
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

	Add(mockF, f1)
}

func TestAddStruct_NoNesting_Ptrs(t *testing.T) {
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

	Add(mockF, f1)
}

func TestAddStruct_NoNesting_PtrsWithNil(t *testing.T) {
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

	Add(mockF, f1)
}
