# quick-gin

## 介绍
    quick-gin是一个基于gin框架的脚手架，基于gin，做了如下一些整理
    
    01. 规划了一套代码组织架构
    02. 对gin.Context进行了扩展，增加很多辅助方法
    03. 集成跨域中间件
    04. 集成请求缓存中间件，
    05. 使用GORM
    06. 封装Jwt
    07. 已经配置好MySQL，Redis 开箱即用
    08. Sentry日志
    09. env全局配置
    10. 自带平滑重启
    11. 完整的代码实例和说明文档
    12. 命令行操作 create_gin_app app_name
    
## 代码组织架构
    Route->Cache->Service->Model
    
    Route：      路由规则和中间件
    Cache：      对需要的请求做缓存
    Service：    验证请求，逻辑处理
    Model：      数据库CURD，数据资源化
    
## 目录说明
    /config/        配置
    /middleware/    中间件
    /model/         模型
    /route/         路由规则
    /service/       业务逻辑
    /util/          工具类
    /config.ini     整个项目的配置文件，采用ini配置方式
    /main.go        负责整个服务的启动，
        
## 相关目录下还有更详细的惯例和说明

    开始对应模块前，请先阅读对应的惯例和约定

## 如何开始上手编写业务
    
    1. 检查config.ini     修改配置参数
    2. 新建route文件       路由映射到service的方法
    3. 新建service文件     用到的数据由model层提供
    4. 新建model文件       负责数据库CURD
    
    上述步骤可以参考自带的实例
    
## 注意事项

    确保go module开启
    推荐设置 goproxy，推荐 https://mirrors.aliyun.com/goproxy/
    请在config.ini中修改mysql和redis的配置
    