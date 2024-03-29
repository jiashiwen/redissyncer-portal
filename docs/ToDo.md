# ToDo

- [x] 发现节点宕机后，修改任务各项状态为broken
- [ ] 迁移任务到其他节点，keepalive 参数为true的任务，当节点宕机时迁移到其他节点

    - [x] 获取节点状态
    - [x] 删除节点，前提是节点已离线
    - [x] 节点健康检查
    - [ ] 节点选择器优化，避免每次轮询任务数量；节点选择指标待讨论

    - [X] portal节点自举，当所有portal节点全部宕机的情况下，启动节点时etcd显示节点alive，测试节点是否真正存活，若已宕机则修改状态后启动

- [ ] 任务api
    - [x] 任务创建
    - [x] 任务停止
    - [x] 任务启动
    - [x] 任务删除
    - [ ] 任务查询
        - [x] 查询所有任务，可自定义pagesize，通过querryid批量返回数据
        - [ ] 根据nodeid查询node上所有任务
        - [x] 根据taskid返回任务状态
        - [x] 根据groupid返回任务状态
        - [x] 根据groupid返回任务状态
        - [ ] 从工作节点获取任务的实时状态，需要直接访问server获取内存中的状态
        - [ ] 综合查询，根据任务节点、任务当前状态返回结果
    - [x] 实现TaskMigrate 任务迁移


- [ ] 补充工具类
    - [x] etcd 游标
    - [x] 在分布式环境下使用游标。具体思路是在etcd中注册游标的queryID，当本地查询不到游标的queryID时，请求etcd中的游标所在位置并发送请求
        - [x] 查询完毕清理本地cursorMap
        - [x] 定时检查cursorMap中超时的游标，最后查询时间戳超过阀值既执行过期清理
    - [ ] 游标在分布式环境中测试
- [ ] 模拟springboot的启动方式，同时支持命令行，环境变量，yaml文件配置启动参数

- [ ] 实现schedule 用于抽象定时执行动作
- [ ] 使用json-iterator 替代 标准库json，增强json解析效率

- [X] 程序后台运行调用go-daemon,为程序添加 '-d' 参数后台运行
- [X] 每次启动将 pid 计入 pid 文件
- [X] 实现命令行start stop status
- [X] 程序退出清理pid文件
- [ ] 实现config文件默认为二进制文件同级

- [ ] 开发测试环境迁移与搭建
    - [X] 基础镜像制作，包括基础编译环境的安装、docker及周边工具安装
    - [X] 开发环境搭建，包括etcd 开发环境，redis单机环境
    - [x] redissycner集群测试安装，包括etcd集群以及portal及server节点安装
    - [X] nginx代理服务器反向代理 portal
    - [ ] 测试环境，包括集群etcd ，redis3-6各个版本的单实例版本*2，redis集群版6套

* 辅助功能及特性
  - [x] http server 跨域访问能力
  - [ ] 权限系统粗粒度实现：拦截器使用checktoken middleware，每次检查token，login 根据用户名和密码发放token
