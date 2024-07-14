package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/tejasp2003/go-url-shortner/routes"
)

func setUpRoutes(app *fiber.App){
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	

 app := fiber.New()

 app.Use(logger.New()) // logger.New() is a middleware that logs all the requests to the console

 setUpRoutes(app)

log.Fatal(app.Listen(os.Getenv("PORT")))

}