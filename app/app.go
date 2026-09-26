package app

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/Ezrabro/subs-check/app/monitor"
	"github.com/Ezrabro/subs-check/assets"
	"github.com/Ezrabro/subs-check/check"
	"github.com/Ezrabro/subs-check/config"
	"github.com/Ezrabro/subs-check/save"
	"github.com/Ezrabro/subs-check/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
)

// App struct manages application state
type App struct {
	configPath string
	interval   int
	watcher    *fsnotify.Watcher
	checkChan  chan struct{} // Channel to trigger checks
	checking   atomic.Bool   // Check status flag
	ticker     *time.Ticker
	done       chan struct{} // Signal to end ticker goroutine
	cron       *cron.Cron    // Crontab scheduler
	version    string
}

// New creates a new app instance
func New(version string) *App {
	configPath := flag.String("f", "", "config文件路径")
	flag.Parse()

	return &App{
		configPath: *configPath,
		checkChan:  make(chan struct{}),
		done:       make(chan struct{}),
		version:    version,
	}
}

// Initialize initializes the application
func (app *App) Initialize() error {
	// Initialize config file path
	if err := app.initConfigPath(); err != nil {
		return fmt.Errorf("initconfig文件路径失败: %w", err)
	}
	// Results snapshot and export cache live beside the config file.
	utils.SetCacheDir(app.configPath)

	// Load configuration file
	if err := app.loadConfig(); err != nil {
		return fmt.Errorf("加载config文件失败: %w", err)
	}

	// init DNS resolver（必须在任何 proxy 连接之前，影响 mihomo 全局 resolver）
	if err := initResolver(); err != nil {
		return fmt.Errorf("init DNS 失败: %w", err)
	}

	// initconfig文件监听
	if err := app.initConfigWatcher(); err != nil {
		return fmt.Errorf("initconfig文件监听失败: %w", err)
	}

	// 从config文件中读取代理，设置代理
	if config.GlobalConfig.Proxy != "" {
		os.Setenv("HTTP_PROXY", config.GlobalConfig.Proxy)
		os.Setenv("HTTPS_PROXY", config.GlobalConfig.Proxy)
	}

	app.interval = func() int {
		if config.GlobalConfig.CheckInterval <= 0 {
			return 1
		}
		return config.GlobalConfig.CheckInterval
	}()

	if config.GlobalConfig.ListenPort != "" {
		if err := app.initHttpServer(); err != nil {
			return fmt.Errorf("initHTTP服务器失败: %w", err)
		}
	}

	if config.GlobalConfig.SubStorePort != "" {
		if runtime.GOOS == "linux" && runtime.GOARCH == "386" {
			slog.Warn("node不支持Linux 32位系统，不启动sub-store服务")
		}
		go assets.RunSubStoreService()
		// 求等吗得，Logs会按预期顺序输出
		time.Sleep(500 * time.Millisecond)
	}

	// 启动内存监控
	monitor.StartMemoryMonitor()

	// 设置信号处理器
	utils.SetupSignalHandler(check.RequestCancel)
	return nil
}

// Run run应用程序主循环
func (app *App) Run() {
	defer func() {
		app.watcher.Close()
		if app.ticker != nil {
			app.ticker.Stop()
		}
		if app.cron != nil {
			app.cron.Stop()
		}
	}()

	// 设置初始定时器模式
	app.setTimer()

	// 仅在cron表达式empty时，首次启动立即执行check
	if config.GlobalConfig.CronExpression != "" {
		slog.Warn("usecron表达式，首次启动不立即执行check")
	} else {
		app.triggerCheck()
	}

	// 在主循环中处理手动触发
	for range app.checkChan {
		go app.triggerCheck()
	}
}

// setTimer 根据config设置定时器
func (app *App) setTimer() {
	// 停止现有定时器
	if app.ticker != nil {
		// 应该先发送停止信号，防止被=nil后panic
		close(app.done)                // 发送停止信号
		app.done = make(chan struct{}) // 创建新通道
		app.ticker.Stop()
		app.ticker = nil
	}

	// 停止现有cron
	if app.cron != nil {
		app.cron.Stop()
		app.cron = nil
	}

	// 检查是否设置了cron表达式
	if config.GlobalConfig.CronExpression != "" {
		slog.Info(fmt.Sprintf("usecron表达式: %s", config.GlobalConfig.CronExpression))
		app.cron = cron.New()
		_, err := app.cron.AddFunc(config.GlobalConfig.CronExpression, func() {
			app.triggerCheck()
		})
		if err != nil {
			slog.Error(fmt.Sprintf("cron表达式 '%s' 解析失败: %v，将use检查间隔时间",
				config.GlobalConfig.CronExpression, err))
			// use间隔时间
			app.useIntervalTimer()
		} else {
			app.cron.Start()
		}
	} else {
		// use间隔时间
		app.useIntervalTimer()
	}
}

// useIntervalTimer use间隔时间模式run
func (app *App) useIntervalTimer() {
	// init定时器
	app.ticker = time.NewTicker(time.Duration(app.interval) * time.Minute)
	done := app.done
	// 启动一个goroutine监听定时器事件
	go func() {
		for {
			select {
			case <-app.ticker.C:
				app.triggerCheck()
			case <-done:
				return // 收到停止信号，退出goroutine
			}
		}
	}()
}

// TriggerCheck 供外部调用的触发check方法
func (app *App) TriggerCheck() {
	select {
	case app.checkChan <- struct{}{}:
		slog.Info("手动触发check")
	default:
		slog.Warn("已有check正在进行，忽略本次触发")
	}
}

// triggerCheck 内部check方法
func (app *App) triggerCheck() {
	// 如果已经在check中，直接返回
	if !app.checking.CompareAndSwap(false, true) {
		slog.Warn("已有check正在进行，跳过本次check")
		return
	}
	defer app.checking.Store(false)

	if err := app.checkProxies(); err != nil {
		slog.Error(fmt.Sprintf("check代理失败: %v", err))
		os.Exit(1)
	}

	// check完成后显示下次检查时间
	if app.ticker != nil {
		// use间隔时间模式
		app.ticker.Reset(time.Duration(app.interval) * time.Minute)
		nextCheck := time.Now().Add(time.Duration(app.interval) * time.Minute)
		slog.Info(fmt.Sprintf("下次检查时间: %s", nextCheck.Format("2006-01-02 15:04:05")))
	} else if app.cron != nil {
		// usecron模式
		entries := app.cron.Entries()
		if len(entries) > 0 {
			nextTime := entries[0].Next
			slog.Info(fmt.Sprintf("下次检查时间: %s", nextTime.Format("2006-01-02 15:04:05")))
		}
	}
	debug.FreeOSMemory()
}

// checkProxies 执行代理check
func (app *App) checkProxies() error {
	slog.Info("开始准备check代理", "进度展示", config.GlobalConfig.PrintProgress)

	// 加载历史可用node到待测队列
	if config.GlobalConfig.KeepDays > 0 {
		if hp := save.LoadHistoryProxies(); len(hp) > 0 {
			config.GlobalProxies = append(config.GlobalProxies, hp...)
		}
	}

	results, err := check.Check()
	if err != nil {
		return fmt.Errorf("check代理失败: %w", err)
	}
	slog.Info("check完成")
	save.SaveConfig(results)
	utils.SendNotify(len(results))
	utils.UpdateSubs()

	// 执行回调脚本
	utils.ExecuteCallback(len(results))

	return nil
}
