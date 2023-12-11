package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var productsList = []Product{}

func main() {
	loadProducts("products.json", &productsList)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET("", GetAllProducts())
		products.GET(":id", GetProduct())
		products.GET("/search", SearchProduct())
	}
	err := r.Run("localhost:8080")
	if err != nil {
		return
	}
}

// loadProducts carga los productos desde un archivo json
func loadProducts(path string, list *[]Product) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &list)
	if err != nil {
		panic(err)
	}
}

// GetAllProducts traer todos los productos almacenados
func GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, productsList)
	}
}

// GetProduct traer un producto por id
func GetProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
		for _, product := range productsList {
			if product.Id == id {
				ctx.JSON(200, product)
				return
			}
		}
		ctx.JSON(404, gin.H{"error": "product not found"})
	}
}

// SearchProduct traer un producto por nombre o categoria
func SearchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("priceGt")
		priceGt, err := strconv.ParseFloat(query, 32)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid price"})
			return
		}
		list := []Product{}
		for _, product := range productsList {
			if product.Price > priceGt {
				list = append(list, product)
			}
		}
		ctx.JSON(200, list)
	}
}
