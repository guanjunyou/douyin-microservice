# 一、项目介绍

> 简单抖音后端项目
>
> 微服务最终版本 ： https://github.com/guanjunyou/douyin-microservice
>
> 早期单体架构版本 ： https://github.com/guanjunyou/douyin 
>
> 文件处理服务 ： https://github.com/guanjunyou/FTPServer/tree/main/ftpServer

# 二、项目分工

| **团队成员** | **主要贡献**                                                 |
| ------------ | ------------------------------------------------------------ |
| 关竣佑       | 负责项目的设计，基础架构组件的搭建，视频 Feed 流接口，文件上传储存服务，点赞优化操作等功能的编写, 项目向微服务演进， 项目性能测试工作，项目后期bug修复工作 |
| 邱祥凯       | 负责评论操作，评论列表，发送消息，聊天记录接口编写，以及评论功能的优化 |
| 王奕丹       | 负责赞操作，喜欢列表接口编写，以及点赞，评论的查询优化       |
| 杨伟宁       | 负责基础架构组件的搭建，关注操作，关注列表，粉丝列表，好友列表接口编写，以及关注操作的优化 |
| 谢声儒       | 负责基础架构组件的搭建，基础接口（除了feed）的编写，点赞功能优化，微服务框架搭建，项目向微服务演进工作 |

## **开发规约**

#### 强制

1. 主体逻辑代码必须放在service层中的impl层，禁止在controller层写过多大的业务的代码，controller层应尽量调用service层的方法实现业务逻辑
2. model 层的函数禁止调用其它model 层相同包下不同 model 的函数
3. 返回给前端的数据若要组装成一个 struct 必须使用 xxxDVO来命名，参见 models.VideoDVO
4. model中 禁止进行sql字符串拼接，避免造成sql注入风险，如需使用参数拼接必须使用 ？ 传参 如

```Go
err := utils.DB.Where("is_deleted != ?", 1).Find(&videolist).Error
```

1. 遇到的所有 error 返回都必须进程处理或返回给上级（如使用 log 输出日志）

```Go
                if err1 != nil {                        log.Printf("Can not get the token!")                }
```

1. 所需用到的参数均放在config.go中，禁止在代码中出现魔法值。（所谓魔法值，是代码中莫名其妙出现的数字，数字意义必须通过阅读其他代码才能推断出来，这样给后期维护或者其他人员阅读代码，带来了极大不便。)如以下代码便出现了魔法值

```Go
// 遍历查询出的审查人对象集合        for(AuditPersonInfoDTO adp : auditPersonInfoDTO){            // 判断审查结果是否为空            if(adp.getStatus()!=null){                // 设置审查状态，status为2代表审核通过，为3代表退回修改                switch (adp.getStatus()){                    case "2" :                        adp.setStatus("审查通过");                        break;                    case "3" :                        adp.setStatus("退回修改");                        break;......
```

1. 每次开发前都必须pull代码！！！不然可能会造成冲突，很难解决。尽量先新建一个分支，测试功能正常后再与main分支合并
2. 禁止对已有文件进行移动（比如说移到其它包内），如需对结构有较大修改请提前说明
3. 每次 push 代码时禁止直接提交到 Master 分支 ！必须新建分支，运行测试正常后再提交分支！合并分支时遇到冲突需慎重解决，不明白的及时提出或让其他人帮忙合并
4. 所有实体类的成员必须使用**首字母大写**的驼峰命名法，Go 语言只用大写首字母才能被其它包访问。
5. 如需更改数据库请提前说明！
6. 如需提交更改后的数据库禁止删掉之前的数据库文件，以 日期-版本号.sql命名 (如：2023-7-21-v1douyin.sql)
7. 分支合并之后必须删除GitHub上的分支，每个人在GitHub上最多拥有一个分支
8. 编写接口时返回的数据一定要按照接口文档要求返回的数据
9. 数据库禁止使用外键

#### **建议**

1. 推荐使用 Goland 进行开发，使用Goland 的 git 图形化工具操作 git

 2.合并分支解决冲突的时候如遇不理解的问题及时提出

1. 开发一个函数后，建议在 test 包下编写测试代码进行测试
2. 如果业务操作间没有太多的关联，建议开启协程，使用 channel 通信。
3. 创建切片数组前，如果能估计大小，建议预先设置好大小，减少后期扩容开销

#### **注意**

1. 请求格式特别是 POST 请求的格式参照原本的代码。它里面有的POST请求不放json而使用拼接URL（我也不知道为什么），这里很坑

# 三、项目实现

### 3.1 技术选型与相关开发文档

#### 3.1.1  技术选型

##### 技术选型

后端框架：gin、go-micro、GORM

中间件：redis、rabbitMq

数据库:MySQL

系统监控: Prometheus , grafana

##### 技术评估

###### 后端框架

- gin:目前进行go-web开发的主流框架，学习成本低且开发效率高
- go-micro:目前go成熟的微服务框架之一，学习成本低且分层明确，支持注册中心可插拔
- GORM:go中最好用的ORM框架之一，覆盖绝大多数的使用场景

###### 中间件

- redis：目前最热门的缓存中间件，基于内存交互可以极大提高相应速度、降低数据库压力、
- rabbitMq：RabbitMQ是一个开源的消息队列系统，用于在应用程序之间传递和存储消息，实现高效的异步通信机制。

###### 数据库

- MySQL: 成熟的关系型数据库，具有广泛的支持和优化工具，适合处理关系型数据。

###### 鉴权

- 登录：登录时使用 jwt 将 username 和 CommonEntity 作为负载生成 token, 然后将 token 存入 redis 中  。 
- 鉴权：操作个人敏感数据或者涉及指定个人的接口时，需要针对用户身份和登录与否进行验证。首先将接收到的 token 进行正确性验证，同时解析出username 等消息，然后从 redis 中查找判断 token 是否过期。

​    (在必须登录的接口，这些操作在网关层的中间件执行，在不登录即可访问的接口若需要获取私人信息则需自行解析 token 鉴权)

##### 技术使用

整体框架采用 go-micro 微服务框架，采用 GROM 与 mysql 数据库进行交互，采用 Redis 作为缓存技术，使用rabbitMq作为消息队列.

目前的 rabbitmq消息队列采取发布订阅模式(Pub/Sub )，可以将消息发送给不同服务的消费者，方便后期模块的扩展

### 3.2 架构设计

本项目目前使用 3台服务器进行部署：分别为微服务主体项目所在服务器，文件服务器，etcd注册中心和rabbitmq 所在服务器。后续若进行微服务分别部署则至少考虑需要5台服务器。

#### 3.2.1 总体架构设计

总体架构图：

暂时无法在飞书文档外展示此内容

##### 微服务架构

1. 采用当今主流的微服务架构进行后端开发，在进行服务拆分的时候考虑到 视频，点赞，评论这三个功能耦合度较高，关注，好友，聊天这三个功能耦合度较高，决定拆分成三个大服务：user ,  video(video, favourite, comment) ,  relation (follow , message) 
2. 采用 ETCD 作为注册中心承担服务发现的功能 ， 使用 RPC 进行服务间的远程调用,  ETCD 使用 docker 在远程服务器上部署
3. 使用proto作为微服务之间传输数据的格式，将请求、响应结构以及远程服务方法编写为proto文件，利用代码生成器生成.pb.go和.pb.micro.go文件，提高开发效率
4. 所有的请求都先请求到 gateway 网关服务，经过鉴权和一系列前置操作后再分发给对应功能的服务
5. 微服务之间的通信采用rpc和消息队列，有效实现了服务之间的解耦合

##### 项目代码结构

每个服务的业务部分代码均采用 controller - service - serviceImpl - model 四级结构

暂时无法在飞书文档外展示此内容

项目总的目录结构如下 （省略部分）

├─app

│  ├─gateway

│  │  ├─cmd   存放 main.go gateway 服务启动入口

│  │  ├─http   网关从HTTP  API 

│  │  ├─middleware 中间件工具

│  │  ├─models 实体

│  │  ├─router  路由控制

│  │  ├─rpc 与其他服务间调用代码

│  │  └─wrappers 微服务之间的调用进行熔断器的封装

│  ├─relation

│  │  ├─cmd

│  │  ├─controller

│  │  ├─models

│  │  ├─mq 消息队列的配置

│  │  ├─rpc

│  │  ├─service

│  │  │  └─impl

│  │  ├─test  测试函数

│  │  └─utils 

│  ├─ 其它 relation 同级的服务

├─config 配置文件

├─idl

│  └─pb  存放 .proto 等文件

├─pkg

│  └─utils 公共工具

├─public

##### 数据库设计

所有数据库的表的结构如下：

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ODA4NWQ2ZmYzYzRmM2VlOWQzZTViNjVjOTg3MjZjZGVfNlNDQmM3cjI3eEw4TDJQMThtdWpleTgyQ0ZBTlBiY01fVG9rZW46TkVsUWJuZzZlb1o5dVl4Rmk1MGNDdFVabjkzXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

三个服务由于数据库层面已经去除耦合，考虑实行 **垂直分库** 

三个服务使用不同的数据库，提升性能和可靠性。分库成**三个数据库** ：**user (user表),**  **video (video 表，like表，comment表) ， relation（message_send_event表，follow 表）**

