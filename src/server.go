package main

import (
	"fmt"
	"gofiber_postgres/src/models"
	"gofiber_postgres/src/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct{
	Author		string	`json:"author"`
	Title			string	`json:"title"`
	Publisher string	`json:"publisher"`
}

type Repository struct{
	DB *gorm.DB
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}
	err := c.BodyParser(&book)

	if err!=nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
			return err
	}

	err = r.DB.Create(&book).Error
	if err!=nil{
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create book"})
			return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book added successfully"} )

	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"could not get books",
			})
			return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"fetched successfully",
			"data": bookModels,
		})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Book{}
	id := c.Params("id")
	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message" : "id cannot be empty",
			})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"could not delete the book",
			})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"Deleted book successfully",
		})
	return nil
}

func (r *Repository) GetBookById(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Book{}
	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"id cannot be empty"})
		return nil
	}

	fmt.Println("the id is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"no book found with the id"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message" : "book found with the id",
			"data" : bookModel,
		})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:	os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSL"),
		DBName:	os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil{
		log.Fatal("could not load the database")
	}
	 
	err = models.MigrateBooks(db)
	if err != nil{
		log.Fatal("could not migrate books")
	}

	r := Repository{
		DB : db,
	}

	app := fiber.New()
	r.SetupRoutes(app)

	app.Listen(":8080")
}