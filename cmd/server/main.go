package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/routes"
	"github.com/RamazanZholdas/KeyboardistSV2/utils"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app, err := app.Intitialize(os.Getenv("MONGO_URI"), os.Getenv("DATABASE_NAME"))
	if err != nil {
		log.Panic(err)
	}

	routes.Setup(app.Fiber)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		utils.LogInfo("Shutting down server...")
		app.Close()
	}()

	if err := app.Fiber.Listen(":8000"); err != nil {
		log.Panic(err)
	}

	utils.LogInfo("Running cleanning up tasks...")
}
