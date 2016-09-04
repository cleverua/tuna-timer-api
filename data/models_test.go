package data


import (
	"testing"
	. "gopkg.in/check.v1"
	"github.com/pavlo/slack-time/utils"
	"log"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ModelsTestSuite struct{
	env *utils.Environment
}

var _ = Suite(&ModelsTestSuite{})

func (s *ModelsTestSuite) SetUpSuite(c *C) {

	e, err := utils.NewEnvironment(utils.TestEnv);

	if err != nil {
		c.Error(err)
	}

	s.env = e

	err = s.env.MigrateDatabase()
	if err != nil {
		c.Error(err)
	}
}

func (s *ModelsTestSuite) TearDownSuite(c *C) {
	s.env.ReleaseResources()
}

func (s *ModelsTestSuite) TestFoo(c *C) {
	log.Println("TestFoo")
	//s.env.OrmDB.AutoMigrate(
		//&Team{},
		//&TeamUser{},
		//&Project{},
		//&Task{},
		//&Timer{},
	//)

	c.Assert(1, Equals, 2)
}