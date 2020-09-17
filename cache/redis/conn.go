package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/17 10:08 上午
 * @Desc:
 */

var (
	//一个连接池对象
	pool *redis.Pool
	//redis连接的相关信息
	redisHost = "127.0.0.1:6379"
	redisPass = "testupload"
)

//创建一个redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300, //以秒为单位，如果超过这个时间都没有被使用我们就会直接进行回收
		Dial: func() (redis.Conn, error) {
			//1. 打开连接
			//第1个参数为协议，第2个参数为Host，第三个参数为
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			//2. 访问认证
			if _, err := c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		//用于定时检查连接是否可用，检查redis - server的状况，如果出问题直接在客户端上关闭了redis的连接
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			//每分钟检测一次可用性
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

//初始化redis连接池
func init() {
	pool = newRedisPool()
}

//对外暴露一个方法用来获取redis连接
func RedisPool() *redis.Pool {
	return pool
}
