package adapter

// Engine 引擎
type Engine interface {

	// Core 核心接口
	Core() Core

	// Manage 管理接口
	Manage() Manage
}
