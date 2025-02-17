package API

import (
	"OPS/module/Consumer"
	"OPS/module/DatabaseConn"
	db "OPS/module/DatabaseConn"

	"OPS/module/Orders"
	"OPS/module/Producer"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var brokers = []string{"localhost:9092"}
var groupID = "order-consumer-group"
var wg sync.WaitGroup

const topic = "orders"

func saveItem(b *gin.Context) {

	o1 := Orders.Order{}

	if err1 := b.ShouldBindJSON(&o1); err1 != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	_, err2 := db.AddRow(DatabaseConn.DbPool, o1)
	if err2 != nil {
		//fmt.Errorf("Facing issues while adding a new row %v", err2)
		log.Fatalf("Error while adding row %v", err2)
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Add Request"})
		return
	}

	b.JSON(http.StatusOK, gin.H{"message": "Successfully Loaded"})
	k, err3 := db.GetAll(DatabaseConn.DbPool)
	if err3 != nil {
		log.Fatalf("Error while fetching the orders %v", err3)
		return
	}

	Orders.OrderList = append(Orders.OrderList, o1)
	fmt.Println(Orders.OrderList)
	b.IndentedJSON(http.StatusOK, k)
	fmt.Println(o1.Product_Id)
	eve, err := db.AddOrder(DatabaseConn.DbPool, o1)
	if err != nil {
		fmt.Errorf("Unable to add an Order, Please check %v ", err)
		b.JSON(http.StatusNotFound, err)
		return
	}

	wg.Add(1)
	go Consumer.ConsumeMessages(brokers, topic, groupID, &wg)

	wg.Add(1)
	go Producer.PublishMessage(brokers, topic, eve, &wg)

	wg.Wait()
}

func StartServer() {
	router := gin.Default()
	router.POST("/Orders", saveItem)
	//router.GET("/Orders", func(b *gin.Context) { getItem(b) })
	router.Run("localhost:8080")
}