三个服务分别连接三个数据库

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YTVhNmNkZDA1MDkwN2IxM2M4OWQzZTNjZDI1NWY0ZDdfcVcwYW4yWmlRUVNFR052ZUxJUHFKU0JHWkxBcWdyWXVfVG9rZW46T29CRGJub2Rsb2d2SEZ4dnhxUGNYaUl3bmhnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 维护数据库和缓存的一致性

1. 更新数据时先更新数据库再更新缓存
2. 删除数据时采用缓存**延迟双删**

删除操作时执行延迟双删确保 redis 中不会出现脏数据：先在 redis 中删除数据后，再从数据库中删除数据，若在数据库删除成功前，另一个线程查询了数据库的没有删除的数据后写回了redis 会导致缓存于数据库不一致

以本项目中的**取消点赞**功能为例介绍缓存延迟双删步骤如下：

1. 从 redis 中对应的 like set 删除该视频 ID 
2. 执行数据库删除等一系列操作
3. 把数据库删除操作写入消息队列
4. 消费消息队列的操作，删除 redis 中对应 like set 中该视频的 ID （本项目中仅实现了第二次删除，未实现消息队列延时删除）

##### 逻辑删除

1. Mysql 的所有表均包含 `is_deleted` 字段，当值为 0 时表示该数据存在，值为 1 时表示该数据被删除。避免误删数据，同时也可以方便地恢复数据
2. 在本项目的所有数据库操作中删除操作均采用 逻辑删除

##### 数据库冗余设计

在微服务的架构下，为了解耦合，不同的表有时候分属不同的服务，导致多表查询变得困难。若涉及到查询其它属于其它服务的数据就要通过RPC调用远程函数，当缓存失效时时间代价很大，会使得用户感觉到明显的时延。故通过数据库表增加冗余字段的`反范式`手段来提高查询性能，在本项目中有如下实践：

1. User 表增加获赞数，作品数，喜欢数，关注数，被关注数等冗余字段，在相关数据更新时主动推送到 user 表，获取用户信息时可立刻返回
2. Video 表增加获赞数，评论数等冗余字段，在相关数据更新时主动推送到 user 表，获取视频信息时可立刻返回

##### 索引的合理设置

###### like表

我们在查询喜欢列表时，考虑到总是会根据当前登录用户的userId去寻找videoId，于是我们为like表的`user_id`建立了普通索引。选择普通索引的原因如下：

1. `user_id`本身是由雪花算法生成的，已经保证了唯一性，因此没必要使用唯一索引
2. 唯一索引与普通索引的查询性能基本没有区别
3. 使用普通索引可以利用change_buff机制对更新语句进行加速，提高交互效率，唯一索引则无法利用这一机制

###### comment表

我们在查询视频的评论列表时，当数据库中没有缓存时，我们经常是根据`videoId`去作为查询条件的，因此我们为`video_id`添加了普通索引

###### follow表

在进行关注数以及粉丝数的查询时会经常用到`user_id`和`follow_user_id`两个字段作为查询条件，因此我们对这两个字段均加了联合索引

##### 独立的文件处理服务

为了使得服务间的解耦和服务的自治，规范化文件的存储，本项目的所有文件（视频，图片）均摈弃把文件存储在服务本地或者使用FTP远程调用的方式存文件。本项目开发一个 fileServer 文件系统来对文件进行统一管理，并暴露 HTTP 接口供其它服务调用，该系统有以下功能：

1. 提供HTTP接口，供 video 服务上传视频文件，使用 ffmpeg 截取第一帧作为图片封面后分别保存到本服务器，返回给 video 服务视频文件，封面图片的URL
2. 供其它服务调用的文件存储接口 （后续还可以对接 minio 等分布式文件系统 ）
3. 文件服务器上使用 Nginx 将视频，图片等文件开放给 APP 前端进行访问

##### 高并发场景解决方案

1. 点赞，关注等场景由于可能在短时间内有较大的并发量，如果任由这些请求立刻操作数据库将会给数据库造成巨大的压力，甚至导致宕机。而且如果走完这一系列的操作再返回给用户，用户将会等待很长的时间，会导致用户的流失。因此，本项目采用 **消息队列****和管道** (channel ) 相结合的方式进行**削峰。**

采用生产者消费者模型进行**异步处理**消费数据，当操作数据成功放置入 消息队列 或者 channel 的时候，即可 

​      返回给用户成功。后续消费数据执行操作确保数据库在平稳的压力下处理，失败率是很低的。倘若出现执行失

​      败的情况，则需要进行**重试**操作（重试操作目前还未实现），重试次数多了之后仍失败就加入失败队列人工介 

​      入处理。

 如果在极端情况下仍然出现后续数据执行失败导致暂时数据不一致的情况，在点赞，关注的功能中也影响极小，牺牲这极小机会的数据不一致来换取用户操作时的快速响应，保证用户体验是值得的。

1. 在点赞 , 关注操作中，在并发量大的情况下，如果恰巧多个点赞请求同时进入，第一个请求未执行完毕，其它请求通过数据库判断未点赞时，会导致连续执行了多次点赞操作。为了保障接口的**`幂等性`**，考虑使用 **Redis** **`分布式锁`**的解决方法。当点赞时尝试获取点赞锁，若获取成功则释放 取消点赞锁 继续执行后续操作。
2. 解决缓存穿透问题：使用**布隆过滤器** 添加数据的 ID ， 或者每次查询到不存在数据时在 redis 中缓存空值        解决缓存雪崩问题：每次生成redis key 的时候 TTL 添加随机值                                             
3. 限流 ：采用熔断器进行限流，后续增加针对某个 ip 进行限流

##### 远程调用重试机制

本项目是微服务架构，服务间存在着许多的远程函数调用，为了避免因网络状况等导致的偶然发生的远程调用失败，在每次调用都设置重试机制，三次都调用的失败的概率很小，若三次调用仍然失败则需要引起重视。

下面是一段重试代码例子：

```Go
var req pb.CheckFollowRequest
req.UserId = userClaim.CommonEntity.Id
req.ToUserId = userId
for i := 0; i < retryLimit; i++ {
    resp, err0 := rpc.RelationClient.CheckFollowForUser(context.Background(), &req)
    if err0 == nil {
       user.IsFollow = resp.IsFollow
       break
    }
}
```

##### 熔断和降级

为了在高并发场景下，保护后端服务，提高系统的可靠性。在 Gateway 层 `wrappers` 包下设置了熔断器（针对 user , video , relation 的服务），以针对 user 服务为例，设置的参数如下：

```Go
func (wrapper *userWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
    cmdName := req.Service() + "." + req.Endpoint()
    config := hystrix.CommandConfig{
        MaxConcurrentRequests:  3000, // 并发数上限
        Timeout:                30000,
        RequestVolumeThreshold: 20,   // 熔断器请求阈值，意思是有20个请求才能进行错误百分比计算
        ErrorPercentThreshold:  50,   // 错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次 开启
        SleepWindow:            1000, // 过多长时间，熔断器再次检测是否开启，单位毫秒ms
    }
    hystrix.ConfigureCommand(cmdName, config)
    return hystrix.Do(cmdName, func() error {
        return wrapper.Client.Call(ctx, req, resp)
    }, func(err error) error {
        return err
    })
}
```

熔断器的参数应该在后续的测试和实践中进行修改

##### 防御恶意攻击

1. **数据库加密**：密码的加密存储使用 bcrypt 算法，由于 bcrypt 算法加入了盐值，盐是一个随机生成的字符串，它与密码一起被哈希。由于盐是随机生成的，因此即使两个用户使用相同的密码，它们的哈希值也不同。这使得攻击者更难以破解密码。校验时，从hash中取出salt，salt跟password进行hash；得到的结果跟保存在DB中的hash进行比对
2.  **IP** **限流与封禁：**使用 Nginx 的 `ngx_http_limit_req_module` 模块对 IP 进行限流，防止恶意 ip 的攻击
   1.  在 Nginx 层设置每个 ip 每秒限制 10 个请求，允许突发达到 20 个每秒

   2.  Nginx.conf 的部分内容如下：

```Go
http {
    include       mime.types;
    default_type  application/octet-stream;


    sendfile        on;
    
    keepalive_timeout  65;

    limit_req_zone $binary_remote_addr zone=mylimit:10m rate=10/s;

    server {
        listen       8010;
        server_name  0.0.0.0;

        upstream targetserver{
        server 127.0.0.1:8080;

    }

        location / {
            limit_req zone=mylimit burst=20nodelay;
            proxy_pass http://targetserver;
        }
        ......
```

在 GateWay 层使用 Redis 用 ip作为 key , 每秒请求的次数作为 value 

##### 系统监控

> 项目运行在 172.28.158.98 服务器上， **`Prometheus`**  **和** **`grafana`**  运行在 172.28.158.130 服务器上，
>
> 这两台服务器均为 内网服务器，之间的带宽为 2G ，装有 centos8 操作系统

本项目采用 **`Prometheus`** **和** **`grafana`** 进行系统监控，使用 grafana 对 Prometheus 搜集到的数据进行可视化监控

下面以监控 user 服务为例子进行介绍

Prometheus.yml 如下 （用于 docker 运行 Prometheus）

```YAML
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    static_configs:
      - targets: ['172.28.158.98:8070']
```

然后在 main.go 中加入以下的部分代码

```Go
func PrometheusBoot() {
    // 创建 HTTP 处理器
    h := promhttp.HandlerFor(
       prometheus.DefaultGatherer,
       promhttp.HandlerOpts{},
    )
    http.Handle("/metrics", h)
    // 启动web服务，监听8085端口
    go func() {
       err := http.ListenAndServe("0.0.0.0:8085", nil)
       if err != nil {
          log.Fatal("ListenAndServe: ", err)
       }
    }()
}
```

