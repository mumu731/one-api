package controller

// mj绘画

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"one-api/model"
	"strconv"
	"time"
)

const ECHOAIKEY = "6e6e88600866799ce190c59e5b3717ca37b16980"

type ReqData struct {
	Prompt string `json:"prompt"`
}

type ImagResData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MessageId string `json:"message_id"`
	} `json:"data"`
}

// Imagine 以文生图
func Imagine(c *gin.Context) {
	id := c.GetInt("id")
	url := fmt.Sprintf("https://gapi.kaiecho.com/api/midj/imagine")
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ECHOAIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
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
	var data ImagResData
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode == http.StatusOK {
		// 插入数据库
		var dbdata ReqData
		err = json.Unmarshal(reqBody, &dbdata)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		date := time.Now().Format("2006-01-02 15:04:05")
		linkMidjourney := model.Midjourney{
			UserId:    id,
			CreatedAt: date,
			Prompt:    dbdata.Prompt,
			Status:    "PROCESSING",
			MessageId: data.Data.MessageId,
		}
		err = linkMidjourney.MidjourneyInsert()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": data})

		// 消费token
		logContent := fmt.Sprintf("Mdjourney绘图固定5000")
		model.RecordConsumeLog(id, 0, 0, "Mdjourney", "Mdjourney以文生图", 5000, logContent)
		model.UpdateUserUsedQuotaAndRequestCount(id, 5000)
		channelId := c.GetInt("channel_id")
		model.UpdateChannelUsedQuota(channelId, 5000)
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
	}
}

// ImagMessage 消息查询
func ImagMessage(c *gin.Context) {
	messageId := c.Query("messageId")
	url := fmt.Sprintf("https://gapi.kaiecho.com/api/midj/message/%s", messageId)
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ECHOAIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
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
	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": data})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
	}
}

// ImagButton 单张放大
func ImagButton(c *gin.Context) {
	url := fmt.Sprintf("https://gapi.kaiecho.com/api/midj/button")
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ECHOAIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
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
	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": data})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
	}
}

type RequestBody struct {
	Data struct {
		MessageId string `json:"message_id"`
		Body      struct {
			ImageUrl string `json:"image_url"`
		}
		Status string `json:"status"`
	}
}

// ImagNotify 消息回调
func ImagNotify(c *gin.Context) {
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	var body RequestBody
	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		return
	}

	// 打印参数
	fmt.Println("数据:", body.Data.MessageId)
	fmt.Println("数据:", body.Data.Body.ImageUrl)
	date := time.Now().Format("2006-01-02 15:04:05")
	linkMidjourney := model.Midjourney{
		ImageUrl:  body.Data.Body.ImageUrl,
		Status:    body.Data.Status,
		MessageId: body.Data.MessageId,
		UpdateAt:  date,
	}

	err = linkMidjourney.MidjourneyUpdate()
	if err != nil {
		return
	}

	return
}

// TextTranslate 翻译
func TextTranslate(c *gin.Context) {
	url := fmt.Sprintf("https://api.175ai.cn/prompt/translate/")
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
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
	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": data})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
	}
}

// GetAllImage 绘图列表
func GetAllImage(c *gin.Context) {
	id := c.GetInt("id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}
	pictures, total, err := model.GetAllPicture(id, page, 10)
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
		"data":    pictures,
		"total":   total,
	})
	return
}
