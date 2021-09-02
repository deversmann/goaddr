package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Contact struct {
	ID        uint   `json:"id" gorm:"primary_key, AUTO_INCREMENT"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zipcode"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

var sampleCon = Contact{FirstName: "Paul", LastName: "Cormir", Address: "100 E. Davie St.", City: "Raleigh", State: "NC", ZipCode: "27601", Phone: "888-RED-HAT-1", Email: "pcormir@redhat.com"}
var db *gorm.DB

func main() {
	r := gin.Default()

	connectDatabase()

	r.GET("/contacts/:id", readContact)
	r.GET("/contacts", readContacts)
	r.POST("/contacts", createContact)
	r.PUT("/contacts/:id", updateContact)
	r.DELETE("/contacts/:id", deleteContact)
	r.Run()
}

func connectDatabase() {
	database, err := gorm.Open("sqlite3", "test.db")

	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&Contact{})
	db = database

	var count int
	db.Model(&Contact{}).Count(&count)
	if count == 0 {
		db.Create(&sampleCon)
	}
}

func readContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "contact with id: " + id + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, con)
}

func readContacts(c *gin.Context) {
	var cons []Contact
	db.Find(&cons)
	c.IndentedJSON(http.StatusOK, gin.H{"contacts": cons})
}

func createContact(c *gin.Context) {
	var newCon Contact
	if err := c.BindJSON(&newCon); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON for contact"})
		return
	}
	db.Create(&newCon)
	c.IndentedJSON(http.StatusCreated, newCon)
}

func updateContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "contact with id: " + id + " not found"})
		return
	}
	origID := con.ID
	if err := c.BindJSON(&con); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON for contact"})
		return
	}
	if origID != con.ID {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Cannot modify ID"})
		return
	}
	db.Save(&con)
	c.JSON(http.StatusOK, con)
}

func deleteContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "contact with id: " + id + " not found"})
		return
	}

	db.Delete(con)
	c.Status(http.StatusNoContent)
}