项目和 Prometheus 启动后 ， 进入 http://172.28.158.130:8060 看到服务已经成功注册

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZDIzMTFhY2RhZGJjMDY3YTNmYTRjMTM0OTI1NDdkYThfZzBxdk1OMkFabG5QbXJ1QnN3MXJEbkYzNjQybmJpWkhfVG9rZW46WVlxSGJhUDRtb1k3N2x4dzI2amNqeTBEbkZlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

运行 grafana 容器 `docker run -d -p 3000:3000 grafana/grafana`

在服务器上运行 **grafana** 容器后，连接对应的   Prometheus  数据源，可在 web 的数据面板看到监控服务的数据

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MWUxZTRhYWEwMmZlYzI2NjAxOGNmMjY5ODlhMjYzNDVfMUVPRUozbDZLN0NqdnJuaTJOQ2E4ZVpLZnBvVllVQ1JfVG9rZW46WVNsS2J3a0Mxb3FGODh4U1ViRWNZQVlNblJmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MzU1MjJjZDI0N2NjNzZjZjkzNTdjYTkyZTJmNzBiODBfTmFuQzM1RGZ3VmtJeVN3YVVnamhxaUZWYmh1N3hSZkpfVG9rZW46WnNYZGJzeEVab05lOWF4YWN0a2MxMUtubkZjXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### 3.2.2 GateWay  及公共组件 功能设计

##### 登录鉴权

使用两层中间件`middleware`对网关收到的所有请求进行预处理，依次为 *`AuthAdminCheck()`**和* *`RefreshHandler`**，用于redis中的token刷新，**`AuthAdminCheck()`**用于登录校验*

*`RefreshHandler`**的具体实现：*

1. *从请求头中获取token*`token := c.Query("token")`，*如果token为空则尝试从body中拿*
2. *判断是否携带token，如果token为空直接放行*
3. *调用**`utils.AnalyseToken(token)`**解析token，将结果保存在**`userClaims`**中*
4. *根据**`userClaims.Name`**查redis，执行*`tokenFromRedis, err := utils.GetTokenFromRedis(userClaims.Name)`
5. `tokenFromRedis`为空则*重建redis缓存*
6. *刷新token的有效期*

*`AuthAdminCheck()`**的具体实现：*

1. 判断请求是否需要登录鉴权，不需要的直接放行
2. 从请求头或者请求体中获取token
3. 使用jwt解析token，从解析结果中获取用户名
4. 根据用户名去查询redis缓存，如果缓存中存在的话放行，不存在则直接阻止该请求

```Go
// 免登录接口列表
var notAuthArr = map[string]string{
    "/douyin/feed/":          "1",
    "/douyin/user/register/": "1",
    "/douyin/user/login/":    "1",
    "/douyin/user/":          "1",
}
```

Router 代码如下 

```Go
import (
    "douyin-microservice/app/gateway/http"
    "douyin-microservice/app/gateway/middleware"
    "github.com/gin-gonic/gin"
)
func InitRouter() *gin.Engine {
    r := gin.Default()
    r.Use(middleware.Cors())
    // public directory is used to serve static resources
    r.Static("/static", "./public")
    //所有请求都需要两层拦截器
    apiRouter := r.Group("/douyin")
    {
        v1 := apiRouter
        v1.Use( middleware.AuthAdminCheck(),middleware.RefreshHandler())
        {
            v1.GET("/feed/", http.FeedHandler)
            v1.POST("/user/register/", http.RegisterHandler)
            v1.POST("/user/login/", http.LoginHandler)
    
            //apiRouter2 := r.Group("/douyin")
            // extra apis - I
            v1.POST("/favorite/action/", http.FavoriteActionHandler)
            v1.GET("/favorite/list/", http.FavoriteListHandler)
            v1.POST("/comment/action/", http.CommentActionHandler)
            v1.GET("/comment/list/", http.CommentListHandler)
            //
            //// extra apis - II
            v1.POST("/relation/action/", http.RelationActionHandler)
            v1.GET("/relation/follow/list/", http.FollowListHandler)
            v1.GET("/relation/follower/list/", http.FollowerListHandler)
            v1.GET("/relation/friend/list/", http.FriendListHandler)
            v1.GET("/message/chat/", http.MessageChatHandler)
            v1.POST("/message/action/", http.MessageActionHandler)
    
            v1.GET("/user/", http.UserInfoHandler)
            v1.POST("/publish/action/", http.PublishHandler)
            v1.GET("/publish/list/", http.PublishListHandler)
        }
    }
    return r
}
```

##### 封装公共实体 CommonEntity

本项目针对数据库中的公共字段 ID ，数据创建时间， 删除标志 封装了一个结构体

```Go
type CommonEntity struct {
    Id         int64     `json:"id,omitempty"`
    CreateDate time.Time `json:"create_date,omitempty"`
    IsDeleted  int64     `json:"is_deleted"`
}

func NewCommonEntity() CommonEntity {
    sf := NewSnowflake()
    return CommonEntity{
       Id:         sf.NextID(),
       CreateDate: time.Now(),
       IsDeleted:  0,
    }
}
```

所有和数据库实体结构体均继承 `CommonEntity` ， 创建时调用 `NewCommonEntity` 函数创建

##### 雪花算法生成分布式ID

由于本项目是分布式系统，而且抖音后端面临的是庞大的用户群体，高并发量以及庞大的数据量需要数据库的分库分表。传统的自增 ID 难以在这种情况下正常运作。因此本项目所有的 ID 均采用雪花算法生成

本项目对雪花算法进行了封装 见项目中的/pkg/utils/Snowflake.go文件

每次 CommonEntity 生成 ID 时都调用函数生成 

##### 敏感词过滤器

本项目中涉及到许多的文本发布功能，如视频标题，评论，聊天消息等等，为了确保文本符合法律法规，没有不允许发布的敏感信息，减少后续人工筛查的工作量，在本项目中设置了敏感词过滤器 `sensitive.Filter`

每当文本上传的时候，若文本中含有敏感词文件中包含的敏感词，则会替换成  `*`  号

#### 3.2.3 User 功能设计

##### 用户注册 /douyin/user/register/

暂时无法在飞书文档外展示此内容

布式系统中，存在多个子系统或服务，这些子系统可能要获取用户身份或者实现登录。通过将Token存储在Redis中，可以实现不同子系统之间的Token共享，从而实现用户在一个子系统登录后，其他子系统无需再次登录。因此考虑将 token 保存 到 redis 中

而因为在这个接口中传入的用户名 `username` 是每个用户唯一的，故使用 username 作为 redis 的 key 

##### 用户登录 /douyin/user/login/

1. 获取传入的 username 和 password 
2. 查询用户是否存在 ，若不存在则返回用户不存在
3. 比较数据库中已经哈希过的密码和用户提供的明文密码是否匹配

```Go
pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
```

1. 生成并保存 token , 返回 token 
2. 登录成功

##### 用户信息 /douyin/user/

1. 首先，函数通过调用`userService.GetUserById(userId)`来获取用户信息。如果在这个过程中发生错误，或者返回的用户不存在，那么函数将返回一个错误信息："用户不存在！"。

在这个步骤中为了提高响应速度，查询用户时先在 redis 中查询，若查不到再从数据库中查询并写入 redis 

用户信息序列化为 `json字符串`的形式存放在 redis 中

1. 如果`token`不为空，那么函数将尝试分析这个`token`。如果分析过程中发生错误，或者返回的用户声明为空，那么函数将直接返回用户信息，不进行后续的检查。

这个步骤的目的是如果用户已经登录，请求获取的目的用户信息就要包含是否关注过的信息

1. 如果token分析成功，那么函数将创建一个`CheckFollowRequest`对象，并设置其`UserId`和`ToUserId`属性。
2. 调用`rpc.RelationClient.CheckFollowForUser(context.Background(), &req)`来检查当前用户是否关注了指定的用户。这个调用的结果将被用来设置`user.IsFollow`属性。

被远程调用 `Relation` 服务的 `CheckFollowForUser` 函数功能是使用  `UserId`和`ToUserId` 查询关注记录是否存在

1. 最后，函数将返回用户信息。

暂时无法在飞书文档外展示此内容

#### 3.2.4 Video 功能设计

##### 视频流接口/douyin/feed/

暂时无法在飞书文档外展示此内容

若传入的 `latestTime` 为 0 ， 则设置为 当前时间

若 token 为空则设置 `userId` 为 -1 

1. 调用`models.GetVideoListByLastTime(latestTime)`方法，根据用户最后观看时间获取视频列表，并将结果存储在`videolist`变量中。同时，初始化一个`size`变量用于存储视频列表的长度。
2. 创建一个`sync.WaitGroup`类型的变量`wg`，用于等待所有并发的协程任务完成。
3. 创建一个长度为`size`的`VideoDVOList`切片，用于存储符合条件的视频信息。
4. 如果`videolist`为空或者在执行过程中出现错误，直接返回`nil`、当前时间和一个错误信息。
5. 遍历`videolist`中的每个视频，对于每个视频：
   1.  a. 获取视频的作者ID（`authorId`）。 

   2.  b. 启动一个并发协程，该协程会执行以下操作：

   3.   i.   通过`videoService`调用`userService`，获取作者的信息。

   4.   ii.  将视频信息复制到一个新的`videoDVO`对象中。 

   5.   iii. 如果`userId`不等于-1，调用`favoriteService.FindIsFavouriteByUserIdAndVideoId(userId, videoDVO.Id)`方法，判断当前用户是否收藏了该视频。否则，将`videoDVO.IsFavorite`设置为`false`。

   6.   iv.  将`videoDVO`对象添加到`VideoDVOList`切片中。 c. 等待所有并发任务完成。
