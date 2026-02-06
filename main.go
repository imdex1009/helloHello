package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Item represents a simple item in our store
type Item struct {
	gorm.Model
	Name        string `json:"name"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

var DB *gorm.DB

func main() {
	var err error
	DB, err = gorm.Open(sqlite.Open("hello.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&Item{})

	// Seed dummy data if empty
	var count int64
	DB.Model(&Item{}).Count(&count)
	if count == 0 {
		items := []Item{
			{
				Name:        "Eco-Friendly Water Bottle",
				Price:       25,
				ImageURL:    "/public/images/water_bottle.png",
				Description: "Stay hydrated and save the planet with this durable, BPA-free water bottle. Made from 100% recycled materials, it keeps your drinks cold for 24 hours.",
			},
			{
				Name:        "Wireless Noise Cancelling Headphones",
				Price:       199,
				ImageURL:    "/public/images/headphones.png",
				Description: "Immerse yourself in music with our industry-leading noise cancelling technology. 30-hour battery life ensures you never miss a beat.",
			},
			{
				Name:        "Nano Banana (Limited Edition)",
				Price:       9999,
				ImageURL:    "/public/images/nano_banana.png",
				Description: "The future of fruit is here. Genetically engineered for perfect sweetness and zero bruising. Comes with a built-in ripeness indicator LED.",
			},
			{
				Name:        "Vintage Polaroid Camera",
				Price:       120,
				ImageURL:    "/public/images/camera.png",
				Description: "Capture moments instantly with this refurbished vintage Polaroid camera. Includes a pack of film to get you started.",
			},
			{
				Name:        "Smart Home Assistant Hub",
				Price:       89,
				ImageURL:    "/public/images/hub.png",
				Description: "Control your entire home with your voice. Compatible with all major smart home devices and features a crystal-clear speaker.",
			},
		}
		DB.Create(&items)
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/public", "./public")

	// GET / - Home (HTML)
	router.GET("/", func(context *gin.Context) {
		var items []Item
		DB.Find(&items)
		context.HTML(http.StatusOK, "index.html", items)
	})

	// GET /products/:id - Product Detail Page (HTML)
	router.GET("/products/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		// gorm.Model ID is uint
		id, err := strconv.Atoi(idParam)
		if err != nil {
			context.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		var item Item
		result := DB.First(&item, id)
		if result.Error != nil {
			context.String(http.StatusNotFound, "Product not found")
			return
		}

		context.HTML(http.StatusOK, "product.html", item)
	})

	// GET /items - List all items (API)
	router.GET("/items", func(context *gin.Context) {
		var items []Item
		DB.Find(&items)
		context.JSON(http.StatusOK, items)
	})

	// POST /items - Create a new item (API)
	router.POST("/items", func(context *gin.Context) {
		var newItem Item
		if err := context.ShouldBindJSON(&newItem); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		DB.Create(&newItem)
		context.JSON(http.StatusCreated, newItem)
	})

	// GET /items/:id - Get a specific item (API)
	router.GET("/items/:id", func(context *gin.Context) {
		idParam := context.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var item Item
		result := DB.First(&item, id)
		if result.Error != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		context.JSON(http.StatusOK, item)
	})

	// Run the server on port 8080
	router.Run(":8080")
}
