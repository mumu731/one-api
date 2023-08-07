package controller

// 数据集
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"one-api/model"
	"strconv"
	"time"
)

// GetAllDatasets 获取所有数据集
func GetAllDatasets(c *gin.Context) {
	id := c.GetInt("id")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	datasets, count, err := model.GetAllDatasets(id, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "获取成功",
		"data":       datasets,
		"totalPages": count,
	})
}

// FilesUploadDatasets 上传文件--中转dify
func FilesUploadDatasets(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer src.Close()

	// 创建一个新的请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 将文件内容复制到请求体中
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = io.Copy(part, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 关闭请求体
	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 创建一个新的请求
	url := "https://cloud.dify.ai/console/api/files/upload" // 替换为目标API的URL
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Cookie", "_ga=GA1.1.1015581394.1688001944; remember_token=de9d5168-fe2c-4a89-a02b-13cbdbab9ce8|3eae104cc60d030d96eb8926e37ec8d91a54ab042f14b8808e8fe7a4db7727788dfd0850f11f999b00fc7fbedf061eba3df2c5bfef08c66f1dac62bcd359a85d; _ga_VDCLZ7W3S5=GS1.1.1688779031.2.0.1688779036.55.0.0; session=553c1f61-ef07-4d2a-ae76-1fc927403d7e.FxUBFH1VhAW7-m78j4kvBrFRvqI; __cuid=917768b431ec40ea8b64d15921748ed7; amp_fef1e8=ff35c648-58cd-439e-8936-7860aeb5d5c2R...1h5mef480.1h5megog9.7.1.8; locale=zh-Hans; _ga_DM9497FN4V=GS1.1.1690954203.22.0.1690954207.56.0.0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	// 发送请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 解析响应JSON
	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 返回响应
	if resp.StatusCode == 201 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": data})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": resp.Status})
	}
}

// EstimateDatasets 上传文件--中转dify
func EstimateDatasets(c *gin.Context) {
	url := fmt.Sprintf("https://cloud.dify.ai/console/api/datasets/indexing-estimate")
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
	req.Header.Set("Cookie", "_ga=GA1.1.1015581394.1688001944; remember_token=de9d5168-fe2c-4a89-a02b-13cbdbab9ce8|3eae104cc60d030d96eb8926e37ec8d91a54ab042f14b8808e8fe7a4db7727788dfd0850f11f999b00fc7fbedf061eba3df2c5bfef08c66f1dac62bcd359a85d; _ga_VDCLZ7W3S5=GS1.1.1688779031.2.0.1688779036.55.0.0; session=553c1f61-ef07-4d2a-ae76-1fc927403d7e.FxUBFH1VhAW7-m78j4kvBrFRvqI; __cuid=917768b431ec40ea8b64d15921748ed7; amp_fef1e8=ff35c648-58cd-439e-8936-7860aeb5d5c2R...1h5mef480.1h5megog9.7.1.8; locale=zh-Hans; _ga_DM9497FN4V=GS1.1.1690954203.22.0.1690954207.56.0.0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

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

//------------------------ 向量数据库操作  ---------------------------------//

const DBURL = "https://fbf0f858-575c-4b19-9946-0a19c379f91f.us-east-1-0.aws.cloud.qdrant.io:6333"
const APIKEY = "rataloXcBhJ6YkRRho-uXAO84ro9yhIL8flwYTNs4LPhDdipq-ID0A"

func CreateCollection(c *gin.Context) {
	name := c.Query("name")
	creatBy := c.Query("creatBy")
	description := c.Query("description")
	wordCount := c.Query("wordCount")

	date := time.Now().Format("2006-01-02 15:04:05")

	url := fmt.Sprintf("%s/collections/%s", DBURL, creatBy)
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "创建失败"})
		return
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "创建失败"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", APIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "创建失败"})
		return
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {

		num, err := strconv.Atoi(wordCount)
		if err != nil {
			fmt.Println("转换失败:", err)
			return
		}
		id := c.GetInt("id")
		connectMDatasets := model.MDatasets{
			Name:        name,
			CreatTime:   date,
			CreatBy:     creatBy,
			Description: description,
			WordCount:   num,
			Remarks:     "",
			UserId:      id,
		}
		err = connectMDatasets.DatasetsInsert()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功"})

	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "创建失败"})
	}
}

func InsertPoints(c *gin.Context) {
	name := c.Query("name")
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", DBURL, name)
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "插入失败100"})
		return
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "插入失败101"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", APIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "插入失败102"})
		return
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "插入成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "插入失败103"})
	}
}

func SearchPoints(c *gin.Context) {
	name := c.Query("name")
	url := fmt.Sprintf("%s/collections/%s/points/search", DBURL, name)
	reqBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "检索失败"})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "检索失败"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", APIKEY)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "检索失败"})
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		// 解析结果
		var result interface{}
		err := json.Unmarshal(body, &result)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "检索失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "检索成功", "data": result})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "检索失败"})
	}
}

func DeleteCollection(c *gin.Context) {
	id := c.Query("id")
	err := model.DeleteDatasetsById(id)
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
	})
	return

}
