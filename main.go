package main

import (
	cx "context"
	"runtime"

	"edgefusion-data-push/api"
	"edgefusion-data-push/plugin/config"
	"edgefusion-data-push/plugin/context"
	"edgefusion-data-push/routers"
	"edgefusion-data-push/service"
	dm "github.com/sineycoder/gorm-dm"
	"gorm.io/gorm"
)

func main() {
	defer context.ClosePlugins()
	runtime.GOMAXPROCS(runtime.NumCPU())
	context.Run(func(ctx context.Context) error {
		var cfg config.Config
		err := ctx.LoadCustomConfig(&cfg)
		if err != nil {
			return err
		}
		db, err := NewDB(cfg)
		if err != nil {
			return err
		}

		newAPI, err := api.NewAPI(&cfg)
		if err != nil {
			return err
		}
		client := service.NewMqttClient(&cfg)

		rs, err := routers.NewServer(&cfg)
		if err != nil {
			return err
		}
		rs.SetAPI(newAPI)
		rs.GetRoute().Use(context.RegisterDatabase(db))
		x := cx.Background()
		x = cx.WithValue(x, "db", db)
		rs.SyncRouter(x, rs.GetRoute())
		context.DatabaseSetHandle(db)
		go rs.Run()
		defer rs.Close()
		go client.Start()
		defer client.Close()
		ctx.Log().Info("init data-push server starting")
		ctx.Wait()
		return nil
	})
}

func NewDB(conf config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dm.Open(conf.Database.Url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
