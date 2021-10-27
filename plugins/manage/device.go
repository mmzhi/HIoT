package manage

import (
	"fmt"
	"github.com/fhmq/hmq/database"
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

// ConfigDeviceController 新建DeviceController
func (e *Engine) ConfigDeviceController(route *gin.RouterGroup) *Engine {
	ctr := DeviceController{e}
	route = route.Group("/device")
	{
		route.POST("/:productId", ctr.add)
		route.GET("/", ctr.list)
		route.GET("/:productId/:deviceId", ctr.get)
		route.POST("/:productId/:deviceId", ctr.update)
		route.DELETE("/:productId/:deviceId", ctr.delete)

		route.POST("/:productId/:deviceId/enable", ctr.enable)
		route.POST("/:productId/:deviceId/disable", ctr.disable)

		route.GET("/:productId/:deviceId/config", ctr.getConfig)
		route.POST("/:productId/:deviceId/config", ctr.updateConfig)

		route.GET("/:productId/:deviceId/topology", ctr.getTopology)
		route.POST("/:productId/:deviceId/topology", ctr.updateTopology)
		route.DELETE("/:productId/:deviceId/topology", ctr.removeTopology)

		route.POST("/:productId/:deviceId/reset", ctr.reset)
	}
	e.deviceController = ctr
	return e
}

// DeviceAddRequest 添加设备请求
type DeviceAddRequest struct {
	DeviceId   *string `json:"deviceId"`
	DeviceName *string `json:"deviceName" binding:"required"`
}

// generateSecret 生成随机密钥
func (ctr *DeviceController) generateSecret() string {
	var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 16)
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
	if product, err = database.Database().Product().Get(productId); err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	if err := database.Database().Device().Add(&model.Device{
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

	err := database.Database().Device().Update(&model.Device{
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

	device, err := database.Database().Device().Get(productId, deviceId)
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
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageCurrent, _ := strconv.Atoi(c.Query("pageCurrent"))

	var productId = c.Query("productId")
	devices, page, err := database.Database().Device().List(model.Page{
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

	err := database.Database().Device().Delete(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}

// enable 启用设备
func (ctr *DeviceController) enable(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	device, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	if device.State != model.DisabledState && device.State != model.InactiveDisabledState {
		c.JSON(http.StatusBadRequest, fail(0, "State is not support"))
		return
	} else if device.State == model.DisabledState { // 禁用
		if err := database.Database().Device().UpdateState(productId, deviceId, model.OfflineState); err != nil {
			c.JSON(http.StatusBadRequest, fail(0, err.Error()))
			return
		}
	} else { // 禁用且未激活
		if err := database.Database().Device().UpdateState(productId, deviceId, model.InactiveState); err != nil {
			c.JSON(http.StatusBadRequest, fail(0, err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, success(nil))
}

// disable 禁用设备
func (ctr *DeviceController) disable(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	device, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	if device.State == model.DisabledState || device.State == model.InactiveDisabledState {
		c.JSON(http.StatusBadRequest, fail(0, "State is not support"))
		return
	} else if device.State == model.InactiveState { // 未激活
		if err := database.Database().Device().UpdateState(productId, deviceId, model.InactiveDisabledState); err != nil {
			c.JSON(http.StatusBadRequest, fail(0, err.Error()))
			return
		}
	} else { // 其余状态
		if err := database.Database().Device().UpdateState(productId, deviceId, model.DisabledState); err != nil {
			c.JSON(http.StatusBadRequest, fail(0, err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, success(nil))
}

// DeviceGetConfigResponse 获取设备配置应答
type DeviceGetConfigResponse struct {
	Config *string `json:"config"`
}

// getConfig 获取配置
func (ctr *DeviceController) getConfig(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	config, err := database.Database().Device().GetConfig(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(DeviceGetConfigResponse{
		Config: config,
	}))
}

// DeviceUpdateConfigRequest 设备配置更新请求
type DeviceUpdateConfigRequest struct {
	Config *string `json:"config"`
}

// updateConfig 更新配置
func (ctr *DeviceController) updateConfig(c *gin.Context) {

	var req DeviceUpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	err := database.Database().Device().UpdateConfig(productId, deviceId, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}
	c.JSON(http.StatusOK, success(nil))
}

// DeviceGetTopologyResponse 获取子设备拓扑信息
type DeviceGetTopologyResponse struct {
	GatewayProductId *string `json:"gatewayProductId"`
	GatewayDeviceId  *string `json:"gatewayDeviceId"`
}

// getTopology 获取子设备拓扑信息
func (ctr *DeviceController) getTopology(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	device, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, success(DeviceGetTopologyResponse{
		GatewayProductId: device.GatewayProductId,
		GatewayDeviceId:  device.GatewayDeviceId,
	}))
}

// DeviceUpdateTopologyRequest 子设备拓扑信息更新请求
type DeviceUpdateTopologyRequest struct {
	GatewayProductId *string `json:"gatewayProductId"`
	GatewayDeviceId  *string `json:"gatewayDeviceId"`
}

// updateTopology 子设备拓扑信息更新
func (ctr *DeviceController) updateTopology(c *gin.Context) {

	var req DeviceUpdateTopologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	subDevice, err := database.Database().Device().Get(productId, deviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	} else if subDevice.ProductType != model.SubDeviceType {
		c.JSON(http.StatusBadRequest, fail(0, "not a sub device"))
		return
	}

	// TODO 下线子设备

	_, err = database.Database().Device().Get(*req.GatewayProductId, *req.GatewayDeviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}

	err = database.Database().Device().UpdateGateway(productId, deviceId, req.GatewayProductId, req.GatewayDeviceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}
	c.JSON(http.StatusOK, success(nil))
}

// removeTopology 子设备拓扑信息移除
func (ctr *DeviceController) removeTopology(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	err := database.Database().Device().UpdateGateway(productId, deviceId, nil, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}
	c.JSON(http.StatusOK, success(nil))
}

// reset 设备密钥重置
func (ctr *DeviceController) reset(c *gin.Context) {

	var productId = c.Param("productId")
	var deviceId = c.Param("deviceId")

	err := database.Database().Device().UpdateSecret(productId, deviceId, ctr.generateSecret())
	if err != nil {
		c.JSON(http.StatusBadRequest, fail(0, err.Error()))
		return
	}
	c.JSON(http.StatusOK, success(nil))
}
