package chat

type ClientServer struct{
	// 在线客户数量
	OnlineCounter int

	// RabbitMQ直接交互器
}

func (cs *ClientServer) Init() {

}

// 创建服务
func (cs *ClientServer) Create() {
	// 创建RabbitMQ消息队列
}

// 加入服务
func (cs *ClientServer) Join() {
}

// 转移客户
func (cs *ClientServer) Transfer() {
}