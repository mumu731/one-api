package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type SearchResult struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Snippet string `json:"snippet"`
}

func search(query string) (map[string]interface{}, error) {

	response, err := http.Get(fmt.Sprintf("https://api.valueserp.com/search?api_key=E6F8AFB946AC47738637FE5E83F5850B&q=%s&hl=zh-cn&gl=cn", url.QueryEscape(query)))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error:", err)
	}

	var data map[string]interface{}
	// 解析JSON数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// 将map转换为JSON字符串
	//jsonStr, err := json.Marshal(data)
	//if err != nil {
	//	return "错误", err
	//}
	//
	//return string(jsonStr), nil

	return data, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type SearchRequest struct {
	Query string `json:"query"`
}

func GetSearchResult(c *gin.Context) {
	var searchRequest SearchRequest
	err := json.NewDecoder(c.Request.Body).Decode(&searchRequest)

	query := searchRequest.Query

	result, err := search(query)
	if err != nil {
		fmt.Println(err)
		return
	}
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
		"data":    result,
	})

}
