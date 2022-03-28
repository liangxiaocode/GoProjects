package main

import (
	"PassageTwo/global"
	"PassageTwo/internal/model"
	"PassageTwo/internal/routers"
	"PassageTwo/pkg/logger"
	"PassageTwo/pkg/setting"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 控制应用程序初始化的流程
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init setupSetting failed,err:%v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init setupLogger failed,err:%v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init setupDBEngine failed,err:%v", err)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		log.Fatalf("init setupSetting failed,err:%v", err)
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		log.Fatalf("init Read Server failed,err:%v", err)
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		log.Fatalf("init Read App failed,err:%v", err)
		return err
	}
	err = setting.ReadSection("Database", &global.DataBaseSetting)
	if err != nil {
		log.Fatalf("init Read Database failed,err:%v", err)
		return err
	}
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}

func setupLogger() error {
	filename := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileText // filename = stroage/logs/app.log
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  filename,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

// 初始化数据库连接
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DataBaseSetting)
	if err != nil {
		return err
	}
	return nil
}

// @title 博客系统
// @version 1.0
// @description Go语言编程之旅:一起用Go做项目
// @termsOfService https://gitee.com/liangluco_worker/projects
func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	// })
	// r.Run(":9000")

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	// 验证配置是否真正映射到配置的结构体上
	// fmt.Printf("global.ServerSetting:%#v\n", global.ServerSetting)
	// fmt.Printf("global.AppSetting:%#v\n", global.AppSetting)
	// fmt.Printf("global.DataBaseSetting:%#v\n", global.DataBaseSetting)

	// global.Logger.Infof("%s: go-programming-tour-book %s", "liangxiaocode", "blog_service")

	s.ListenAndServe()
}
