package main

import (
	"github.com/UrielJaloto/Rar-Cracker/internal/config"
)

func main() {
	config := config.ParseConfig()

	if err := config.ValidateFields(); err != nil {
		println(err.Error())
		return
	}

	config.PrintFields()
}