6. 返回`VideoDVOList`切片、下一个视频的创建时间（即`videolist`中最后一个视频的创建时间）以及可能的错误信息。

VideoDVO 如下：

```Go
type VideoDVO struct {
    utils.CommonEntity
    Author        User   `json:"author"`
    PlayUrl       string `json:"play_url"`
    CoverUrl      string `json:"cover_url"`
    FavoriteCount int64  `json:"favorite_count"`
    CommentCount  int64  `json:"comment_count"`
    IsFavorite    bool   `json:"is_favorite"`
    Title         string `json:"title,omitempty"`
}
```

##### 投稿接口 /douyin/publish/action/

暂时无法在飞书文档外展示此内容

1. 从`title`中过滤敏感词汇，将结果赋值给`replaceTitle`。
2. 调用`utils.UploadToServer`方法，将视频文件上传到服务器，并将返回的封面文件名赋值给`coverName`。如果上传过程中出现错误，直接返回错误信息。
3. 创建一个`models.Video`结构体实例，并设置其属性值。其中，`CommonEntity`、`AuthorId`、`PlayUrl`和`CoverUrl`分别设置为新创建的通用实体、作者ID、视频播放地址和封面图片地址。`Title`属性设置为过滤后的标题。
4. 调用`models.SaveVideo`方法，将视频信息保存到数据库中。如果保存过程中出现错误，返回错误信息。
5. 更新用户的发布作品数，将其加1。
6. 远程调用 User 服务的方法，更新用户信息。如果更新过程中出现错误，返回错误信息。
7. 如果以上步骤都执行成功，返回发布成功

在这个接口中，曾经考虑过使用异步处理的方式，让后面的保存等操作在返回给用户之后再进行，提高响应速度，但是后面考虑到这个接口不能出现用户发布视频失败但是不知情的严重情况，故牺牲速度来确保用户成功上传视频

##### 发布列表 /douyin/publish/list/

1. 调用`models.GetVediosByUserId`方法，根据用户ID获取用户发布的所有视频信息，并将结果赋值给`videoList`。如果获取过程中出现错误，直接返回错误信息。
2. 获取`videoList`的长度，并将其赋值给变量`size`。
3. 创建一个长度为`size`的`models.VideoDVO`切片，并将其赋值给变量`VideoDVOList`。
4. 创建一个同步等待组`wg`，用于等待所有协程完成。
5. 定义一个`error`类型的变量`err0`，用于接收协程产生的错误。
6. 使用`for`循环遍历`videoList`，对于每个视频： a. 将协程并发数加1。b. 启动一个协程，在其中执行以下操作： i.   获取当前视频的作者ID，并将其赋值给变量`userId`。 ii.  创建一个`pb.UserRequest`结构体实例，并设置其`UserId`属性为`userId`。 iii. 调用`rpc.UserClient.GetUserById`方法，根据用户请求获取用户信息，并将结果赋值给变量`userResp`。如果获取过程中出现错误，将错误赋值给变量`err1`。 iv.  创建一个`models.VideoDVO`结构体实例，并使用`copier.Copy`方法将其复制到新创建的结构体实例中。如果复制过程中出现错误，将错误赋值给变量`err`。 v.   将用户信息转换为`BuildUser`函数返回的用户对象，并将其赋值给`videoDVO.Author`。 vi.  将`videoDVO`添加到`VideoDVOList`切片中。 c.  协程执行完毕后，将`wg.Done()`作为等待组的结束信号。
7. 调用`wg.Wait()`方法，等待所有协程完成。
8. 检查协程内是否存在错误，如果有，则返回错误信息。
9. 返回`VideoDVOList`切片和`nil`错误。

**文件类型检查:** 在接口中，用户需要上传视频保存到服务器。但是若有人利用此接口上传恶意文件或者木马，则会造成严重的危害，因此考虑使用`h2non/filetype` 包进行文件类型检查

在 impl 文件中加入了如下代码进行文件类型检查

```Go
if filetype.IsVideo(data) {
    return errors.New("视频格式异常！")
}
```

#### 3.2.5 Favourite 功能设计

##### 赞操作 /douyin/favorite/action/

1. 尝试获取Redis分布式锁,该分布式锁基于`setnx`命令实现，分为两种锁，锁的key如下：
   1. ```Go
      lockKey := config.LikeLock + userIdStr + videoIdStr 
      unLikeLockKey := config.UnLikeLock + userIdStr + videoIdStr
      ```

   2.  对于点赞动作类型进行不同的处理：

a.点赞：尝试获取以lockKey为key的锁，i. 获取失败则直接返回`errors.New("-1")` ii. 获取成功，释放以unlockKey为key的锁

b.取消点赞: 尝试获取以unlockKey为key的锁，i. 获取失败则直接返回`errors.New("-1")` ii. 获取成功，释放以lockKey为key的锁

以上操作用来避免同一个用户重复点赞或取消点赞

1. 以`userLikeKey := config.LikeKey + userIdStr`为key从redis中相应的set结构查询有没有点赞的`videoId`，如果没有进一步从数据库中查询
2.  对于点赞操作，如果查询到直接返回“用户已点赞”，查不到则可以进一步调用`models.GetVideoById(videoId)`查询具体video数据，从中获取`authorId`

i.封装`mqData := models.LikeMQToUser{UserId: userId, VideoId: videoId, ActionType: actionType, AuthorId: authorId}`序列化为json后发送给`LikeRMQ`队列

ii.发送到对应管道`mq.LikeChannel <- mqData`

iii.主程序直接返回

iv.Video模块的*`func LikeConsumer`*`(ch <-`*`chan `*`models.LikeMQToUser) `异步消费管道中的消息，将相应video的`FavoriteCount`++，并将videoId添加到`userId`对应的set集合中

v.user模块的*func* `(userService UserServiceImpl) likeConsume(message <-`*`chan `*`amqp.Delivery)`异步消费队列中的消息，对点赞用户执行`user.FavoriteCount = user.FavoriteCount + 1`，对视频作者执行`user.TotalFavorited++`，此时如果是视频作者给自己点赞必须在同一条update语句更新这两个字段。如果在两条update语句中更新同一条记录，会因为update语句的redolog文件会被覆盖，导致只有后一条更新生效

1. 对于取消点赞操作，如果查询不到则直接返回`未找到要取消的点赞记录`的错误，查到则可以进一步调用`models.GetVideoById(videoId)`查询具体video数据，从中获取`authorId`

i.封装`mqData := models.LikeMQToUser{UserId: userId, VideoId: videoId, ActionType: actionType, AuthorId: authorId}`序列化为json后发送给`LikeRMQ`队列

ii.发送到对应管道`mq.LikeChannel <- mqData`

iii.主程序直接返回

iv.Video模块的*`func LikeConsumer`*`(ch <-`*`chan `*`models.LikeMQToUser) `异步消费管道中的消息，将相应video的`FavoriteCount`--，并将videoId从`userId`对应的set集合中删除

v.user模块的*func* `(userService UserServiceImpl) likeConsume(message <-`*`chan `*`amqp.Delivery)`异步消费队列中的消息，对点赞用户执行`user.FavoriteCount = user.FavoriteCount - 1`，对视频作者执行`user.TotalFavorited`--，如果是作者取消点赞也要保证使用一条update语句

示意图如下：

暂时无法在飞书文档外展示此内容

##### 喜欢列表  /douyin/favorite/list/、

1. 拼接`likeKey := config.LikeKey + strconv.FormatInt(userId, 10)`，以该key从redis相应的set结构中找到所有的videoId，如果找不到则进一步从数据库中的like表中根据`userId`查询相应的`videoId`，将查询到的加入到`likeIdsSet`中
2. 创建一个同步等待组`wg`，用于等待所有协程完成,创建*`var `*`res []models.LikeVedioListDVO`用于保存待返回信息
3. for循环遍历`likeIdsSet`，对于每次循环，在开始前`wg`的计数器+1，并开启一个协程，协程内进行如下处理

a.根据videoId从数据库中查询详细的`video`记录，从记录中获取`AuthorId`

b.调用`rpc.UserClient.GetUserById`方法从远程user服务中查询作者信息，保存在`author`中

c.创建*`var `*`likeVideoListDVO models.LikeVedioListDVO`，将`video`和`author`封装进去

d.执行`res = `*`append`*`(res, likeVideoListDVO)`，将数据添加到切片res中

e.`wg`的计数器-1

1. *现在* *`res`* *包含了所有视频的作者和视频信息,直接返回**`res`*

LikeVedioListDVO的结构：

```Go
type LikeVedioListDVO struct {
    Video
    Author *User json:"author" gorm:"foreignKey:AuthorId"
}
```

#### 3.2.6 Comment 功能设计.

##### 评论操作 /douyin/comment/action/

当action_type=1即发表评论时：

1. `将comment.User.Id`封装到*`var `*`req pb.UserRequest`中，调用`rpc.UserClient.GetUserById`从user服务查询具体用户信息
2. 将评论的id加入到布隆过滤器中
3. 封装`models.CommentMQToVideo`结构体，并将该结构体发送到相应管道中`mq.CommentChannel <- toMQ`
4. 主程序直接返回
5. `commentActionConsumer()`异步消费`mq.CommentChannel`中的消息流程如下：

