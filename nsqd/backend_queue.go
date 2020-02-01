package nsqd

// BackendQueue represents the behavior for the secondary message
// storage system
//消息(存储)队列的接口定义
type BackendQueue interface {
	Put([]byte) error	//推入消息
	ReadChan() chan []byte // this is expected to be an *unbuffered* channel //读取消息的channel
	Close() error	//关闭队列
	Delete() error	//删除队列
	Depth() int64	//队列深度
	Empty() error	//清空队列
}
