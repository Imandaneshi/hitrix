package controller

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/dto/indexes"
	"github.com/coretrix/hitrix/pkg/dto/list"
	"github.com/coretrix/hitrix/pkg/entity"
	errorhandling "github.com/coretrix/hitrix/pkg/error_handling"
	"github.com/coretrix/hitrix/pkg/errors"
	accountModel "github.com/coretrix/hitrix/pkg/model/account"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/coretrix/hitrix/pkg/view/requestlogger"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type MenuItem struct {
	Label string
	URL   string
	Icon  string
}

type DevPanelController struct {
}

//func (controller *DevPanelController) GetActionListAction(c *gin.Context) {
//	actions := []*MenuItem{
//		{
//			Label: "Clear Cache",
//			URL:   "/dev/clear-cache/",
//			Icon:  "mdiCached",
//		},
//	}
//
//	c.JSON(200, actions)
//}

func (controller *DevPanelController) GetSettingsAction(c *gin.Context) {
	appService := service.DI().App()

	response.SuccessResponse(c, gin.H{
		"AppMode": appService.Mode,
	})
}

func (controller *DevPanelController) CreateDevPanelUserAction(c *gin.Context) {
	passwordService := service.DI().Password()

	ormService := service.DI().OrmEngine()

	form := &accountModel.LoginDevForm{}
	if err := binding.ShouldBindQuery(c, form); err != nil {
		fieldError, ok := (err).(errors.FieldErrors)
		if ok {
			response.ErrorResponseFields(c, fieldError, nil)

			return
		}

		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	adminEntity := service.DI().App().DevPanel.UserEntity

	passwordHash, err := passwordService.HashPassword(form.Password)
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)
	}

	adminTableSchema := ormService.GetRegistry().GetEntitySchemaForEntity(adminEntity)
	response.SuccessResponse(
		c,
		fmt.Sprintf(`INSERT INTO %s (Email, Password) VALUES('%s', '%s')`, adminTableSchema.GetTableName(), form.Username, passwordHash))
}

func (controller *DevPanelController) PostLoginDevPanelAction(c *gin.Context) {
	loginForm := accountModel.LoginDevForm{}
	token, refreshToken, err := loginForm.Login(c)

	errType, ok := err.(errors.FieldErrors)

	if ok && errType != nil {
		response.ErrorResponseFields(c, errType, nil)

		return
	}

	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) PostGenerateTokenAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	devPanelUserEntity := c.MustGet(account.LoggedDevPanelUserEntity).(app.IDevPanelUserEntity)

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		response.ErrorResponseGlobal(c, err, nil)

		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (controller *DevPanelController) GetClearCacheAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	redisService := ormService.GetRedis()

	redisService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetClearRedisStreamsAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Stream == "" {
		panic("stream pool is not defined")
	}

	redisStreamsService := ormService.GetRedis(appService.RedisPools.Stream)
	redisStreamsService.FlushDB()

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) DeleteRedisStreamAction(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Stream == "" {
		panic("stream pool is not defined")
	}

	redisStreamService := ormService.GetRedis(appService.RedisPools.Stream)

	name := c.Param("name")
	if name == "" {
		panic("provide stream name")
	}

	redisStreamService.XTrim(name, 0)

	c.JSON(200, gin.H{})
}

