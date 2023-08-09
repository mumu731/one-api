package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"one-api/common"
)

var DB *gorm.DB

func createRootAccountIfNeed() error {
	var user User
	//if user.Status != common.UserStatusEnabled {
	if err := DB.First(&user).Error; err != nil {
		common.SysLog("no user exists, create a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		rootUser := User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
			AccessToken: common.GetUUID(),
			Quota:       100000000,
		}
		DB.Create(&rootUser)
	}
	return nil
}

func CountTable(tableName string) (num int64) {
	DB.Table(tableName).Count(&num)
	return
}

func InitDB() (err error) {
	var db *gorm.DB
	//if os.Getenv("SQL_DSN") != "" {
	//	// Use MySQL
	//	common.SysLog("using MySQL as database")
	//	db, err = gorm.Open(mysql.Open(os.Getenv("SQL_DSN")), &gorm.Config{
	//		PrepareStmt: true, // precompile SQL
	//	})
	//} else {
	//	// Use SQLite
	//	common.SysLog("SQL_DSN not set, using SQLite as database")
	//	common.UsingSQLite = true
	//	db, err = gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{
	//		PrepareStmt: true, // precompile SQL
	//	})
	//}
	common.SysLog("using MySQL as database")
	db, err = gorm.Open(mysql.Open("dev_openai:3ZSpSnMYC38CZw24@tcp(45.154.13.71:3306)/dev_openai"), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
	//db, err = gorm.Open(mysql.Open("root:ul25oU1db964@tcp(zeabur-gcp-asia-east1-1.clusters.zeabur.com:30011)/dev_openai"), &gorm.Config{
	//	PrepareStmt: true, // precompile SQL
	//})
	common.SysLog("database connected")
	if err == nil {
		DB = db
		if !common.IsMasterNode {
			return nil
		}
		err := db.AutoMigrate(&Channel{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Token{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Option{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Redemption{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Ability{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Log{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&MDatasets{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Order{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Midjourney{})
		if err != nil {
			return err
		}
		common.SysLog("database migrated")
		err = createRootAccountIfNeed()
		return err
	} else {
		common.FatalLog(err)
	}
	return err
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
