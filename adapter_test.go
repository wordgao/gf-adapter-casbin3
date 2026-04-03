package adapter_test

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	adapter "github.com/wordgao/gf-adapter-casbin3"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func TestNewAdapter(t *testing.T) {
	// 测试空数据库实例
	_, err := adapter.NewAdapter(adapter.Options{})
	if err != adapter.ErrDBRequired {
		t.Errorf("Expected ErrDBRequired, got %v", err)
	}
}

func TestNewAdapterByDB(t *testing.T) {
	// 测试空数据库实例
	_, err := adapter.NewAdapterByDB(nil)
	if err != adapter.ErrDBRequired {
		t.Errorf("Expected ErrDBRequired, got %v", err)
	}
}

func TestAdapterBasic(t *testing.T) {
	// 跳过需要真实数据库的测试
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// 假设已配置数据库
	db := g.DB()
	if db == nil {
		t.Skip("Database not configured")
	}

	// 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("test_casbin_rule"),
	)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx := context.Background()

	// 测试添加策略
	err = a.AddPolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
	if err != nil {
		t.Errorf("Failed to add policy: %v", err)
	}

	// 测试删除策略
	err = a.RemovePolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
	if err != nil {
		t.Errorf("Failed to remove policy: %v", err)
	}
}

func TestAdapterWithCasbin(t *testing.T) {
	// 跳过需要真实数据库的测试
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// 假设已配置数据库
	db := g.DB()
	if db == nil {
		t.Skip("Database not configured")
	}

	// 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("test_casbin_rule"),
	)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	// 创建模型
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	// 创建执行器
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		t.Fatalf("Failed to create enforcer: %v", err)
	}

	// 添加策略
	ok, err := e.AddPolicy("alice", "data1", "read")
	if err != nil {
		t.Errorf("Failed to add policy: %v", err)
	}
	if !ok {
		t.Error("AddPolicy returned false")
	}

	// 检查权限
	ok, err = e.Enforce("alice", "data1", "read")
	if err != nil {
		t.Errorf("Failed to enforce: %v", err)
	}
	if !ok {
		t.Error("Enforce should return true")
	}

	// 检查无权限
	ok, err = e.Enforce("bob", "data1", "read")
	if err != nil {
		t.Errorf("Failed to enforce: %v", err)
	}
	if ok {
		t.Error("Enforce should return false")
	}

	// 删除策略
	ok, err = e.RemovePolicy("alice", "data1", "read")
	if err != nil {
		t.Errorf("Failed to remove policy: %v", err)
	}
	if !ok {
		t.Error("RemovePolicy returned false")
	}
}

func TestBatchOperations(t *testing.T) {
	// 跳过需要真实数据库的测试
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// 假设已配置数据库
	db := g.DB()
	if db == nil {
		t.Skip("Database not configured")
	}

	// 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("test_casbin_rule"),
	)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx := context.Background()

	// 批量添加
	rules := [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"charlie", "data3", "delete"},
	}

	err = a.AddPoliciesCtx(ctx, "p", "p", rules)
	if err != nil {
		t.Errorf("Failed to add policies: %v", err)
	}

	// 批量删除
	err = a.RemovePoliciesCtx(ctx, "p", "p", rules)
	if err != nil {
		t.Errorf("Failed to remove policies: %v", err)
	}
}

func TestUpdateOperations(t *testing.T) {
	// 跳过需要真实数据库的测试
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// 假设已配置数据库
	db := g.DB()
	if db == nil {
		t.Skip("Database not configured")
	}

	// 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("test_casbin_rule"),
	)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx := context.Background()

	// 添加初始策略
	err = a.AddPolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
	if err != nil {
		t.Errorf("Failed to add policy: %v", err)
	}

	// 更新策略
	err = a.UpdatePolicyCtx(ctx, "p", "p",
		[]string{"alice", "data1", "read"},
		[]string{"alice", "data1", "write"},
	)
	if err != nil {
		t.Errorf("Failed to update policy: %v", err)
	}

	// 清理
	err = a.RemovePolicyCtx(ctx, "p", "p", []string{"alice", "data1", "write"})
	if err != nil {
		t.Errorf("Failed to remove policy: %v", err)
	}
}

func TestTransaction(t *testing.T) {
	// 跳过需要真实数据库的测试
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// 假设已配置数据库
	db := g.DB()
	if db == nil {
		t.Skip("Database not configured")
	}

	// 创建适配器
	a, err := adapter.NewAdapterByDB(db,
		adapter.WithTableName("test_casbin_rule"),
	)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx := context.Background()

	// 开始事务
	ta, err := a.BeginTransaction(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// 在事务中添加策略
	err = ta.AddPolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
	if err != nil {
		ta.Rollback()
		t.Errorf("Failed to add policy in transaction: %v", err)
	}

	// 提交事务
	err = ta.Commit()
	if err != nil {
		t.Errorf("Failed to commit transaction: %v", err)
	}

	// 清理
	err = a.RemovePolicyCtx(ctx, "p", "p", []string{"alice", "data1", "read"})
	if err != nil {
		t.Errorf("Failed to remove policy: %v", err)
	}
}

func TestFieldName(t *testing.T) {
	// 测试自定义字段名
	fieldName := &adapter.FieldName{
		PType: "p_type",
		V0:    "v0",
		V1:    "v1",
		V2:    "v2",
		V3:    "v3",
		V4:    "v4",
		V5:    "v5",
	}

	if fieldName.PType != "p_type" {
		t.Errorf("Expected p_type, got %s", fieldName.PType)
	}
}

func TestGetCreateTableSQL(t *testing.T) {
	fieldName := adapter.DefaultFieldName()
	sql := adapter.GetCreateTableSQL("casbin_rule", fieldName)

	if sql == "" {
		t.Error("SQL should not be empty")
	}

	if !contains(sql, "CREATE TABLE") {
		t.Error("SQL should contain CREATE TABLE")
	}

	if !contains(sql, "casbin_rule") {
		t.Error("SQL should contain table name")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
