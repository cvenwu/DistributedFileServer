package oss

import (
	cfg "fileSystem/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/21 6:09 下午
 * @Desc:
 */

var ossCli *oss.Client

//创建ossClient对象
func Client() *oss.Client {
	//说明之前创建过
	if ossCli != nil {
		return ossCli
	}

	ossCli, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccesskeyID, cfg.OSSAccessKeySecret)
	if err != nil {
		log.Println("oss创建连接失败-------------------->", err)
		return nil
	}

	return ossCli
}


//用于获取bucket对象
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil {
			log.Println("oss创建连接失败-------------------->", err)
			return nil
		}
		return bucket
	}
	return nil
}

//DownloadURL：临时授权下载的url

func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return signedURL
}