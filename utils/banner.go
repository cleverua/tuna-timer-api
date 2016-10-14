package utils

import (
	"log"
)

// PrintBanner just a bit of fun
func PrintBanner(env *Environment) {
	log.Println("Tuna Timer")
	log.Printf("Version: %s\n", env.AppVersion)
	log.Printf("Environment: %s\n", env.Name)
	log.Println("--------------------------------------------")
}
