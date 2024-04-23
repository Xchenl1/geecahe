package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	// 记录自己地址，包括主机名/IP 和端口
	self string
	// 作为节点间通讯地址的前缀，默认是 /_geecache/
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPoll serving unexpected path:" + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupname := parts[0]
	key := parts[1]
	// 我们约定访问路径格式为 /<basepath>/<groupname>/<key>，通过 groupname 得到 group 实例
	group := GetGroup(groupname)
	if group == nil {
		http.Error(w, "no such group"+groupname, http.StatusNotFound)
		return
	}
	// 再使用 group.Get(key) 获取缓存数据
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 最终使用 w.Write() 将缓存值作为 httpResponse 的 body 返回
	w.Header().Set("Content_type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
