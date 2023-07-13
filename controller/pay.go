package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"one-api/model"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GetPayUrl 易支付-支付请求
func GetPayUrl(c *gin.Context) {
	orderNo := c.Query("orderNo")
	urls := "https://www.u3b.net/mapi.php"
	params := url.Values{}
	for key, value := range FormPayQuery(orderNo) {
		params.Add(key, value)
	}
	resp, err := http.PostForm(urls, params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var data map[string]interface{}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
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
		"data":    data,
	})

	return
}

// FormPayQuery 易支付-支付参数
func FormPayQuery(orderNo string) map[string]string {
	// 定义参数
	params := map[string]string{
		"pid":          "1182",
		"type":         "alipay",
		"out_trade_no": SearchOrderss(orderNo).OrderNo,
		"notify_url":   "http://182.44.52.201:15600/api/notify_url",
		"return_url":   "http://127.0.0.1:1002",
		"name":         "Tokens",
		"money":        strconv.FormatFloat(SearchOrderss(orderNo).Price, 'f', -1, 64),
		"clientip":     "192.168.1.100",
		"device":       "pc",
	}
	key := "QgWCn26z7uGQJEUwgXZn8rZ2gCiiiC7c"
	// 调用签名函数
	sign := GenerateSign(params, key)
	// 将签名添加到参数中
	params["sign"] = sign
	params["sign_type"] = "MD5"

	return params
}

// GenerateSign 易支付-生成签名
func GenerateSign(params map[string]string, key string) string {
	// 将参数名按ASCII码从小到大排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接参数字符串
	var paramStr string
	for _, k := range keys {
		paramStr += k + "=" + params[k] + "&"
	}
	paramStr = strings.TrimRight(paramStr, "&")
	// 拼接字符串与商户密钥进行MD5加密
	sign := fmt.Sprintf("%x", md5.Sum([]byte(paramStr+key)))
	return sign
}

// GetOrderNo 生成单号 当前日期+随机的5位数字
func GetOrderNo() string {
	rand.Seed(time.Now().UnixNano())
	return time.Now().Format("20060102150405") + fmt.Sprintf("%05v", rand.Int31n(100000))
}

// NotifyHandler 易支付-支付通知处理
func NotifyHandler(c *gin.Context) {
	// 读取请求参数
	tradeNo := c.Query("trade_no")
	outTradeNo := c.Query("out_trade_no")

	// 打印参数
	fmt.Println("易支付订单号:", tradeNo)
	fmt.Println("商户订单号:", outTradeNo)
	if tradeNo == "" || outTradeNo == "" {
		c.String(http.StatusOK, "fail")
		return
	}
	UpdateOrder(outTradeNo)
	UpdateBalance(outTradeNo)
	c.String(http.StatusOK, "success")
	return
}

// Order GetPayUrl 支付请求
type Order struct {
	Id         int     `json:"id"`
	OrderNo    string  `json:"orderNo" gorm:"unique;index"`
	CreatTime  string  `json:"creatTime" gorm:"not null"`
	PayTime    string  `json:"payTime" gorm:"default:null"`
	UpdateTime string  `json:"updateTime" gorm:"default:null"`
	State      int     `json:"state" gorm:"default:0"`
	UserId     int     `json:"userId" gorm:"default:0"`
	Price      float64 `json:"price" gorm:"default:0"`
	Remarks    string  `json:"remarks" gorm:"default:null"`
}

func SearchOrderss(orderNo string) Order {
	order, err := model.SearchOrders(orderNo)
	if err != nil {
		fmt.Println(err)
		return Order{}
	}
	fmt.Println(order)
	return Order(order)
}

// UpdateOrder 更新订单
func UpdateOrder(orderNo string) {
	date := time.Now().Format("2006-01-02 15:04:05")
	model.Update(orderNo, 1, date)

}

// GetPayUrl 易支付-查询支付状态
func GetPayAct(c *gin.Context) {
	orderNo := c.Query("orderNo")
	urls := "https://www.u3b.net/api.php?act=order&pid=1182&key=QgWCn26z7uGQJEUwgXZn8rZ2gCiiiC7c&out_trade_no=" + orderNo
	resp, err := http.Get(urls)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var data map[string]interface{}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
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
		"data":    data,
	})

	return
}

// UpdateBalance 更新用户余额
func UpdateBalance(orderNo string) {

	order, err := model.SearchOrders(orderNo)
	if err != nil {
		fmt.Println(err)
		return
	}

	//查询用户剩余金额
	user, err := model.GetUserById(order.UserId, false)

	balance := user.Balance + order.Price

	model.UpdateUserBalance(order.UserId, balance)

}
