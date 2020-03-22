##模型说明

    本项目未使用ORM库,而是采用了近似sql的sqlx
    https://github.com/jmoiron/sqlx

    model的定义不需要严格按照数据库表定义

    比如user表的字段 [id,name,age,created_at...],有可能你当前这个项目只用到了[id,name]，
    那你的model只需要定义用到的字段

    type User struct {
        ID          uint64     `db:"id" `
        Name        *string    `db:"user_name"`
    }
    由于我们使用的sqlx本不是ORM，定义struct的意义只在于接受结果，而不一定要和数据表一一映射，
    只要你定义的struct字段包含sqlx结果集返回的字段就行

##惯例和约定

    比如我要新建一个User模型，在model中新建 UserModel文件夹【文件夹名大写驼峰】，
    然后在UserModel下新建 UserModel文件【文件名大写驼峰】，
    可能有一些字段我们并不想返回给前端，所以进行资源化很有必要，通常在还需要
    新建一个Source文件


##迁移文件

    项目中最好包含数据表的迁移文件，在model下面新建migrate.sql文件，
    包含创建数据表的sql代码