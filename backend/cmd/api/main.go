package main

import (
    "log"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
    app := fiber.New()
    app.Use(cors.New())

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok", "api": "cubearchitect"})
    })

    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    log.Fatal(app.Listen(":" + port))
}