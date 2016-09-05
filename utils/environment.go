package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/olebedev/config"

	"github.com/jinzhu/gorm"
	// PosgreSQL driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tanel/dbmigrate"
)

const (
	// ProductionEnv - a value that indicates about production env
	ProductionEnv = "production"
	// TestEnv - a value that indicates about test env
	TestEnv = "test"
)

const (
	// ConfigFile - path to YML config file
	ConfigFile string = "../config.yml"

	// MigrationsFolder - the folder to look migration SQLs in
	MigrationsFolder string = "../data/migrations/"
)

// Environment is a thing that holds env. specific stuff
type Environment struct {
	OrmDB *gorm.DB
	RawDB *sql.DB // for unit tests
}

// NewEnvironment creates a new environment
func NewEnvironment(environment string) (*Environment, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	cfg.Env() // test this out!
	cfg, err = cfg.Get(environment)

	env := &Environment{}
	connection, err := connectToDatabase(cfg)
	if err != nil {
		return nil, err
	}
	env.OrmDB = connection
	env.RawDB = env.OrmDB.DB()
	return env, nil
}

// ReleaseResources - supposed to be called in the end of application/test suite lifecycle
func (env *Environment) ReleaseResources() {
	env.OrmDB.Close()
}

// MigrateDatabase - performs database migrations
func (env *Environment) MigrateDatabase() error {
	log.Println("Migrating database...")

	err := dbmigrate.Run(env.RawDB, MigrationsFolder)
	if err != nil {
		return err
	}

	log.Println("Database migrated!")
	return nil
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {

	log.Println("Connecting to database:")
	log.Printf("database.name: %s", cfg.UString("database.name"))
	log.Printf("database.host: %s", cfg.UString("database.host"))
	log.Printf("database.port: %s", cfg.UString("database.port"))
	log.Printf("database.user: %s", cfg.UString("database.user"))
	log.Print("database.pass: ***********")

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
		log.Println("Failed to connect!!")
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected successfully!")
	return db, nil
}

func readConfig() (*config.Config, error) {
	cfg, err := config.ParseYamlFile(ConfigFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
