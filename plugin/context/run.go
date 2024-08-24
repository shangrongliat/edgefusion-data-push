package context

import (
	"edgefusion-data-push/plugin/logs"
	"flag"
	"os"
	"runtime/debug"
)

// Run service
func Run(handle func(Context) error) {

	var h bool
	var c string
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&c, "c", "etc/conf.yml", "the configuration file")
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	ctx := NewContext(c)
	defer func() {
		if r := recover(); r != nil {
			ctx.Log().Error("service is stopped with panic", logs.Any("panic", r), logs.Any("stack", string(debug.Stack())))
		}
	}()

	pwd, _ := os.Getwd()
	ctx.Log().Info("service starting", logs.Any("pwd", pwd))
	err := handle(ctx)
	if err != nil {
		ctx.Log().Error("service has stopped with error", logs.Error(err))
	} else {
		ctx.Log().Info("service has stopped")
	}
}
