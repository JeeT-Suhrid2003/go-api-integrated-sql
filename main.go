package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Book struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var db *gorm.DB

func initDB() {
	// connect to MySQL running in Docker
	dsn := "root:root@tcp(127.0.0.1:3306)/booksdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// auto-migrate schema (creates table if not exists)
	db.AutoMigrate(&Book{})
}

func getBooks(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	db.Create(&newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found!"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight (OPTIONS request)
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found!"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not available!"})
		return
	}

	book.Quantity--
	db.Save(&book)
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameter"})
		return
	}

	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found!"})
		return
	}

	book.Quantity++
	db.Save(&book)
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	initDB()
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)

	router.Run("localhost:8080")

}
