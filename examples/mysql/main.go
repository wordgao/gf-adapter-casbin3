package main

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v3"
	adapter "github.com/wordgao/gf-adapter-casbin3"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	// 1. 通过 GoFrame 配置文件创建适配器
	// config.yaml:
	// database:
	//   default:
	//     link: "mysql:root:password@tcp(127.0.0.1:3306)/casbin"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 添加单条策略
	ok, err := e.AddPolicy("alice", "data1", "read")
	if err != nil {
		panic(fmt.Sprintf("Failed to add policy: %v", err))
	}
	fmt.Printf("Add policy result: %v\n", ok)

	// 添加多条策略
	rules := [][]string{
		{"bob", "data2", "write"},
		{"charlie", "data3", "delete"},
	}
	ok, err = e.AddPolicies(rules)
	if err != nil {
		panic(fmt.Sprintf("Failed to add policies: %v", err))
	}
	fmt.Printf("Add policies result: %v\n", ok)

	// 6. 检查权限
	ok, _ = e.Enforce("alice", "data1", "read")
	fmt.Printf("alice can read data1: %v\n", ok)

	ok, _ = e.Enforce("bob", "data2", "write")
	fmt.Printf("bob can write data2: %v\n", ok)

	ok, _ = e.Enforce("alice", "data1", "write")
	fmt.Printf("alice can write data1: %v\n", ok)

	// 7. 更新策略
	ok, err = e.UpdatePolicy(
		[]string{"alice", "data1", "read"},
		[]string{"alice", "data1", "write"},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to update policy: %v", err))
	}
	fmt.Printf("Update policy result: %v\n", ok)

	// 8. 删除策略
	ok, err = e.RemovePolicy("bob", "data2", "write")
	if err != nil {
		panic(fmt.Sprintf("Failed to remove policy: %v", err))
	}
	fmt.Printf("Remove policy result: %v\n", ok)

	// 9. 保存策略
	err = e.SavePolicy()
	if err != nil {
		panic(fmt.Sprintf("Failed to save policy: %v", err))
	}
	fmt.Println("Policy saved successfully")

	// 10. 使用事务
	ta, err := a.BeginTransaction(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to begin transaction: %v", err))
	}

	err = ta.AddPolicyCtx(ctx, "p", "p", []string{"dave", "data4", "read"})
	if err != nil {
		ta.Rollback()
		panic(fmt.Sprintf("Failed to add policy in transaction: %v", err))
	}

	err = ta.Commit()
	if err != nil {
		panic(fmt.Sprintf("Failed to commit transaction: %v", err))
	}
	fmt.Println("Transaction committed successfully")
}
