package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/deversmann/goaddr/log"
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
var (
	dialect = "sqlite"
	dsn     = "file::memory:?cache=shared"
	port    = "8080"
)

func main() {
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
	log.Info.Println("Server running and listening on port ", port)
	r.Run(fmt.Sprintf(":%s", port))
}

func initConfig() {
	logLevel := getenv("GOADDR_LOGLEVEL", "")
	switch logLevel {
	case "DEBUG":
		log.Debug.SetOutput(os.Stderr)
	case "NONE":
		log.Info.SetOutput(ioutil.Discard)
	}
	log.Debug.Println("Debugging log on")
	log.Info.Println("Info log on")

	dialect = getenv("GOADDR_DBDIALECT", dialect)
	dsn = getenv("GOADDR_DBDSN", dsn)
	port = getenv("GOADDR_PORT", port)
	log.Debug.Println("Using dialect =", dialect)
	log.Debug.Println("Using dsn     =", dsn)
	log.Debug.Println("Using port    =", port)
}

func initDB() {
	log.Debug.Println("Connecting to database")
	var dialector gorm.Dialector
	switch dialect {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgresql":
		dialector = postgres.Open(dsn)
	default:
		log.Info.Fatalf("Unknown/unimplemented database dialect: %s", dialect)
	}

	database, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}) // turning off GORM's internal logging
	if err != nil {
		log.Info.Fatal("Unrecoverable error opening database:", err)
	}
	db = database
	if err := db.AutoMigrate(&Contact{}); err != nil {
		log.Info.Fatal("Unrecoverable error migrating database:", err)
	}

	var count int64
	db.Model(&Contact{}).Count(&count)
	if count == 0 {
		log.Info.Println("Empty database detected, creating sample entry.")
		db.Create(&sampleCon)
	}
}

func readContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact
	if err := db.First(&con, id).Error; err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	log.Debug.Println("Successfully retrieved contact with id: ", id)
	c.IndentedJSON(http.StatusOK, con)
}

func readContacts(c *gin.Context) {
	var cons []Contact
	tx, err := parseQuery(c, db)
	if err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	if tx.Find(&cons).RowsAffected == 0 {
		log.Debug.Println("No results returned")
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no contacts found"})
		return
	}
	log.Debug.Printf("Successfully retrieved %d contacts", len(cons))
	c.IndentedJSON(http.StatusOK, gin.H{"contacts": cons})
}

func createContact(c *gin.Context) {
	var newCon Contact
	if err := c.BindJSON(&newCon); err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON for contact"})
		return
	}
	db.Create(&newCon)
	log.Debug.Println("Successfully created contact with id: ", newCon.ID)
	c.IndentedJSON(http.StatusCreated, newCon)
}

func updateContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	origID := con.ID
	if err := c.BindJSON(&con); err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON for contact"})
		return
	}
	if origID != con.ID {
		log.Debug.Println("Record ID mismatch")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Cannot modify ID"})
		return
	}
	db.Save(&con)
	log.Debug.Println("Successfully updated contact with id: ", con.ID)
	c.JSON(http.StatusAccepted, con)
}

func deleteContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		log.Debug.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	db.Delete(con)
	log.Debug.Println("Successfully deleted contact with id: ", con.ID)
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

func parseQuery(c *gin.Context, tx *gorm.DB) (*gorm.DB, error) {
	// sorting first
	if sortBy := c.Query("sort_by"); sortBy != "" {
		log.Debug.Printf("- sort_by: %s\n", sortBy)
		for _, chunk := range strings.Split(sortBy, ",") {
			asRunes := []rune(chunk)
			if asRunes[0] == '-' {
				log.Debug.Printf("- Order: %s %s\n", string(asRunes[1:]), "desc")
				tx = tx.Order(fmt.Sprintf("%s %s", string(asRunes[1:]), "desc"))
			} else {
				log.Debug.Printf("- Order: %s\n", chunk)
				tx = tx.Order(chunk)
			}
		}
	}
	// page size
	if limit := c.Query("limit"); limit != "" {
		iLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Debug.Printf("invalid value for limit: %s", limit)
			return tx, fmt.Errorf("invalid value for limit: %s", limit)
		}
		log.Debug.Printf("- Limit: %s", limit)
		tx = tx.Limit(iLimit)
	}
	// page number
	if offset := c.Query("offset"); offset != "" {
		iOffset, err := strconv.Atoi(offset)
		if err != nil {
			log.Debug.Printf("invalid value for offset: %s", offset)
			return tx, fmt.Errorf("invalid value for offset: %s", offset)
		}
		log.Debug.Printf("- Offset: %s", offset)
		tx = tx.Offset(iOffset)
	}
	return tx, nil
}
