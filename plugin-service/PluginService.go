package plugin_service

import (
	"github.com/jsen-joker/goos/core/support-plugin"
	"github.com/jsen-joker/goos/core/support-plugin/manager"
	"github.com/jsen-joker/goos/plugin-service/controller"
)

type PluginService struct {
}
func (p *PluginService) Init()  {

	manager.RegisterRouter(&manager.Route{Name: "echo",     Method:"GET",   Pattern:"/serviceEcho", HandlerFunc: controller.ServiceEcho})

}
func (p *PluginService) Start()  {
}

func CreatePlugin() *support_plugin.Plugin {
	return &support_plugin.Plugin{
		PluginMeta: support_plugin.PluginMeta{
			Name: "service",
			Version: "",
		},
		PluginBoot: new(PluginService),
	}
}