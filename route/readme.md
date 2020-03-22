##路由说明

    在route下面新建路由目录，比如api
    在api目录下新建路由，路由按功能拆分，
    比如文章板块article.go，用户板块user.go
    
##惯例和约定

    路由目录和文件名小写驼峰

    路由文件包含一个加载函数，例如：
    
    func LoadArticleRoute(r *gin.Engine) {
    	g := r.Group("/article")
    
    	g.POST("/add", func(c *gin.Context) {
    		ArticleService.Add(request.New(c))
    	})
    }

    
##加载路由
    
    在route目录下的route.go中调用所有的路由加载函数