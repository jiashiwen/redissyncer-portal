# ToDo

- [x] 发现节点宕机后，修改任务各项状态为
- [ ] 编写nodeservice
    - [ ] 获取节点状态
    - [ ] 删除节点，前提是节点已离线
    - [ ] 节点健康检查
    - [ ] portal节点自举，当所有portal节点全部宕机的情况下，启动节点时etcd显示节点alive，测试节点是否真正存活，若已宕机则修改状态后启动

- [ ] 任务api
    - [X] 任务创建
    - [X] 任务停止
    - [ ] 任务启动
    - [X] 任务删除
    - [ ] 任务查询
        - [ ] 查询所有任务，可自定义pagesize，通过querryid批量返回数据
        - [ ] 根据nodeid查询node上所有任务
        - [ ] 根据taskid返回任务状态
        - [ ] 根据groupid返回任务状态
        - [ ] 根据groupid返回任务状态

- [ ] etcd补充工具类
    - [X] etcd 游标