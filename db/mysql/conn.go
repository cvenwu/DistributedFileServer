package mysql

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/15 4:08 下午
 * @Desc: 专门用于创建mysql连接的
 */

import (
	"database/sql"                     //golang提供的标准接口
	_ "github.com/go-sql-driver/mysql" //一般我们不会直接使用这个驱动的方法，而是通过golang提供的接口来访问，因此可以加入一个_，作用：当导入这个驱动的时候这个驱动就会进行初始化并且会将自己注册到database/sql 里面去
	"log"
	"os"
)

var db *sql.DB

func init() {
	//注意：一般生产环境中不建议使用root用户来进行数据库连接
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3301)/fileserver?charset=utf8")

	//设置同时活跃的连接数
	db.SetMaxOpenConns(1000)

	//测试一下数据库连接是否成功
	err := db.Ping()
	if err != nil {
		//说明连接失败，直接打印错误信息即可。
		log.Printf("Failed to connect to mysql, err：", err.Error())
		//强制进程退出
		os.Exit(1)
	}

	log.Println("---------------------------------数据库连接测试成功---------------------------------")
}

//外部提供一个可以访问的方法
//返回我们创建的数据库的连接对象
func DBConn() *sql.DB {
	return db
}
