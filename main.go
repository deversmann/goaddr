package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	viper.SetDefault("DBDialect", "sqlite")
	viper.SetDefault("DBDSN", "file::memory:?cache=shared")
	viper.SetDefault("Port", "8080")
	viper.SetConfigName("goaddr-config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("goaddr")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	for k, v := range viper.AllSettings() {
		fmt.Printf("  %s=%s\n", k, v)
	}

	r := gin.Default()

	connectDatabase()

	gin.EnableJsonDecoderDisallowUnknownFields() // error on undefined JSON fields
	v1 := r.Group("/api/v1/contacts")
	{
		v1.POST("", createContact)
		v1.GET("", readContacts)
		v1.GET("/:id", readContact)
		v1.PUT("/:id", updateContact)
		v1.DELETE("/:id", deleteContact)
	}
	r.Run(fmt.Sprintf(":%s", viper.Get("Port")))
}

func connectDatabase() {
	var dialector gorm.Dialector
	dialect := viper.Get("DBDialect").(string)
	dsn := viper.Get("DBDSN").(string)

	switch dialect {
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "postgresql":
		dialector = postgres.Open(dsn)
	default:
		panic(fmt.Sprint("Unknown/unimplemented database dialect: ", dialect))
	}

	database, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&Contact{})
	db = database

	var count int64
	db.Model(&Contact{}).Count(&count)
	if count == 0 {
		db.Create(&sampleCon)
	}
}

func readContact(c *gin.Context) {
	id := c.Param("id")
	var con Contact

	if err := db.First(&con, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, con)
}

func readContacts(c *gin.Context) {
	var cons []Contact
	if db.Find(&cons).RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no contacts found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"contacts": cons})
}

func createContact(c *gin.Context) {
	var newCon Contact
	if err := c.BindJSON(&newCon); err != nil {
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}
	origID := con.ID
	if err := c.BindJSON(&con); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON for contact"})
		return
	}
	if origID != con.ID {
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "contact with id: " + id + " not found"})
		return
	}

	db.Delete(con)
	c.Status(http.StatusNoContent)
}
