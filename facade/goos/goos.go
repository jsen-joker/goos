package main

import (
	"github.com/jsen-joker/goos/core/env"
	"github.com/jsen-joker/goos/core/support-plugin"
	"github.com/jsen-joker/goos/facade/goos/lifecycle"
)

//var plugins = &support_plugin.Plugins{}
//
//func checkVersion(plugin *support_plugin.Plugin) {
//	if utils.Empty(&plugin.Version) {
//		plugin.Version = APP_VERSION
//	}
//}


func init() {

	// 初始化环境变量
	env.InitEnv()

	life := lifecycle.GoosLifecycle{
		GoosPlugins: &support_plugin.Plugins{},
		AppVersion: env.GoosVersion,
	}
	life.Init()

	defer life.BeforeDestroy()
	//// 初始化环境变量
	//env.InitEnv()
	//
	//// 初始化数据库
	//support_db.Init()
	//defer support_db.Close()
	//
	//
	//// 插件初始化
	//plugins.Init()
	//var plugin = support_plugin.ReflectCreatePlugin(plugin_security.CreatePlugin)
	//checkVersion(plugin)
	//plugins.Register(plugin.Name, plugin)
	//plugin = support_plugin.ReflectCreatePlugin(plugin_config.CreatePlugin)
	//checkVersion(plugin)
	//plugins.Register(plugin.Name, plugin)
	//plugin = support_plugin.ReflectCreatePlugin(plugin_goos_ui.CreatePlugin)
	//checkVersion(plugin)
	//plugins.Register(plugin.Name, plugin)
	//plugin = support_plugin.ReflectCreatePlugin(plugin_service.CreatePlugin)
	//checkVersion(plugin)
	//plugins.Register(plugin.Name, plugin)
	//// 完成插件注册
	//
	//// 初始化插件
	//for _, p := range plugins.PluginList() {
	//	p.Init()
	//}
	//// 创建router
	//router := manager.CreateRouter()
	//log.Fatal(http.ListenAndServe(":8080", router))
	//// 启动插件
	//for _, p := range plugins.PluginList() {
	//	p.Start()
	//}
}

func main() {

}
