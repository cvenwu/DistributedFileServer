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
6. 第4章节中，老师专门新建了一个util/resp.go中对我们返回的响应信息以结构体的形式做了封装，使得我们更够更容易的返回json数据作为响应。可以改进的点：但是如果我们可以自己改进一下老师的util/resp.go，遵循restful规范就好了，同时最好能够把用户的登录和注册页面以及用户的home.html能够使用vue就更加好了


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



## 第4章：账号系统与鉴权

互联网大部分服务需要账号鉴权才可以操作。因此我们开发一个账号系统的功能：
1. 支持用户注册以及登录，
2. 持久化用户会话的session或者token，用户登录之后且在session失效或者用户登出之前访问其他的api功能接口，例如上传和下载文件（这就是一个鉴权的过程）
3. 可以进行用户数据资源隔离，每个用户之间的数据互不干扰，很多用户上传同一个文件，但是云端只存储了一份数据，但是都可以访问这个文件，并且其中一个用户的删除操作并不会影响其他用户的访问操作


![GuzDfO](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/GuzDfO.png)用户所有的操作都会首先经过用户关口这个模块，只有经过这个授权之后才可以继续访问用户的数据



本章架构图：![8URaUc](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/8URaUc.png)


### 4-2 
1. 建立用户表 ![FTTxiv](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/FTTxiv.png)
2. 在db文件夹下新建一个user.go，里面新建一个函数（UserSignUp）用于处理用户注册插入到数据库的数据这个流程
3. 在handler文件夹下新建user.go里面新建一个handler（SignupHandler）用于处理用户发送过来的请求并将其进行用户注册的整个流程
4. 在static/view下插入signup.html(老师已经提前写好)用于展示注册页面
5. 之后在main.go文件中映射我们刚才定义的处理用户请求对应的handler函数以及对应的请求路径。


自己在第4章上可以做的：
1. 注册页面中的密码显示是采用明文，![VxLTBN](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/VxLTBN.png)，我们自己可以改进一下
2. 注册成功显示的页面：![6NoH5s](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/6NoH5s.png)


### 4-3
用户的登录
原理：将用户登录输入的密码与自己设置的盐值拼接之后进行加密与我们数据库中的密码进行比对，

登录的后端逻辑分为3步：
1. 验证用户的用户名以及密码是否正确
2. 校验密码通过之后生成一个访问的凭证，是其他API接口访问的标志。
    两种实现：
        1. 基于Token的验证，token生成之后发送给客户端，客户端每次请求的时候都需要token来进行验证
        2. 基于session和cookie的方式来进行验证
    我们在这里采用基于token的验证方式
3. 需要返回一些登录成功的信息（比如用户信息），或者重定向到主页。我们重定向到主页


编码步骤：
1. 在db/user.go文件中编写验证用户名与密码的逻辑操作，如果用户名和密码正确返回true，否则返回false
2. 在我们的handler/user.go中编写用户登录的逻辑处理步骤(按照上面写到的3个步骤)
    2.1 用户名以及密码进行校验
    2.2 通过之后，生成token
        2.2.1 通过在handler/user.go中定义一个函数GenToken用来生成token，需要传入一个用户名，返回生成的token
                生成规则(自定义)：使用md5加密字符串(用户名 拼接 时间戳 拼接 "_tokensalt" ) 拼接 时间戳的后8位。因为md5加密之后生成的是32位，我们想要40位的token字符串
        2.2.2 生成之后，我们需要在db/user.go实现一个写token到数据库的操作。需要自己提前建立token对应的表
                    ```sql
                    CREATE TABLE `tbl_user_token` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                      `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
                      `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
                        PRIMARY KEY (`id`),
                      UNIQUE KEY `idx_username` (`user_name`)
                    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
                    ```
    2.3 登录之后使得用户重定向到我们的home.html页面
        
        
### 4-4 实现登录后用户信息查询
 
注意：我们只有注册和登录不需要token验证，其他接口都需要token验证。
因此我们需要在登录接口成功之后，需要将token带上给客户端或浏览器，然后将其缓存到本地，每次都可以调用

所以我们首先要修改登录接口，登录成功之后，同时返回我们的token，因为返回的东西较多，所以建议使用json作为响应的body

这里我们将转换json的过程封装到了util/resp.go中，已经写好了

