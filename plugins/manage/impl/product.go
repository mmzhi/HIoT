package impl

import (
	"fmt"
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/plugins/manage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
	"net/http"
)

// ProductController 产品控制器
type ProductController struct {
	*Engine
}

// NewProductController 新建ProductController
func NewProductController(e *Engine) manage.IManage {
	return &ProductController{e}
}

// Run 运行产品Controller
func (ctr *ProductController) Run() {
	route := ctr.Group("/api/v1/product")
	{
		route.POST("/", ctr.add)
		route.POST("/:productId", ctr.update)
		route.GET("/:productId", ctr.get)
		route.GET("/", ctr.list)
		route.DELETE("/:productId", ctr.delete)
	}
}

// ProductAddRequest 添加产品 请求
type ProductAddRequest struct {
	ProductId   *string               `json:"productId"`
	ProductType *database.ProductType `json:"productType" binding:"min=1,max=3"`
	ProductName *string               `json:"productName" binding:"required"`
}

// ProductAddResponse 添加产品 应答
type ProductAddResponse struct {
	ProductId string `json:"productId"`
}

// add 添加产品
func (ctr *ProductController) add(c *gin.Context) {
	var req ProductAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var product = database.Product{
		ProductType: *req.ProductType,
		ProductName: *req.ProductName,
	}

	if req.ProductId == nil {
		product.ProductId = ksuid.New().String()
	} else {
		product.ProductId = *req.ProductId
	}

	err := ctr.database.Product().Add(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(ProductAddResponse{
		ProductId: product.ProductId,
	}))
}

// ProductUpdateRequest 	修改产品信息请求
type ProductUpdateRequest struct {
	ProductName *string `json:"productName" binding:"required"`
}

// update 修改产品信息
func (ctr *ProductController) update(c *gin.Context) {
	var req ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var productId = c.Param("productId")

	err := ctr.database.Product().Update(&database.Product{
		ProductId:   productId,
		ProductName: *req.ProductName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
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
