package main

import (
   "fmt"
   "log"
   "os"
	 "strings"
	 "database/sql" 
	 "github.com/lib/pq" 

   "github.com/gofiber/fiber/v2"
)

type Node struct {
	Id				int `json:"id"`
	Info 			string `json:"info"`
	Children 	[]string `json:"children"`
}


func rootHandler(c *fiber.Ctx, db *sql.DB) error {
	rows, err:= db.Query("SELECT * FROM public.\"Node3\"")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("an unexpected error occured")
	}

	var nodes []Node
	for rows.Next() {
		var a int
		var b string
		var c []string
		rows.Scan(&a, &b, pq.Array(&c))
		fmt.Println(c)
		nodes = append(nodes, Node{Id: a, Info: b, Children: c})
	}

	return c.JSON(nodes)
}

func createHandler(c *fiber.Ctx, db *sql.DB) error {
	newNode:= Node{}
	{
		err := c.BodyParser(&newNode)
		if err != nil {
			log.Printf("An error occured: %v", err)
			return c.SendString(err.Error())
		}
	}

	_, err := db.Exec("INSERT into public.\"Node3\" (info, children) VALUES ($1, $2)", newNode.Info, "{" + strings.Join(newNode.Children, ",") + "}")
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)
	}
	return c.SendStatus(200)
}

func updateHandler(c *fiber.Ctx, db *sql.DB) error {
	id:= c.Query("id")
	newNode:= Node {}
	{
		err:= c.BodyParser(&newNode)
		if err != nil {
			log.Printf("An error occured: %v", err)
			return c.SendString(err.Error())
		}
	}

	_, err := db.Exec("UPDATE public.\"Node3\" SET info=$1, children=$2 WHERE id=$3", newNode.Info, "{" + strings.Join(newNode.Children, ",") + "}", id)
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)
	}
	return c.SendStatus(200)
}


func main() {
	uri := "postgresql://postgres:test@localhost/dynalist?sslmode=disable"
	db, err := sql.Open("postgres", uri)

	if err != nil {
		log.Fatal(err)
	}
   
	app := fiber.New()
   
	 
	app.Get("/", func(c *fiber.Ctx) error {
		return rootHandler(c, db)
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		return createHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return updateHandler(c, db)
	})
  	 
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}