package websocket

import (
	"fmt"
	"log"
	"os"
)

type Router struct {
	// 路由器名字
	Name 			string

	// 缓存大小
	CacheSize		int
	// 缓存
	Cache 			chan *Pack

	// 进程数量
	ProcNum			int
	// 回调
	Proc 			func(pack *Pack)
	// 进程描述
	ProcDesc		string
}

func (r *Router) Run() {
	defer func() {
		fmt.Printf("路由器：%s启动成功\n",r.Name)
	}()

	if r.ProcNum < 1 {
		log.Println("路由器的进程数量不能少于1")
		os.Exit(0)
	}
	if r.CacheSize < 1 {
		log.Println("路由器的Channel缓存大小，不能少于1")
		os.Exit(0)
	}

	for i:=0;i<r.ProcNum;i++ {
		go func(r *Router) {
			for {
				select {
				case data := <- r.Cache:
					r.Proc(data)
				}
			}
		}(r)
	}
}

func (r *Router) Push(pack *Pack) {
	r.Cache <- pack
}