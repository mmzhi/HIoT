package manage

import (
	"fmt"
	"github.com/fhmq/hmq/database"
	"github.com/fhmq/hmq/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mcuadros/go-defaults"
	"github.com/segmentio/ksuid"
	"net/http"
	"strconv"
)

// ProductController 产品控制器
type ProductController struct {
	*Engine
}

// ConfigProductController 新建ProductController
func (e *Engine) ConfigProductController(route *gin.RouterGroup) *Engine {
	ctr := ProductController{e}
	route = route.Group("/product")
	{
		route.POST("/", ctr.add)
		route.POST("/:productId", ctr.update)
		route.GET("/:productId", ctr.get)
		route.GET("/", ctr.list)
		route.DELETE("/:productId", ctr.delete)
	}
	e.productController = ctr
	return e
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
	ProductId   *string            `json:"productId"`
	ProductType *model.ProductType `json:"productType" binding:"min=1,max=3"`
	ProductName *string            `json:"productName" binding:"required"`
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

	var product = model.Product{
		ProductType: *req.ProductType,
		ProductName: *req.ProductName,
	}

	if req.ProductId == nil {
		product.ProductId = ksuid.New().String()
	} else {
		product.ProductId = *req.ProductId
	}

	err := database.Database().Product().Add(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
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

	err := database.Database().Product().Update(&model.Product{
		ProductId:   productId,
		ProductName: *req.ProductName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}

// ProductListItemResponse 产品列表子项应答
type ProductListItemResponse struct {
	ProductId   *string            `json:"productId"`
	ProductType *model.ProductType `json:"productType" binding:"min=1,max=3"`
	ProductName *string            `json:"productName" binding:"required"`
}

// ProductListResponse 产品列表应答
type ProductListResponse struct {
	List []ProductListItemResponse `json:"list"`
	Page Page                      `json:"page"`
}

// list 获取产品列表
func (ctr *ProductController) list(c *gin.Context) {

	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	pageCurrent, _ := strconv.Atoi(c.Param("pageCurrent"))

	products, page, err := database.Database().Product().List(model.Page{
		Current: pageCurrent,
		Size:    pageSize,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var productsResp []ProductListItemResponse

	for _, v := range products {
		productsResp = append(productsResp, ProductListItemResponse{
			ProductId:   &v.ProductId,
			ProductType: &v.ProductType,
			ProductName: &v.ProductName,
		})
	}

	c.JSON(http.StatusOK, success(ProductListResponse{
		List: productsResp,
		Page: Page{
			Current: page.Current,
			Pages:   page.Pages,
			Size:    page.Size,
			Total:   page.Total,
		},
	}))
}

// ProductGetResponse 获取指定产品信息请求
type ProductGetResponse struct {
	ProductId   *string            `json:"productId"`
	ProductType *model.ProductType `json:"productType"`
	ProductName *string            `json:"productName"`
}

// get 获取指定产品信息
func (ctr *ProductController) get(c *gin.Context) {
	var productId = c.Param("productId")

	product, err := database.Database().Product().Get(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(ProductGetResponse{
		ProductId:   &product.ProductId,
		ProductType: &product.ProductType,
		ProductName: &product.ProductName,
	}))
}

// delete 删除产品
func (ctr *ProductController) delete(c *gin.Context) {
	var productId = c.Param("productId")

	err := database.Database().Product().Delete(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}
