package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Snippet string `json:"snippet"`
}

func search(query string) ([]SearchResult, error) {
	// 构建搜狗搜索 URL
	searchURL := fmt.Sprintf("https://www.so.com/s?q=%s", url.QueryEscape(query))

	// 发送 HTTP GET 请求并获取响应
	response, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	// 解析搜索结果
	results := make([]SearchResult, 0)
	doc.Find(".result .res-title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		linkElement := s.Find("a[href]")
		href, _ := linkElement.Attr("href")
		snippet := s.Parent().Find(".res-desc").Text()

		result := SearchResult{
			Title:   strings.TrimSpace(title),
			Href:    strings.TrimSpace(href),
			Snippet: strings.TrimSpace(snippet),
		}
		results = append(results, result)
	})

	return results, nil
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
