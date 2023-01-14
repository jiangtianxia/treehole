## 功能模块

### 登录业务接口

- [x] 使用md5对密码进行加密存储数据库
- [x] 使用jwt生成token进行鉴权登录
- [x] 雪花算法生成唯一ID
- [x] 邮箱注册、密码登录、找回密码

### 用户业务接口

- [x] 修改用户信息

### 公共业务接口

- [x] 文件本地上传
- [ ] 文件分片上传（目前只能上传图片，要上传视频等大文件可做）
- [ ]   文件上传至云服务（OSS）

### 帖子业务接口（核心）

- [x] 创建帖子
- [x] 查询帖子
- [x] 获取发布帖子信息
- [x] 删除帖子
- [x] 获取帖子详细信息
- [x] 热度或发布时间排行榜
- [x] 修改帖子
- [x] 点赞帖子
- [ ] 社区分类

### 评论业务接口

- [x] 发送评论
- [x] 获取文章评论
- [x] 获取评论记录
- [x] 删除评论

### 聊天室功能

- [x] websocket通讯
- [x] 心跳检测
- [x] 获取当前聊天室在线人数
- [ ] 获取聊天记录

### 高并发（核心）

- [x] 令牌桶限流策略
- [x] 跨域问题（前端）
- [ ] 跨域问题（后端）
- [x] 热度排名算法
- [x] 缓存问题（待优化）
- [x] 缓存数据与数据库数据不一致问题(待优化)

### 部署上线

- [ ] 使用docker-compose部署Nginx + GoLand + Redis + Mysql





### 笔记

#### 限流

防止同一时间大量请求访问服务器从而导致服务器宕机。

解决：

计数器、滑动窗口、漏桶、令牌桶（常用）、Redis + Lua分布式限流

#### 热度算法

![热度算法](E:\心灵树洞\热度算法.png)

前时间->将数据存入数据库->将数据库存入redis(事务操作)。

> 热度时间为7天，redis存放七天后，点赞sorted set过期，点赞只加点赞数，不加热度，同时判断自己该帖子热度是否为前2000名，如果不是，则回存数据到数据库。

**为什么redis数据要存回数据库，不能永久存放？**

redis存储在内存当中，而内存资源紧缺，存过多会导致没内存，数据库则存放到磁盘当中，因此需要存回数据库，这也是为什么redis读取速度比mysql快的原因，redis读取数据指向读取内存即可获取数据，而数据库读取数据需要先读取磁盘，再从磁盘中将数据写入内存获取数据。

![硬盘数据库读取流程](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a6673d3836a64d7eb61373521fe94d8d~tplv-k3u1fbpfcp-zoom-in-crop-mark:4536:0:0:0.awebp)

![内存数据库读取流程](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2369bd758cfe470bace91d560a2fcc9f~tplv-k3u1fbpfcp-zoom-in-crop-mark:4536:0:0:0.awebp)

图片转载：https://juejin.cn/post/7049148028875178020

**redis命令的使用**：

**创建帖子：**

1、创建有序集合

```redis
zadd treehole:note:time timescore(float64) noteid
zadd treehole:note:score score(float64) noteid
zadd treehole:note:voted:noteid votedscore(-1或1)(float64) userid   
--treehole:note:voted:noteid设置有效期为7天，当用户点赞时判断是否过了七天，没过继续点赞，过了，则获取帖子排名，判断是否在前2000名以内，不在则删除对应的哈希集合数据和分数数据
```

2、创建哈希集合

```redis
Hmset noteInfoKey noteInfo
```

**查询帖子（待优化）：**

到数据库模糊查询数据，直接将数据返回，不对redis进行操作。

**获取发布帖子信息（待优化）：**

到数据库操作查询数据，直接将数据返回，不对redis进行操作。

**删除帖子：**

1、删除哈希集合的key和有序集合中的成员投票的key

```go
del noteInfoKey
del treehole:voted:noteid
```

2、删除时间、和分数有序集合中的成员

```go
zrem treehole:time noteid
zrem treehole:score noteid
```

3、删除数据库及本地图片信息

**获取帖子详细信息**：

流程：

先到redis哈希集合中查询是否存在该帖子，如果存在，则使用author_identity，去查询作者信息。

不存在，则到数据库中查询数据，同时将数据放入redis哈希集合，设置过期时间为半天，作为缓存。

1、获取哈希集合中的数据，不存在长度为0

```go
hgetall noteid 
```

2、创建哈希集合，并设置过期时间

```go
Hmset noteInfoKey noteInfo
expire noteInfoKey time
```

3、访问问量加1

```redis
hincrby treehole:noteinfo:noteid visit 1
```

4、判断当前用户是否有点赞，如果数据库没有数据，则创建数据。

```redis
zscore treehole:note:voted:noteid userid
```

**按照热度或时间获取帖子信息：**

1、获取前100个分数最高的帖子id（从高到低递减）

```redis
ZRevRange treehole:note:score 0 100
```

2、根据noteid，获取哈希集合中的用户信息

```go
Hgetall noteid  
```

3、获取有序集合的成员数

```redis
zcard key
```

**修改帖子**：

1、删除缓存数据

2、修改数据库数据

3、休眠几秒之后再次删除缓存数据

**点赞帖子**：

判断缓存中是否存在该缓存投票key

