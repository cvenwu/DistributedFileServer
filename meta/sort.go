package meta

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/15 10:18 上午
 * @Desc: 按照文件上传时间进行排序
 */

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByUploadTime []FileMeta

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}
