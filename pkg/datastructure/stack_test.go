package datastructure_test

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestStackTestSuite(t *testing.T) {
	suite.Run(t, new(StackTestSuite))
}

type StackTestSuite struct {
	suite.Suite
}

func (s *StackTestSuite) TestStackMustPopCorrectly() {
	underTest := datastructure.Stack[string]{"A", "B", "C"}

	// Pop correctly
	var result *string
	var ok bool

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("C", *result)

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("B", *result)

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("A", *result)

	result, ok = underTest.Pop()
	s.False(ok)
	s.Nil(result)
}

func (s *StackTestSuite) TestStackMustPushCorrectly() {
	underTest := datastructure.Stack[string]{}

	underTest.Push("A")
	underTest.Push("B")
	underTest.Push("C")

	// Pop correctly.
	var result *string
	var ok bool

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("C", *result)

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("B", *result)

	result, ok = underTest.Pop()
	s.True(ok)
	s.NotNil(result)
	s.Equal("A", *result)

	result, ok = underTest.Pop()
	s.False(ok)
	s.Nil(result)
}