1、没有。判断分数缓存中是否存在该文章key，存在则判断排名是否在2000名以内，不在则删除key，在则修改投票数据。

2、有。修改投票数据和分数数据。

判断key是否存在

```redis
exisit key
```

读取成员的排名

```redis
ZREVRANK key member
```

#### 评论缓存

![评论缓存](E:\心灵树洞\评论缓存.png)

**发送评论：**

1、将数据插入数据库。

2、将数据放入redis，同时设置过期时间为7天。

在list前插数据

```redis
lpush key vlaue
```

**获取文章评论：**

1、判断列表中是否存在该key

2、不存在，则查找数据库数据，并将其加入缓存，设置过期时间为1天。

3、存在，则查询数据返回即可。

```redis
exisit key    -- 不存在，返回0
llen key
lrange key stat stop
hgetall key
```

**获取评论记录：**

直接操作数据库即可，无需操作缓存。

**删除评论：**

1、先删除缓存数据

```redis
// 删除哈希集合数据
del key

// 删除list中的该用户的数据
lrem key count VALUE
```

2、删除数据库数据

3、休眠几秒后再次删除缓存数据

#### 聊天室

![聊天室](E:\心灵树洞\聊天室.png)

**聊天流程：**

1、建立websocket连接，将用户信息加入到哈希缓存中，用户信息过期时间为2天，开启监听器。

2、Conn.ReadMessage()读取到前端发送的数据。

3、判断是否为心跳检测类型，如果是，则判断是否还有心跳，即判断连接是否已经断开，断开直接返回，不断开则回应。

4、不是则将信息发送，同时将消息缓存至有序集合和哈希集合当中，同时有序集合的长度大小为500条，哈希集合有效期为7天。

**缓存消息和缓存用户信息：**

```redis
// 插入在线集合缓存之前，先移除该用户信息，再插入
lrem key count VALUE

// 保留后500条聊天记录
Ltrim key 0 500
```

**获取在线人数：**

1、获取用户信息的在线list信息

2、判断缓存中用户信息是否已经过期，不过期则查询

3、过期，则到数据库提取用户信息，同时将期存入缓存，过期时间为1天

**清除超时连接：**

使用定时任务清理超时连接，即清理没有心跳的连接。

**获取聊天记录：**

1、判断查询的查询的大小，判断查询大小是否在500以内。

2、500以内获取缓存数据，返回。

3、500以外查询数据库，返回。

#### 缓存击穿 

同一时间大量数据访问一个不存在的帖子，请求会先查询redis，不存在之后查询数据库，由于数据库同一时间请求过多，导致数据库宕机问题。

解决：

1、增加接口校验，过滤掉不合理的请求。（感觉不是很好）

2、redis主动缓存null结果，并设置较短的过期时间。

#### 缓存雪崩

在高峰期时，redis机器出现故障（宕机），导致完全访问不了，大量请求访问数据库导致数据库宕机。

解决：

1、redis集群，部署多个redis。（不考虑）

2、本地缓存。（前端）

3、限流&降级。（令牌桶限流）

4、设置热点数据永不过期。（热点排名）

5、缓存数据的过期时间设置随机，防止同一时间大量数据过期现象发生。（实际很难实现，还会引发新问题）

#### 缓存击穿

缓存中大量的热点key同时过期，正好有大量请求访问，因为缓存中没有，所以全部打响了DB数据库，数据库抗不住挂了，宕机了。

解决：

1、设置热点数据永不过期。（热点排名）

2、第一时间去数据库获取数据填充到redis中，并且在请求数据库时需要加锁，避免所有的请求都直接访问数据库，一旦有一个请求正确查询到数据库时，将结果存入缓存中，其他的线程获取锁失败后，让其睡眠一会，之后重新尝试查询缓存，获取成功，直接返回结果；获取失败，则再次尝试获取锁，重复上述流程。

#### 缓存数据与数据库数据不一致问题

先删缓存后更新数据库+延迟双删（一般为写一个数据的时间，3~5秒）

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4b1816efac5d44b4a088a8e0ec63eba9~tplv-k3u1fbpfcp-zoom-in-crop-mark:4536:0:0:0.awebp?)

图片转载：https://juejin.cn/post/7156237969202872350#heading-31



### 参考资料

- 基于雪花算法生成用户ID

https://www.yuque.com/docs/share/e50bbca1-e019-45e2-b77b-a9ba01fbede3?# 

- 傻傻分不清之 Cookie、Session、Token、JWT

https://juejin.cn/post/6844904034181070861

- Redis为什么这么快？

https://juejin.cn/post/7049148028875178020

- 新来个技术总监，把限流实现的那叫一个优雅，佩服！

https://juejin.cn/post/7145435951899574302

- 面试突击81：什么是跨域问题？如何解决？

https://juejin.cn/post/7140618562351136804

- 高并发下秒杀商品，你必须知道的9个细节

https://juejin.cn/post/7140900177438572574

- 聊一聊作为高并发系统基石之一的缓存，会用很简单，用好才是技术活

https://juejin.cn/post/7151937376578142216

- 聊一聊缓存和数据库不一致性问题的产生及主流解决方案以及扩展的思考

https://juejin.cn/post/7156237969202872350

- 为什么有HTTP协议，还要有websocket协议？

https://juejin.cn/post/7144161126652051464