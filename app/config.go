package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Ezrabro/subs-check/config"
	"github.com/Ezrabro/subs-check/utils"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// initConfigPath initconfig文件路径
func (app *App) initConfigPath() error {
	if app.configPath == "" {
		execPath := utils.GetExecutablePath()
		configDir := filepath.Join(execPath, "config")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("创建config目录失败: %w", err)
		}

		app.configPath = filepath.Join(configDir, "config.yaml")
	}
	return nil
}

// loadConfig 加载config文件
func (app *App) loadConfig() error {
	yamlFile, err := os.ReadFile(app.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return app.createDefaultConfig()
		}
		return fmt.Errorf("读取config文件失败: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, config.GlobalConfig); err != nil {
		return fmt.Errorf("解析config文件失败: %w", err)
	}

	slog.Info("config文件读取成功")
	return nil
}

// createDefaultConfig 创建默认config文件
func (app *App) createDefaultConfig() error {
	slog.Info("config文件不存在，创建默认config文件")

	if err := os.WriteFile(app.configPath, []byte(config.DefaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("写入默认config文件失败: %w", err)
	}

	slog.Info("默认config文件创建成功")
	slog.Info(fmt.Sprintf("请editconfig文件: %s", app.configPath))
	os.Exit(0)
	return nil
}

// initConfigWatcher initconfig文件监听
func (app *App) initConfigWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败: %w", err)
	}

	app.watcher = watcher

	// 防抖定时器，防止vscode等软件先临时创建文件在覆盖，会产生两次write事件
	var debounceTimer *time.Timer
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if absPath, _ := filepath.Abs(app.configPath); event.Name != absPath {
					continue
				}
				// 兼容容器外修改
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					// 如果定时器存在，重置它
					if debounceTimer != nil {
						debounceTimer.Stop()
					}

					// 创建新的定时器，延迟100ms执行
					debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
						slog.Info("config文件发生变化，正在重新加载")
						oldCronExpr := config.GlobalConfig.CronExpression
						oldInterval := app.interval

						if err := app.loadConfig(); err != nil {
							slog.Error(fmt.Sprintf("重新加载config文件失败: %v", err))
							return
						}

						// 检查cron表达式或check间隔是否变化
						if oldCronExpr != config.GlobalConfig.CronExpression ||
							oldInterval != config.GlobalConfig.CheckInterval {

							app.interval = func() int {
								if config.GlobalConfig.CheckInterval <= 0 {
									return 1
								}
								return config.GlobalConfig.CheckInterval
							}()
							slog.Warn("check设置发生变化，重新config定时器")

							// usesetTimer方法重新设置定时器
							app.setTimer()
						}
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error(fmt.Sprintf("config文件监听错误: %v", err))
			}
		}
	}()

	// Start监听config文件目录
	if err := watcher.Add(filepath.Dir(app.configPath)); err != nil {
		return fmt.Errorf("添加config文件监听失败: %w", err)
	}

	slog.Info("config文件监听已启动")
	return nil
}