func (controller *DevPanelController) GetAlters(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	alters := ormService.GetAlters()
	result := make([]string, len(alters))

	force := c.Query("force")
	if force != "" {
		redisService := ormService.GetRedis()
		redisService.FlushDB()
	}

	for i, alter := range alters {
		if force != "" {
			alter.Exec(ormService.Engine)
		} else {
			result[i] = alter.SQL
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisStreams(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	ormService.GetEventBroker().GetStreamsStatistics()

	stats := ormService.GetEventBroker().GetStreamsStatistics()
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Stream < stats[j].Stream
	})
	response.SuccessResponse(c, stats)
}

// GetRedisStatistics TODO: check if this is missing with Lukasz
func (controller *DevPanelController) GetRedisStatistics(_ *gin.Context) {
	//ormService := service.DI().OrmEngineForContext(c.Request.Context())

	//stats := tools.GetRedisStatistics(ormService)
	//sort.Slice(stats, func(i, j int) bool {
	//	return stats[i].RedisPool < stats[j].RedisPool
	//})
	//response.SuccessResponse(c, stats)
}

func (controller *DevPanelController) GetRedisSearchStatistics(c *gin.Context) {
	response.SuccessResponse(c, service.DI().OrmEngineForContext(c.Request.Context()).GetRedisSearchStatistics())
}

func (controller *DevPanelController) GetRedisSearchAlters(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	altersSearch := ormService.GetRedisSearchAlters()
	result := make([]map[string]string, len(altersSearch))

	force := c.Query("force")
	for i, alter := range altersSearch {
		if force != "" {
			alter.Execute()
		} else {
			result[i] = map[string]string{
				"Query":   alter.Query,
				"Changes": strings.Join(alter.Changes, " | "),
			}
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) GetRedisSearchIndexes(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	stats := ormService.RedisSearchEngine.GetRedisSearchStatistics()

	indexList := make([]indexes.Index, len(stats))

	for i, stat := range stats {
		indexList[i] = indexes.Index{
			Name:      stat.Index.Name,
			TotalDocs: stat.Info.NumDocs,
			TotalSize: uint64(
				stat.Info.DocTableSizeMB +
					stat.Info.KeyTableSizeMB +
					stat.Info.SortableValuesSizeMB +
					stat.Info.InvertedSzMB +
					stat.Info.OffsetVectorsSzMB,
			),
		}
	}

	result := indexes.ResponseDTOList{
		Indexes: indexList,
	}
	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) PostRedisSearchForceReindex(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	indexName := c.Param("index")
	if indexName == "" {
		response.ErrorResponseGlobal(c, "index is required", nil)

		return
	}

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	ormService.ForceReindex(indexName)
	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostRedisSearchForceReindexAll(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	rsIndexes := ormService.ListIndices()

	concurrently := c.Query("concurrently")
	if concurrently != "" {
		ormService := service.DI().OrmEngineForContext(c.Request.Context())

		wg := sync.WaitGroup{}
		wg.Add(len(rsIndexes))

		for _, index := range rsIndexes {
			go func(index string) {
				defer func() {
					if r := recover(); r != nil {
						service.DI().ErrorLogger().LogError(r)
					}
				}()

				ormService.ForceReindex(index)
				wg.Done()
			}(index)
		}

		wg.Wait()
	} else {
		for _, index := range rsIndexes {
			ormService.ForceReindex(index)
		}
	}

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostRedisSearchIndexInfo(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	indexName := c.Param("index")
	if indexName == "" {
		response.ErrorResponseGlobal(c, "index is required", nil)

		return
	}

	appService := service.DI().App()
	if appService.DevPanel == nil || appService.RedisPools.Search == "" {
		panic("stream pool is not defined")
	}

	response.SuccessResponse(c, ormService.Info(indexName))
}

func (controller *DevPanelController) GetFeatureFlags(c *gin.Context) {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := redisearch.NewRedisSearchQuery()
	var featureFlagEntities []*entity.FeatureFlagEntity

	ormService.RedisSearch(query, beeorm.NewPager(1, 1000), &featureFlagEntities)

	type feature struct {
		Name       string
		Registered bool
		Enabled    bool
	}

	result := make([]*feature, len(featureFlagEntities))

	for i, featureFlagEntity := range featureFlagEntities {
		result[i] = &feature{
			Name:       featureFlagEntity.Name,
			Registered: featureFlagEntity.Registered,
			Enabled:    featureFlagEntity.Enabled,
		}
	}

	response.SuccessResponse(c, result)
}

func (controller *DevPanelController) PostEnableFeatureFlag(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.ErrorResponseGlobal(c, "name is required", nil)

		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := redisearch.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)

	if !found {
		response.ErrorResponseGlobal(c, "feature is missing", nil)

		return
	}

	featureFlagEntity.Enabled = true
	ormService.Flush(featureFlagEntity)

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) PostDisableFeatureFlag(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.ErrorResponseGlobal(c, "name is required", nil)

		return
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	query := redisearch.NewRedisSearchQuery()
	query.FilterString("Name", name)

	featureFlagEntity := &entity.FeatureFlagEntity{}
	found := ormService.RedisSearchOne(featureFlagEntity, query)

	if !found {
		response.ErrorResponseGlobal(c, "feature is missing", nil)

		return
	}

	featureFlagEntity.Enabled = false
	ormService.Flush(featureFlagEntity)

	response.SuccessResponse(c, nil)
}

func (controller *DevPanelController) GetEnvValues(c *gin.Context) {
	response.SuccessResponse(c, map[string]interface{}{"list": os.Environ()})
}

func (controller *DevPanelController) PostRequestsLogger(c *gin.Context) {
	request := list.RequestDTOList{}

	err := binding.ShouldBindJSON(c, &request)
	if errorhandling.HandleError(c, err) {
		return
	}

	res, err := requestlogger.RequestsLogger(c.Request.Context(), request)
	if errorhandling.HandleError(c, err) {
		return
	}

	response.SuccessResponse(c, res)
}
