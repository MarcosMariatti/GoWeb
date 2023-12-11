package main

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const datePattern = `^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$`

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyRequestCreate struct {
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
		products.POST("", CreateProduct())
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

// CreateProduct adds a new product to the slice productsList
func CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body BodyRequestCreate
		err := ctx.ShouldBind(&body)
		if err != nil {

			// incremental id
			id := len(productsList) + 1

			// check if the date has a valid format
			if validateDate(body.Expiration, datePattern) {
				date := strings.Split(body.Expiration, "/")

				// check if the day is between 1 and 31
				day, dayErr := strconv.Atoi(date[0])
				if dayErr != nil || (day < 0 || day > 31) {
					ctx.JSON(http.StatusBadRequest, map[string]any{"message": "invalid day"})
					return
				}

				// check if the day is between 1 and 31
				month, monthErr := strconv.Atoi(date[1])
				if monthErr != nil || (month < 0 || month > 12) {
					ctx.JSON(http.StatusBadRequest, map[string]any{"message": "invalid month"})
					return
				}

				year, yearErr := strconv.Atoi(date[2])
				if yearErr != nil || year < 0 {
					ctx.JSON(http.StatusBadRequest, map[string]any{"message": "invalid year"})
					return
				}

			} else {
				ctx.JSON(http.StatusBadRequest, map[string]any{"message": "invalid expiration date"})
				return
			}
			// check for unique code value
			for _, prod := range productsList {
				if body.CodeValue == prod.CodeValue {
					ctx.JSON(http.StatusBadRequest, map[string]any{"message": "code value is used"})
					return
				}
			}
			prod := Product{
				Id:          id,
				Name:        body.Name,
				Quantity:    body.Quantity,
				CodeValue:   body.CodeValue,
				IsPublished: body.IsPublished,
				Expiration:  body.Expiration,
				Price:       body.Price,
			}
			productsList = append(productsList, prod)

			ctx.JSON(http.StatusCreated, map[string]any{
				"message": "product created",
				"data": Product{
					Id:          prod.Id,
					Name:        prod.Name,
					Quantity:    prod.Quantity,
					CodeValue:   prod.CodeValue,
					IsPublished: prod.IsPublished,
					Expiration:  prod.Expiration,
					Price:       prod.Price,
				},
			})
		}
		ctx.JSON(http.StatusBadRequest, map[string]any{"message": "invalid request body"})
		return
	}

}

// validateDate check if a date string match a given pattern
func validateDate(date string, pattern string) bool {
	match, _ := regexp.MatchString(pattern, date)
	return match
}
