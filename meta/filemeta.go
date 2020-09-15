package meta

import "sort"

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/14 8:34 下午
 * @Desc: 文件元信息结构体
 */

//文件元信息结构体
type FileMeta struct {
	//通过sha1作为文件的唯一标志
	//也可以使用md5
	FileSha1 string
	//文件名
	FileName string
	//文件大小
	FileSize int64
	//存在本地的路径
	Location string
	//存储时间戳(时间格式化后的字符串)
	UploadAt string
}

//定义一个对象存储所有上传文件的元信息
//key可以唯一标志文件的，也就是FileSha1
var fileMetas map[string]FileMeta

func init() {
	//做初始化工作
	fileMetas = make(map[string]FileMeta)
}

//新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//获取文件元信息：通过sha1值获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))

	//为了避免用户输入的数量大于我们总文件数量，避免panic
	if len(fMetaArray) >= count {
		return fMetaArray[0:count]
	} else {
		return fMetaArray[0:len(fMetaArray)]
	}
}

//删除文件元信息
//生产环境中我们需要做一些安全的判断，比如delete操作会不会引起线程同步的问题，如果多线程必须保证map安全
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