1. 修改我们的登录handler，在登录成功之后，我们返回用户一个我们自己写好的json数据，![n1vn7L](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/n1vn7L.png)
2. 编写用户信息查询接口 UserInfoHandler()
    2.1 解析请求参数
    2.2 验证token是否有效，其中我们又编写了单独的一个方法用来验证token是否有效【老师课上没有写这个函数的代码 】
        2.2.1 首先获取我们token的后8位(因为后8位我们就是截取的时间戳的后8位)，
        2.2.2 根据我们自己设置的token有效时间(例如我们自己设置的是1天或者几天)验证是否在有效期内
        2.2.3 如果在有效期内则去数据库中查询并验证与我们的token是否一致
        2.2.4 最后返回一致或不一致即可
    2.3 如果token，我们直接查询用户信息，此时我们在db/user.go文件中新建一个获取用户信息的方法GetUserInfo，
    2.4 将查询到的用户信息使用我们的RespMsg封装作为一个json格式的数据发送给用户
3. 去main.go中添加我们对应的路由以及处理函数

为什么要编写用户信息查询函数：因为我们设置登录之后跳转到home.html页面，此时home.html加载的时候会调用请求路径/user/info查询用户的信息，所以我们要编写一个用户信息查询接口


### 4-5 接口梳理小结

用户注册接口：![y2exeg](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/y2exeg.png)
用户登录接口：![yT9heG](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/yT9heG.png)
用户信息查询接口：![kxNls1](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/kxNls1.png)

**拦截器验证token**：
原因：因为每个API都会校验token（除了登录和注册）。如果每个接口都校验，会造成很多代码的重复，因此我们使用拦截器(服务端接收到用户请求之后，在转发给对应的Handler之前，将请求拦截下来验证用户名以及token)验证token
步骤：![UtuHSe](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/UtuHSe.png)
 
 
### 4-6 访问鉴权接口的实现（拦截器）
1. 在handler下新建一个auth.go，里面编写一个拦截器的方法
2. 在main.go中路由处理函数前面加上拦截器 ![UQAlAj](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/UQAlAj.png)
3. 将UserInfoHandler函数中写的验证token是否有效就可以注释掉来避免重复校验  ![86GMLx](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/86GMLx.png) 
4. 此时我们可以尝试一下登录，发现登录之后可以跳转到主页面，同时我们试着触发异常，将浏览器的cookie清空，当前停留的页面(home.html)刷新页面出现403forbidden


本章小结：
1. mysql用户表的字段设计以及校验token字段的表的设计，具体token验证并没有详细讲解。
2. 实现了用户注册、登录、用户信息查询的接口，同时编写了对应的html界面（老师直接给了代码，没有讲解）
3. 实现了一个简单的验证token的拦截器(针对http请求，如果其他方式的请求，例如rpc需要实现其他方式的拦截器)，使用拦截器的好处：用户在访问非注册和登录的其他接口的时候都需要统一鉴权，避免了大量的代码冗余，使得代码复用性更高，性能提升。

> 其实拦截器不止有验证token的作用，还有其他好处：过滤ip白名单，角色权限控制。

### 第4章自己可以做的
1. 可以自己改进一下老师的util/resp.go，遵循restful规范就好了，最后如果使用graphQL就更加美好了
2. 同时最好能够把用户的登录和注册页面以及用户的home.html能够使用vue就更加好了
3. 自己可以尝试使用session/cookie来进行认证,加深对安全方面的理解


## 第5章 Hash计算与秒传功能的实现


### 5-1 本章内容介绍以及本章架构和技术介绍

> 主要是基于文件的hash值计算与秒传功能的实现

文件的校验值算法（3个主流）：
CRC32生成的是32位的校验值，CRC64生成的是64位的校验值
MD5 是16个字节108位
SHA1 20个字节也就是160位

CRC我们一般叫做校验码，而MD5/SHA1一般是哈希值和散列值，主要体现了这3个算法原理的不同
CRC使用的是多项式的除法
比如一个文件分成两块，MD5会先计算出第一块的md5值，然后基于第一块的md5值以及第2块文件的内容来算出整个文件的md5值，

安全级别：比如文件内容不同，而算出来的hash值一样，这就是不安全的，crc是最弱的，sha1是最高的，md5居中，当然有更安全的方法比如sha256或sha512

