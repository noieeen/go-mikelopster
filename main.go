package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"Price"`
	Category string `json:"category"`
	Quantity int    `json:"quantity"`
}

var db *sql.DB

func main() {
	app := fiber.New()
	db = setupDatabase()
	defer db.Close()

	// Setup route
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// api/health
	api.Options("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	// api/v1/products
	v1.Post("/products", createProduct)

	log.Fatal(app.Listen(":3000"))

	// Init Table
	// if err := initDB(db); err != nil {
	// 	log.Fatal(err)
	// }
}

func setupDatabase() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected!")
	return db
}

func initDB(db *sql.DB) error {
	// products
	if _, err := db.Exec(`DROP TABLE IF EXISTS public.products`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS public.products
	(
		id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
		name text COLLATE pg_catalog."default",
		price integer DEFAULT 0,
		category text COLLATE pg_catalog."default",
		quantity integer DEFAULT 0
	)`); err != nil {
		return err
	}

	return nil
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)

	if err := c.BodyParser(product); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO public.products(
		name, price, category, quantity)
		VALUES ($1, $2, $3, $4);`, product.Name, product.Price, product.Category, product.Quantity); err != nil {
		return err
	}

	return c.SendString("Create " + product.Name + " Successfully!")
}
