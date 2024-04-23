package geecache

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// Get 从对应group中查找缓存值，PeerGetter 就对应于上述流程中的 HTTP 客户端
	Get(group string, key string) ([]byte, error)
}
