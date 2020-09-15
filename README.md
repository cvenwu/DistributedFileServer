# DistributedFileServer
1. 基于Golang实现的一个分布式文件上传服务
2. 重点结合开源存储(Ceph)以及公有云(阿里OSS)支持断点续传以及秒传功能。
3. 微服务化以及容器化部署，从而实现一个分布式文件的存储服务。


## 自己下去要做的一些事情：
1. 了解私有云，公有云，容器编排的概念
2. 

## 预期收获
1. redis/rabbidMq：实现将服务模块中的同步逻辑转换成异步来执行，结合进行业务开发
2. docker/kubernets(一个容器编排工具)
3. 分布式对象存储(Ceph)
4. 阿里云OSS对象存储服务

## 收货哪些干活
1. 将会熟悉文件分块断点上传以及秒传的概念
2. 文件对象从私有云迁移到公有云的经验，从cephq迁移到阿里云的OSS的经验


## 课程安排
1. 2-6章：构建一个基础版的文件上传服务，基本功能可以使用，例如文件上传下载，但是性能架构方面自己还是需要进行升级和优化
2. 7-11章：架构逐步升级，搭建一个完整优化的分布式服务
3. 12章：课程总结

## 自己的建议：
1. 如果到时候看完了想要写到简历上，自己可以回头看第一章里面有讲具体过程，
2. 老师的接口并没有遵循restful接口，自己做完这个项目可以尝试自己改进，使用restful接口来实现，如果有机会可以通过graphQL来实现这个
3. 自己可以通过vue来实现页面，将内容更加丰富
4. 更合理的是在文件上传保存之后，在下载那一端部署一个反向代理，然后将文件作为一个静态资源来处理，例如nginx，下载的时候后端服务会提供一个接口，用于构造下载文件的url，客户端获取url之后，就去下载，下载的时候会经过nginx，
nginx再做一次静态资源访问将文件download下来，一些限流以及权限访问都可以在nginx做，可以减轻golang实现后端的压力。
5. 自己觉得我们定义的文件元信息，没有给json对应序列化加上一个tag，以为golang中属性字母都是大写开头为了其他包可以访问，因此自己觉得最好在对应的字段后面加上一个json的tag



自己在第2章内容上可以做的：

1. 因为上传同一个文件，sha1值都是根据上传文件的内容生成，所以不同用户上传同一文件但是上传时的文件名如何处理？比如A,B都上传同一个文件，但是A上传的时候叫做a.txt而B上传的时候叫做b.txt
2. 另外一个是不同用户上传同一个文件，但是我们如果使用计算sha1并更新，将会不断重新上传，我们其实可以判断sha1是否已经存在，避免不必要的上传，来节省带宽。
3. 自己在用户查询最近上传文件的时候做了判断，因为用户需要传入一个limit参数表示查最近上传的几个，所以如果文件总数小于limit的时候，直接切片切[0:limit]将会越界panic。



## 最初的文件系统（第2章完成之后实现的效果）

