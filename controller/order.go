package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"one-api/model"
	"strconv"
	"time"
)

func GetAllOrder(c *gin.Context) {
	id := c.GetInt("id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}
	orders, total, err := model.GetAllOrders(id, page, 10)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    orders,
		"total":   total,
	})
	return
}

func AddOrder(c *gin.Context) {
	id := c.GetInt("id")
	order := model.Order{}
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	orderNumber := generateOrderNumber()

	date := time.Now().Format("2006-01-02 15:04:05")

	cleanOrder := model.Order{
		OrderNo:   orderNumber,
		CreatTime: date,
		State:     0,
		UserId:    id,
		Price:     order.Price,
	}
	err = cleanOrder.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    cleanOrder,
	})
	return
}

// 生成订单号
func generateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())

	// 获取当前日期
	date := time.Now().Format("20060102")

	// 生成8位随机数
	randomNum := fmt.Sprintf("%08d", rand.Intn(100000000))

	// 拼接订单号
	orderNumber := date + randomNum

	return orderNumber
}

// 更新订单
//func UpdateOrder(orderNo string) {
//	order := model.Order{
//		OrderNo: orderNo,
//		State:   1,
//	}
//
//	err := order.Update()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//}
