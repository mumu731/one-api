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
	"sort"
	"strings"
	"time"
)

// GetPayUrl 支付请求
func GetPayUrl(c *gin.Context) {
	urls := "https://www.u3b.net/mapi.php"
	params := url.Values{}
	for key, value := range FormPayQuery() {
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

// FormPayQuery 支付参数
func FormPayQuery() map[string]string {
	// 定义参数
	params := map[string]string{
		"pid":          "1182",
		"type":         "alipay",
		"out_trade_no": GetOrderNo(),
		"notify_url":   "http://112.8.204.167:15678/api/notify_url",
		"return_url":   "http://112.8.204.167:15678/api/notify_url",
		"name":         "Tokens",
		"money":        "1",
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

// GenerateSign 生成签名
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

// NotifyHandler 支付通知处理
func NotifyHandler(c *gin.Context) {
	// 读取请求参数
	tradeNo := c.Query("trade_no")
	outTradeNo := c.Query("out_trade_no")

	// 打印参数
	fmt.Println("易支付订单号:", tradeNo)
	fmt.Println("商户订单号:", outTradeNo)

	c.String(http.StatusOK, "success")
	return
}
