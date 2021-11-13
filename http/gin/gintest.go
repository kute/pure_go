package main

import (
	"fmt"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)
import "github.com/gin-gonic/gin"

type GinConfig struct {
	/**
	是否禁用控制台输出颜色
	*/
	ifDisableConsoleColor bool

	// 日志记录文件
	ifLogFile bool
}

func main() {

	fmt.Println("github.com/gin-gonic/gin begin")
	GinConfig{
		ifDisableConsoleColor: false,
		ifLogFile:             false,
	}.globalSettings()

	// default()默认会自带logger和recovery中间件
	router := gin.Default()

	routerSettings(router)

	customeLogFormat(router)

	quickStart(router)

	//createRouterTest()
}

func routerSettings(router *gin.Engine) {
	// 自定义recovery组件，将任何panic错误转换为500
	router.Use(gin.CustomRecovery(func(context *gin.Context, err interface{}) {
		if recovery, ok := err.(string); ok {
			context.String(http.StatusInternalServerError, fmt.Sprintf("error:%s", recovery))
		}
		context.AbortWithStatus(http.StatusInternalServerError)
	}))
}

func customeLogFormat(router *gin.Engine) {

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// your custom format
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		SkipPaths: []string{
			"/ping2",
		},
	}))
}

func (ginConfig GinConfig) globalSettings() {
	// 禁用控制台颜色
	if ginConfig.ifDisableConsoleColor {
		gin.DisableConsoleColor()
		//gin.ForceConsoleColor()
	}
	// 日志写文件
	if ginConfig.ifLogFile {
		if f, err := os.Create("gin.log"); err == nil {
			gin.DefaultWriter = io.MultiWriter(f)
		}
	}
}

func quickStart(router *gin.Engine) {

	router.GET("/ping", func(c *gin.Context) {
		// json响应
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// action参数可为空
	router.GET("/user/:name/*action", func(context *gin.Context) {
		var name = context.Param("name")
		println(context.FullPath())
		context.String(http.StatusOK, "response name=%s", name)
	})

	/**
	POST  /post_test?id=abcd&name=kute&choice[a]=1&choice[b]=2

	post body: message=hello&names[first]=37&names[second]=38
	*/
	router.POST("/post_test", func(context *gin.Context) {
		id := context.Query("id")
		name := context.DefaultQuery("name", "not found")
		choices := context.QueryMap("choices")
		message := context.PostForm("message")
		names := context.PostFormMap("names")
		context.JSON(http.StatusOK, gin.H{
			"id":      id,
			"name":    name,
			"choices": choices,
			"message": message,
			"names":   names,
		})
	})

	// 文件上传
	router.MaxMultipartMemory = 8 << 20 // 8M, 默认20M
	router.POST("/upload", func(context *gin.Context) {
		// 单个文件
		// curl -X POST http://localhost:8080/upload -F "file=@/Users/appleboy/test.zip" -H "Content-Type: multipart/form-data"
		if file, err := context.FormFile("file"); err == nil {
			println("upload filename=%s", file.Filename)
			println(file.Size)
		}
		// 多个文件
		// curl -X POST http://localhost:8080/upload -F "upload[]=@/Users/appleboy/test1.zip" -F "upload[]=@/Users/appleboy/test2.zip" -H "Content-Type: multipart/form-data"
		form, _ := context.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
		}
	})

	// 划分组，统一前缀，初始化时可以默认多个handler
	v1 := router.Group("/api/v1/hello", func(context *gin.Context) {
		// 默认这里的handler都会被执行到
		context.String(http.StatusOK, "default handler 1\n")
	}, func(context *gin.Context) {
		// 默认这里的handler都会被执行到
		context.String(http.StatusOK, "default handler 2\n")
	})
	{
		// response如下：
		//default handler 1
		//default handler 2
		//S1
		v1.GET("/s1", func(context *gin.Context) {
			context.String(http.StatusOK, "S1\n")
		})

		v1.GET("/s2", func(context *gin.Context) {
			// xml响应
			context.XML(http.StatusOK, gin.H{
				"code":    1,
				"message": "ok",
			})
		})

		// 内嵌组
		nestedGroup := v1.Group("/nest")
		{
			nestedGroup.GET("/content", func(context *gin.Context) {
				// ymal响应
				context.YAML(http.StatusOK, gin.H{
					"code":    1,
					"message": "ok",
				})
			})

			nestedGroup.GET("/protoBuf", func(context *gin.Context) {
				reps := []int64{int64(1), int64(2)}
				label := "test"
				// The specific definition of protobuf is written in the testdata/protoexample file.
				data := &protoexample.Test{
					Label: &label,
					Reps:  reps,
				}
				context.ProtoBuf(http.StatusOK, data)
			})
		}
	}

	// listen 8080
	_ = router.Run(":8080")
}

func createRouterTest() {

	router := gin.Default()

	fmt.Println(router.AppEngine)
}
