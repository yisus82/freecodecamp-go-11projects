package main

import (
	"08-fiber-crm/database"
	"08-fiber-crm/lead"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/leads", lead.GetLeads)
	app.Get("/leads/:id", lead.GetLead)
	app.Post("/leads", lead.CreateLead)
	app.Delete("/leads/:id", lead.DeleteLead)
}

func main() {
	app := fiber.New()
	db := database.GetDatabase()
	setupRoutes(app)
	app.Listen(":3000")
	defer db.Close()
}
