package API

import (
	db "OPS/module/DatabaseConn"
	"OPS/module/Orders"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func saveItem(b *gin.Context) {
	dbPool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Login Issue")
		return
	}

	o1 := Orders.Order{User_Id: 30, Product_Id: 30, Quantity: 39, Total_Price: 23.25}

	if err1 := b.ShouldBindJSON(&o1); err1 != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	_, err2 := db.AddRow(dbPool, o1)
	if err2 != nil {
		fmt.Errorf("Facing issues while adding a new row")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Add Request"})
		return
	}

	b.JSON(http.StatusOK, gin.H{"message": "Successfully Loaded"})
	k, err3 := db.GetAll(dbPool)
	if err3 != nil {
		log.Fatalf("Error while fetching the orders")
		return
	}

	b.IndentedJSON(http.StatusOK, k)

}

func StartServer() {
	router := gin.Default()
	router.POST("/Orders", func(b *gin.Context) { saveItem(b) })
	//router.GET("/Orders", func(b *gin.Context) { getItem(b) })
	router.Run("localhost:8080")
}
