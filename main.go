package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Contact struct {
	ID        uint   `json:"id" gorm:"primary_key, AUTO_INCREMENT"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zipcode   string `json:"zipcode"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

var sampleCon = Contact{
	Firstname: "Paul",
	Lastname:  "Cormir",
	Address:   "100 E. Davie St.",
	City:      "Raleigh",
	State:     "NC",
	Zipcode:   "27601",
	Phone:     "888-RED-HAT-1",
	Email:     "pcormir@redhat.com",
}

var db *gorm.DB

// default values will be overwritten as needed
var dialect = "sqlite"
var dsn = "file::memory:?cache=shared"
var port = "8080"

func main() {
	log.SetPrefix("[GOADDR] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	initConfig()
	initDB()

	r := gin.Default()
	gin.EnableJsonDecoderDisallowUnknownFields() // error on undefined JSON fields
	apiRoot := r.Group("/api")
	v1 := apiRoot.Group("/v1")
	contactsRoute := v1.Group("/contacts")
	{
		contactsRoute.POST("", createContact)
		contactsRoute.GET("", readContacts)
		contactsRoute.GET("/:id", readContact)
		contactsRoute.PUT("/:id", updateContact)
		contactsRoute.DELETE("/:id", deleteContact)
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func initConfig() {
	dialect = getenv("GOADDR_DBDIALECT", dialect)
	dsn = getenv("GOADDR_DBDSN", dsn)
	port = getenv("GOADDR_PORT", port)
	log.Println("Using dialect =", dialect)
	log.Println("Using dsn     =", dsn)
	log.Println("Using port    =", port)
}

func initDB() {
	var dialector gorm.Dialector
	switch dialect {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgresql":
		dialector = postgres.Open(dsn)
	default:
		log.Fatalf("Unknown/unimplemented database dialect: %s", dialect)
	}

	database, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}) // turning off GORM's internal logging
	if err != nil {
		log.Fatal(err)
	}
	db = database
	if err := db.AutoMigrate(&Contact{}); err != nil {
		log.Fatal(err)
	}

	var count int64
	db.Model(&Contact{}).Count(&count)
	if count == 0 {
		log.Println("Empty database detected, creating sample entry.")
		db.Create(&sampleCon)
	}
}

func readContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact
	if err := db.First(&con, id).Error; err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, con)
}

func readContacts(c *gin.Context) {
	var cons []Contact
	if db.Find(&cons).RowsAffected == 0 {
		log.Println("No results returned")
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no contacts found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"contacts": cons})
}

func createContact(c *gin.Context) {
	var newCon Contact
	if err := c.BindJSON(&newCon); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON for contact"})
		return
	}
	db.Create(&newCon)
	c.IndentedJSON(http.StatusCreated, newCon)
}

func updateContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	origID := con.ID
	if err := c.BindJSON(&con); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON for contact"})
		return
	}
	if origID != con.ID {
		log.Println("Record ID mismatch")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Cannot modify ID"})
		return
	}
	db.Save(&con)
	c.JSON(http.StatusAccepted, con)
}

func deleteContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	db.Delete(con)
	c.Status(http.StatusNoContent)
}

// returns the environment variable or a fallback if not set
func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return fallback
	}
	return val
}
