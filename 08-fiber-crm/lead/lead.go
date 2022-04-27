package lead

import (
	"08-fiber-crm/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

var db *gorm.DB

func init() {
	db = database.GetDatabase()
	db.AutoMigrate(&Lead{})
}

func GetLeads(c *fiber.Ctx) error {
	var leads []Lead
	err := db.Find(&leads).Error
	c.JSON(leads)
	return err
}

func GetLead(c *fiber.Ctx) error {
	var lead Lead
	db.First(&lead, c.Params("id"))
	if lead.ID == 0 {
		c.Status(404).SendString("Lead not found")
		return nil
	}
	c.JSON(lead)
	return nil
}

func CreateLead(c *fiber.Ctx) error {
	var lead Lead
	err := c.BodyParser(&lead)
	if err != nil || !checkLeadRequiredFields(lead) {
		c.Status(400).SendString("Error parsing body")
		return nil
	}
	db.Create(&lead)
	c.Location("/api/v1/leads/" + strconv.FormatUint(uint64(lead.ID), 10))
	c.Status(201).JSON(lead)
	return err
}

func DeleteLead(c *fiber.Ctx) error {
	db.Delete(&Lead{}, c.Params("id"))
	return nil
}

func checkLeadRequiredFields(lead Lead) bool {
	if lead.Name == "" || lead.Company == "" || lead.Email == "" || lead.Phone == "" {
		return false
	}
	return true
}
