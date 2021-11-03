package core

import "errors"

// parseSysTopic 解析topic，获取目标 productId 和 deviceId
func parseSysTopic(topic string) (productId string, deviceId string, err error) {
	if params := sysTopicRegexp.FindStringSubmatch(topic); len(params) != 3 {
		return "", "", errors.New("not a right topic")
	} else {
		return params[1], params[2], nil
	}
}
