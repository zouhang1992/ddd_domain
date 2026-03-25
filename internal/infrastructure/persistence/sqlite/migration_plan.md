# 数据迁移方案

## 概述

此文档详细描述了从旧的 `house.db` SQLite 数据库到新的 DDD 架构的 `ddd.db` 数据库的完整数据迁移方案。

## 迁移策略

采用增量式迁移策略，分为多个阶段：

1. **备份阶段**：确保源数据安全
2. **验证阶段**：检查源数据库结构
3. **迁移阶段**：将数据导入新数据库
4. **验证阶段**：检查迁移后的数据完整性
5. **清理阶段**：处理临时文件和日志

## 源数据库结构

**旧的 `house.db` 数据库（来自 /Users/bytedance/go/src/code.byted.org/zouhang/house 系统）：**

```
-- 主要表
landlords (id, name, phone, note)
rooms (id, location_id, room_number, tags)
leases (id, room_id, landlord_id, tenant_name, tenant_phone, start_date, end_date, rent_amount, status, last_charge_at)
bills (id, lease_id, type, status, rent, water, electric, other, amount, paid_at, note)
deposits (id, lease_id, amount, status, note)

-- 位置表（可能的结构）
locations (id, short_name, detail)
```

## 目标数据库结构

**新的 `ddd.db` 数据库：**

```
-- 已在 connection.go 中定义的 schema

-- 房间表（现有的）
CREATE TABLE IF NOT EXISTS rooms (
    id TEXT PRIMARY KEY,
    location_id TEXT NOT NULL,
    room_number TEXT NOT NULL,
    tags TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

-- 位置表（现有的）
CREATE TABLE IF NOT EXISTS locations (
    id TEXT PRIMARY KEY,
    short_name TEXT NOT NULL,
    detail TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- 新增的表
CREATE TABLE IF NOT EXISTS landlords (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS leases (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    landlord_id TEXT,
    tenant_name TEXT NOT NULL,
    tenant_phone TEXT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    rent_amount INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    note TEXT,
    last_charge_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (landlord_id) REFERENCES landlords(id)
);

CREATE TABLE IF NOT EXISTS bills (
    id TEXT PRIMARY KEY,
    lease_id TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    amount INTEGER NOT NULL DEFAULT 0,
    rent_amount INTEGER NOT NULL DEFAULT 0,
    water_amount INTEGER NOT NULL DEFAULT 0,
    electric_amount INTEGER NOT NULL DEFAULT 0,
    other_amount INTEGER NOT NULL DEFAULT 0,
    paid_at DATETIME,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (lease_id) REFERENCES leases(id)
);

CREATE TABLE IF NOT EXISTS deposits (
    id TEXT PRIMARY KEY,
    lease_id TEXT NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'collected',
    refunded_at DATETIME,
    deducted_at DATETIME,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (lease_id) REFERENCES leases(id)
);
```

## 迁移脚本

### 1. 备份脚本

```bash
#!/bin/bash

# 备份脚本
SOURCE_DB="/Users/bytedance/go/src/code.byted.org/zouhang/house/house.db"
TARGET_DIR="/Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$TARGET_DIR/backup_house_$TIMESTAMP.db"

mkdir -p $TARGET_DIR

if [ -f "$SOURCE_DB" ]; then
    echo "备份源数据库 $SOURCE_DB 到 $BACKUP_FILE"
    cp "$SOURCE_DB" "$BACKUP_FILE"
    if [ $? -eq 0 ]; then
        echo "备份成功"
    else
        echo "备份失败"
        exit 1
    fi
else
    echo "源数据库 $SOURCE_DB 不存在"
    exit 1
fi
```

### 2. 迁移脚本

```go
// 数据迁移脚本（使用 Go 语言）
package main

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    sourceDB := "/Users/bytedance/go/src/code.byted.org/zouhang/house/house.db"
    targetDB := "/Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db"

    // 打开源数据库
    srcDB, err := sql.Open("sqlite3", sourceDB)
    if err != nil {
        fmt.Printf("打开源数据库失败: %v\n", err)
        os.Exit(1)
    }
    defer srcDB.Close()

    // 打开或创建目标数据库
    tgtDB, err := sql.Open("sqlite3", targetDB)
    if err != nil {
        fmt.Printf("打开目标数据库失败: %v\n", err)
        os.Exit(1)
    }
    defer tgtDB.Close()

    // 创建目标数据库结构（使用 connection.go 中的 schema）
    createSchema(tgtDB)

    // 执行迁移
    if err := migrateData(srcDB, tgtDB); err != nil {
        fmt.Printf("数据迁移失败: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("数据迁移完成")
}

func createSchema(db *sql.DB) {
    schema := `
