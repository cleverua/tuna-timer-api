package utils

import (
	"testing"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type STSuite struct{
	env *Environment
}

var _ = Suite(&STSuite{})


func (s *STSuite) TestNewEnvironment(c *C) {

	env, err := NewEnvironment(TestEnv)
	if err != nil {
		c.Error(err)
	}

	c.Assert(env.OrmDB, NotNil)
	c.Assert(env.RawDB, NotNil)
}

