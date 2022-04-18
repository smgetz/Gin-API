package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

//book struct, formatted with lowercase for json data
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

//Moc Data creating a slice of book
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

//Query id param
//check id or not if not reutrn message
//get book by id if exists
//check quant if 0 not available
//else return okay, book
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID Query Param"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

//Same as checkout but adding to quant
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID Query Param"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// Gin Context is info of request and allows to send request
//books is data
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

//Helper that will actually get book by id
//gin.H custom json error message
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

//Get book by ID
//Returning a pointer to that book
//Range over b to identify by i, if found return i else not found
//Error if book does/doesnt exist
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book Not Found")
}

//Creates book via POST request in main, appended to the slice books of book struct
func createBook(c *gin.Context) {
	var newBook book

	//bind payload to newBook of book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// Create Router from Gin for diff routes
//Path = updating
func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
