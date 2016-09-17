package utils

import (
	"database/sql"
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
	MigrationsFolder string = "models/migrations/"

	// gormLogSQL - Whether GORM SQL logging is enabled or not
	gormLogSQL bool = false
)

// Environment is a thing that holds env. specific stuff
type Environment struct {
	AppVersion string
	Name       string
	CreatedAt  time.Time
}

// NewEnvironment creates a new environment
func NewEnvironment(environment string, appVersion string) (*Environment, *gorm.DB) {
	cfg, err := readConfig(environment)
	if err != nil {
		log.Fatal(err) //no way to launch the app without an Environment, fatal!
	}

	cfg, err = cfg.Get(environment)
	cfg.Env() // test this out!

	env := &Environment{Name: environment, AppVersion: appVersion, CreatedAt: time.Now()}
	connection, err := connectToDatabase(cfg)
	if err != nil {
		log.Fatal(err) //no way to launch the app without a DB, fatal!
	}

	connection.LogMode(gormLogSQL)
	return env, connection
}

// MigrateDatabase - performs database migrations
func (env *Environment) MigrateDatabase(db *sql.DB) error {
	log.Println("Migrating database...")

	err := dbmigrate.Run(db, adjustPath(env.Name, MigrationsFolder))
	if err != nil {
		log.Printf("Failed to migrate database! Error was: %s\n", err)
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
