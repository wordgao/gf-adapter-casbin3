package main

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	adapter "github.com/wordgao/gf-adapter-casbin3"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	// 1. SQLite 配置
	// config.yaml:
	// database:
	//   default:
	//     link: "sqlite:./casbin.db"
	//     debug: true

	db := g.DB()

	// 2. 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("casbin_rule"),
		adapter.WithAutoCreate(true),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create adapter: %v", err))
	}

	// 3. 创建 Casbin 执行器
	e, err := casbin.NewEnforcer("model.conf", a)
	if err != nil {
		panic(fmt.Sprintf("Failed to create enforcer: %v", err))
	}

	// 4. 加载策略
	err = e.LoadPolicy()
	if err != nil {
		panic(fmt.Sprintf("Failed to load policy: %v", err))
	}

	// 5. 添加策略
	ok, err := e.AddPolicy("alice", "data1", "read")
	if err != nil {
		panic(fmt.Sprintf("Failed to add policy: %v", err))
	}
	fmt.Printf("Add policy result: %v\n", ok)

	// 6. 检查权限
	ok, _ = e.Enforce("alice", "data1", "read")
	fmt.Printf("alice can read data1: %v\n", ok)
}
