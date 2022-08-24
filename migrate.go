package db_importer

// Author: L. Jiang <l.jiang.1024@gmail.com>

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/ismdeep/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type s struct {
	Dialect  string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Skips    []string
	CLI      string
}

// DBImporterMigrate model
type DBImporterMigrate struct {
	ID        string `gorm:"type:varchar(255);primary_key;"`      // fileName as ID
	Status    string `gorm:"type:varchar(20);not null;default:0"` // status: pending, success, failed, skipped
	FailedMsg string `gorm:"type:text"`                           // failed msg
	ExecStart time.Time
	ExecEnd   time.Time
}

const (
	// MigrateStatusSuccess success
	MigrateStatusSuccess = "SUCCESS"
	// MigrateStatusSkipped skipped
	MigrateStatusSkipped = "SKIPPED"
)

func isSkipped(fileName string, skippedList []string) bool {
	for _, s2 := range skippedList {
		if s2 == fileName {
			return true
		}
	}

	return false
}

func Migrate(configRoot string) error {
	var conf s

	v := viper.New()
	v.AddConfigPath(configRoot)
	v.SetConfigName("db-importer")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(&conf); err != nil {
		return err
	}

	var db *gorm.DB
	switch conf.Dialect {
	case "mysql", "mariadb":
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local&charset=utf8mb4,utf8",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "zzzz_",
			},
		})
		if err != nil {
			return err
		}
		db = conn
	default:
		panic(fmt.Errorf("unsupported db dialect. [%v]", conf.Dialect))
	}

	if err := db.AutoMigrate(&DBImporterMigrate{}); err != nil {
		return err
	}

	// 1. get sql file names (dict asc sort)
	files, err := ioutil.ReadDir(fmt.Sprintf("%v/sql", configRoot))
	if err != nil {
		return err
	}
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	sort.Strings(fileNames)

	for _, fileName := range fileNames {
		// 1. check if is already migrated
		var cnt int64
		if err := db.Model(&DBImporterMigrate{}).Where("id=?", fileName).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt >= 1 {
			//log.Warn("migrate", log.String("info", "already migrated"), log.String("name", fileName))
			continue
		}

		// 2. skipped
		if isSkipped(fileName, conf.Skips) {
			db.Create(&DBImporterMigrate{
				ID:        fileName,
				Status:    MigrateStatusSkipped,
				FailedMsg: "",
				ExecStart: time.Now(),
				ExecEnd:   time.Now(),
			})
			continue
		}

		// 3. run import command
		startTime := time.Now()
		sqlF, err := os.Open(fmt.Sprintf("%v/sql/%v", configRoot, fileName))
		if err != nil {
			return err
		}
		cmd := exec.Command(conf.CLI,
			fmt.Sprintf("-h%v", conf.Host),
			fmt.Sprintf("-P%v", conf.Port),
			fmt.Sprintf("-u%v", conf.Username),
			fmt.Sprintf("-p%v", conf.Password),
			conf.Database)
		cmd.Stdin = sqlF
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error("migrate",
				log.String("info", "migrate failed"),
				log.String("name", fileName),
				log.FieldErr(err))
			return err
		}
		db.Create(&DBImporterMigrate{
			ID:        fileName,
			Status:    MigrateStatusSuccess,
			FailedMsg: "",
			ExecStart: startTime,
			ExecEnd:   time.Now(),
		})
	}

	return nil
}
