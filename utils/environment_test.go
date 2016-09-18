package utils

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type EnvironmentTestSuite struct {
	env *Environment
}

var _ = Suite(&EnvironmentTestSuite{})

func (s *EnvironmentTestSuite) TestNewEnvironment(c *C) {
	env := NewEnvironment(TestEnv, "1")
	c.Assert(env.AppVersion, Equals, "1")
	c.Assert(env.CreatedAt, NotNil)
}
