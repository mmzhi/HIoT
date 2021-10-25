package broker

// Kick 剔除客户端
func (b *Broker) Kick(productId string, deviceId string) {
	clientId := productId + ":" + deviceId
	cli, ok := b.clients.Load(clientId)
	if ok {
		conn, succss := cli.(*client)
		if succss {
			conn.Close()
		}
	}
}
