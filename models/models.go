package models

import (
	"baby/settings"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

//数据表结构

type Types struct {
	gorm.Model
	Firsts  string `json:"firsts" gorm:"type:varchar(255)"`
	Seconds string `json:"seconds" gorm:"type:varchar(255)"`
}

type Commodities struct {
	gorm.Model
	Name     string    `json:"name" gorm:"type:varchar(255)"`
	Sizes    string    `json:"sizes" gorm:"type:varchar(255)"`
	Types    string    `json:"types" gorm:"type:varchar(255)"`
	Price    float64   `json:"price" `
	Discount float64   `json:"discount" `
	Stock    int64     `json:"stock" `
	Sold     int64     `json:"sold" `
	Likes    int64     `json:"likes" `
	Created  time.Time `json:"created" `
	Img      string    `json:"img" gorm:"type:varchar(255)"`
	Details  string    `json:"details" gorm:"type:varchar(255)"`
}

type Users struct {
	gorm.Model
	Username  string    `json:"username" gorm:"type:varchar(255);unique"`
	Password  string    `json:"password" gorm:"type:varchar(255)"`
	IsStaff   int64     `json:"isStaff" gorm:"default:0"`
	LastLogin time.Time `json:"lastLogin"`
}

type Carts struct {
	gorm.Model
	Quantity    int64       `json:"quantity"`
	CommodityId int64       `json:"commodityId"`
	Commodities Commodities `gorm:"foreignKey:CommodityId"`
	UserId      int64       `json:"userId"`
	Users       Users       `json:"-" gorm:"foreignKey:UserId"`
}

type Orders struct {
	gorm.Model
	Price   string `json:"price" gorm:"type:varchar(255)"`
	PayInFo string `json:"payInFo" gorm:"type:varchar(255)"`
	UserId  int64

	Users Users `json:"-" gorm:"foreignKey:UserId"`
	State int64 `json:"state" `
}

type Records struct {
	gorm.Model
	CommodityId int64       `json:"commodityId"`
	Commodities Commodities ` gorm:"foreignKey:CommodityId"`
	UserId      int64       `json:"userId"`
	Users       Users       `json:"-" gorm:"foreignKey:UserId"`
}

type Jwts struct {
	gorm.Model
	Token  string    `json:"token" gorm:"type:varchar(1000)"`
	Expire time.Time `json:"expire"`
}

//func (u *Users) BeforeSave(db *gorm.DB) error {
//	m := md5.New()
//	m.Write([]byte(u.Password))
//	u.Password = hex.EncodeToString(m.Sum(nil))
//	return nil
//}

var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	settings.MySQLSetting.User,
	settings.MySQLSetting.Password,
	settings.MySQLSetting.Host,
	settings.MySQLSetting.Name)

var DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
	//禁止创建数据表的外键约束
	DisableForeignKeyConstraintWhenMigrating: true,
})

//数据库初始化

func Setup() error {
	// 检查 DB 是否初始化（原代码中 err 未定义，这里补充校验）
	if DB == nil {
		return errors.New("数据库连接未初始化（DB 为 nil）")
	}

	// 定义需要迁移的模型列表
	models := []interface{}{
		&Types{},
		&Commodities{},
		&Users{},
		&Carts{},
		&Orders{},
		&Records{},
		&Jwts{}, // 存储 JWT 黑名单的模型
	}

	// 循环迁移并处理错误
	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型 %T 失败: %w", model, err)
		}
	}

	// 获取底层 SQL 连接池，处理可能的错误
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 配置连接池参数（根据实际业务调整）
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	// 启动 JWT 过期 Token 定时清理任务（每天凌晨 2 点执行）
	startJwtCleanupTask()

	fmt.Println("数据库模型迁移和连接池配置成功")
	return nil
}

// startJwtCleanupTask 启动定时清理任务（每天凌晨 3 点执行）
func startJwtCleanupTask() {
	// 程序启动时先清理一次历史过期数据
	cleanupExpiredJwts()

	// 计算首次执行时间（今天凌晨 3 点，若已过则次日执行）
	now := time.Now()
	firstRun := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if now.After(firstRun) {
		firstRun = firstRun.Add(24 * time.Hour)
	}
	delay := firstRun.Sub(now)

	// 启动定时器：首次延迟 delay 后，每 24 小时执行一次
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		// 等待首次执行
		time.Sleep(delay)
		cleanupExpiredJwts()

		// 定时循环执行
		for range ticker.C {
			cleanupExpiredJwts()
		}
	}()
}

// cleanupExpiredJwts 清理过期的 JWT 记录（Expire < 当前时间）
func cleanupExpiredJwts() {
	now := time.Now()
	// 执行删除：删除 Expire 早于当前时间的记录
	result := DB.Where("expire < ?", now).Unscoped().Delete(&Jwts{})

	// 处理结果
	if result.Error != nil {
		fmt.Printf("[%s] 清理过期 JWT 失败: %v\n", now.Format("2006-01-02 15:04:05"), result.Error)
		return
	}
	fmt.Printf("[%s] 清理过期 JWT 成功，共删除 %d 条记录\n", now.Format("2006-01-02 15:04:05"), result.RowsAffected)
}
