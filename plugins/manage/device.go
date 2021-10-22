package manage

import (
	"fmt"
	"github.com/fhmq/hmq/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"math/rand"
	"net/http"
)

// DeviceController 设备控制器
type DeviceController struct {
	*Engine
}

// NewDeviceController 新建DeviceController
func NewDeviceController(e *Engine) IManage {
	return &DeviceController{e}
}

// Run 运行产品Controller
func (ctr *DeviceController) Run() {
	route := ctr.Group("/api/v1/device")
	{
		route.POST("/:productId", ctr.add)
		route.POST("/:productId/:deviceId", ctr.update)
		route.GET("/:productId/:deviceId", ctr.get)
		route.GET("/:productId", ctr.list)
		route.DELETE("/:productId/:deviceId", ctr.delete)
	}
}

// DeviceAddRequest 添加设备请求
type DeviceAddRequest struct {
	DeviceId   *string `json:"deviceId"`
	DeviceName *string `json:"deviceName" binding:"required"`
}

// generateSecret 生成随机密钥
func (ctr *DeviceController) generateSecret() string {
	var letters = []rune("0123456789ABCDEF")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// add 添加设备
func (ctr *DeviceController) add(c *gin.Context) {
	var req DeviceAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}
	var productId = c.Param("productId")
	var (
		product *database.Product
		err     error
	)
	if product, err = ctr.database.Product().Get(productId); err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	if err := ctr.database.Device().Add(&database.Device{
		ProductId: product.ProductId,
		DeviceId:  *req.DeviceId,

		ProductType:  product.ProductType,
		DeviceName:   *req.DeviceName,
		DeviceSecret: ctr.generateSecret(),

		State: database.InactiveState,
	}); err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}

// DeviceUpdateRequest 	修改设备信息请求
type DeviceUpdateRequest struct {
	DeviceName *string `json:"deviceName" binding:"required"`
}

// update 更新设备信息
func (ctr *DeviceController) update(c *gin.Context) {
	var req DeviceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	err := ctr.database.Device().Update(&database.Device{
		ProductId:  productId,
		DeviceId:   deviceId,
		DeviceName: *req.DeviceName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}

// DeviceGetResponse 获取设备信息应答
type DeviceGetResponse struct {
	ProductId string `json:"productId"`
	DeviceId  string `json:"deviceId"`

	ProductType  database.ProductType `json:"productType"`
	DeviceName   string               `json:"deviceName"`
	DeviceSecret string               `json:"deviceSecret"`

	FirmwareVersion *string              `json:"firmwareVersion"`
	IpAddress       *string              `json:"ipAddress"`
	State           database.DeviceState `json:"state"`

	OnlineTime  *Datetime `json:"onlineTime"`
	OfflineTime *Datetime `json:"offlineTime"`

	CreatedAt Datetime `json:"createdAt"`
	UpdatedAt Datetime `json:"updatedAt"`
}

// get 设备详情
func (ctr *DeviceController) get(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	device, err := ctr.database.Device().Get(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(DeviceGetResponse{
		ProductId: device.ProductId,
		DeviceId:  device.DeviceId,

		ProductType:  device.ProductType,
		DeviceName:   device.DeviceName,
		DeviceSecret: device.DeviceSecret,

		FirmwareVersion: device.FirmwareVersion,
		IpAddress:       device.IpAddress,

		State: device.State,

		OnlineTime:  PDatetime(device.OnlineTime),
		OfflineTime: PDatetime(device.OfflineTime),

		CreatedAt: Datetime{device.CreatedAt},
		UpdatedAt: Datetime{device.UpdatedAt},
	}))
}

func (ctr *DeviceController) list(c *gin.Context) {

}

func (ctr *DeviceController) delete(c *gin.Context) {

}
