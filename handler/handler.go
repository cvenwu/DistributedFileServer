package handler

import (
	"encoding/json"
	dblayer "fileSystem/db"
	"fileSystem/meta"
	"fileSystem/store/ceph"
	"fileSystem/util"
	"gopkg.in/amz.v1/s3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/14 7:15 下午
 * @Desc: 专门用于处理上传的接口
 */

//实现一个用于上传文件的接口
//处理文件上传
//第1个参数：用于向用户返回数据的ResponseWriter对象
//第2个参数：用于接收用户请求的request对象指针
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//判断用户请求的方法是什么？
	//如果是get就返回上传文件的html页面
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("/Users/yirufeng/实习/项目/DistributedFileServer/static/view/index.html")
		if err != nil {
			log.Println(err)
			io.WriteString(w, "Internal Server Error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" { //如果是post则要接收用户要上传的文件流以及存储到本地目录
		//接收文件流以及存储到本地目录,
		//因为客户端通过form表单提交文件
		//返回3个参数：文件句柄，文件头，错误信息
		file, head, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to get data, err：%s\n", err.Error())
			return
		}
		//同样在退出之前记得关闭文件资源
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建一个本地的文件流来接收
		//接收一个参数表示创建文件的路径以及名称
		//./表示当前目录，也就是项目根目录
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			log.Printf("Failed to create file，err:%s\n", err.Error())
		}
		//同样在退出之前记得关闭文件资源
		defer newFile.Close()

		//第3步：将内存中的文件流拷贝到buffer中
		//第1个参数表示已经写入的字节长度，第2个参数为错误信息
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			log.Printf("Failed to save data into file，err:%s\n", err.Error())
			return
		}

		//计算sha1之前需要将newFile句柄移动到最前面0的位置
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		log.Println("上传文件的sha1值为-------------------->", fileMeta.FileSha1)


		//7-4视频进行代码的改动：同时将文件写入ceph进行存储
		newFile.Seek(0, 0)
		//读取全部文件内容
		data, _ := ioutil.ReadAll(newFile)
		//获取指定的桶(我们自己提前创建好桶了)
		bucket := ceph.GetCephBucket("userfile")
		//写入到ceph的路径，加上文件的sha1保证唯一性
		cephPath := "/ceph/" + fileMeta.FileSha1
		//第1个参数为存放的路径，第2个参数是存放什么数据，第3个参数是数据存放的类型(我们使用octet-stream表明是二进制的数据)，第4个是权限，生产环境中需要设置成私有来让用户进行认证
		_ = bucket.Put(cephPath, data, "octet-stream", s3.PublicRead)
		//将文件元信息存储的路径修改为ceph存储的路径
		fileMeta.Location = cephPath


		//meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetaDB(fileMeta)

		//同时将该记录更新到文件用户表中
		//更新用户文件表，也就是将用户上传的文件同步到用户文件表中。
		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			log.Println("------------------------上传文件成功---------------------------")
			//可以向用户返回一个成功的信息或页面
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		} else {
			log.Println("------------------------上传文件失败---------------------------")
			w.Write([]byte("Upload file failed...."))
		}
	}
}

//表示成功上传文件的信息：
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished...")
}

//获取文件的元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {

	//解析客户端请求的参数
	r.ParseForm()

	//获取参数

	//假设客户端上传的文件参数为filehash
	//r.Form["filehash"]返回的是一个数组，默认是取第1个
	filehash := r.Form["filehash"][0]
	//通过客户上传的filehash来获取对应文件的信息
	//fmeta := meta.GetFileMeta(filehash)
	fmeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//自己在这里加的代码，
	//说明没有找到对应的文件
	if fmeta.FileName == "" {
		log.Println("Resource Not Found")
		w.Write([]byte("Resource Not Found"))
		w.WriteHeader(http.StatusOK)
		return
	}

	//将结构对象序列化json返回给客户端
	//返回两个参数：
	//第1个：转换后的byte数组类型的数据
	//第2个：相关的error信息
	data, err := json.Marshal(fmeta)
	if err != nil {
		log.Println("序列化json格式数据失败------------------------")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	//解析
	r.ParseForm()
	//解析之后获取filehash值
	fsha1 := r.Form.Get("filehash")
	//获取文件对应的元信息
	fm := meta.GetFileMeta(fsha1)
	//服务端通过文件元信息的位置读取文件到内存然后返回给客户端
	f, err := os.Open(fm.Location)
	//有错误
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//如果可以访问到文件，记得退出函数的时候要关闭资源
	defer f.Close()

	//因为我们现在是小文件，所以可以使用这个方法
	//当文件很大的时候，我们需要使用流，就是每次读一小部分数据返回给客户，然后刷新缓存继续读取文件末尾为止
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//理论上工作是已经做完了的，但是为了让浏览器做一个演示，我们需要将一个http的响应头，让浏览器识别出来就可以当成一个文件进行下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")

	//到这里直接将我们内存中的数据返回
	w.Write(data)
}

//自己参照老师最后的项目代码复原了一下
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	//fileMetas := meta.GetLastFileMetas(limitCnt)

	username := r.Form.Get("username")
	fileMetas, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		log.Println("用户上传文件列表查询失败-------------------")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMetas)
	if err != nil {
		log.Println("批量查询-------------->结果序列化失败------------>", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//修改文件元信息(这里只涉及到了文件重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//解析客户端请求的参数列表,客户会携带3个参数
	//1.第1个参数表示操作类型，0表示重命名，1代表其他操作
	//2.文件的唯一标识,sha1
	//3.更新后的文件名
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	//如果操作类型不是0，直接返回一个403错误
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//如果客户端不是post请求方法
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//获取当前文件的元信息
	curFileMeta := meta.GetFileMeta(fileSha1)
	//修改文件名
	curFileMeta.FileName = newFileName
	//重新保存到map中
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		log.Println("文件重命名---------序列化失败-----------》", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件，需要用户传入一个filesha1
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//获取用户需要删除的文件的filesha1
	filesha1 := r.Form.Get("filehash")

	//获取文件元信息
	fMeta := meta.GetFileMeta(filesha1)
	//删除文件，这里可能会失败，但是我们先忽略，只要在map中我们删除了，就当做文件已经删除了
	os.Remove(fMeta.Location)

	//从文件元信息map中删除
	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}

//客户尝试秒传
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {

	//1. 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize := r.Form.Get("filesize")

	//2. 从文件表中查询相同hash的文件的记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		log.Println("查询文件失败----------------->")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//3. 如果查不到记录说明秒传失败
	if fileMeta == nil {
		log.Println("秒传失败------------------------------")
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	//4. 上传过则将文件信息同步到用户文件表中，说明秒传成功
	fz, _ := strconv.Atoi(filesize)
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fz))
	if suc {
		log.Println("秒传成功------------------------------")
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		log.Println("秒传失败------------------------------")
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}
}
