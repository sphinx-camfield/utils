package booter

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type BooterTest struct {
	suite.Suite
}

func (suite *BooterTest) TestBoot() {
	// Given
	boots := []BootFunc{
		func(c *Container) CleanUpFunc {
			return func() {}
		},
	}

	// When
	clean := Boot(boots)

	// Then
	suite.NotNil(clean)
}

func TestBoot(t *testing.T) {
	suite.Run(t, new(BooterTest))
}
