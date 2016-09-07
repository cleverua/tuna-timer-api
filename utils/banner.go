package utils

import (
	"log"
)

// PrintBanner just a bit of fun
func PrintBanner(env *Environment) {
	log.Println(banner)
	log.Printf("Version: %s\n", env.AppVersion)
	log.Printf("Environment: %s\n", env.Name)
	log.Println("--------------------------------------------")
}

const banner string = `-------------------------------------------- 
  _, ,  _    _,,  ,    ___,___, , ,  _, 
 (_, | '|\  /  |_/    ' | ' |  |\/| /_, 
  _)'|__|-\'\_'| \      |  _|_,| '|'\_
`
