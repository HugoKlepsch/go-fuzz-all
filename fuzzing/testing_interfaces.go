package fuzzing

import "testing"

//go:generate mockgen -destination ../internal/mocks/testinlFMock.go -package mocks github.com/hugoklepsch/go-fuzz-all/fuzzing TestingF
type TestingF interface {
	Add(...any)
	Fuzz(any)
}

//go:generate mockgen -destination ../internal/mocks/testingTMock.go -package mocks github.com/hugoklepsch/go-fuzz-all/fuzzing TestingT
type TestingT interface {
	testing.TB
}
