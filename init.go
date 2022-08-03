package db_importer

import (
	"github.com/ismdeep/log"
	"os"
)

func init() {
	if configRoot := os.Getenv("CONFIG_ROOT"); configRoot != "" {
		log.Info("init", log.String("configRoot", configRoot))
		if err := Migrate(configRoot); err != nil {
			panic(err)
		}
	}
}
