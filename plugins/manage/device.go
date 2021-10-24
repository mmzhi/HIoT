package manage

import (
	"fmt"
	"github.com/fhmq/hmq/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"math/rand"
	"net/http"
	"strconv"
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
		product *model.Product
		err     error
	)
	if product, err = ctr.database.Product().Get(productId); err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	if err := ctr.database.Device().Add(&model.Device{
		ProductId: product.ProductId,
		DeviceId:  *req.DeviceId,

		ProductType:  product.ProductType,
		DeviceName:   *req.DeviceName,
		DeviceSecret: ctr.generateSecret(),

		State: model.InactiveState,
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

	err := ctr.database.Device().Update(&model.Device{
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

	ProductType  model.ProductType `json:"productType"`
	DeviceName   string            `json:"deviceName"`
	DeviceSecret string            `json:"deviceSecret"`

	FirmwareVersion *string           `json:"firmwareVersion"`
	IpAddress       *string           `json:"ipAddress"`
	State           model.DeviceState `json:"state"`

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

// DeviceListItemResponse 获取设备信息应答
type DeviceListItemResponse struct {
	ProductId string `json:"productId"`
	DeviceId  string `json:"deviceId"`

	ProductType  model.ProductType `json:"productType"`
	DeviceName   string            `json:"deviceName"`
	DeviceSecret string            `json:"deviceSecret"`

	FirmwareVersion *string           `json:"firmwareVersion"`
	IpAddress       *string           `json:"ipAddress"`
	State           model.DeviceState `json:"state"`

	OnlineTime  *Datetime `json:"onlineTime"`
	OfflineTime *Datetime `json:"offlineTime"`

	CreatedAt Datetime `json:"createdAt"`
	UpdatedAt Datetime `json:"updatedAt"`
}

type DeviceListResponse struct {
	Page Page                     `json:"page"`
	List []DeviceListItemResponse `json:"list"`
}

// list 设备列表
func (ctr *DeviceController) list(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	pageCurrent, _ := strconv.Atoi(c.Param("pageCurrent"))

	var productId = c.Param("productId")
	devices, page, err := ctr.database.Device().List(model.Page{
		Current: pageCurrent,
		Size:    pageSize,
	}, &model.Device{
		ProductId: productId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var devicesResp []DeviceListItemResponse

	for _, v := range devices {
		devicesResp = append(devicesResp, DeviceListItemResponse{
			ProductId: v.ProductId,
			DeviceId:  v.DeviceId,

			ProductType:  v.ProductType,
			DeviceName:   v.DeviceName,
			DeviceSecret: v.DeviceSecret,

			FirmwareVersion: v.FirmwareVersion,
			IpAddress:       v.IpAddress,

			State: v.State,

			OnlineTime:  PDatetime(v.OnlineTime),
			OfflineTime: PDatetime(v.OfflineTime),

			CreatedAt: Datetime{v.CreatedAt},
			UpdatedAt: Datetime{v.UpdatedAt},
		})
	}

	c.JSON(http.StatusOK, success(DeviceListResponse{
		List: devicesResp,
		Page: Page{
			Current: page.Current,
			Pages:   page.Pages,
			Size:    page.Size,
			Total:   page.Total,
		},
	}))
}

// delete 删除设备
func (ctr *DeviceController) delete(c *gin.Context) {
	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	err := ctr.database.Device().Delete(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}
