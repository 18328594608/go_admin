package main

import (
	"context"
	"flag"
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/logger"
	"go_web_app/route"
	"go_web_app/setting"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	filePath := ""
	// 第二个参数 就是命令行的-x  第三个是默认值 第四个是-h的时候的说明
	flag.StringVar(&filePath, "file", "", "配置文件地址")
	flag.Parse()
	// 加载配置文件
	if err := setting.Init(filePath); err != nil {
		fmt.Printf("init settings failed:%s \n", err)
		return
	}
	// 初始化日志f
	if err := logger.Init(setting.Config.LogConfig, setting.Config.Mode); err != nil {
		fmt.Printf("init settings failed:%s \n", err)
		return
	}
	zap.L().Debug("logger init success")
	defer zap.L().Sync()
	// 不要遗漏2个 db的close

	// 初始化mysql
	if err := mysql.Init(setting.Config.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed:%s \n", err)
		return
	}
	zap.L().Debug("mysql init success")

	mysql.LoadConfig()
	defer mysql.Close()
	// 注册路由
	r := route.Setup(setting.Config.Mode)
	// 启动服务 （优雅关机）

	zap.L().Debug("server listen", zap.String("port", strconv.Itoa(setting.Config.Port)))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Config.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown server")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
		zap.L().Error("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
