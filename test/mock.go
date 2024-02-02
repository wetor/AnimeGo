package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func serverInit() *gin.Engine {
	return gin.Default()
}

func randPort(min, max int) int {
	return rand.Intn(max-min) + min
}

func serverStart(ctx context.Context, r *gin.Engine) string {
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
		return
	})

	addr := fmt.Sprintf("127.0.0.1:%d", randPort(10000, 20000))

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error: %v\n", err)
			}
		}()
		// 监听context的关闭信号
		<-ctx.Done()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v\n", err)
		}
	}()
	host := fmt.Sprintf("http://%s", addr)

	for {
		resp, err := http.Get(host + "/ping")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("HTTP server is ready")
			break
		}
		fmt.Println("HTTP server is not ready yet, retrying...")
		time.Sleep(250 * time.Millisecond)
	}

	return host
}

func MockMikanStart(ctx context.Context) string {
	r := serverInit()
	// 定义路由和处理函数
	r.GET("/Home/Episode/:hash", func(c *gin.Context) {
		testData, has := c.GetQuery("testdata")
		if !has {
			testData = "mikan"
		}
		hash := c.Param("hash")

		data, err := GetData(testData, hash)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "text/html", data)
		return
	})
	r.GET("/Home/bangumi/:mikan_id", func(c *gin.Context) {
		testData, has := c.GetQuery("testdata")
		if !has {
			testData = "mikan"
		}
		mikanId := c.Param("mikan_id")
		data, err := GetData(testData, mikanId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "text/html", data)
		return
	})
	return serverStart(ctx, r)
}

func MockBangumiStart(ctx context.Context) string {
	r := serverInit()
	// 定义路由和处理函数
	r.GET("/v0/subjects/:id", func(c *gin.Context) {
		testData, has := c.GetQuery("testdata")
		if !has {
			testData = "bangumi"
		}
		id := c.Param("id")
		data, err := GetData(testData, id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	})
	return serverStart(ctx, r)
}

func MockThemoviedbStart(ctx context.Context) string {
	r := serverInit()
	// 定义路由和处理函数
	r.GET("/3/discover/tv", func(c *gin.Context) {
		testData, has := c.GetQuery("testdata")
		if !has {
			testData = "themoviedb"
		}
		query := c.Query("with_text_query")
		data, err := GetData(testData, query)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	})
	r.GET("/3/tv/:id", func(c *gin.Context) {
		testData, has := c.GetQuery("testdata")
		if !has {
			testData = "themoviedb"
		}
		id := c.Param("id")
		data, err := GetData(testData, id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/json", data)
		return
	})
	return serverStart(ctx, r)
}
