package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/21 4:33 下午
 * @Desc:
 */

var cephConn *s3.S3

//暴露接口：外部方法可以通过该方法来获取一个S3的链接
//获取ceph连接
func GetCephConnection() *s3.S3{

	//0. 加入一个判断逻辑，避免重复初始化
	if cephConn != nil {
		return cephConn
	}


	//1. 初始化ceph的信息
	//包括入口的host, bucket等
	auth := aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}
	//配置region
	curRegion := aws.Region{
		Name: "default",
		EC2Endpoint: "http://127.0.0.1:9080", //根据docker中的9080端口映射出的一个，映射到docker容器中的80端口
		S3Endpoint: "http://127.0.0.1:9080",
		S3BucketEndpoint: "",
		S3LocationConstraint: false,
		S3LowercaseBucket: false,
		Sign: aws.SignV2,
	}

	//2. 创建S3类型的连接
	return s3.New(auth, curRegion)
}

//获取指定的bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	return GetCephConnection().Bucket(bucket)
}