架构如图所示：
![Y5wuE2](https://gitee.com/yirufeng/images/raw/master/uPic/Y5wuE2.png)

## 接口列表

![JomueR](https://gitee.com/yirufeng/images/raw/master/uPic/JomueR.png)



获取文件的信息：http://localhost:8080/file/meta?filehash=上传文件的sha1哈希值



## 第2章：步骤


### 文件上传
2-1，2-2做的内容：
1. 项目根目录下新建一个main.go以及handler文件夹同时在handler文件夹下新建一个handler.go文件(专门用于处理上传的接口)
2. 项目根目录下新建一个static文件夹用于存放静态资源文件，同时新建一个view文件夹，用于存储对应的html文件
3. 在handler.go下我们自己实现了一个文件上传函数，同时在上传成功后我们重定向到一个文件上传成功的处理函数中

2-3的内容：
不是文件上传上去就完成了，有时候我们要去查询，如何做呢？
就是把每次上传文件的内容记录下来，例如文件的哈希md5以及sha1作为文件的id，上传时间，文件名，文件大小，文件存储路径，这些内容都要保留下来，便于查询接口返回给用户
上面这些内容我们都封装到了util.go 老师已经写好的，直接用就可以了
1. 新建一个文件夹util同时将老师的util.go放进去
2. 新建文件夹meta，同时新建filemeta.go用来存储文件元信息，
        同时我们在这个里面新建了一个全局变量(采用键值对，键为文件的sha1，值就是对应的文件元信息结构体)用来存储所有上传文件的元信息，
        我们新增了两个方法一个是用来增加或者更新文件的元信息，另外一个是获取文件的元信息
3. 修改我们之前的处理文件上传逻辑，在上传之前新建一个文件元结构体，并且同时更新之前的处理文件上传逻辑。
    3.1 新建一个文件的结构体
    3.2 文件结构体的FileSize字段为io.Copy(newFile, file)返回的第1个结果，也就是已经写入的字节长度
    3.3 之后我们需要获取文件的sha1，但是获取之间，我们需要将newFile也就是我们写入到目的地的文件句柄在刚才写完之后，我们要将其重新移动到最开始的位置来计算sha1，之后我们就可以给文件结构体的sha1进行赋值
    3.4 这时我们可以进行上传文件元结构信息的操作(也就是将文件映射到我们前面定义的map中便于查询)

文件上传演示总结：

![R1QW2Q](https://gitee.com/yirufeng/images/raw/master/uPic/R1QW2Q.png)
真正的生产环境中会保存在redis或mysql中


### 文件元信息查询

2-4 文件元信息查询接口：
原理：通过文件的hash值来查询，

实现的功能：
1. 单个查询（通过指定的hash值来查询）
2. 批量查询（例如，最近上传的文件对应的信息）


1. 首先编写一个获取文件元信息的处理函数，从表单中获取文件的hash值，之后去我们定义的方法中根据对应的hash值查询对应的文件元信息，将文件元信息序列化后返回给用户
2. 将该函数添加到main函数中的映射中
3. 通过命令`sha1sum 刚才上传的文件`直接获取对应文件的sha1值，
4. 启动运行之后，发现通过`http://localhost:8080/file/meta/?filehash=刚才上传文件的hash`可以获取到我们序列化后的json文件元信息


自己百度的资料：
1. `sha1sum命令补充：`sha1sum命令用于生成和校验文件的sha1值。它会逐位对文件的内容进行校验。是文件的内容，与文件名无关，也就是文件内容相同，其sha1值相同。[参考](https://www.linux265.com/course/5071.html)

### 文件下载接口
原理：根据用户提供的filehash值来将文件对应的内容返回给用户。
更合理的是在文件上传保存之后，在下载那一端部署一个反向代理，然后将文件作为一个静态资源来处理，例如nginx，下载的时候后端服务会提供一个接口，用于构造下载文件的url，客户端获取url之后，就去下载，下载的时候会经过nginx，
nginx再做一次静态资源访问将文件download下来，一些限流以及权限访问都可以在nginx做，可以减轻golang实现后端的压力。


思路：服务端通过文件元信息的位置读取文件到内存然后返回给客户端
1. 首先解析用户请求中携带的filehash值，
2. 根据我们自己写的方法从map中查询到文件元信息之后，就可以获取到文件存储的位置
3. 之后打开文件，将其读到内存，然后将数据写入到给用户返回的响应中
4. 其实前3步我们就已经做完了下载，但是为了让浏览器有一个下载文件的效果，我们又在返回的响应头部加入两个字段，使得浏览器有一个下载文件的效果
5. 启动运行之后，发现通过`http://localhost:8080/file/download?filehash=刚才上传文件的hash`便可以通过浏览器下载获取的文件

### 文件Meta更新(重命名)接口

#### 修改文件元信息，以文件重命名为例

需要用户上传两个值，一个是sha1(文件唯一的哈希值)，另外一个是重命名后的文件名

步骤：
1. 解析用户请求的参数列表，用户携带3个参数，1个操作类型，1个文件的唯一标识，1个更新后的文件名
2. 根据用户请求的参数列表获取用户的操作类型，(如果操作类型不是0，那么代表其他操作，我们就返回一个403)
3. 如果用户请求的类型不是post，直接返回一个405
4. 获取用户请求中携带的文件唯一标识，直接修改文件元信息map集合中的文件名为更新后的文件名即可
5. 将修改后的文件元信息序列化返回给用户


### 文件删除接口
主要是实现两个删除操作
1. 删除索引数据，也就是从我们的map中删除
2. 删除文件系统中的数据，也就是物理删除

步骤：
1. 解析用户请求的参数
2. 获取用户请求中的文件filesha1。
3. 根据文件的fielsha1获取文件对应的元信息中的文件存储路径，进行删除
3. 从存储文件元信息的map集合删除文件对应的键值对
4. 返回给用户操作成功

注意：目前所有文件所有元信息都存储在内存中，所以一旦重启所有元信息都没了，在真正的生产环境中我们存放到redis或者mysql数据库，便于增删改查以及安全性。



这个是存放在meta下的filemeta.go文件中，老师视频中没有讲到：

### 有一个视频这里没有，关于文件信息查询，url中使用了query

1. 首先在meta下面新建一个sort.go，按照上传时间排序
2. 之后再fileMeta.go中编写一个获取批量文件元信息列表的方法![48YP7u](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/48YP7u.png)
3. 编写对应的handler处理函数FileQueryHandler, 解析用户请求的参数，获取用户传入的参数limit,表示查最近上传的几个文件，之后调取我们meta中编写的方法，新建一个存储文件元信息的切片，同时根据最近上传时间查询之后添加到切片中并序列化后返回
4. 最后在main.go中映射对应的url以及handler处理函数


### 本章小结
![ZbUE8T](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/ZbUE8T.png)


## v1文档列表

> v1表示已经搭建好的最初文件系统，只实现了文件上传，下载，重命名，删除，以及查询最近上传文件的功能。

接口：

1. 上传文件：到http://localhost:8080/file/upload传入要上传的文件
2. 查询上传文件的元信息(get)：http://localhost:8080/file/meta?filehash=文件的sha1值
3. 查询最近上传的几个文件(get)：http://localhost:8080/file/query/?limit=一个值
4. 文件修改(post)(只是修改了map中的映射文件名，实际上物理存储的文件名没修改)：http://localhost:8080/file/update?op=0&filehash=文件的sha1值&filename=新文件名(如111.png)
5. 文件删除(get)：http://localhost:8080/file/delete?filehash=文件sha1值
6. 文件下载(get)：http://localhost:8080/file/download?filehash=文件sha1值


## 第3章


### 3-1 本章介绍以及mysql的主从复制架构
之前存储的文件元信息都存储在内存中，一旦遇到特殊情况就会丢失，因此我们这里使用mysql(为啥不使用mongodb呢，其实都可以，但是选择最适合自己的即可)

相比于之前的架构图，我们加入了mysql之后的架构图，如图所示：![6kS7BR](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/6kS7BR.png)


单点模式：一旦mysql进程挂掉整个数据库服务就会被强制终止，因此我们需要更稳定的架构
主从模式：一个主节点(master)，多个从节点(slave)组成的数据库集群。主从架构的原理图如图所示：![IBTDmR](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/IBTDmR.png)

图中的集群是一个master节点和一个slave节点，工作原理如下：
1. master的mysql实例所有的增删改查所有操作都会写入到bin log的日志文件中
2. slave节点的mysql会启动一个io线程实时读取master的bin log会将其同步到relay log日志文件中，同时mysql会另外启动一个sql线程来专门读取relay log日志并将所有操作重放一遍。
3. 最终实现的效果就是在master发生的所有数据修改过程都会按序在slave节点重放一遍，这样保持了两个节点的数据一致性，并且如果从节点有多个的时候也会类似这样的流程保持一致性。一旦主节点挂掉，从节点还可以提供读数据的服务，留出宝贵的时间启动主服务或者将从节点切换成主节点。

多主模式：比主从更强大的架构。服务部署在不同地区的机房里面，每个机房都运行数据库服务，并且要求每个机房的数据库都可以读写，并且要数据同步。这时就可以使用多主模式。


### 3-2 实操快速部署一个mysql主从模式

为了更方便的演示在本机上通过docker启动两个mysql容器 `sudo docker ps`查看容器
![MzkRqN](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/MzkRqN.png)


【推荐】自己百度到的资料：如何通过docker搭建一个mysql主从架构的集群(http://www.fall.ink/post/22)

【推荐】推荐自己去查看一下docker的命令[菜鸟教程](https://www.runoob.com/docker/docker-rm-command.html)

如何搭建一个主从架构的集群
1. 登录docker中两个mysql的节点 `mysql -uroot -h127.0.0.1 -P端口号 - p` 之后输入密码，(因为两个docker镜像中的mysql只是映射到本机的端口不同，所以要指定对应的mysql端口号)
2. 假设做成主节点的mysql(假设为master)控制台，要做成从节点的mysql(假设为slave)控制台，
3. 找到将要做成主节点的binlog信息，`show master status;`
4. 回到从节点，配置一下master的信息，也就是告诉从节点将要从哪里读取master的binlog 。最后一个参数(MASTER_LOG_POS)指定从哪里开始复制，为0表示从日志最开始的地方进行复制 `change master to MASTER_HOST='master的ip地址',MASTER_USER='主master的实例的用户名(自己要提前建)',MASTER_PASSWORD='对应实例的用户名的密码',MASTER_LOG_FILE='填上刚才在master节点执行show master status结果的File对应的值',MASTER_LOG_POS=0;`
5. 在从节点上执行`start slave;` 并执行`show slave status\G;` 
6. 从节点上查看 ![piaO8j](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/piaO8j.png)
7. 测试一下主从数据的同步，
    7.1 首先在master上创建一个数据库，![4DevYp](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/4DevYp.png) 此时切换到从节点，我们使用`show databases;`查看从节点是否有对应的数据库
    7.2 切换到master，创建一张数据表，![oItkNy](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/oItkNy.png) 之后去从节点上看一下是否有对应的数据表
    7.3 切换到master，向表中插入数据，![RHsZjc](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/RHsZjc.png)，之后去slave查看一下对应表是否有数据
    7.4 首先去从节点查看一下游标的位置`show slave status\G;`，![D4qZrL](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/D4qZrL.png)去master查看一下binlog中游标的位置，![j300Xq](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/j300Xq.png)发现与我们的从节点的master的位置一样
8. 




### 3-3 文件表的设计以及创建

我们删除记录并不是真正的删除，而是加入一个标记字段，避免物理操作误删除，以及可以通过修改该字段的值来修改，同时可以减少删除造成的数据库文件的空洞和碎片，

0. 创建一个数据库叫做`create database fileserver default character set utf8;`
1. 按照这个sql文件创建数据表：![at1zrj](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/at1zrj.png) ,status列标识文件的状态
    1.1 延伸：一个表中的数据量太大的时候我们需要进行分库分表，分表有水平分表和垂直分表，垂直分表（将一个表分成多个表，例如一个表的有16列，
                将第1个表设计为9列，第2个表设计为9列，因为这两个表都是有通用的一列来进行连接，也就是一条记录一刀切成两半，这两个表通过唯一键进行关联）是把不同字段切分开来 
          水平分表意思就是每个表的结构一样，无非不过就是将其中一个表的数据切开放到多个表中进行放置。
    1.2 这里以水平分表为例，因为filesha1的值后两位都是16进制，所以一共有16*16种可能也就是256个表，如图所示：![5S9at9](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/5S9at9.png)
    1.3 如果一个文件的一个filesha值以ff结尾那么将会被存放到`tbl_ff`这个表中，如果一个文件表的filesha1值的后两位作为表的区分，256张表有一个公共的前缀tbl_。
    1.4 优点：水平分表(横着将表分成两个部分)的方式非常简单，只是根据一个规则 缺点：如果基于原来的分表进行拓展则需要改变这个规则，此时需要移动数据到对应的表中，比较繁琐。
    
    
作者设计的表的结构 
```sql
CREATE TABLE `tbl_file` ( `id` int(11) not null auto_increment, `file_sha1` CHAR(40) not null DEFAULT '' COMMENT 'hash', `file_name`VARCHAR(256) not null DEFAULT '' COMMENT '', `file_size` BIGINT(20) DEFAULT '0' COMMENT '',
`file_addr` VARCHAR(1024) not null default '' comment '', `create_at` datetime default NOW() comment '', `update_at`datetime default NOW() on update CURRENT_TIMESTAMP() COMMENT '', `status` int(11) not null DEFAULT '0' COMMENT '(//)', `ext1` int(11) DEFAULT '0' COMMENT '1', `ext2` text comment '2', PRIMARY KEY (`id`), unique key `idx_file_hash` (`file_sha1`), key `idx_status` (`status`) ) ENGINE=INNODB DEFAULT CHARSET=utf8;
```

### 3-4 持久化元数据到数据库中

步骤
1. 在项目目录下新建一个db目录，同时在db目录下新建一个mysql子文件夹，专门用于创建mysql连接的，新建一个文件conn.go在mysql文件夹下，
    1.1 在该文件下init函数下进行数据库的连接，并且提供一个方法用于返回创建的数据库连接的方法
    1.2 同时在db目录下新建一个file.go文件，新建一个向数据库插入(使用了`insert into`语句)记录的函数（通过预编译操作，避免sql注入），也就是文件上传成功之后将文件信息插入到数据库中。
2. 在db.go定义一个全局变量初始化数据库连接以及设置
3. 在file.go新建一个函数来告诉我们上传文件成功的时候应该要干些什么，


关于insert into语句与insert语句的区别请参考：[insert](https://blog.csdn.net/qq_30715329/article/details/79363761)


### 3-5 从文件表获取元数据


1. 在filemeta.go新建一个函数，在update文件元信息的时候， 直接写入到mysql的表中。
2. 在handler文件中，修改之前的uploadFileMeta将`meta.UpdateFileMeta(fileMeta)`修改为`meta.UpdateFileMetaDB(fileMeta)`
3. 在file.go的文件中，增加一个接口用于 ：查询(修改文件的元信息)接口，同时为了返回的数据我们还新建了一个tbfile结构体。
4. 在filemeta.go文件中，新建一个GetFileMetaDB() 用于从数据库查询文件的元信息并且返回，
5. 修改handler.go中的GetFileMetaHandler() 将`fmeta := meta.GetFileMeta(filehash)`修改为`fmeta, err := meta.GetFileMetaDB(filehash)
                                                                                	if err != nil {
                                                                                		w.WriteHeader(http.StatusInternalServerError)
                                                                                		return
                                                                                	}`

### 本章小结

![nXs1IU](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/nXs1IU.png)

在完成对数据库文件表的插入以及查询操作之后，做一个使用mysql的小结
1. 通过官方提供的sql.DB来管理数据库连接对象
2. 通过sql.Open()方法来创建协程安全的sql.DB对象，不需要频繁调用open以及close方法
3. 优先使用prepared statement 来进行预编译防止sql的注入攻击，比手动拼接的字符串更有效

本章小结：![Yy09Rm](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/Yy09Rm.png)
1. 讲了mysql的特点应用场景，为啥在项目中选择mysql以及本项目使用Mysql的好处
2. 讲了主从架构的工作原理以及从实际操作情况下搭建一个主从架构，同时基于当前文件上传做了一个文件表的创建
3. 通过golang访问mysql进行插入和查询


### 第3章自己实践
1. 目前批量查询还是从内存中查询，无法从数据库中读取数据，需要进一步的优化
2. 自己有一个疑问就是为啥修改文件元信息的时候我们参数从get方法中获取，但是又要判断如果请求的方法不是post将会返回一个服务器内部错误呢？

## 额外阅读和补充
1. 【推荐】自己百度到的资料：如何通过docker搭建一个mysql主从架构的集群(http://www.fall.ink/post/22)
2. 【推荐】推荐自己去查看一下docker的命令[菜鸟教程](https://www.runoob.com/docker/docker-rm-command.html)