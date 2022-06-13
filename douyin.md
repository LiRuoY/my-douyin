# 项目

## 技术

- gin
- jwt4
- gorm
- viper配置文件
- mysql

## 包
![image-20220613102408437](https://user-images.githubusercontent.com/98474722/173273838-3c8b7f4e-ac08-4524-9b1c-e1c1ddbda2e3.png)


## 架构
![image-20220613100845934](https://user-images.githubusercontent.com/98474722/173273754-e1af5a69-e420-447b-b782-a0dc47ddd137.png)


## 技术点

- 数据库表的设定
    - 一对一
        - user----userLogin
    - 一对多
        - user----videos
        - user----comments
        - video----comments
    - 多对多
        - user-----favorVideo-----表favor_videos
        - user----follows(关注的人)----表follows

- 视频存放
    - static：存放所有用户的视频
    - 文件夹：每个用户在static下都有一个名为其id的文件夹，创建时机为第一次上传视频时
    - 文件名：采取年月日+filename的方式
    - 数据库存放该视频在本地的url，如http://192.168.93.7:8080/static/1/202206091.mp4

- 时区
    - time包解析字符串--->时间，如果不用gorm约束的CreatedAt和UpdatedAt要格外注意从数据库解析时间字符串
        - Parse函数，需要layout和timestr都有时区，默认返回utc时间
        - ParseInLocation函数，使用本地时间，或者指定时区
    - gorm的时间配置parseTime=True&loc=Asia%2fShanghai
        - 使用约束的CreatedAt和UpdatedAt，自动把当前时间存入数据库，把数据库时间字符串转为time
- 密码
    - 采用标准库的md5加密
- dao操作
    - 使用Preload和Association方法操作表关系
- token
    - 按顺序从query，form，header中拿去token，有一个有就可以
    - 按顺序从query，form拿取userId，用来和token的useId校验，没有就不用校验
- 功能的细节
    - user
        - name和password长度都是  0<len<=32
    - 点赞
        - 可以点赞自己的视频
        - 未点赞才可以点赞，点赞了才可以取消点赞，否则返回err
    - 评论
        - 可以发多条评论
        - 对title进行长度限制，0<title<=32



# jwt

- 不用jwt提供的middlerware了，自己写了一个

# 功能

1. 用户注册✔️
2. 用户登录✔️
3. 用户信息✔️
4. 投稿接口✔️
5. 发布列表✔️
6. 视频流接口✔️
7. 赞操作✔️
8. 点赞列表✔️
9. 评论操作✔️
10. 评论列表✔️
11. 关注操作✔️
12. 关注列表✔️
13. 粉丝列表✔️

# cmd

## api

### handlers

- 一个handler专门放在一个文件，用来SendResoponse，处理service返回的err
    - common-response-----status_code   status_msg
    - 由于不确定到底能不能在Resoponse里把所有数据都写成data，所以各自定义
- 
- 完蛋，因为想用jwt的middlerware练练，现在搞得有点乱，好像这个middlerware是多余的一样
    - middlerware主要用在登录，所有验权操作，其实也可以自己写一个这个很小的框架就3个操作ReleaseToken，ParseToken，Authenticate
    - 因为项目的一些要求，又自己写了ReleaseToken，ParseToken，Authenticate
        - register要返回token，ReleaseToken
        - 推送视频用表单，这个中间件不能从表单自动验权，ParseToken，Authenticate



1. 都先解析request-json，调用service处理request，调用SendResoponse发送resoponse-json
    1. register   post   ✔️
    
    2. login   post   ✔️
    
    3. user-info   get  query  ✔️
        1. 只用query user_id,再把转为int64的id传给service
        
    4. 投稿接口 post    multipart/form-data✔️
    
        1. 因为token在表单，所以要用自定义的方法来验证，jwt不能用表单来验证
        2. 检查视频格式
        3. 保存视频在服务端本地，在static所在的目录
            1. 要先创建专属于这个user的文件夹
            2. 规则都写在静态资源里面了
        4. 再调用srevice
    
    5. 发布列表  get   query✔️
    
    5. 视频流接口 get✔️
    
    5. 赞操作  post✔️
    
    5. 点赞列表  get ✔️
    
       1. 因为response和video-list是一样的，所以用video里的Send
       
    9. 评论操作   post✔️
    
    9. 评论列表 get✔️
    
    9. 关注操作  post ✔️
    
    9. 关注列表 get✔️
    
    13. 粉丝列表 get✔️
    
         
    
         

## db

### dal

#### service

1. register		✔️——————因为jwt中间件没有register功能，所以要自己写一个
    1. request-name,password     
    2. response--id，token   
    3. 参数校验，0=<name<=128  password!=0
    4. 在db根据name查询是否存在
        1. 不存在就可以在db创建，保存id
        2. 存在返回已存在错误	
    5. 颁发token

2. login  ✔️——————因为用了jwt中间件的Loginhandler所以不用这个了
    1.  request	name,password   
    2.   response    id，token
    3. 参数校验，0=<name<=128  password!=0
    4. 在db根据name查询，再验证密码
        1. 成功返回id
        2. 失败直接返回错误	
    5. 颁发token




这下面都是身份认证后的操作，即request都带token

1. user-info  ✔️
    1.  request	  name,password  
    2. response    返回user的id,name,follow_count,follower_count,isfollow，即user表
    3. 在user表查询这个id就可以了，再返回
2. 投稿接口✔️   
    1. 这里很恶心，因为token在form，要自定义一个parseToken来验证
    2. request	data--file,token,title
    3. response     只返回common-response
    3. 关系————这里是User has many Videos
    4. 把在这个file的路径和title保存在数据库就好了
        1. 在本地的保存形式，在static文件夹下，创建以id开头的文件夹，再创建video文件夹，picture文件夹，把用户的视频，图片存在这里
3. 发布列表✔️
    1. request	token，user_id
    2. response   user_id的所有video
    2. 关系————这里是User has many Videos
    3. 根据user_id去db拿视频就好了，把视频切片单独拿出来，再把user放入每个视频的Author
    3. publish-list 要根据自己有没有对自己的视频点赞来设置isFavorite为true
4. 视频流接口✔️
    1. request	latest_time--int64 （--精确到秒的时间戳，如果没有这个，就是当前时间）,token
    2. response  video_list(最多30个video,按投稿时间倒序--晚-早)和next_time(返回视频中最早视频的时间--精确到秒的时间戳,int64)
    3. 先对created_time进行降序，拿走>=latest_time的视频，最多30个
        1. 用map对在videos里拿到的user_id去重，value为struct{}，节省空间
        2. 再把user全取出来，再把user放入video的author里
    4. feed不应该有isFavorite
5. 赞操作         ✔️             
    1. request		user_id,token,video_id,action_type(1-点赞,2-取消点赞)
    2. response    commonResponse
    3. 关系————这里是User和Videos是many2many，表为favor_videos
    4. 直接用Association给User和Video在favor_videos加关系，再给video的FavoriteCount加1就好了，
        1. 但是，根据要求，当在返回videoList是有IsFavorite字段的
            1. 处理方法：Video的IsFavorite在数据库可以不用管，当要IsFavorite的时候，把video在数据库取出来，再由判断FavoriteCount是否为0，来决定是否填上true就好了
        2. 已点赞，不给再点，未点赞不给取消点赞
            1. 先在favor_videos判断是否有这个user_id,video_id
                1. 试了很久，发现Find，row，rows（但row，rows的scan会返回取不到值的err）不会返回找不到err，而first，take会，
                2. 因为没有表的结构，所以用table不用model，所以只能用take了
            2. 在favor函数里是用swith做favor和cancelfavor的
        3. 暂时自己可以给自己点赞
    5. 另外，点赞和取消赞应该开启事务
        1. 因为如果对点赞+1成功，但添加在favor_videos的关系失败时，就应该全部失败，反过来也一样
6. 点赞列表  ✔️
    1. request	 user_id,token
    2. response   用户点过赞的列表，注意isFavorite应该为true
    3.   直接拿user_id去db查询，取出来后再一个一个把isFavorite设为true
    4. 把video的author也找出来填上去
7. 评论操作✔️
    1. request	 user_id,token,video_id,action_type(1评论，2删除),comment_text,comment_id
    2. response     返回当前操作的评论，包含id，user，content，create_dare日期为mm-dd，
    3. 再直接在db里的comment创建或者删除，并返回这个comment
        1. 创建时，因为user_id，video_id都是外键，所以，如果没有这个user_id，video_id，那么创建失败
            1. 创建完再查询author
            2. 如果在comment里已经有了，那么根据需求自定义可以评论几次
            3. 创建后会自动返回id
        2. 删除时，也一样，如果有这个comment就删除，如果没有这个comment，应该返回err
            1. 先查询发过这条comment的作者，如果没有就直接返回（因为delete空是不返回err的）
            2. 所以得先查再删
        3. 以上都得把当前的comment返回
        4. 因为video有commentCount，因此要开启事务
    4. 再把user塞进comment里
    5. 可以发布多条评论（暂时不限制），删除只能删除一次（所以得先查询）
8. 评论列表✔️
    1. request	 video_id，token
    2. response     返返回视频所有评论按时间倒序---晚--早
    3. comment里也包含了user
9. 关注操作 ✔️
    1. request  user_id,token,to_user_id,action_type(1-关注，2-取消)
    2. response    common-response
    3. 这里是many2many关系，表为follower_relations
    4. 由于user里有FollowCount，所以要开启事务
    5. 这里基本完全和点赞差不多

10. 关注列表✔️
    1. 这里要注意is_follow要为true

11. 粉丝列表✔️
     1. 这里，如果user对粉丝也有关注，也要显示为true





# bug

- 想在判断user是否存在，找不到时返回err，又少写语句
    - DB.Where("name=?", name).Error，无论如何都不会返回err
    - DB.Debug().First(&o, "name=?", "asd").Error，这样才会返回err，当找不到时也会
    - 上面方案都不行，改为DB.Model(&UserLogin{}).Where("name=?", name).Count(&c)，比较好
        - 因为不需要err，也不要不存在时返回err

- 

- 绑定query的bug，绑定query应该用form标签，否则不能把query的数据绑定到结构体里

- 

- **所有数字都是float64，除非用反射拿到对应的类型**

    - 在用jwt4.MapClaims拿token.Claims里面的MapClaims值时，因为这只是map，所以所有数字为float64，
    - 所以拿id时注意转为int64
    - 不要忘了c.Set(mw.IdentityKey, id)，不然可能不太好取值，当然，直接返回也是可以的

- 

- gorm设置为Local或者Asia%2fShanghai（%2f是gorm要求的），即中国

    - Parse函数返回UTC，或者在layout和timestr都有时区，返回timestr的时区
        - 这很麻烦，

    - ParseInLocation函数则只返回传入的时区，
        - 可以是time.Local，系统时间
        - 也可以是自定义，如chinaZone, err := time.LoadLocation("Asia/Shanghai")
        - 时区文件可以在go的root目录下的 lib\time\zoneinfo.zip   查看

    - 当时把用Parse解析的拿到请求了，所以错了，这个是UTC时间戳，多了8小时
        - 现在用ParseInLocation就没问题了
        - 教训：解析字符串为时间用Parse函数和ParseInLocation函数要格外注意，其他倒不用注意


- 打开文件
    - config.yml在项目下，main.go也在项目下
    - config的init函数在项目下的config包下
    - 在main里打开文件传相对路径给init函数可以用config.yml
    - 在init函数打开文件使用相对路径必须用../config.yml





# 麻烦事

- 自己写了一个jwt中间件

- ❌当用jwt中间件时
    - 要额外自定义Unauthorized，LoginResponse
        - LoginResponse返回id和token
            - token就是LoginResponse参数message
            - 但是id就不好拿了，id可以用设置IdentityKey为"id"，这个可以通过ExtractClaims，c.Keys，c.Get()拿取，但是，这3个都是在Keys这个map里
                - 而且只有在middlewareImpl调用时，才会 c.Set("JWT_PAYLOAD", claims)，  c.Set(mw.IdentityKey, identity)
                - 而且middlewareImpl是GinJWTMiddleware的handler方法，是要创建GinJWTMiddleware对象后才能通过调用其MiddlewareFunc()方法返回一个函数，这个函数只调用middlewareImpl方法，这是用来做中间件方法--handler的
                - 还可以通过GetClaimsFromJWT()拿取，但这个要用GinJWTMiddleware对象调用，在设置LoginResponse时都还没生成呢
                - 于是我在Authenticator里用CheckUser返回的userID来调用c.Set(constants.IdentityKey, userID)，因为Authenticator是最先调用的，所以在LoginResponse用没问题，当用middlewareImpl时也只是重新设为userID而已
        
    - 不能从表单拿token，太恶心了，因为推送视频的request的token在表单上，为什么？？。。。
    
    - 有一个大大的bug
    
        - id和token的id不对称也会被验证成功，解决方法：使用Authorizator判断
    
        - Authorizator是对user的再次验证，func(data interface{}, c *gin.Context) bool
    
        - data是Identity的value，即id
    
        - 如果request上传了user_id,必须在这里判断这个user_id是否等于Identity
    
        - ```go
            Authorizator: func(data interface{}, c *gin.Context) bool {
                //如果没有userId则不用进行这种隐私校验
                //本来，如果没有配置这个函数，也是直接返回true
                qid := c.Query(constants.UserIdQuery)
                if qid == "" {
                    return true
                }
                queryId, err := strconv.ParseInt(qid, 10, 64)
                if err != nil {
                    return false
                }
            
                //Authorizator,c.keys["id"]=%!d(float64=1) float64
                //这是因为MapClaims只是map，因此从token字符串解析的是float64
                id, ok := data.(float64)
                if !ok {
                    return false
                }
                tokenId := int64(id)
            
                return tokenId == queryId
            },
            ```



# 领悟

- 事务
    - 做多个修改，要么都成功，要么都失败
- service使用Flow，让Flow这个结构体保存参数，不要一直在函数里传
    - 当然，这是要对参数做很多处理（可能由这个参数又参数其他参数），所以保存在Flow里，减少传参，还有方便调用
    - 当然，如果对参数做的处理极少，那么可以不要Flow，反而更好
- 查询数据库时，可以选择把对象地址给dao层，这样就不用再返回这个对象了
    - 最好还是只在service和dao层之间这么做，不要在handler和service之间也这么做







# 关于app的使用的bug

- app一直用不了，但是6-9的晚上我因为用java的tomcat查看了一下8080端口，竟然发现一直有一个applicationwebserver.exe这个应用占用8080端口，于是把它关掉，就和大家说的一样，用无线局域网适配器 WLAN的 IPv4 地址+自己开放的端口号当domain就好了，再把这个domain放在app的高级设置里
- 点赞，发评论，关注都没有user_id
    - 强行在middleware写了一个解析token获取id函数

- publish-list和favorite-list一直用不了，闪退  
- 前端点赞逻辑有问题
