package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type SearchResult struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Snippet string `json:"snippet"`
}

func search(query string) (string, error) {

	region := "zh-cn"
	page := 1

	response, err := http.Get(fmt.Sprintf("https://duckduckgo.com/?q=%s", url.QueryEscape(query)))
	if err != nil {
		fmt.Println("Error:", err)
		return "错误", err
	}
	defer response.Body.Close()

	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "错误", err
	}

	regex := regexp.MustCompile(`vqd=["']([^"']+)["']`)
	match := regex.FindStringSubmatch(string(html))
	var vqd string
	if len(match) > 1 {
		vqd = strings.ReplaceAll(strings.ReplaceAll(match[1], `"`, ""), "'", "")
	}

	safeSearchBase := map[string]int{"On": 1, "Moderate": -1, "Off": -2}
	PAGINATION_STEP := 25

	res, err := http.Get(fmt.Sprintf("https://links.duckduckgo.com/d.js?q=%s&l=%s&p=%d&s=%d&df=%d&o=json&vqd=%s",
		url.QueryEscape(query), region, safeSearchBase["On"], max(PAGINATION_STEP*(page-1), 0), getTimeMillis(), vqd))
	if err != nil {
		fmt.Println("Error:", err)
		return "错误", err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error:", err)
		return "错误", err
	}

	referenceResults := make([][]interface{}, 0)
	if results, ok := result["results"].([]interface{}); ok {
		for _, r := range results {
			row := r.(map[string]interface{})
			if n, ok := row["n"]; !ok || n == nil {
				if body, ok := row["a"].(string); ok {
					referenceResults = append(referenceResults, []interface{}{body, row["u"]})
					if len(referenceResults) > 2 {
						break
					}
				}
			}
		}
	}

	resultData, err := json.Marshal(referenceResults)
	if err != nil {
		fmt.Println("Error:", err)
		return "错误", err
	}
	return string(resultData), nil
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
