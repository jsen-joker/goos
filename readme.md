# Goos

### Api设计原则
    
    在良好的区分度的情况下，尽可能的精简Api，（例如，create和update可以合并为一个Api，在统一范式的情况下，这是一种可读性较高的设计方法）
    尽可能地暴露最少的接口
    尽可能的对接口进行权限验证（紧的权限策略）
    Api命名规范 /api/{PluginName}/{Resource}，对应的CRUD用（POST、GET、POST、DELETE）这些方法赋予语义
    Resource具体规范，一般get one 直接用Resource名字即可，如果是get一个list，则用{Resource}List（加List后缀表示）
    Api对应方法命名规范：由于Resource缺少CRUD语义，因此方法的命名要在Resource前面加上对应的Get，Delete，Post方法使用默认Resource
    
    