a.调用`models.SaveComment`将评论数据保存到数据库中

b.以评论的`videoId`为key,将评论的id保存到相应的zset结构中，zset的score为评论的创建时间，成员为评论的id

c.以评论的id为构造key`commentExistKey := "comment:" + strconv.Itoa(int(commentDB.Id))`，将评论进行json序列化后保存到对应的string结构

当action_type=2即删除评论时：

1. 使用布隆过滤器初步判断待删除的评论id是否存在，如果不存在直接返回error
2. 构造`commentExistKey := "comment:" + strconv.Itoa(int(commentId))`判断redis中待删除的评论是否存在，不存在直接返回error
3. 从redis中删除相应的缓存
4. 封装`models.CommentMQToVideo`结构体，并将该结构体发送到相应管道中`mq.CommentChannel <- toMQ`
5. 主程序直接返回
6. `commentActionConsumer()`异步消费`mq.CommentChannel`中的消息流程如下：

a.调用`models.DeleteComment`将评论数据从数据库中删除

b.以评论的`videoId`为key,将评论的id从zset中的成员删除

c.以评论的id为构造key`commentExistKey := "comment:" + strconv.Itoa(int(commentDB.Id))`，将redis中保存的评论序列化字符串再次删除一次，双删保证数据库与缓存的一致性

CommentMQToVideo的结构：

```Go
type CommentMQToVideo struct {
    utils.CommonEntity
    ActionType int    `json:"action_type"`
    UserId     User   `json:"user"`
    VideoId    int64  `json:"video_id"`
    Content    string `json:"content"`
    CommentID  int64  `json:"id"`
}
```

流程图如下：

暂时无法在飞书文档外展示此内容

##### 评论列表 /douyin/comment/list/

1. 先使用布隆过滤器判断请求的视频id是否在缓存中，如果不在的话直接返回空
2. 以`videoId`为key从redis相应的zset结构中取出数据，如果不存在的话需进行查数据库并进行缓存的重构，维护以`videoId`为key、`commentId`为成员、创建时间为score的zset结构和以`commentId`为key、值为评论json序列化字符串的string结构
3. *创建**`var `*`comments []models.Comment`用于保存待返回的数据
4. 遍历取得的缓存，对于每一个取得的每一个评论Id，根据`commentId`从缓存中查询评论的json字符串，分为如下情况：

a. 查询不到则从数据库中查询并维护到redis中,然后添加到`comments`中

b. 查询到了就直接反序列为`models.Comment`的对象并添加到`comments`中

1. 返回`comments`

`models.Comment`的结构如下：

```Go
type Comment struct {
    utils.CommonEntity
    User    User   `json:"user"`
    Content string `json:"content,omitempty"`
}
```

#### 3.2.7 Relation 功能设计

##### 关注/取消关注接口  /douyin/relation/action/

1. 首先，检查userId是否等于`toUserId`，如果相等，则返回错误信息"你不能关注(或者取消关注)自己"。

2. 定义两个分布式锁key，一个用于关注操作，另一个用于取消关注操作。

3. 根据`actionType`的值进行不同的操作：

   1. 1.  如果`actionType为1`（关注操作）：

         1. 在Redis中设置一个分布式锁，锁的过期时间为`config.FollowLockTTL * time.Second`。如果设置成功 继续执行；否则，返回"已关注"的错误信息。

         ​       关注分布式锁的 key 为  `lockKey := config.FollowLock + userIdStr + toUserIdStr`

         1. 删除Redis中的取消关注锁。
         2. 检查Redis中是否存在用户关注的集合，如果存在，则检查该集合中是否有`toUserId`。如果已经关注了`toUserId`，则返回"已关注"的错误信息；否则，从数据库中查询关注记录。
         3. 如果缓存中没有找到关注记录，则从数据库中查询关注记录。如果找到了关注记录，则将`isExists`设置为true。
         4. 如果`isExists`为true，则返回"该用户已关注"的错误信息；否则，继续执行。
         5. 创建一个关注消息，将其加入消息队列，并将消息序列化为JSON格式。
         6. 返回nil表示操作成功。

   2. 1.  如果`actionType为2`（取消关注操作）：

         1. 在Redis中设置一个分布式锁，锁的过期时间为`config.UnFollowLockTTL * time.Second`。如果设置成功 继续执行；否则，返回"已取消关注"的错误信息。

         ​       取消关注的分布式锁 key 为 `unFollowLockKey := config.UnFollowLock + userIdStr + toUserIdStr`

         1. 删除Redis中的关注锁。
         2. 检查Redis中是否存在用户关注的集合，如果存在，则检查该集合中是否有`toUserId`。如果没有关注`toUserId`，则返回"未找到要取消的关注记录"的错误信息。
         3. 如果上一步中在缓存找到对应的集合， 把缓存中对应的集合 `follow` 和 `follower` 中的ID删除
         4. 如果缓存中没找到了关注记录，则从数据库中查询关注记录。如果未找到关注记录或关注记录的ID为0，则返回错误信息。
         5.  发送消息到对应channel  `mq.FollowChannel <- mqData`
         6. 创建一个取消关注消息，将消息序列化为JSON格式,  将其加入 RabbitMQ 消息队列

           为什么通知 User 更新数据使用消息队列而不使用协程？因为User 服务和 Relation 服务属于不同的服务，故使用消息队列来进行异步处理。

         1. 返回nil表示操作成功。

在上述步骤中发送给 channel 和 消息队列的消息结构如下

```Go
type FollowMQToUser struct {
    UserId       int64 `json:"user_id"`
    FollowUserId int64 `json:"follow_user_id"`
    ActionType   int   `json:"action_type"`
}
```

1. Relation 服务对 channel  `FollowChannel` 的消费流程：
   1. 启动多协程并发对 channel 进行监听消费
   2. 取出数据
      1. 如果`actionType为1`（关注操作）：将关注记录持久化到 mysql 数据库中，然后往 Redis 对应的 `userId` 的 `follow`集合 和 `FollowUserId` 的 `follower`集合插入 ID
      2. 如果`actionType为2`（取消关注操作）：从数据库中删除关注记录后，在Redis 对应的 `userId` 的 `follow`集合 和 `FollowUserId` 的 `follower`集合删除 ID 【缓存延迟双删】
2. User 服务对 RabbitMQ `FollowMQ` 的消费流程：
   1. 启动多协程并发消费 FollowMQ 中的数据
   2. 更新 user 表中相关用户的 关注数 和 被关注数 （同时更新 redis）

示意图如下：

暂时无法在飞书文档外展示此内容

##### 关注列表/douyin/relation/follow/list/

1. 定义变量 follows 和 users 两个切片 
2.  查询指定用户的关注列表，并将结果存储到 follows  变量中。如果查询出错，则返回错误信息。

若 Redis 中找到对应的 `follow` 集合则从集合中取出关注的 ID 

1. 定义协程并发更新函数，对上面查询出的每个关注 ID 远程调用 user 服务查询对应的 user 组装成 users 切片
2. 使用 wg.Wait() 等待所有协程完成并发更新操作。
3. 重建缓存，并返回关注列表

##### 粉丝列表  /douyin/relation/follower/list/

1. 定义变量 followers 和 users 两个切片 
2.  查询指定用户的粉丝列表，并将结果存储到 followers  变量中。如果查询出错，则返回错误信息。

若 Redis 中找到对应的 `follower` 集合则从集合中取出粉丝的 ID 

1. 定义协程并发更新函数，对上面查询出的每个粉丝 ID 远程调用 user 服务查询对应的 user 组装成 users 切片
2. 使用 wg.Wait() 等待所有协程完成并发更新操作。
3. 重建缓存，并返回粉丝列表

##### 好友列表 /douyin/relation/friend/list/

1. 查询 Redis 缓存中是否存在该用户的 `follow` 和 `follower` 集合， 若不存在则重建这两个集合的缓存
2. 使用 `SInter()` 函数 ， 求这两个集合的交集即可获得好友的 ID 列表
3. 定义协程并发更新函数，对上面查询出的每个好友 ID 远程调用 user 服务查询对应的 user 组装成 users 切片
4. 使用 wg.Wait() 等待所有协程完成并发更新操作。
5. 返回好友列表

####  3.2.8 Message 功能设计

##### 发送消息 /douyin/message/action/

1. 调用 `utils.AnalyseToken(token)` 函数对 token 进行分析，获取用户信息并存储在 userClaim 中，并以此获取用户的 ID 。如果分析出错，则返回错误信息。
2. 将发送者用户 ID 其存储在 userId 变量中。
3. 调用 `models.SaveMessageSendEvent(&models.MessageSendEvent{...})` 函数保存消息发送事件对象，并将返回的错误信息存储在 err 变量中。如果保存消息发送事件失败，则返回错误信息。

##### 聊天记录 /douyin/message/chat/

1. 调用 `utils.AnalyseToken(token)` 函数对 token 进行分析，获取用户信息并存储在 userClaim 中，并以此获取用户的 ID 。如果分析出错，则返回空的 []models.MessageDVO 和错误信息。
2. 判断传入的时间戳 `pre_msg_time` 的值 分别 处理
   1. ```Go
      var preTime time.Time
      if preMsgTime != "0" {
          me, _ := strconv.ParseInt(preMsgTime, 10, 64)
          preTime = time.Unix(me, 0)
          if preTime.Year() > 9999 {
             preTime = time.Unix(me/1000, 0)
          }
      } else {
          preTime = time.Unix(0, 0)
      }
      ```
