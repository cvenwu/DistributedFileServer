package db

import (
	mydb "fileSystem/db/mysql"
	"log"
	"time"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/16 6:15 下午
 * @Desc: 用户文件表相关的db操作
 */

//定义一个结构体与用户文件表的基本结构是相同的
//用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

//实现更新用户文件表的方法
//用来表示插入成功与否
//更新文件用户表
func OnUserFileUploadFinished(username string, filehash string, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare("Insert into tbl_user_file(`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) values (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return false
	}

	//记得关闭资源
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		log.Println(err)
		return false
	}
	//说明插入成功
	return true
}

//批量获取用户文件信息，也就是查询一个用户所有上传的文件，返回该用户上传的所有文件组成的列表
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where username = ? limit ?")
	if err != nil {
		log.Println("stmt初始化是啊比")
		return nil, err
	}

	defer stmt.Close()

	//查询用户文件的信息，会返回查询成功之后所有符合记录的列表
	ret, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	//用于存放结果
	var userFiles []UserFile

	for ret.Next() {
		uFile := UserFile{}
		err := ret.Scan(&uFile.FileHash, &uFile.FileName, &uFile.FileSize, &uFile.UploadAt, &uFile.LastUpdated)
		//如果scan失败我们直接跳出循环
		if err != nil {
			log.Println(err.Error())
			break
		}
		userFiles = append(userFiles, uFile)
	}

	return userFiles, nil
}
