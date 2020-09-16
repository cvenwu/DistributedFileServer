package db

import (
	mydb "fileSystem/db/mysql"
	"log"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/16 9:31 上午
 * @Desc:
 */

//通过用户名以及密码完成user表的注册操作
//用户注册
//插入数据成功返回true，否则返回false
func UserSignUp(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) values(?, ?)")
	if err != nil {
		log.Println("Failed to insert，err：", err)
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		log.Println("Failed to insert，err：", err)
		return false
	}

	//进一步校验是否插入：重复注册也算失败
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

//判断用户名和密码是否正确
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where username = ? limit 1")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err.Error())
		return false
	} else if rows == nil { //判断返回的记录是否为空
		log.Println("username not found：", username)
		return false
	}

	//将我们查询到的rows转换一下格式，转换为元素为map类型的数组
	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false
}

//刷新用户登录的token
func UpdateToken(username string, token string) bool {
	//这里之所以用到replace是因为token可以重复生成的并且可以插入的。而user_name是一个unique的key，所以有旧的我们通过最新的值进行覆盖
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`, `user_token`) values(?, ?)")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true

}

//定义一个用户结构体，与Mysql中user表的结构一一对应
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("select user_name, signup_at from tbl_user where username = ? limit1")
	if err != nil {
		log.Println(err.Error())
		return user, err
	}

	//即使关闭资源
	defer stmt.Close()

	//执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}