3. 调用  `models.FindMessageSendEventByUserIdAndToUserId(userId, toUserId, preTime )`  函数查找消息发送事件表中指定用户发送给指定用户的记录，并将结果存储在· messageSendEvents· 变量中。如果查找失败，则返回空的  []models.MessageDVO 和错误信息。（若`pre_msg_time` 不为 "0" 则不调用此函数，因为后续新增的消息需要剔除本人发的）
4. 调用  `models.FindMessageSendEventByUserIdAndToUserId(toUserId, userId, preTime )`  函数查找消息发送事件表中指定用户接收到指定用户发送的消息的记录，并将结果存储在 `messageSendEventsOpposite`  变量中。如果查找失败，则返回空的 []models.MessageDVO 和错误信息。
5. 将 `messageSendEvents` 和 `messageSendEventsOpposite` 合并成一个列表
6. 用多协程并发将排序后的消息数组组装成 `MessageDVO` 数组
7. 对`MessageDVO` 数组按照 `CreateTime`时间戳升序排序

MessageDVO 的结构如下 

```Go
type MessageDVO struct {
    Id         int64  `json:"id,omitempty"`
    ToUserId   int64  `json:"to_user_id,omitempty"`
    UserId     int64  `json:"from_user_id,omitempty"`
    Content    string `json:"content,omitempty"`
    CreateTime int64  `json:"create_time,omitempty"`
}
```

> 聊天记录功能的缓存优化方向 （暂未实现）：
>
> 1. 使用 redis 的 Sorted Set 数据结构， key 为 userId 和 toUserId  的组合，表示两个人间的聊天记录，value存储聊天记录的集合，score 取每条聊天记录的时间戳，每个集合设置过期时间
> 2. 分页请求加载聊天记录，redis 的集合中只储存部分聊天记录，redis 中的聊天记录取完再从数据库中查询，区别对待冷热数据，由于时间久远的聊天记录使用的频率较少，可以不存放在redis中
> 3. 新增聊天记录时往 redis 集合中 SADD 聊天记录，若对应的 key 不存在，则需要取出最新的20条聊天记录重建这个缓存集合

  

# 四、测试结果

## 功能测试

#### User接口

##### 用户注册 /douyin/user/register/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MDUyNDE0MWFlMmRiMDU1Zjg1ZDNlMGM5NDA3OTYzYzRfYUlZSURxMGZQcmlBYm1XcTRsUjI2ajFFWW96NVd2cXBfVG9rZW46SWNUNGJVV0pXb0xEb2p4aThBZmNmaUVMbnloXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 用户登录 /douyin/user/login/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NjUyMDU5ZDdkNGE0YTE3NmQ4MmVhNjlhNTAyNjMzYTRfN1VEcVd2ckIwSk1wUDFzcEdnWWRxZzVaOTdFR3lZVU1fVG9rZW46S2hwMWJ3WUwxb0RXa3J4V0Z2TGMzV002bjBUXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 用户信息 /douyin/user/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjBjNjA1OTBiYjM3MzEwNGZjNGJhYTJmNGYyOGVjNTBfT0JKTVdRdFdMbExtakFnWW1wMFBrbFhRQ1lkMVNNckxfVG9rZW46VHozZ2JGb0Vjb3hsMTl4UVRhRmNuSTM5bldmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### Video接口

##### 视频流接口/douyin/feed/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NThhZmM1NTFiMDQwZTY1ZDBlYjg0OTgwODM3ZTcxYjFfU3R2TjFzRUdBa3lQZEV3SFlNWTVDdHlJQmdSa0hCN0tfVG9rZW46SnoyN2IyRjdjb1E3UEx4T1A0YWNENm9WbkhlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 投稿接口 /douyin/publish/action/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=Mjg5YmVjMjAwNTY0MmIwMTE2ZjU2YjYxZDg0OTFhYmFfWHdOVXRrcWZHdUE4NFh6OWVsT0lJc0o3YkF2bkptcEZfVG9rZW46V2Q5aGJTTEhhb0RHWUV4S1lpZmM4eEQ3bmdnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 发布列表 /douyin/publish/list/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZWMwYmVmZGU4MWM4YTVjZmQ1ZDAwMGEwMGNlMzE4NjdfNU1WaFFZaHRySHpPQjZZQUMxZ3NxZXpyYWVPeHlJa1RfVG9rZW46SHRiaGJJSUxvbzNzQkx4UHpFd2NYdG0yblNoXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### Favourite 功能

##### 赞操作 /douyin/favorite/action/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=N2ZlMjg2YjI1MTE2MzE3NWZlM2QxMjM5OWMxMWM4ZDBfRmtHVGZHZVRaVDJWUllYajlpcERlcFZIZ2xrb1JkTFFfVG9rZW46QVhyMmJuZ3hVb1RJYnV4eTFQVmNJWDNsbkVjXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 喜欢列表  /douyin/favorite/list/

只有该接口的status_code为string类型，其他接口均为int类型，于是我们统一全为int型

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NzUxOTE5OGViMTAzZTBiZDhkNGQwMmIyN2JkZjljNTFfcE9iUkZYa0ZweVg3MjNFanYyZ2tpWnZldnRkcnpkQjFfVG9rZW46QTc1V2JGRlFwb0RjQzd4dGtleWNmY293bkxlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### Comment 功能

##### 评论操作 /douyin/comment/action/

- 发表评论

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YTFkZDM5NDViYTBmNWQwMTczNjljMzQzZDlkYzYwNzNfcmR0ajJEa1JUelJKbExyT3JLRXFJOWZOSUNWOWhGTGlfVG9rZW46WmNxOWJCWlpFb2J3R1p4ZTJObmNmSTNxbm5lXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

- 删除评论

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YzQyYmE5MTkzODdmMjNkNTcwYTRjNjc1OTc1YjQ3ODlfakVueU1zNTBjNG5FZWdHVFh4YTFRdXRhNkpXczJvaW5fVG9rZW46S2ZndWJMd2JGbzFuaDZ4eXlET2NIQ2M1blBmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 评论列表 /douyin/comment/list/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=OTQ4ZmRlZGIwYzllNjIwZjdmYjU3NGE2Y2RkMTU1ODlfWnd1QjQ2MWkya1JVOGoyS1ZacThHTVhVWFNLOVVGbjRfVG9rZW46UWgydGI1V3VZb240bmF4Y05uSmNFV1FzbkljXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### Relation 功能

##### 关注/取消关注接口  /douyin/relation/action/

- 关注

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NGNkOWVhY2JhOGYxNWVkNWFkYTk3NWNlZDI2MWJhMTFfR1J3ZkNoZ2JKUFNWcUkxdkdUYTRqdFExNm9US2pRcUZfVG9rZW46T0p6cGJsdDNSbzY4a0F4UTI1V2M0V1hIbldoXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

- 取消关注

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YjljMjA4Y2NlYjg1NmIxZWJjNzEyMjg0NDU5YzIxMDFfUllMMUNveHlYb040MzNqWmRuTnNhRTJQZEpZbnJxYXdfVG9rZW46WGNIY2JFRzFvb2lWaFJ4ZDE1emN5MUFFbmhmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 关注列表/douyin/relation/follow/list

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MjdiMzQwOTU3ZTQyM2Q4ZjMwNTAxMWNiYTRiYzczZGVfdEpYVHdhYzU2VkNSYmtpaUtCaFVZcXhNWEU4Nlo1NHZfVG9rZW46UlhwZmJ2d0Qxb2FYRkh4ODR2WWNXRnRnbjVwXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 粉丝列表  /douyin/relation/follower/list/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZmNkMzMzZmNlNzhlYzQyOTQ3ZGY2MDhmMTQ0NGU2YjRfS05ENm9mWnZ5ZDF4Z0luUmhCTzN1NTh3V00wbmJINWxfVG9rZW46T1pWUmJUbkhqb1F4eTJ4SGdtSWNyaGdBbnVoXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 好友列表 /douyin/relation/friend/list/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NWQyZmFjMjE4NTA1M2Q5OTdlZTU0MzgwYzBjZGE1YjdfWklwUFZ5V0NvVGI2SlRmZ1N3UmRCcElhRGV6Uk9ESkZfVG9rZW46QTNKdWJvVVk2b1VTSFV4MU80ZGNCN2pabjRlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### Message 功能

##### 发送消息 /douyin/message/action/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MjM3NjU5Yjc4ZWU1MWU4ZGJjODRkY2M1MTE0ZGUwZTBfSkhUZlJRUlRHRnVkd2NUYTVUaTM4Vng3am9jaHliUmhfVG9rZW46VWZqTmI2czh2b0hlQVN4THJqN2NGWDJYbm5jXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

##### 聊天记录 /douyin/message/chat/

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NjViNTg4NmFkZjNiY2YzZTY3NTg4ZDEwZmRlYjhhZmJfMUF4RFJxRU5PM0dUV3FpWHIwWlZBZHRCWXE5TGNjOVNfVG9rZW46RFl6YmJwamdYbzhuMEF4UWtNQWMzM1ltbmxiXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

## 性能测试

测试硬件条件：

(**本项目，****etcd****)**，**redis ，****jmeter** **，（****mysql****,rabbitmq）**分别部署在**四台****云服务器**上，均为 centos8 操作系统

