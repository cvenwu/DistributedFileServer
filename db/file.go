package mydb

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/15 4:13 下午
 * @Desc:
 */

import (
	"database/sql"
	mydb "fileSystem/db/mysql"
	"log"
)

//将上传文件的元信息同步到数据库中
//返回一个bool值：如果操作成功将会返回一个true否则返回false
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	//使用预编译的语句来进行数据库的操作
	//好处：可以防止sql的注入攻击，防止外部用户的恶意sql拼接

	//statis在insert的时候默认为1
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_file(file_sha1, file_name, file_size, file_addr, status) values(?, ?, ?, ?, 1)")
	log.Println(stmt, err)

	if err != nil {
		log.Println("Failed to prepare statement，err ：" + err.Error())
		return false
	}

	//如果成功预编译，我们需要通过defer来进行资源的预关闭
	defer stmt.Close()

	//通过exec方法真正执行对应的sql语句
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println("执行了对应的sql语句------------------------------>")
	//最后一步要判断一下数据是否重复插入，避免数据重复插入的情况
	if ret, err := ret.RowsAffected(); nil == err {
		//说明没有产生新增记录，但是sql表已经生效了
		if ret <= 0 {
			log.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}

	//说明插入失败
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//查询对应的元信息
//1. 首先定义一个文件元信息结构体来表示该表的所有字段
//2. 从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_addr, file_name, file_size from tbl_file where file_sha1 = ? and status = 1 limit 1")

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	//实现资源的关闭
	defer stmt.Close()

	tfile := TableFile{}
	stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &tfile, nil

}