CREATE TABLE IF NOT EXISTS sagas (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    state TEXT NOT NULL,
    current_step INTEGER NOT NULL DEFAULT 0,
    error TEXT,
    data BLOB,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
    id TEXT PRIMARY KEY,
    short_name TEXT NOT NULL,
    detail TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    id TEXT PRIMARY KEY,
    location_id TEXT NOT NULL,
    room_number TEXT NOT NULL,
    tags TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE TABLE IF NOT EXISTS landlords (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS leases (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    landlord_id TEXT,
    tenant_name TEXT NOT NULL,
    tenant_phone TEXT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    rent_amount INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    note TEXT,
    last_charge_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (landlord_id) REFERENCES landlords(id)
);

CREATE TABLE IF NOT EXISTS bills (
    id TEXT PRIMARY KEY,
    lease_id TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    amount INTEGER NOT NULL DEFAULT 0,
    rent_amount INTEGER NOT NULL DEFAULT 0,
    water_amount INTEGER NOT NULL DEFAULT 0,
    electric_amount INTEGER NOT NULL DEFAULT 0,
    other_amount INTEGER NOT NULL DEFAULT 0,
    paid_at DATETIME,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (lease_id) REFERENCES leases(id)
);

CREATE TABLE IF NOT EXISTS deposits (
    id TEXT PRIMARY KEY,
    lease_id TEXT NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'collected',
    refunded_at DATETIME,
    deducted_at DATETIME,
    note TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (lease_id) REFERENCES leases(id)
);
`

    _, err := db.Exec(schema)
    if err != nil {
        fmt.Printf("创建数据库结构失败: %v\n", err)
        os.Exit(1)
    }
}

func migrateData(srcDB, tgtDB *sql.DB) error {
    // 迁移 landlords
    if err := migrateLandlords(srcDB, tgtDB); err != nil {
        return err
    }

    // 迁移 locations（假设源数据库使用 locations 表）
    if err := migrateLocations(srcDB, tgtDB); err != nil {
        return err
    }

    // 迁移 rooms
    if err := migrateRooms(srcDB, tgtDB); err != nil {
        return err
    }

    // 迁移 leases
    if err := migrateLeases(srcDB, tgtDB); err != nil {
        return err
    }

    // 迁移 bills
    if err := migrateBills(srcDB, tgtDB); err != nil {
        return err
    }

    // 迁移 deposits
    if err := migrateDeposits(srcDB, tgtDB); err != nil {
        return err
    }

    return nil
}

func migrateLandlords(srcDB, tgtDB *sql.DB) error {
    rows, err := srcDB.Query("SELECT id, name, phone, note FROM landlords")
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id, name, phone, note string
        if err := rows.Scan(&id, &name, &phone, &note); err != nil {
            return err
        }

        now := time.Now()
        _, err := tgtDB.Exec(`
            INSERT OR REPLACE INTO landlords (id, name, phone, note, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, id, name, phone, note, now, now)
        if err != nil {
            return err
        }
    }

    return nil
}