其中 redis ， jmeter , 本项目 所在的服务器均为 **32G** **内存** **， 8核**，这三台服务器间的带宽为 2G 。

Mysql 所在的服务器配置是 **2G****内存** **，2 核**

所有测试均使用 centos8 下的 jmeter 发请求

（由于受到本地Windows系统的配置，带宽和网络情况影响，故使用和项目部署的云服务器同一个子网下的一台 32G 内存 8核 的centos8 服务器上部署 jmeter）

#### 测试项目1：高并发随机点赞/取消点赞测试

**目的**：在高并发下测试点赞功能是否还能保证幂等性，是否会出现数据不一致

1. 使用python脚本生成 一个 csv 文件，每一行是 1 或者 2 （随机出现）（代表请求的点赞或者取消点赞），jmeter 每个测试线程依次读取然后拼接得到 http 请求。

jmx文件部分内容

```YAML
          </elementProp>
          <stringProp name="HTTPSampler.domain">ip</stringProp>
          <stringProp name="HTTPSampler.port">8080</stringProp> 
          <stringProp name="HTTPSampler.protocol">http</stringProp>
          <stringProp name="HTTPSampler.contentEncoding">utf-8</stringProp>
          <stringProp name="HTTPSampler.path">/douyin/favorite/action/?video_id=7097084071494288384&amp;action_type=${type}&amp;token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NzA5MDMwNjQxMDkzOTk0MTg4OCwiY3JlYXRlX2RhdGUiOiIyMDIzLTA3LTI3VDIwOjI2OjIwWiIsImlzX2RlbGV0ZWQiOjAsIm5hbWUiOiIyMDIwMjIzMTAxNEAxNjMuY29tIn0.Ja5dz5k47VSON2pgsGrMEzbKUCka6j_4p-ytoga1iRE</stringProp>
          <stringProp name="HTTPSampler.method">POST</stringProp>
          <boolProp name="HTTPSampler.follow_redirects">true</boolProp>
          <boolProp name="HTTPSampler.auto_redirects">false</boolProp>
          <boolProp name="HTTPSampler.use_keepalive">true</boolProp>
          <boolProp name="HTTPSampler.DO_MULTIPART_POST">false</boolProp>
          <boolProp name="HTTPSampler.BROWSER_COMPATIBLE_MULTIPART">true</boolProp>
          <stringProp name="HTTPSampler.embedded_url_re"></stringProp>
          <stringProp name="HTTPSampler.implementation">HttpClient4</stringProp>
          <stringProp name="HTTPSampler.connect_timeout"></stringProp>
          <stringProp name="HTTPSampler.response_timeout"></stringProp>
        </HTTPSamplerProxy>
```

1. 使用 jmeter 进行一分钟持续压测

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YjY2NzE2Njk5ZjM4OGZkMWY0NTU2OWIzNGM4OGM5NWRfTE9BZmN6cEFNY3JLbEY4MENGV0cwaVNUMzlKMVIyWnpfVG9rZW46U1JmRmJZN3VZb0x3T1p4Q1FnMmNJY3lvbkl4XzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

1. 得到测试报告如下

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NmU5MWIwMDRkY2E5MjI1MjVjOTkwMWQ5NGM0NTFmOTlfSHVmbnR5a28yTzJnQ3VWUDFRdmJySXBrYlRVcWk5QjZfVG9rZW46TmtsemJ5b1Z1b3lSRVp4QVBoRGNXRG9ubkFnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

可见在每秒吞吐量 18000 以上的巨大压力下，平均响应时间仅为 25ms 左右，可见 消息队列 的异步操作保障了用户的流畅体验感。对比测试前后的 user ， video ,  like 表的数据，并对测试的csv文件计算得到的数据进行比较

发现，三张表数据一致性没有被破坏，而且没有出现连续多次点赞/取消点赞的执行，结果保持了正确性

#### 测试项目2：高并发的查询请求下 执行取消点赞操作是否会造成 redis 和 mysql 数据不一致现象

部分 jmx 文件内容

```YAML
</elementProp>
          <stringProp name="HTTPSampler.domain">ip</stringProp>
          <stringProp name="HTTPSampler.port">8080</stringProp>
          <stringProp name="HTTPSampler.protocol">http</stringProp>
          <stringProp name="HTTPSampler.contentEncoding">utf-8</stringProp>
          <stringProp name="HTTPSampler.path">/douyin/favorite/action/?user_id=7090306410939941888&amp;token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NzA5MDMwNjQxMDkzOTk0MTg4OCwiY3JlYXRlX2RhdGUiOiIyMDIzLTA3LTI3VDIwOjI2OjIwWiIsImlzX2RlbGV0ZWQiOjAsIm5hbWUiOiIyMDIwMjIzMTAxNEAxNjMuY29tIn0.Ja5dz5k47VSON2pgsGrMEzbKUCka6j_4p-ytoga1iRE</stringProp>
          <stringProp name="HTTPSampler.method">POST</stringProp>
          <boolProp name="HTTPSampler.follow_redirects">true</boolProp>
          <boolProp name="HTTPSampler.auto_redirects">false</boolProp>
          <boolProp name="HTTPSampler.use_keepalive">true</boolProp>
          <boolProp name="HTTPSampler.DO_MULTIPART_POST">false</boolProp>
          <boolProp name="HTTPSampler.BROWSER_COMPATIBLE_MULTIPART">true</boolProp>
          <stringProp name="HTTPSampler.embedded_url_re"></stringProp>
          <stringProp name="HTTPSampler.implementation">HttpClient4</stringProp>
          <stringProp name="HTTPSampler.connect_timeout">4000</stringProp>
          <stringProp name="HTTPSampler.response_timeout">4000</stringProp>
        </HTTPSamplerProxy>
```

1. 开启 jmeter 进行一分钟压测，同时使用 apifox 发送多次 点赞/取消点赞请求
2. Jmeter 压测报告如下

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZWYxYmYyOWMyNmQ4NWI4YmQwNGQ4ZmQwNzgzYWQxODlfU1hTck1SdVF4MmhsRnllblNSWERnU3B1MXJFNjNZSEpfVG9rZW46VjJQbmJidWhXb1BXOW14M0RRVmNYdlNDbkVnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

在每秒超过 17000 的吞吐量请求下 ， 观察 mysql 和 redis 中的点赞/取消点赞情况保持一致，没有出现数据不一致和脏数据

#### 测试项目3：高并发查询用户信息压测

Jmeter 压测报告如下

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=YmI4YjU5N2UyYTczZGY0ZWNiNGRmOWJmM2MzNTE0OTNfaGlDNnV0NDlBV29UOUh6cDI1M1JTazJJUFRPZjdVWjJfVG9rZW46VmRRb2JNVjRVb1Q2WER4aG00TWNKS2NhbkdjXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

#### 测试项目4：评论功能压测与优化

初始的评论部分代码是这样的

```Go
func (commentService CommentServiceImpl) PostComments(comment models.Comment, video_id int64) error {
     var req pb.UserRequest
     req.UserId = comment.User.Id
     resp, err := rpc.UserClient.GetUserById(context.Background(), &req)
     if err != nil {
      return err
     }
     user := BuildUser(resp.User)
    bloomFilter.BloomFilter.Add([]byte(strconv.Itoa(int(comment.Id))))

    toMQ := models.CommentMQToVideo{
        CommonEntity: comment.CommonEntity,
        ActionType:   1,
        UserId:       user,
        VideoId:      video_id,
        Content:      comment.Content,
        CommentID:    -1,
    }
    mq.CommentChannel <- toMQ
    return nil
}
```

进行了压测后结果非常不理想，如下图所示

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZGY0NzQ0Y2MxN2U5NzRmNDM5YjA3YjY0NWZmNDg4NDVfUWVKNFFuR3RDVnJVd1Bab1dQWWhuSG9jTjhtRVc0NXdfVG9rZW46RmR2YWJTd0xtbzQ3Y1J4WUhRMGNVTHkybmRlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

当每秒处理请求数到 150左右时就出现了 error ， 而且平均时延高达 1.6秒，最大时延超过了30秒 。这显然是难以接受的，这将给用户带来极差的体验。

后来经过很多次的排查和压测，发现了速度瓶颈在于  mq.CommentChannel <- toMQ ，将消息放入 channel 的这行代码。在请求高峰下，并发向 channel 添加数据导致性能急剧下降，而且消费速率远低于生产速率，channel 的内的消息数很快到达缓冲容量，导致阻塞发生。

后来考虑使用 创建协程 异步加入 channel 的方法来降低用户时延，使得这段高时延的情况不被用户感知

```Go
func channelPut(msg models.CommentMQToVideo) { 
    mq.CommentChannel <- msg
}

func (commentService CommentServiceImpl) PostComments(comment models.Comment, video_id int64) error {
     var req pb.UserRequest
     req.UserId = comment.User.Id
     resp, err := rpc.UserClient.GetUserById(context.Background(), &req)
     if err != nil {
      return err
     }
     user := BuildUser(resp.User)

    bloomFilter.BloomFilter.Add([]byte(strconv.Itoa(int(comment.Id))))

    toMQ := models.CommentMQToVideo{
        CommonEntity: comment.CommonEntity,
        ActionType:   1,
        UserId:       user,
        VideoId:      video_id,
        Content:      comment.Content,
        CommentID:    -1,
    }
    go channelPut(toMQ)
    return nil
}
```

