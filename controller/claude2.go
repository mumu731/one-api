package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	claude "github.com/all-in-aigc/claude-webapi"
	"github.com/gin-gonic/gin"
	"net/http"
	"one-api/common"
	"strings"
)

type ClauRequestBody struct {
	Text        string                   `json:"text"`
	Attachments []map[string]interface{} `json:"attachments"`
}

var (
	baseUri    string = "https://claude-proxy.f91.dev"
	orgid      string = "f597cb2e-9643-4d03-82ae-22e24bc8fd03"
	sessionKey string = "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"
	userAgent  string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	debug      bool   = true
)

// new claude-webapi client
func getClient() *claude.Client {
	cli := claude.NewClient(
		claude.WithBaseUri(baseUri),
		claude.WithSessionKey(sessionKey),
		claude.WithOrgid(orgid),
		claude.WithUserAgent(userAgent),
		claude.WithDebug(debug),
	)

	return cli
}

func ProxyClaude2(c *gin.Context) {

	cli := getClient()
	var requestBody ClauRequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text := requestBody.Text
	attachments := requestBody.Attachments

	println(attachments)
	conversationId := "781e5ed7-a2e9-40d1-8074-3d0fa6c362f9"

	params := map[string]interface{}{
		"attachments": attachments,
		"completion": map[string]interface{}{
			"incremental": true,
			"model":       cli.GetModel(), // default model is "claude-2"
			"prompt":      "",
			"timezone":    "Asia/Shanghai", // your custom timezone
		},
		"organization_uuid": cli.GetOrgid(),
		"conversation_uuid": conversationId,
		"text":              text,
	}
	res, err := cli.GetChatStream(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
	}

	c.Status(200)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no")

	for v := range res.Stream {
		println(v.Get("completion").String())
		if v.Get("stop_reason").String() == "stop_sequence" {
			break
		}
		reply := v.Get("completion").String()

		// 将回复逐行输出给前端
		_, _ = c.Writer.WriteString(reply + "\n")
	}

}

type ClaudeResponses struct {
	Completion string `json:"completion"`
}

func ProxyClaude(c *gin.Context) {
	cli := getClient()
	conversationId := "e4dd001c-5391-44a3-a744-420402370bf0"

	var requestBodys ClauRequestBody

	if err := c.ShouldBindJSON(&requestBodys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text := requestBodys.Text
	attachments := requestBodys.Attachments

	client := &http.Client{}
	parms := map[string]interface{}{
		"completion": map[string]interface{}{
			"prompt":      text,
			"timezone":    "Asia/Shanghai",
			"model":       "claude-2",
			"incremental": true,
		},
		"organization_uuid": cli.GetOrgid(),
		"conversation_uuid": conversationId,
		"text":              text,
		"attachments":       attachments,
	}

	requestBody, err := json.Marshal(parms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	req, err := http.NewRequest(c.Request.Method, "https://claude-proxy.f91.dev/api/append_message", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "sessionKey=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")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	println(6600)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 设置响应头
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no")

	// 将gin.Context.Writer转换为http.ResponseWriter并获取Flusher
	w := c.Writer.(http.ResponseWriter)
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 使用Scanner逐行读取响应数据并实时推送给前端
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		data := scanner.Text()
		println(data)
		// 发送数据到前端
		data = strings.TrimSuffix(data, "\r")
		c.Render(-1, common.CustomEvent{Data: data})

		// 刷新缓冲区，将数据推送给前端
		flusher.Flush()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取响应失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

}