计算效率：crc效率最高，

应用场景：客户端服务端文件传输都会计算crc校验值对比，而md5和sha1则用于文件和数据签名，sha1用于签名，如果安全性更高可以使用sha256或sha512


秒传--------------------------------


在云存储中，哪些场景会遇到秒传：
1. 用户上传：当我们上传一个很大的文件的时候秒传也可以发挥作用，因为之前有用户上传过相同文件，当我们上传的时候云服务就会检测出来，并且将上传状态立马置为完成，这样免去传输的过程。
2. 离线下载：有些大文件可以瞬间完成下载。
3. 好友分享：

要实现秒传的关键点：
1. 记录每一个文件的hash值，我们一般会用md5或sha1以及sha256。每上传一个文件到云存储服务，云服务都会计算并且记录下来文件的hash值，用户下次上传的时候如果hash值相同则避免重复上传。如果客户端不可以完成hash计算则无法触发秒传
2. 用户文件的关联，除了与用户无关的唯一文件表，还需要创建一个用户文件表，这样可以基于用户实现逻辑资源的隔离。也就是一个文件的记录在与用户无关的文件表中只会存在一个记录，但是又可以由多个用户来实现分享，因此在用户文件表中，一个文件可能会有多条记录。

架构图：


唯一文件表：之前我们创建过的文件表，一个文件只在表中存储一条记录，以文件的hash值作为唯一主键
用户文件表：存了每一个用户的所有元数据，比如一个用户上传了100个文件，那么这个表就会有100条记录，无论文件重复与否
hash计算可以内嵌在上传server里面，做一个内部逻辑模块的存在，也可以单独的抽离出来做一个独立的微服务 ，对外部提供接口来调用。做成独立微服务的好处：可以减轻server的负载，同时使得整个架构更加灵活以及可扩展。

完整流程：
1. 用户上传文件数据，上传server获取到用户上传文件流，一边将流存储到本地中一边通知hash计算的模块，通知它去计算当前上传文件内容的hash值，等到文件上传完成之后，同时计算完成文件的hash值，之后将上传的文件写入到唯一文件表中，最后一步将文件元信息关联到用户文件表中，下次可以通过用户文件表查询相关信息
2. 如果发生了秒传， 就没有实际的文件传输了（图片中的①），因为上传server比对hash发现相同，就没有图中的②③④步了，直接将文件的元信息关联到用户文件表中。


CRC应用场景：主要应用在①步中用户分段上传文件，每段文件都会计算一个crc校验值两边进行比对。


### 5-2 用户文件表结构
 ![Lv58N9](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/Lv58N9.png)
 
 
在mysql中创建文件表即可


### 5-3 改造上传接口

之前我们已经实现了文件哈希计算的逻辑，但是可以优化，因为上传文件的时候，计算哈希时间比较长，所以可以把它抽离出来变成一个独立的微服务，然后就可以进行异步的处理，

1. 首先在db文件夹下新建一个userfile.go里面新建一个方法，专门用户向用户文件表插入记录
2. 修改上传文件的handler，在上传文件到唯一文件表中的后面我们加上几个逻辑，主要用来使得用户上传成功之后，可以将用户以及对应上传的文件记录同步到用户文件表中。

### 5-4 文件上传列表查询接口的实现
> 前面我们将文件上传成功之后会跳转到home页面，里面有一个文件列表，但是上传成功之后并没有任何显示，因此我们需要编写一个接口用于显示用户成功上传的文件列表。

步骤
1. 新建一个函数在userfile.go中专门用于处理用户上传成功的文件列表查询。
2. 在我们对应的handler中的处理函数需要注释掉我们之前使用的获取最近上传文件的函数，而是使用我们刚才写的用户文件列表查询函数





接下来我们尝试一下文件秒传的实现，在秒传实现之后，我们上传一个上传过的文件，但是会重命名成其他文件上传，所以正常情况下，用户文件表中会多一个记录，而在文件表中只有一个记录

### 5-5 golang实现秒传判断接口
秒传接口的具体逻辑过程：![cclxnZ](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/cclxnZ.png)


如果秒传成功用户不需要上传文件，而秒传失败用户需要重新上传文件进行真正的传输。
最主要的原理就是计算上传文件的hash，如果之前已经在了文件表(唯一文件表)说明我们可以秒传。

