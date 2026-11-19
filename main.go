package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID     int64   `json:"id" binding: "required"`
	Name   string  `json: "name" binding: "required"`
	Dept   string  `json: "dept" binding: "required"`
	Salary float64 `json: "salary" binding: "required"`
}

type EmployeeRequest struct {
	Name   string  `json: "name" binding: "required"`
	Dept   string  `json: "dept" binding: "required"`
	Salary float64 `json: "salary" binding: "required"`
}

var (
	db     = make(map[int64]Employee)
	lock   sync.RWMutex
	nextId = 1
)

func main() {

	r := gin.Default()

	r.GET("/health", getHealth)
	r.GET("/employee", getAll)
	r.POST("/employee", create)

	r.Run(":8000")

}

func getHealth(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "healthy, up and running"})

}

func getAll(c *gin.Context) {

	lock.Lock()
	defer lock.Unlock()

	count := len(db)
	if count == 0 {
		fmt.Println("Count is Zero")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data is Null/Empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": db})

}

func create(c *gin.Context) {
	var request EmployeeRequest

	err :=
		c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request Error"})
		return
	}
	var newEmp Employee
	newEmp.Dept = request.Dept
	newEmp.ID = int64(nextId + 1)
	newEmp.Name = request.Name
	newEmp.Salary = request.Salary

	db[newEmp.ID] = newEmp

	c.JSON(http.StatusCreated, gin.H{"message": "created", "data": newEmp})

}
