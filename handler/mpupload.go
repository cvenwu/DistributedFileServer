package handler

import (
	rPool "fileSystem/cache/redis"
	"fileSystem/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/17 10:24 上午
 * @Desc:
 */

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string //用于唯一标志每一次上传的动作
	ChunkSize  int    //（除了最后一个分块的）每一个分块的大小
	ChunkCount int    //文件分块的数量
}

//用于分块上传之前的初始化
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求信息
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "param invalid", nil).JSONBytes())
	}

	//2. 获得redis连接
	rConn := rPool.RedisPool().Get()

	//不要忘记关闭连接
	defer rConn.Close()

	//3. 生成分块上传的初始化信息，将初始化信息封装成一个struct
	upInfo := MultipartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		//生成规则：当前用户名+当前时间戳
		UploadId:  username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: 5 * 1024 * 1024, //每一块的最大大小是5M
		//文件大小除以每个块的大小，之后向上取整获得分块的个数
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	//4. 将生成的初始化信息写入到redis缓存中
	//通过hset命令将这些初始化信息存进去
	//我们约定key是前缀(MP_)+当前上传的id
	//其实这里可以通过hmset命令一次性将所有内容写进去
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadId, "filesize", upInfo.FileSize)

	//5. 将响应初始化数据返回给客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

//上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求参数
	r.ParseForm()
	//username := r.Form.Get("username")
	//说明当前分块是属于哪个uploadid的请求
	uploadID := r.Form.Get("uploadid")
	//当前文件的分块，如果为3也就是说当前是文件分块中的第3 块
	chunkIndex := r.Form.Get("index")

	//2. 获得redis连接池的连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//3. 根据当前用户以及上传的分块序号获得文件句柄，用于存储当前文件块的内容
	//关于文件上传的存储目录，我们这里约定将文件上传到data目录下，以uploadId作为一个子文件夹，将不同的upload请求分隔开来
	//每个文件对应的是一个目录
	//返回的是一个文件的句柄以及错误信息


	//如果通过os.Create()创建一个之前文件目录不存在的时候，将会报错，因此我们需要创建目录
	fpath := "/data/" + uploadID + "/" + chunkIndex
	//path.Dir()将目录提取出来，
	//第2个参数为权限，设置该目录权限为0744，除了当前用户拥有7的权限其他用户只有4的权限
	os.MkdirAll(path.Dir(fpath), 0744)
	f, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed.", nil).JSONBytes())
		return
	}

	//不要忘记关闭当前的文件资源
	defer f.Close()

	//将得到的文件分块内容存储起来，通过一个for循环每次只读1M或10M
	buf := make([]byte, 1024*1024) //每次读取1M
	for {
		//通过request中的body的read()方法将读取到的内容写到buf中
		n, err := r.Body.Read(buf)
		//通过write将buf中的内容写到文件中
		f.Write(buf[:n])
		//如果读到文件最后，此时err不为nil，我们需要退出
		if err != nil {
			break
		}
	}

	//4. 完成上传当前文件分块之后，更新redis中的缓存数据，表明当前块已经上传完成
	//第3个参数为当前文件分块的序号，
	//第4个参数为1
	//好处：通过查询上传进度的时候，判断到所有的idx上传完成则表示当前文件分块上传成功。
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	//5. 返回处理结果给客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")
	//2. 获得redis连接池中的连接
	conn := rPool.RedisPool().Get()
	defer conn.Close()

	//3. 通过uploadId查询并判断所有分块是否上传完成
	//因为我们之前每成功上传一块我们都会在redis中有一条记录
	//所以这里如果记录的个数等于我们的块数，说明上传完成
	//将查询出的结果转换成interface的array。通过redis.Values()转换
	data, err := redis.Values(conn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0 //应该上传的总数量
	chunkCount := 0 //实际上传的数量
	//这里为什么跳转2，是因为通过hgetall查出来的所有结果，key和value都是在同一个array里面
	//所以每次循环里面我们要同时越过key和value
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" { //如果找到hash中有chkidx_开头并且v为1
			chunkCount++
		}
	}
	//判断两个的值是否一致
	if chunkCount != totalCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	//4. 所有分块上传完成之后，我们要进行合并，暂时跳过

	//5. 更新唯一文件表以及用户文件表，与我们之前普通上传的逻辑是一样的
	fsize, _ := strconv.Atoi(filesize)
	//最后一个参数为fileaddr因为我们合并分块还没有做，所以这里没法确定，先没写
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	//6. 响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