func migrateLocations(srcDB, tgtDB *sql.DB) error {
    // 检查源数据库是否有 locations 表
    var count int
    err := srcDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='locations'").Scan(&count)
    if err != nil {
        return err
    }

    if count > 0 {
        rows, err := srcDB.Query("SELECT id, short_name, detail FROM locations")
        if err != nil {
            return err
        }
        defer rows.Close()

        for rows.Next() {
            var id, shortName, detail string
            if err := rows.Scan(&id, &shortName, &detail); err != nil {
                return err
            }

            now := time.Now()
            _, err := tgtDB.Exec(`
                INSERT OR REPLACE INTO locations (id, short_name, detail, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?)
            `, id, shortName, detail, now, now)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func migrateRooms(srcDB, tgtDB *sql.DB) error {
    rows, err := srcDB.Query("SELECT id, location_id, room_number, tags FROM rooms")
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id, locationID, roomNumber, tags string
        if err := rows.Scan(&id, &locationID, &roomNumber, &tags); err != nil {
            return err
        }

        now := time.Now()
        _, err := tgtDB.Exec(`
            INSERT OR REPLACE INTO rooms (id, location_id, room_number, tags, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `, id, locationID, roomNumber, tags, now, now)
        if err != nil {
            return err
        }
    }

    return nil
}

func migrateLeases(srcDB, tgtDB *sql.DB) error {
    rows, err := srcDB.Query(`
        SELECT id, room_id, landlord_id, tenant_name, tenant_phone,
               start_date, end_date, rent_amount, status, last_charge_at
        FROM leases
    `)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id, roomID, landlordID, tenantName, tenantPhone, status string
        var startDate, endDate, lastChargeAt time.Time
        var rentAmount int64

        if err := rows.Scan(&id, &roomID, &landlordID, &tenantName, &tenantPhone,
            &startDate, &endDate, &rentAmount, &status, &lastChargeAt); err != nil {
            return err
        }

        now := time.Now()
        _, err := tgtDB.Exec(`
            INSERT OR REPLACE INTO leases (
                id, room_id, landlord_id, tenant_name, tenant_phone, start_date, end_date,
                rent_amount, status, last_charge_at, created_at, updated_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, id, roomID, landlordID, tenantName, tenantPhone, startDate, endDate,
            rentAmount, status, lastChargeAt, now, now)
        if err != nil {
            return err
        }
    }

    return nil
}

func migrateBills(srcDB, tgtDB *sql.DB) error {
    rows, err := srcDB.Query(`
        SELECT id, lease_id, type, status, rent, water, electric, other, amount, paid_at, note
        FROM bills
    `)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id, leaseID, type, status, note string
        var rent, water, electric, other, amount int64
        var paidAt time.Time

        if err := rows.Scan(&id, &leaseID, &type, &status, &rent, &water, &electric,
            &other, &amount, &paidAt, &note); err != nil {
            return err
        }

        now := time.Now()
        _, err := tgtDB.Exec(`
            INSERT OR REPLACE INTO bills (
                id, lease_id, type, status, amount, rent_amount, water_amount, electric_amount,
                other_amount, paid_at, note, created_at, updated_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, id, leaseID, type, status, amount, rent, water, electric, other, paidAt,
            note, now, now)
        if err != nil {
            return err
        }
    }

    return nil
}

func migrateDeposits(srcDB, tgtDB *sql.DB) error {
    rows, err := srcDB.Query(`SELECT id, lease_id, amount, status, note FROM deposits`)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var id, leaseID, status, note string
        var amount int64

        if err := rows.Scan(&id, &leaseID, &amount, &status, &note); err != nil {
            return err
        }

        now := time.Now()
        _, err := tgtDB.Exec(`
            INSERT OR REPLACE INTO deposits (id, lease_id, amount, status, note, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `, id, leaseID, amount, status, note, now, now)
        if err != nil {
            return err
        }
    }

    return nil
}
```

### 3. 使用方法

```bash
# 1. 编译迁移脚本
cd /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain
CGO_ENABLED=1 go build -o data/migrate data/migrate.go

# 2. 运行备份脚本
chmod +x data/backup.sh
./data/backup.sh

# 3. 运行迁移
./data/migrate

# 4. 验证迁移结果
./data/migrate verify
```

## 验证方法

### 1. 简单验证

```bash
# 检查是否有数据迁移
ls -la /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/
sqlite3 /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db "SELECT count(*) FROM landlords"
```

### 2. 复杂验证

```bash
# 比较行数
echo -e "\n--- 数据行数比较 ---"
echo "源数据库 landlords count: $(sqlite3 /Users/bytedance/go/src/code.byted.org/zouhang/house/house.db "SELECT count(*) FROM landlords")"
echo "目标数据库 landlords count: $(sqlite3 /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db "SELECT count(*) FROM landlords")"
echo "源数据库 leases count: $(sqlite3 /Users/bytedance/go/src/code.byted.org/zouhang/house/house.db "SELECT count(*) FROM leases")"
echo "目标数据库 leases count: $(sqlite3 /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db "SELECT count(*) FROM leases")"
echo "源数据库 bills count: $(sqlite3 /Users/bytedance/go/src/code.byted.org/zouhang/house/house.db "SELECT count(*) FROM bills")"
echo "目标数据库 bills count: $(sqlite3 /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db "SELECT count(*) FROM bills")"
echo "源数据库 rooms count: $(sqlite3 /Users/bytedance/go/src/code.byted.org/zouhang/house/house.db "SELECT count(*) FROM rooms")"
echo "目标数据库 rooms count: $(sqlite3 /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/ddd.db "SELECT count(*) FROM rooms")"
```

## 回滚计划

**如果迁移失败，可以使用以下回滚步骤：**

1. **停止迁移进程**
2. **检查当前状态**
3. **恢复备份**（如果需要）

```bash
# 从备份中恢复源数据
cp /Users/bytedance/go/src/github.com/zouhang1992/ddd_domain/data/backup_house_*.db /Users/bytedance/go/src/code.byted.org/zouhang/house/house.db
```

4. **删除目标数据库**
5. **重新开始迁移过程**

## 生产环境迁移

**生产环境迁移需要额外的预防措施：**

1. **使用事务**：确保迁移的原子性
2. **分阶段部署**：
   - 在测试环境验证
   - 在生产环境进行小规模测试
   - 在低流量时间窗口进行
3. **监控**：
   - 监控系统资源使用情况
   - 确保 API 响应正常
   - 检查错误率

## 后续优化

1. **索引优化**：根据查询模式创建适当的索引
2. **数据清理**：删除不再需要的旧数据
3. **性能优化**：
   - 使用分区
   - 优化查询
   - 调整配置

## 支持和帮助

如在迁移过程中遇到问题，请：

1. 检查日志文件
2. 验证源数据库结构
3. 确保磁盘空间充足
4. 联系开发团队