1. 在handler.go中编写一个文件秒传的函数![cclxnZ](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/cclxnZ.png)
2. 在main.go中添加对应的url逻辑规则以及handler处理函数 ![Gy39dU](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/Gy39dU.png)
3. 使用postman模拟登录，同时秒传一个文件，但是使用不同的文件名，看一下响应信息，同时看一下tbl_file以及tbl_user_file表的信息


一个衍生的问题：假如有这么一个应用场景，同时有多个用户同时上传同一个文件，这个时候云端逻辑如何处理呢？
参考思路：解决方法不唯一，只有结合实际的应用场景才行。
方法：
1. 允许不同用户同时上传同一个文件
2. 先完成上传的先入库
3. 后上传的只更新用户文件表，并且删除已经上传的文件




### 本章小结

hash算法的应用场景与对比（crc，md5，sha1）
秒传（说白了就是后端做了文件的共享，相同的文件只会上传1次，数据库中只会保存1份，所以实际存储的文件数目与用户看到的文件数目是一对多的关系）的原理和简单实现
关键点：通过hash唯一表示一个文件。


## 第6章 文件的分块上传与断点续传 

### 6-1 分块上传与断点续传

两个概念：
1. 分块上传：文件切成多块，独立传输，上传完成后合并。（文件上传之后会进行文件完整性校验）
2. 断点续传（基于分块传输的机制来实现，云端会将传输后的每一块的信息都缓存好了，下一次重传会查询没有传输的文件块，从最前面的没有传输的文件块开始传送，也就是offset，之后开始重传） 传输暂停或者异常中断之后，可基于原来进度重传，例如当用户点击暂停之后，下次可以接着当前暂停的进度传输而没有必要重新开始传输

几点说明：
1. 小文件不建议分块上传（例如几M的文件）。即使用了分块也只分成了1块，此时还会请求分块以及合并接口，所以会浪费部分资源，
2. 可以并行上传分块，并且可以无序传输，
3. 只要合理设置并行数量（），就可以利用分块上传机制极大的提高传输效率
4. 可以减少传输失败后重试的流量以及时间（因为有断点续传，失败后传输过的文件块不需要重新传输）


具体流程：![QEujyz](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/QEujyz.png)
1. 云端初始化一些上传信息，例如提示客户端文件要切分成几块，每个文件多大
2. 客户端执行上传文件块的任务，（这个过程可以并行执行）
3. 通知客户端上传完成

在这个过程中我们还可以加入两个功能：比如用户上传的过程中进行上传取消(upload abort)，或者用户查询自己的上传进度（也可以查到文件是否需要分块上传等，还可以用进度条显示用户自己的上传进度）(upload query)


服务架构变迁： ![HPax4j](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/HPax4j.png)
几个变动点：
1. 将原来的普通文件上传改为分块上传（但是其实也是支持普通上传的，对于小文件我们会使用普通上传，对于大文件我们采用分块上传）
2. 新增一个redis缓存服务（用于存储每一个文件已经上传的每一块的元信息，比如每一块的序列化，大小，起始位置）
    为什么使用redis：因为只是在用户上传的过程中保留这些上传成功的文件块的信息，信息量不会太大，主要取决于同时有多少个文件在上传。另外这部分数据操作频繁，要求效率高以保证我们云端和客户端之间沟通响应比较及时


### 6-2 文件上传通用接口

基于上面的逻辑架构，我们实现了如下的几个文件上传通用接口：![4qsgIS](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/4qsgIS.png)

通知分块上传完成与取消分块上传这两个接口是一个互斥的关系，


接口：上传初始化 ![1eJCSi](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/1eJCSi.png)
步骤如下：
1. 判断文件是否已经上传过，如果上传过直接触发秒传
2. 生成唯一的上传id（一串固定的随机字符串，可能混淆了时间戳等信息），注意：一个用户前后上传同一个文件，这个id也是不一样的
3. 缓存分块初始化信息（按照约定的规则初始化分块的信息，例如我们要分成多少块以及每块的大小，还有上传的有效时间等，会被缓存到Redis中） 发送给客户端


为了本机的演示，作者已经提前安装好了redis。redis-server默认登录是不需要密码验证的，但是配置之后需要进行验证， 

一般在Golang中如果频繁使用redis，我们最好使用一个redis连接池进行操作。这里使用的是redi-go

