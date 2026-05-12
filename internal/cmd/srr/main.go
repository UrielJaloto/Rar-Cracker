package main

import (
	"github.com/UrielJaloto/surgical-rar-recovery/internal/services"
)

func main() {
	application := services.NewApplication()
	application.Run()
}
