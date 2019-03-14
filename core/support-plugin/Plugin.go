package support_plugin

type PluginMeta struct {
	// 插件名字（保证唯一性）
	Name string
	// 插件版本
	Version string
}

type PluginBoot interface {
	Init()
	Start()
}

type Plugin struct {
	PluginMeta
	PluginBoot
}

type Plugins struct {
	pluginList map[string] *Plugin
}

func (p *Plugins) Init() {
	p.pluginList = make(map[string] *Plugin)
}
func (p *Plugins) Register(name string, plugin *Plugin) {
	p.pluginList[name] = plugin
}
func (p *Plugins) PluginList() map[string] *Plugin {
	return p.pluginList
}