步骤：
1. 项目目录下新建一个cache目录，在cache下新建一个redis文件夹，在redis新建conn.go，用来管理redis的连接池、
        1.1 新建一个函数用来初始化redis的连接池
        1.2 我们将我们初始化redis连接池的函数放置在init函数中
        1.3 暴露一个方法用来返回redis连接池中的连接

### 6-3 实现初始化分块上传接口

步骤
1. 在handler下新建一个mpupload.go文件，之后新建一个函数以及为了便于返回我们自己定义了一个结构体用于放置我们存放的文件分块结构体信息
2. 实现分块上传的接口之初始化上传：InitialMultipartUploadHandler
    
    
    
### 6-5 实现上传具体文件分块的handler

1. 在mpupload.go文件中编写UploadPartHandler函数，用于处理文件分块上传

其实在这个函数中，我们可以进一步优化，进行一个分块哈希的校验，每一个客户端都需要客户端上传本地计算好的哈希值，
然后服务端接收到分块之后，再计算一次两者进行比较，如果一致则表示当前内容是没有篡改和丢失的。否则当前上传无效


### 6-5 实现通知上传合并接口
步骤
1. 在mpupload.go中新建一个函数专门用于通知上传合并



### 6-6 分块上传的场景测试和小结
1. 将上面的几个函数映射到我们的url中，
2. 进行分块上传的测试，页面上不好进行测试，我们通过一个测试脚本进行测试，test目录下的test_mpupload.go用来模拟客户端进行分块上传流程
3. 启动测试脚本之后发现上传之后确实分块存放，此时我们进一步的测试验证文件的完整性：
    3.1 将分块存放的文件进行合并
    3.2 将合并之后的文件计算哈希值与客户端上传文件之前计算文件的哈希值进行校验
    3.3 流程如图所以：![sOFp9L](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/sOFp9L.png)


进行小结：
![Aadz1i](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/Aadz1i.png)

1. 我们实现了3个接口的主要逻辑：初始化分块信息，上传分块，通知上传分块完成 
2. 还剩下两个接口：取消上传分块，查看分块上传的整体状态。本章不讲解了，因为比较简单
    2.1 取消上传：![O0Hd2f](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/O0Hd2f.png)
        2.1.1 云端接收到取消请求，首先检查upload中的请求是否有效，如果有效就删除已经上传过的分块文件
        2.1.2 删除redis中的缓存状态：根据用户名以及uploadid，找到记录并删除掉
        2.1.3 更新Mysql表的状态，
    2.2 查看分块上传的整体状态：![5I11WV](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/5I11WV.png)
        核心的就是：通过用户名以及upload查到所有还没有上传的数据，返回给用户
        例如当一个用户上传的文件被划分成100块，已经上传了50块，还有50块没有上传，此时用户看到的进度就是50%，并且将所有没有上传的块返回给客户端
    
本章小结：
1. 讲解了分块上传以及断点续传的概念以及原理
2. 分块上传的流程做了讲解（分块上传的初始化，上传分块，通知分块上传完成，取消上传分块，查看分块上传进度）
3. 抽取了几个重要的接口进行代码实现以及具体的测试， 

断点续传实现：分块上传中断一下，客户端获取上传进度，来得到自己还需要上传哪些分块（得到分块序号），得到消息后继续上传从而完成断点续传。



## 第7章

### 7-1 
1. ceph是什么：一种分布式存储系统，也是redhat旗下的开源存储产品。
2. ceph主要用于解决什么问题：为了更好的解决数据分布式存储，相对于其他存储，能够更加充分利用存储节点的计算能力，
在存储数据的时候能够计算得到某一个数据存储的位置，从而尽量将数据分布均匀。因为自带hash算法，使得不会出现单点故障，理论上无限拓展节点和扩容。
3. ceph的历史和现状：openstack私有云后端存储的标配。公有云用很少的代码可以实现存储，
4. ceph的特点：
    4.1 部署简单，可以直接用docker快速搭建，
    4.2 可靠性高，多副本隔离存储，数据强一致性。
    4.3 性能高：由文件系统决定
    4.4 分布式：可扩展性强
    4.5 客户端支持多语言接入。
    4.6 开源
5. ceph的体系架构：（自己具体没有做笔记）![3caqGv](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/3caqGv.png)

接入ceph对象存储

