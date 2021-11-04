package manage

import (
	"github.com/gin-gonic/gin"
)

// MessageController 消息控制器
type MessageController manage

// add 添加设备
func (ctr *MessageController) add(c *gin.Context) {
	// 获取 body
	//body, err := c.GetRawData()
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, failWithError(model.ErrInvalidFormat))
	//	return
	//}
	//
	//topic := c.Param("topic")
	//productId := c.Query("productId")
	//deviceId := c.Query("deviceId")
	//
	//var (
	//	product *model.Product
	//	err     error
	//)
	//if product, err = repository.Database().Product().Get(productId); err != nil {
	//	c.JSON(http.StatusBadRequest, failWithError(err))
	//	return
	//}
	//
	//if err := repository.Database().Device().Add(&model.Device{
	//	ProductId: product.ProductId,
	//	DeviceId:  *req.DeviceId,
	//
	//	ProductType:  product.ProductType,
	//	DeviceName:   *req.DeviceName,
	//	DeviceSecret: ctr.generateSecret(),
	//
	//	State: model.InactiveState,
	//}); err != nil {
	//	c.JSON(http.StatusBadRequest, failWithError(err))
	//	return
	//}
	//
	//c.JSON(http.StatusOK, success(nil))
}
