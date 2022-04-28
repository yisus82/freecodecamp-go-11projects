package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	client *mongo.Client
	db     *mongo.Database
}

var mongoInstance MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017/" + dbName

type Employee struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Salary float64            `json:"salary" bson:"salary"`
	Age    int                `json:"age" bson:"age"`
}

func connect() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	mongoInstance.client = client
	mongoInstance.db = client.Database(dbName)

	return nil
}

func checkRequiredFields(employee Employee) bool {
	if employee.Name == "" || employee.Salary == 0 || employee.Age == 0 {
		return false
	}
	return true
}

func main() {
	if err := connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employees", func(c *fiber.Ctx) error {
		employees := make([]Employee, 0)
		collection := mongoInstance.db.Collection("employees")
		query := bson.D{}

		cursor, err := collection.Find(c.Context(), query)
		if err != nil {
			return err
		}

		err = cursor.All(c.Context(), &employees)
		if err != nil {
			return err
		}

		return c.JSON(employees)
	})

	app.Get("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var employee Employee
		collection := mongoInstance.db.Collection("employees")

		employeeID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(400).SendString("Invalid ID")
		}

		query := bson.D{primitive.E{Key: "_id", Value: employeeID}}
		result := collection.FindOne(c.Context(), query)
		if result.Err() != nil {
			return c.Status(404).SendString("Employee not found")
		}

		err = result.Decode(&employee)
		if err != nil {
			return err
		}

		return c.JSON(employee)
	})

	app.Post("/employees", func(c *fiber.Ctx) error {
		var employee Employee

		err := c.BodyParser(&employee)
		if err != nil || !checkRequiredFields(employee) {
			return c.Status(400).SendString("Error parsing body")
		}

		collection := mongoInstance.db.Collection("employees")
		result, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return err
		}

		employee.ID = result.InsertedID.(primitive.ObjectID)
		c.Location("/employees/" + employee.ID.Hex())
		return c.Status(201).JSON(employee)
	})

	app.Put("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var employee Employee

		err := c.BodyParser(&employee)
		if err != nil || !checkRequiredFields(employee) {
			return c.Status(400).SendString("Error parsing body")
		}

		collection := mongoInstance.db.Collection("employees")
		employeeID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(400).SendString("Invalid ID")
		}

		employee.ID = employeeID

		query := bson.D{primitive.E{Key: "_id", Value: employeeID}}
		update := bson.D{
			primitive.E{
				Key: "$set",
				Value: bson.D{
					primitive.E{Key: "name", Value: employee.Name},
					primitive.E{Key: "salary", Value: employee.Salary},
					primitive.E{Key: "age", Value: employee.Age},
				},
			},
		}

		result, err := collection.UpdateOne(c.Context(), query, update)
		if err != nil {
			return err
		}

		if result.MatchedCount == 0 {
			return c.Status(404).SendString("Employee not found")
		}

		return c.JSON(employee)
	})

	app.Delete("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		collection := mongoInstance.db.Collection("employees")

		employeeID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(400).SendString("Invalid ID")
		}

		query := bson.D{primitive.E{Key: "_id", Value: employeeID}}

		_, err = collection.DeleteOne(c.Context(), query)
		if err != nil {
			return err
		}

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
