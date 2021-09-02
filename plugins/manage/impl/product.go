package impl

import (
	"github.com/fhmq/hmq/adapter"
	"github.com/fhmq/hmq/database"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"net/http"
)

// ProductController 产品控制器
type ProductController struct {
	database database.IDatabase // 数据库功能
	handler  adapter.IHandler   // broker扩展方法
}

// ProductAddRequest 添加产品 请求
type ProductAddRequest struct {
	ProductId   string               `json:"productId"`
	ProductType database.ProductType `json:"productType" binding:"required"`
	ProductName string               `json:"productName" binding:"required"`
}

// ProductAddResponse 添加产品 应答
type ProductAddResponse struct {
	ProductId string `json:"productId"`
}

// add 添加产品
func (ctr *ProductController) add(c *gin.Context) {
	var req ProductAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, fail(0, ""))
		return
	}

	if req.ProductId == "" {
		req.ProductId = ksuid.New().String()
	}

	var product = database.Product{}
	err := ctr.database.Product().Add(&product)
	if err != nil {

		return
	}
}

// update 修改产品信息
func (ctr *ProductController) update(c *gin.Context) {

}

// list 获取产品列表
func (ctr *ProductController) list(c *gin.Context) {

}

// get 获取指定产品信息
func (ctr *ProductController) get(c *gin.Context) {

}

// delete 删除产品
func (ctr *ProductController) delete(c *gin.Context) {

}