优化后的压测，用户感知的性能得到很大的提升

也起到了在评论流量洪峰下对于用户体验的保障作用

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=Yjg1OWM1NTZmNTIzNDVkNDY0YjgyZDFhNmYxMDQ2ZDZfbkNGa0Z0VFduZWROMTEwSFM3NGt3SXZoTWdhTFJPZmhfVG9rZW46SzR5Y2IwdFg3b3RJV1F4SEVMS2NOdDMybjVmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

经过了比较长时间后，数据库才更新完毕数据，由于有了 channel 的削峰，数据库能在低水平的压力下运行。避免了数据库的压力突增导致的数据储存失败甚至宕机的出现。

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=OTJmZWJlZTIyNmZlOTU1MGM1NjExMDliNWZlOThkZjBfSTE3Rk12R0cwQkVvMDJnc096UXBnWHhEZThTY3FkVTFfVG9rZW46TmhCVWJpbVFXb0luRVV4ZDNIeGMxWkg0blBnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

查看数据库 所在服务器  (2G 2核) ， 并没有造成巨大的压力

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=NmZiYTIwZjI5ZDY0YWQ1NmJjMmFmZjIxYzI5NDlkYTNfQmt2QWZVMFM5RUh0UUhoV3JkR0Zyd2xmVURMbXlucXFfVG9rZW46UHhrYmIzaEl3b2pOSGN4bGtrTGNUNGFFbktlXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

在线上运行时，若遇到评论量突增，用户评论的体验稳定性也能得到保障，为后续的扩容处理争取时间。在检测到 channel 数据压力增大时，可动态地对数据库集群进行平滑扩容从容应对。避免了突增的请求直接将数据库击垮的危险。

#### 测试项目5 ：索引优化查询性能对比

在关注相关的接口中需要调用 `FindByUserIdAndVedioId()` 函数，该函数通过 GORM 构造的一个 SQL 查询例子是

```SQL
SELECT * FROM follow WHERE user_id = 517753093834674727 AND follow_user_id = 264013038577388680 AND is_deleted = 0;
```

在经历了多次高并发测试后， follow 表内有至少 **110 万行**数据 。对该表执行以上 SQL 语句后发现耗时超过 0.4s

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=MzMwYjBhZmY0ZjA2OTBlYmUwZTQ3YzhkZWVhZWJlZTRfcWhEZFo2RGNFRVFJbFI0MHFnNHRxZmllbXpmcEFCemlfVG9rZW46SzZxR2IxWjFZb2M1dHB4bHFqdWNFZWpPbkNnXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

当给该表的 `user_id` 和 `follow_user_id` 添加了联合索引 user_and_follow （数据结构为B+树）后再次执行该 SQL 语句测试，耗时仅为 0.016秒。速度**提升了近30倍**。显然经过了我们的索引优化后，在大数据量下性能大大提升。

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=OWU4M2UzM2M2ODdlNjhjY2M0OWUyNzA1ZGI0M2Y2NzVfcHJKZTUxQXJhRm9QekxNV0NKZlNJZU9KVEdTUkJMMGpfVG9rZW46UWdaeGJRb0Jwb296ZlR4ZTFPUGNkY1Rqbng2XzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

通过 explain 查询执行计划时，也证明该 SQL 使用了联合索引

![img](https://da0sq9guct3.feishu.cn/space/api/box/stream/download/asynccode/?code=ZGY0OTNmMWRhOGJlYTRiZWIzY2I5MmU3MDQ5Yzk3Y2VfUVRJTUJvVkNHODhGWk5xMVQzeDNTcmFTWU9KRXkyb21fVG9rZW46UjZid2JTNG5Wb1ZsUWh4RjE5SGNuUFZGbjBmXzE2OTM1NDA0NTQ6MTY5MzU0NDA1NF9WNA)

# 五、Demo 演示视频 

**请开启声音观看**

暂时无法在飞书文档外展示此内容

**下面是单独演示聊天功能的视频**

暂时无法在飞书文档外展示此内容

# 六、项目总结与反思

## 目前存在的问题和优化项

1. 目前使用 `setnx` 实现分布式锁在 redis 分片集群或者主从集群中可能会失效，并不是一个高可用的分布式锁。当两个 redis 结点没有及时进行数据同步可能导致另一个线程在另一个结点再次获取了 锁 引发并发问题。因此后续考虑采用 RedLock 算法改进或者使用 Zookeeper 等其它组件实现分布式锁
2. 本项目中可能存在的 `大key问题` . 如在本项目的点赞功能中，如一个用户点赞了成千上万个视频，那 redis 缓存中的 like set 会出现有数千个元素的集合，由于 redis 是单线程的，因此要避免大 key 的出现 。改进方案是进行大 key 的拆分。可以把一个 like set 拆成多个集合。 而且很久以前点赞过的视频可能再也不会刷，这部分数据变成了冷数据，不应该在 redis 中占用空间。所以要考虑将这部分的点赞记录不存放于缓存中。
3. 本项目中未采用链路追踪技术，而路径追踪有助于分析和监控微服务间的通信
4. 本项目中，服务内涉及的异步削峰操作没有使用 rabbitmq 进行，而是使用 channel ，虽然性能开销较小，但是却导致并发量受限和存在无法可靠地存储大数据的问题。Channel 没有提供类似于消息持久化、重试、消息确认等高级的消息队列功能，在数据可靠性上不如专门的消息队列。由于在本项目中，点赞，评论，关注等更新本服务的数据为一次性操作且失败率较低，对于性能的追求高于可靠性，故没有像服务间那样使用 rabbitmq 消息队列，而是使用 channel .
5. 使用 Kafka 和 ES 搜集分析日志数据
6. 未实现自动化编译部署方案

## 架构演进的可能性

1. 在针对多机房分布式部署的情景下，采用 Redis 分片集群的方式部署
2. 由于时间原因，目前的 mysql 数据库没有进行水平分表，  随着用户数量的增长，数据量会变得很庞大，因此要考虑进行水平分表和读写分离。
3. 后期在项目中可以引入 大数据分析功能，使用消息队列对数据进行离线处理，分析用户画像，并训练使用推荐算法针对不同用户推送个性化的视频
4. 引入 Hbase 存储好友关系和社交网络。HBase 支持范围查询和前缀查询等功能，这对于查询某个用户的好友列表、共同好友等功能是非常有帮助
5. 由于时间原因， 本项目中未采用负载均衡技术。但是在真实场景中，抖音后端将面临庞大请求压力，故需要使用 Nginx 进行多级负载均衡。
6. 考虑使用 `channel 池`的思路对使用 channel 进行削峰进行优化，一开始创建若干 channel ，生产数据的时候放入其中一个 channel , 以此提高稳定性，避免出现单一 channel 数据积压导致阻塞等问题
7. 后续考虑引入 GPT （如：讯飞星火大模型） ， 通过对讯飞星火 API 的调用，实现给用户提供**语音识别输入功能（讯飞星火语言识别API），视频合规鉴别功能（讯飞星火视频合规API），智能生成视频标题，评论功能（GPT按照关键字生成）， 根据用户对视频的评论的褒贬程度影响后续视频推荐 （讯飞星火****情感分析****API）**

## 项目过程中的反思与总结

**项目过开发程中遇到的问题**

1. 在 user 的点赞操作消费函数中，存在一种情况是 用户自己点赞自己，由于开启了事务，而且需要更新 TotalFavorited 和 FavoriteCount 两个字段。起初发现这样同时更新同一行事务提交后只记录最后一次的更新。之后发现，是因为事务的默认隔离级别是 `可重复读` ，由于MySQL 的 MVCC 机制，两次更新前的读取都只读取到事务开始执行时的快照。
   1.  修改后的函数如下，设置了对是否更新同一行的判断

   2. ```Go
      tx := utils.GetMysqlDB().Begin()
      //获得当前用户
      user, err := model.GetUserById(userId)
      
      //查询视频作者
      author, err2 := model.GetUserById(data.AuthorId)
      if err2 != nil {
          panic(err2)
      }
      actionType := data.ActionType
      
      if actionType == 1 {
          //喜欢数量+一
          user.FavoriteCount = user.FavoriteCount + 1
          //如果是同一个作者，在同一个事务中必须保证针对同一行的操作只出现一次
          if user.Id == author.Id {
             user.TotalFavorited++
          }
          err = model.UpdateUser(tx, user)
          if err != nil {
             log.Println("err:", err)
             tx.Rollback()
             panic(err)
          }
          if user.Id != author.Id {
             //总点赞数+1
             author.TotalFavorited = author.TotalFavorited + 1
             err = model.UpdateUser(tx, author)
             if err != nil {
                log.Println("err:", err)
                tx.Rollback()
                panic(err)
             }
          }
      
      } else {
          //喜欢数量-1
          user.FavoriteCount = user.FavoriteCount - 1
          //如果是同一个作者，在同一个事务中必须保证针对同一行的操作只出现一次
          if user.Id == author.Id {
             user.TotalFavorited--
          }
          err = model.UpdateUser(tx, user)
          if err != nil {
             log.Println("err:", err)
             tx.Rollback()
             panic(err)
          }
          //总点赞数-1
          if user.Id != author.Id {
             author.TotalFavorited = author.TotalFavorited - 1
             err = model.UpdateUser(tx, author)
             if err != nil {
                log.Println("err:", err)
                tx.Rollback()
                panic(err)
             }
          }
      }
      tx.Commit()
      ```