原有的上传都是将文件存储在上传节点的本地的，为了提高性能以及可靠性，一般需要使用分布式存储系统，使用ceph搭建私有云是一个不错的选择。

### 7-2 
ceph 基础组件
![BgR3MR](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/BgR3MR.png)
1. OSD（object storage , 用户数据以及对象存储的守护进程）：用于集群中所有数据与对象的存储：存储(所有用户存储数据到物理盘都会经过OSD)/复制(OSD得到一个数据的副本然后存储到另外一个地址，做成多副本)/平衡（集群规模变动，数据复制与迁移）/恢复数据（某个节点的硬盘坏了，将该节点的数据恢复到其他盘）等等，其他操作如，发心跳到监控节点，维持集群的健康
2. Monitor：集群的监视器，负责维护集群的健康状态以及元数据，比如集群中所有节点的属性以及关系，OSD会获取映射表与数据，计算出对象最终的存储地址
3. MDS(meta data server)：保存文件系统服务的元数据
4. GW：提供兼容的gateway服务


作者在本机也通过docker安装了ceph集群，![025dlH](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/025dlH.png)
查看其中一个monitor节点的健康状态，![c1bZIb](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/c1bZIb.png)

因为ceph兼容aws的s3接口
![i6dc5m](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/i6dc5m.png)
对于一个对象来说，数据就是就是文件中的数据，而元数据就是文件的描述信息 

### 7-3 

加入ceph之后的架构变化：![mIFnh3](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/mIFnh3.png)

变化主要在于：存储这一块文件上传完之后，存储到服务器本地之后， 会将本地存储的文件转移到ceph集群中。之后用户下载文件的时候直接从ceph读取并下载。

在本地存储没有转移到ceph之前，用户还是需要从本地访问文件并下载。
而这个过程是同步的，也就是说用户上传文件首先会保存到本地，之后本地就会将文件转移到ceph，
这个过程就是同步的，有个不好的地方就是这个过程需要写两次，会增加用户上传时间，所以后面章节会将这个过程转换成异步的逻辑，异步明显的好处就是用户并不会增加等待上传文件的时间，避免用户体验差
不过异步【第9章节讲解】带来任务的复杂性，会将任务写入到队列中。

EndPoint：对外暴露的微服务，存储服务的入口，web服务入口点的url


文件写入到ceph之后，数据库表也会有一定修改，会将相关写入的存储地址写入到mysql中


操作：
1. 在项目目录下新建store文件夹，在store新建ceph文件夹，新建ceph_conn.go用于客户端向ceph进行连接
            1.1 在该文件下编写一个函数用来进行ceph的连接，连接之前我们需要进行一些配置，我们使用CMD进行配置：![GUVIZX](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/GUVIZX.png)
                执行命令之后的结果：![DXP1nO](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/DXP1nO.png)
                之后我们就可以根据我们配置的这两个key进行上传以及下载
            1.2 在该文件下编写一个函数用来获取指定的bucket对象
            1.3 在项目目录下的test文件夹编写一个测试文件test_ceph.go ![PqKOq5](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/PqKOq5.png) ![VWcm2X](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/VWcm2X.png)
                        
### 7-4 实现ceph的文件上传下载

步骤：
1. 找到我们之前写的UploadHandler的函数，插入一段代码，在用户文件内容写入到本地存储之后，并且在更新文件表之前。在这个范围内将这个文件写入到ceph中


自己可以做的：
1. 老师配置好ceph之后，上传文件，传到了ceph中，![KhK6ev](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/KhK6ev.png)，老师建议学生可以在后面加上一个下载按钮并同时开发下载功能进行下载
       1.1 ![wE9sw7](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/wE9sw7.png)
       
       
### 本章小结（7-4视频后半部分讲解了）
![VaRMw8](https://cdn.jsdelivr.net/gh/sivanWu0222/ImageHosting@master/uPic/VaRMw8.png)

老师建议多实践来实现功能模块加深对ceph的理解程度。对其他开源云存储也会有帮助。



## 额外阅读和补充
1. 【推荐】自己百度到的资料：如何通过docker搭建一个mysql主从架构的集群(http://www.fall.ink/post/22)
2. 【推荐】推荐自己去查看一下docker的命令[菜鸟教程](https://www.runoob.com/docker/docker-rm-command.html) 