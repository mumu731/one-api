package controller

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ReadTxt(c *gin.Context) {
	fileURL := "https://diandi-app.oss-cn-hangzhou.aliyuncs.com/other/beiying.txt" // 替换为你的网络文件URL

	// 发送GET请求获取网络文件内容
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("无法获取文件:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无法获取文件",
		})
		return
	}
	defer resp.Body.Close()

	// 创建一个Scanner来读取文件内容
	scanner := bufio.NewScanner(resp.Body)

	// 存储文本行的切片
	var lines []string

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 检查是否有错误发生
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时发生错误:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无法获取文件",
		})
		return
	}

	// 将文本行连接起来，并移除空格和换行符
	compressedText := strings.Join(lines, "")
	compressedText = strings.ReplaceAll(compressedText, " ", "")
	compressedText = strings.ReplaceAll(compressedText, "\n", "")

	// 分割文本
	var splitText []string
	for i := 0; i < len(compressedText); i += 200 {
		end := i + 200
		if end > len(compressedText) {
			end = len(compressedText)
		}
		splitText = append(splitText, compressedText[i:end])
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    splitText,
	})
}
