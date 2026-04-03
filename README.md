# gf-adapter-casbin3

GoFrame 2.x 的 Casbin v3 适配器，支持完整的 Casbin v3 接口。

[![Go](https://img.shields.io/github/go-mod/go-version/wordgao/gf-adapter-casbin3)](https://github.com/wordgao/gf-adapter-casbin3)
[![License](https://img.shields.io/github/license/nkwd/gf-adapter-casbin3)](LICENSE)

## 特性

- ✅ **完整实现 Casbin v3 全部接口**
  - `Adapter` - 基础接口
  - `ContextAdapter` - 支持 context 超时控制
  - `UpdatableAdapter` - 支持策略更新
  - `BatchAdapter` - 支持批量操作
  - `FilteredAdapter` - 支持过滤加载
  - `TransactionalAdapter` - 支持事务

- ✅ **基于 GoFrame 2.10.0 ORM**
  - 利用 GF 的链式操作
  - 支持 GF 的配置管理
  - 支持 GF 的事务机制

- ✅ **多数据库支持**
  - MySQL
  - PostgreSQL
  - SQLite
  - Oracle
  - SQL Server

- ✅ **自动建表**
  - 首次使用自动创建 casbin_rule 表
  - 支持自定义表名和字段名

- ✅ **高性能**
  - 批量操作优化
  - 事务支持
  - 索引优化

## 安装

```bash
go get github.com/wordgao/gf-adapter-casbin3
```

## 快速开始

### 基本使用

```go
package main

import (
    "github.com/casbin/casbin/v3"
    adapter "github.com/wordgao/gf-adapter-casbin3"
    "github.com/gogf/gf/v2/frame/g"
)

func main() {
    // 通过 GoFrame 配置创建适配器
    // config.yaml:
    // database:
    //   default:
    //     link: "mysql:root:password@tcp(127.0.0.1:3306)/casbin"
    
    db := g.DB()
    
    a, err := adapter.NewAdapterByDB(db)
    if err != nil {
        panic(err)
    }
    
    // 创建 Casbin 执行器
    e, err := casbin.NewEnforcer("model.conf", a)
    if err != nil {
        panic(err)
    }
    
    // 加载策略
    e.LoadPolicy()
    
    // 检查权限
    ok, _ := e.Enforce("alice", "data1", "read")
    if ok {
        // 允许访问
    }
}
```

### 自定义配置

```go
a, err := adapter.NewAdapter(adapter.Options{
    GDB: db,
    TableName: "my_casbin_rule",
    FieldName: &adapter.FieldName{
        PType: "p_type",
        V0: "v0",
        V1: "v1",
        V2: "v2",
        V3: "v3",
        V4: "v4",
        V5: "v5",
    },
    AutoCreate: true,
})
```

### 使用 Context

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 带超时的加载策略
err := a.LoadPolicyCtx(ctx, model)
```

### 批量操作

```go
rules := [][]string{
    {"alice", "data1", "read"},
    {"bob", "data2", "write"},
    {"charlie", "data3", "delete"},
}

// 批量添加
err := a.AddPolicies("p", "p", rules)

// 批量删除
err := a.RemovePolicies("p", "p", rules)
```

### 事务支持

```go
// 开始事务
ta, err := a.BeginTransaction(ctx)
if err != nil {
    panic(err)
}

// 在事务中操作
err = ta.AddPolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
if err != nil {
    ta.Rollback()
    panic(err)
}

// 提交事务
err = ta.Commit()
if err != nil {
    panic(err)
}
```

### 过滤加载

```go
filter := &adapter.Filter{
    P: []string{"", "domain1"},  // p 策略过滤条件
    G: []string{"", "", "domain1"}, // g 策略过滤条件
}

err := a.LoadFilteredPolicy(model, filter)
```

## 数据库配置

### MySQL

```yaml
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/casbin"
```

### PostgreSQL

```yaml
database:
  default:
    link: "postgresql:root:password@127.0.0.1:5432/casbin"
```

### SQLite

```yaml
database:
  default:
    link: "sqlite:./casbin.db"
```

## 表结构

默认表名：`casbin_rule`

```sql
CREATE TABLE casbin_rule (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL DEFAULT '',
    v0 VARCHAR(100) NOT NULL DEFAULT '',
    v1 VARCHAR(100) NOT NULL DEFAULT '',
    v2 VARCHAR(100) NOT NULL DEFAULT '',
    v3 VARCHAR(100) NOT NULL DEFAULT '',
    v4 VARCHAR(100) NOT NULL DEFAULT '',
    v5 VARCHAR(100) NOT NULL DEFAULT '',
    UNIQUE KEY unique_index (ptype, v0, v1, v2, v3, v4, v5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API 文档

### Adapter

```go
// 创建适配器
func NewAdapter(opts Options) (*Adapter, error)
func NewAdapterByDB(db gdb.DB, opts ...Option) (*Adapter, error)
func NewAdapterByGroup(group string, opts ...Option) (*Adapter, error)

// 基础操作
func (a *Adapter) LoadPolicy(model model.Model) error
func (a *Adapter) SavePolicy(model model.Model) error
func (a *Adapter) AddPolicy(sec, ptype string, rule []string) error
func (a *Adapter) RemovePolicy(sec, ptype string, rule []string) error
func (a *Adapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error

// Context 操作
func (a *Adapter) LoadPolicyCtx(ctx context.Context, model model.Model) error
func (a *Adapter) SavePolicyCtx(ctx context.Context, model model.Model) error
func (a *Adapter) AddPolicyCtx(ctx context.Context, sec, ptype string, rule []string) error
func (a *Adapter) RemovePolicyCtx(ctx context.Context, sec, ptype string, rule []string) error
func (a *Adapter) RemoveFilteredPolicyCtx(ctx context.Context, sec, ptype string, fieldIndex int, fieldValues ...string) error

// 更新操作
func (a *Adapter) UpdatePolicy(sec, ptype string, oldRule, newRule []string) error
func (a *Adapter) UpdatePolicies(sec, ptype string, oldRules, newRules [][]string) error
func (a *Adapter) UpdateFilteredPolicies(sec, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error)

// 批量操作
func (a *Adapter) AddPolicies(sec, ptype string, rules [][]string) error
func (a *Adapter) RemovePolicies(sec, ptype string, rules [][]string) error

// 事务
func (a *Adapter) BeginTransaction(ctx context.Context) (*TransactionalAdapter, error)

// 过滤
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行短测试（跳过数据库测试）
go test -short ./...

# 运行特定测试
go test -v -run TestAdapterBasic
```

## 依赖

- [Casbin v3](https://github.com/casbin/casbin) - 权限管理库
- [GoFrame v2.10.0](https://github.com/gogf/gf) - Go 开发框架

## 许可证

MIT License
