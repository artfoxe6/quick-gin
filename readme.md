# quick-gin
基于```gin```的一个轻量```api```骨架，可以帮助你快速开始业务开发

#### 概括
之前做过其他语言的Web开发，切换到gin后会感觉少了点什么，最明显的一点他没有提供一个推荐或默认的目录结构<br>
比如常见的```controller```，```model```，```service```，```repository```等等<br>
这会造成新手不知道如何去组织和划分代码<br>
还有就是没有提供一些开箱即用的Web通用组件，比如数据库，缓存等

#### 目录结构
  ```bash
    project
    ├── build 二进制文件打包目录
    ├── cmd 应用目录，可能存放多个应用
    │   └── app 具体的应用，名字根据业务来取就行
    │       └── main.go 应用的入口，只负责启动服务，不应该涉及业务代码
    ├── config 配置目录
    │   ├── config.ini 这个文件一般不应该提交
    │   └── config.ini.example
    ├── go.mod
    ├── go.sum
    ├── internal 由于业务项目一般是不对外公开的，也不存在被其他项目导入的可能，所以放到internal中
    │   ├── app 对应cmd里面的具体应用
    │   │   ├── api => controller
    │   │   │   ├── base.go 
    │   │   ├── model 这里面只定义模型和关联关系，不建议在这里写CRUD
    │   │   │   ├── base.go
    │   │   ├── router.go 路由
    │   │   └── service
    │   │       ├── student.go
    │   └── pkg
    ├── LICENSE
    ├── log
    │   ├── req20230513.log
    │   └── runtime.log
    ├── Makefile
    ├── readme.md
    └── test
  ```

#### 开始使用
- 克隆仓库到本地
- 进入根目录执行初始化命令
  ```
    ./init       Linux
    ./init_mac   Mac
    ./init.exe   Windows
  ```
- 启动服务
  ```
  make
  ```
    详情查看Makefile文件，里面定义几种常用命令<br>
    make build 打包应用<br>
    make stop 停止服务<br>
    make run 启动服务<br>
    make restart 等于 build+stop+run，不带参数时默认执行这个<br>
    之所以用make方式是因为现在网络上很多开源Golang项目都在用<br>
    如果不喜欢你依然可以使用 go run cmd/main.go 这种方式<br>


#### 温馨提示
- 根目录结构参考自 https://github.com/golang-standards/project-layout
- 业务层次结构参考了一些开源项目，然后结合一些自己的经验摸索出来的，大家可以参考借鉴，不一定适用所有项目，<br>
  如果有更加优雅的组织方式，欢迎留言讨论
- 每一个业务层具体干什么这里在啰嗦一下 
  >controller层 参数验证，返回处理结果
  > 
  >model层 推荐只定义模型和关联关系，不涉及CRUD,
  但是如果你的项目极小，可以考虑在这里写CRUD，砍掉service层，
  相当于把service层的功能分摊到model和controller层了
  > 
  >service层 业务逻辑，调用各种pkg,数据库CRUD
  但是如果你的项目很复杂，可能需要复用或者缓存CURD，这个时候你可以在抽一个repository层出来专门处理数据库,
  相当于把service在细拆


    

    