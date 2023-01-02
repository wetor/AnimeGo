package api

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/wetor/AnimeGo/internal/store"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/cache"
)

// BoltList godoc
//  @Summary 获取Bolt数据库的Bucket列表或key列表
//  @Description 获取Bolt数据库的Bucket列表，或指定Bucket下的key列表
//  @Tags bolt
//  @Accept  json
//  @Produce  json
//  @Param type query webModels.BoltListRequest true "获取bolt数据库列表"
//  @Success 200 {object} webModels.Response{data=webModels.BoltListResponse}
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/bolt [get]
func BoltList(c *gin.Context) {
	var request webModels.BoltListRequest
	if !checkRequest(c, &request) {
		return
	}
	var db *cache.Bolt
	if request.DB == "bolt" {
		db = store.Cache
	} else if request.DB == "bolt_sub" {
		db = store.BangumiCache
	} else {
		c.JSON(webModels.Fail("参数错误，未找到数据库"))
		return
	}

	var list []string
	if request.Type == "bucket" {
		list = db.ListBucket()
	} else if request.Type == "key" {
		if len(request.Bucket) == 0 {
			c.JSON(webModels.Fail("参数错误，type为 " + request.Type + " 时，需要 bucket 参数"))
			return
		}
		list = db.ListKey(request.Bucket)
	} else {
		c.JSON(webModels.Fail("参数错误，不支持的type：" + request.Type + "，目前仅支持 bucket 和 key"))
		return
	}
	c.JSON(webModels.Succ("列表", webModels.BoltListResponse{
		Type:   request.Type,
		Bucket: request.Bucket,
		Data:   list,
	}))
}

// Bolt godoc
//  @Summary 获取Bolt数据库的值
//  @Description 获取Bolt数据库指定Bucket和key所储存的值
//  @Tags bolt
//  @Accept  json
//  @Produce  json
//  @Param type query webModels.BoltGetRequest true "获取bolt数据库值"
//  @Success 200 {object} webModels.Response{data=webModels.BoltGetResponse}
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/bolt/value [get]
func Bolt(c *gin.Context) {
	var request webModels.BoltGetRequest
	if !checkRequest(c, &request) {
		return
	}
	var db *cache.Bolt
	if request.DB == "bolt" {
		db = store.Cache
	} else if request.DB == "bolt_sub" {
		db = store.BangumiCache
	} else {
		c.JSON(webModels.Fail("参数错误，未找到数据库"))
		return
	}

	ttl, val, err := db.GetValue(request.Bucket, request.Key)
	if err != nil {
		c.JSON(webModels.Fail("查询失败，" + err.Error()))
		return
	}
	m := make(map[string]any)
	err = jsoniter.Unmarshal([]byte(val), &m)
	if err != nil {
		str := ""
		err = jsoniter.Unmarshal([]byte(val), &str)
		if err != nil {
			c.JSON(webModels.Fail("转换失败，" + err.Error()))
			return
		}
		c.JSON(webModels.Succ("查询结果", webModels.BoltGetResponse{
			Bucket: request.Bucket,
			Key:    request.Key,
			Value:  str,
			TTL:    ttl,
		}))
		return
	}
	c.JSON(webModels.Succ("查询结果", webModels.BoltGetResponse{
		Bucket: request.Bucket,
		Key:    request.Key,
		Value:  m,
		TTL:    ttl,
	}))
}

// BoltDelete godoc
//  @Summary 删除Bolt数据库的值
//  @Description 删除Bolt数据库指定Bucket和key所储存的值
//  @Tags bolt
//  @Accept  json
//  @Produce  json
//  @Param type query webModels.BoltDeleteRequest true "删除bolt数据库值"
//  @Success 200 {object} webModels.Response{data=webModels.BoltGetResponse}
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/bolt/value [delete]
func BoltDelete(c *gin.Context) {
	var request webModels.BoltDeleteRequest
	if !checkRequest(c, &request) {
		return
	}
	var db *cache.Bolt
	if request.DB == "bolt" {
		db = store.Cache
	} else {
		c.JSON(webModels.Fail("参数错误，只能删除 bolt 数据库中的数据"))
		return
	}

	err := db.Delete(request.Bucket, request.Key)
	if err != nil {
		c.JSON(webModels.Fail("删除失败，" + err.Error()))
		return
	}
	c.JSON(webModels.Succ("删除成功"))
}
