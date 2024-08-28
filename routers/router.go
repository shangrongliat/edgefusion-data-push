package routers

import (
	"context"
	"embed"
	"net/http"
	"time"

	"edgefusion-data-push/api"
	"edgefusion-data-push/plugin/config"
	log "edgefusion-data-push/plugin/logs"
	"edgefusion-data-push/plugin/persist"
	"github.com/gin-gonic/gin"
)

const (
	DefaultAPICacheDuration = time.Second * 2
)

type Server struct {
	ExternalHandlers []gin.HandlerFunc
	APICache         persist.CacheStore

	cfg    *config.Config
	router *gin.Engine
	server *http.Server
	api    *api.API
	log    *log.Logger
}

func NewServer(config *config.Config) (*Server, error) {
	router := gin.New()
	server := &http.Server{
		Addr:           config.Server.Port,
		Handler:        router,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{
		cfg:      config,
		router:   router,
		server:   server,
		APICache: persist.NewInMemoryStore(DefaultAPICacheDuration),
		log:      log.L().With(log.Any("server", "EdgeFusionDataPushRouter")),
	}, nil
}

var File embed.FS

func (s *Server) SyncRouter(ctx context.Context, engine *gin.Engine) {
	s.router.NoRoute(NoRouteHandler)
	s.router.NoMethod(NoMethodHandler)
	s.router.GET("/health", Health)
	s.router.Use(RequestIDHandler)
	s.router.Use(LoggerHandler)
	s.router.Use(s.AuthHandler)
	s.router.Use(s.ExternalHandlers...)
	engine.Use(config.Cors())
	v1 := engine.Group("v1/hook")
	node := v1.Group("/dvr")
	{
		node.POST("callback", config.Wrapper(s.api.DvrFinish))
		node.GET("test", config.Wrapper(s.api.GetInfluxData))
		node.GET("Testw", config.Wrapper(s.api.WInfluxData))
	}
	live := v1.Group("/live")
	{
		live.POST("start", config.Wrapper(s.api.DvrFinish))
		live.POST("stop", config.Wrapper(s.api.DvrFinish))
	}
}

// GetRoute get router
func (r *Server) GetRoute() *gin.Engine {
	return r.router
}

func (s *Server) Run() {
	log.L().Debug("admin  server start: ", log.Any("addr", s.cfg.Server.Port))
	if err := s.server.ListenAndServe(); err != nil {
		log.L().Info("admin server stopped", log.Error(err))
	}
}

func (s *Server) SetAPI(api *api.API) {
	s.api = api
}

// Close server
func (s *Server) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTime)
	s.server.Shutdown(ctx)
}

// auth handler
func (s *Server) AuthHandler(c *gin.Context) {
	//cc := common.NewContext(c)
	//err := s.Auth.Authenticate(cc)
	//if err != nil {
	//	s.log.Error("request authenticate failed",
	//		log.Any(cc.GetTrace()),
	//		log.Any("namespace", cc.GetNamespace()),
	//		log.Any("authorization", c.Request.Header.Get("Authorization")),
	//		log.Error(err))
	//	common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied, common.Field("error", err)), true)
	//}
}
