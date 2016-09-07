package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/olebedev/config"

	"github.com/jinzhu/gorm"
	// PosgreSQL driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tanel/dbmigrate"
)

const (
	// ProductionEnv - a value that indicates about production env
	ProductionEnv = "production"
	// DevelopmentEnv - a value that indicates about development env
	DevelopmentEnv = "development"
	// TestEnv - a value that indicates about test env
	TestEnv = "test"
)

const (
	// ConfigFile - path to YML config file
	ConfigFile string = "config.yml"

	// MigrationsFolder - the folder to look migration SQLs in
	MigrationsFolder string = "data/migrations/"

	// gormLogSQL - Whether GORM SQL logging is enabled or not
	gormLogSQL bool = false
)

// Environment is a thing that holds env. specific stuff
type Environment struct {
	AppVersion string
	Name       string
	CreatedAt  time.Time
	OrmDB      *gorm.DB
}

// NewEnvironment creates a new environment
func NewEnvironment(environment string, appVersion string) *Environment {
	cfg, err := readConfig(environment)
	if err != nil {
		log.Fatal(err) //no way to launch the app without an Environment, fatal!
	}
	cfg.Env() // test this out!
	cfg, err = cfg.Get(environment)

	env := &Environment{Name: environment, AppVersion: appVersion, CreatedAt: time.Now()}
	connection, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatal(err) //no way to launch the app without a DB, fatal!
	}
	env.OrmDB = connection
	env.OrmDB.LogMode(gormLogSQL)

	return env
}

// ReleaseResources - supposed to be called in the end of application/test suite lifecycle
func (env *Environment) ReleaseResources() {
	log.Println("Releasing resources...")
	env.OrmDB.Close()
	log.Println("Done releasing resources")
}

// MigrateDatabase - performs database migrations
func (env *Environment) MigrateDatabase() error {
	log.Println("Migrating database...")

	err := dbmigrate.Run(env.OrmDB.DB(), adjustPath(env.Name, MigrationsFolder))
	if err != nil {
		return err
	}

	log.Println("Database migrated!")
	return nil
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf("sslmode=disable dbname=%s host=%s port=%s user=%s password=%s",
			cfg.UString("database.name"),
			cfg.UString("database.host"),
			cfg.UString("database.port"),
			cfg.UString("database.user"),
			cfg.UString("database.pass"),
		))

	if err != nil {
		return nil, err
	}
	return db, nil
}

func readConfig(environmentName string) (*config.Config, error) {
	cfg, err := config.ParseYamlFile(adjustPath(environmentName, ConfigFile))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// a hack to walk around this issue:
// http://stackoverflow.com/questions/23847003/golang-tests-and-working-directory
// does it have a nicer solution?
func adjustPath(environmentName string, resource string) string {
	if environmentName == TestEnv {
		return "../" + resource
	}
	return resource
}
