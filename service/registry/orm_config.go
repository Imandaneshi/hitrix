package registry

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"fmt"
	"reflect"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"

	"github.com/latolukasz/beeorm"
)

type ORMRegistryInitFunc func(registry *beeorm.Registry)

func ServiceProviderOrmRegistry(init ORMRegistryInitFunc) *service.DefinitionGlobal {
	var defferFunc func()
	var err error
	var ormConfig beeorm.ValidatedRegistry
	var appService *app.App
	var configService config.IConfig

	return &service.DefinitionGlobal{
		Name: service.ORMConfigService,
		Build: func(ctn di.Container) (interface{}, error) {
			appService = ctn.Get(service.AppService).(*app.App)
			configService = ctn.Get(service.ConfigService).(config.IConfig)

			registry := beeorm.NewRegistry()

			configuration, ok := configService.Get("orm")
			if !ok {
				return nil, errors.New("no orm config")
			}

			yamlConfig := map[string]interface{}{}
			for k, v := range configuration.(map[interface{}]interface{}) {
				yamlConfig[fmt.Sprint(k)] = v
			}

			registry.InitByYaml(yamlConfig)
			init(registry)

			if appService.IsInTestMode() {
				overwriteORMConfig(appService, configService, registry, yamlConfig)
			}

			ormConfig, defferFunc, err = registry.Validate()
			return ormConfig, err
		},
		Close: func(obj interface{}) error {
			if appService.IsInTestMode() && os.Getenv("PARALLEL_TESTS") == "true" {
				removeDBs(appService, configService)
			}

			defferFunc()
			return nil
		},
	}
}

func removeDBs(appService *app.App, configService config.IConfig) {
	mysqlConnection := strings.Split(configService.MustString("orm.default.mysql"), "/")
	db, err := sql.Open("mysql", mysqlConnection[0]+"/?multiStatements=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	newDBName := "t_" + appService.ParallelTestID

	_, err = db.Exec("DROP DATABASE `" + newDBName + "`")

	if err != nil {
		panic(err)
	}
}

func overwriteORMConfig(appService *app.App, configService config.IConfig, registry *beeorm.Registry, yamlConfig map[string]interface{}) {
	mysqlConnection := strings.Split(configService.MustString("orm.default.mysql"), "/")
	db, err := sql.Open("mysql", mysqlConnection[0]+"/?multiStatements=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	newDBName := "t_" + appService.ParallelTestID
	color.Blue("DB name: %s", newDBName)

	if os.Getenv("PARALLEL_TESTS") == "true" {
		_, err = db.Exec("DROP DATABASE IF EXISTS `" + newDBName + "`; CREATE DATABASE `" + newDBName + "`")
	} else {
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + newDBName + "`")
	}

	if err != nil {
		panic(err)
	}

	registry.RegisterMySQLPool(mysqlConnection[0] + "/" + newDBName)

	for pool, value := range yamlConfig {
		if _, ok := value.(map[interface{}]interface{})["sentinel"]; ok {
			for masterConf, sentinelConnections := range value.(map[interface{}]interface{})["sentinel"].(map[interface{}]interface{}) {
				sentinelConn := make([]string, 0)
				sentinelConnValues := reflect.ValueOf(sentinelConnections)

				for i := 0; i < reflect.ValueOf(sentinelConnections).Len(); i++ {
					sentinelConn = append(sentinelConn, fmt.Sprint(sentinelConnValues.Index(i)))
				}

				settings := strings.Split(fmt.Sprint(masterConf), ":")
				dbIndex, _ := strconv.Atoi(settings[1])

				registry.RegisterRedisSentinel(settings[0], appService.ParallelTestID, dbIndex, sentinelConn, fmt.Sprint(pool))
			}
		}
	}
}
