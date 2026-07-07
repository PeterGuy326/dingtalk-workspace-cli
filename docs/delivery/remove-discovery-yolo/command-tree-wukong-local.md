# wukong-local visible command tree

- binary: `/tmp/dws-wukong-local`
- generated_at: `2026-07-07T06:03:59Z`
- command_count: `974`
- help_errors: `0`

```text
dws
  agoal  # Agoal 管理
  aiapp  # AI 应用创建 / 查询 / 修改
  aidesign  # AI 设计（文生图/图生图/编辑/超分/抠图）
  aisearch  # AI 搜问
  aitable  # AI 表格操作
  attendance  # 考勤打卡 / 排班 / 统计
  blackboard  # 企业公告管理
  calendar  # 日历日程 / 会议室 / 闲忙
  chat  # 群聊 / 消息 / 机器人
  conference  # 视频会议：发起/预约/邀请入会/会中控制
  contact  # 通讯录 / 用户 / 部门 / 角色 / 人员关系
  credit  # 企业信用查询 (芝麻企业信用)
  devdoc  # 开放平台文档搜索
  ding  # DING 消息 / 发送 / 撤回
  doc  # 钉钉文档管理
  docparse  # 文档解析（PDF/图片转 Markdown）
  drive  # 钉盘文件管理
  finance  # 财务 / 发票 / 凭证 / 银行
  law  # 法律咨询与检索
  live  # 直播列表 / 信息
  mail  # 邮箱 / 邮件收发
  minutes  # AI 听记 / 会议纪要
  oa  # OA 审批 / 同意 / 拒绝 / 撤销
  pat  # 行为授权管理
  report  # 钉钉日志（OA 周报应用 / 日志模版填报）
  sheet  # 钉钉表格管理
  todo  # 待办任务管理
  wiki  # 知识库 / 空间管理 / 节点管理 / 成员管理
  yida  # 宜搭（应用 / 表单 / 流程审批）
  dws  # <service> [command] [flags]
  api  # 调用钉钉 OpenAPI (Raw HTTP)
  auth  # 认证管理
  completion  # 生成 Shell 自动补全脚本
  config  # 配置管理
  dev  # 开放平台开发者能力
  help  # 查看任意命令的帮助信息
  plugin  # 插件管理
  profile  # 组织 profile 管理
  recovery  # 错误恢复辅助命令
  schema  # 查看有限的本地 Schema（静态端点模式）
  skill  # 技能管理
  version  # 显示版本信息
    contract  # 经营合约管理
    report  # 周月报管理
    scorecard  # 计分卡管理
    strategy  # 战略解码管理
    user  # 用户目标管理
    create  # 创建 AI 应用
    modify  # 修改 AI 应用
    query  # 查询 AI 应用
    edit  # 编辑图片 — 通过文本指令修改已有图片
    generate  # 文生图 — 根据 prompt 生成图片
    generate-with-image  # 参考图生图 — 根据参考图 + prompt 生成图片
    generate-with-template  # 模板生图 — 根据参考图 + 模板生成图片
    isolate  # 抠图 — 去除背景提取主体
    upscale  # 超分辨率 — 2倍放大提升图片清晰度
    behavior  # 搜索明确的发送/创建/接收等行为记录
    enterprise  # 搜索企业内部知识内容和相关消息
    person  # 搜索企业人员
    advperm  # 高级权限管理（开关 / 角色查看与删除）
    attachment  # 附件管理
    base  # Base 管理
    chart  # 图表管理
    create  # 创建 AI 表格（dws aitable base create 的别名）
    dashboard  # 仪表盘管理
    export  # 数据导出
    field  # 字段管理
    form  # 表单管理
    import  # 数据导入
    info  # 获取 AI 表格信息（dws aitable base get 的别名）
    list  # 获取 AI 表格列表（dws aitable base list 的别名）
    record  # 记录管理
    search  # AI 表格搜索（dws aitable base search 的别名）
    section  # 文件夹与节点管理
    table  # 数据表管理
    template  # 模板搜索
    view  # 视图管理
    workflow  # 自动化工作流管理（启停 / 查看 / 列表）
    adjustment  # 补卡规则
    approve  # 审批单查询
    boss-check  # BOSS 改签打卡记录
    check  # 打卡查询
    checkin  # 签到管理
    class  # 班次规则
    globalsetting  # 全局规则设置（仅管理员）
    group  # 考勤组
    overtime  # 加班规则
    record  # 考勤记录
    report  # 查询考勤报表和结果
    rules  # 查询考勤组与考勤规则
    schedule  # 排班管理
    selfsetting  # 个人规则设置
    shift  # 班次查询
    summary  # 查询某个人的考勤统计摘要
    vacation  # 假期管理
    create  # [危险] 创建并发送公告(全员，不可撤回)
    list  # 查询用户公告列表
    acl  # 管理我的日历访问权限（共享给他人）
    attachment  # 日程附件管理
    attendee  # 日程参会人管理
    book  # 日历本管理（我能看哪些日历）
    busy  # 闲忙查询
    event  # 日程管理
    room  # 会议室管理
    bot  # 机器人管理
    category  # 会话分组管理
    chmod  # 授予 chat 高风险操作权限
    clear-all-red-point  # 清除所有会话红点（全部已读）
    clear-messages  # 清空当前用户指定会话的聊天记录
    clear-red-point  # 清除会话红点
    conversation-info  # 获取会话基础信息
    data-auth  # 授予 chat 数据读取权限
    file  # 会话文件上传
    group  # 群组管理
    group-mute  # 全员禁言 / 取消全员禁言
    group-mute-member  # 指定群成员禁言 / 取消禁言
    group-role  # 群身份管理
    hide  # 会话列表中隐藏会话
    list-all-conversations  # 分页获取当前用户的全部会话列表
    list-top-conversations  # 拉取置顶会话列表
    mark-read  # 标记消息已读
    mark-unread  # 标记会话为未读
    message  # 会话消息管理
    mute  # 会话消息免打扰
    mute-at-all  # 关闭/开启 @所有人消息提醒
    mute-red-envelope  # 关闭/开启红包消息提醒
    search  # 根据关键词搜索群聊
    search-common  # 搜索共同群（查询指定人共同所在的群聊）
    set-top  # 会话置顶 / 取消置顶（支持单聊/群聊）
    meeting  # 会议管理
    member  # 成员管理
    dept  # 部门查询
    label  # 角色查询
    relation  # 人员关系查询
    user  # 人员查询
    annual  # 工商年报
    bidding  # 招投标信息
    branch  # 分支机构信息
    cert-info  # 资质证照信息
    change  # 工商变更信息
    equity  # 股权信息
    info  # 企业基本工商信息
    ip  # 知识产权信息
    kp  # KP联系人信息
    license  # 行政许可信息
    member  # 企业成员信息
    risk  # 企业风险信息
    search  # 企业名称搜索
    article  # 文档文章
    message  # DING 消息管理
    block  # 块级编辑
    comment  # 文档评论 / 评论管理
    create  # 创建文档
    export  # 导出在线文档 (支持 docx / markdown / pdf)
    file  # 文件管理
    info  # 获取文档元信息
    media  # 文档媒体 / 附件管理
    read  # 读取文档内容 (Markdown)
    template  # 文档模板管理
    update  # 更新文档内容
    version  # 文档历史版本管理
    convert  # 解析文件内容为 Markdown
    commit  # 提交文件上传
    copy  # 复制文件/文档到指定位置
    delete  # 删除文件/文件夹到回收站
    download  # 下载钉盘文件到本地
    info  # 获取文件元数据信息
    list  # 获取文件/文件夹列表（统一入口）
    list-spaces  # 获取钉盘空间列表 (deprecated → dws wiki space list --type orgSpace/mySpace)
    mkdir  # 创建文件夹
    move  # 移动文件/文档到指定位置
    permission  # 文档节点权限管理
    publish  # 文件互联网公开发布管理
    recent  # 获取最近访问/编辑的文档列表
    recycle  # 钉盘回收站管理
    rename  # 重命名文件/文档
    search  # 搜索文件（聚合钉盘+文档空间）
    upload  # 上传本地文件到钉盘或文档空间
    upload-info  # 获取文件上传信息
    account  # 企业账户管理
    bank  # 银行交易明细
    category  # 收支类别管理
    company  # 主体管理
    customer  # 客户管理
    digital-invoice  # 数电发票管理
    gather  # 自定义经营报表数据采集
    invoice  # 发票管理
    journal  # 现金日报
    payment  # 支付 / 付款管理
    process  # 审批单管理
    receipt  # 付款单/收款单管理
    supplier  # 供应商管理
    voucher  # 会计凭证
    case  # 案例检索
    consult  # 法律咨询
    deliadvice  # 法律咨询 (得力)
    delicasesearch  # 案例检索 (得力)
    delisearch  # 法规检索 (得力)
    search  # 法规/案例统一检索
    stream  # 直播流管理
    attachment  # 邮件附件管理
    auto-reply  # 邮件自动回复管理
    contact  # 邮件联系人管理
    draft  # 草稿管理
    folder  # 邮件文件夹管理
    mailbox  # 邮箱地址管理
    message  # 邮件管理
    rule  # 收信规则管理
    tag  # 邮件标签管理
    template  # 邮件模板管理
    thread  # 邮件会话管理
    user  # 邮箱用户管理
    get  # 获取听记内容
    hot-word  # 个人热词管理
    list  # 听记列表
    mind-graph  # 思维导图管理
    permission  # 听记成员权限管理
    record  # 控制听记录音
    replace-text  # 查找替换段落和纪要中匹配的文字
    speaker  # 发言人管理
    tag  # 听记标签/分组管理
    update  # 更新听记信息
    upload  # 文件上传管理
    approval  # 审批管理
    chmod  # 授予指定权限
    create  # [deprecated] 已废弃，请改用 `dws report entry submit`
    created  # [deprecated] 已废弃，请改用 `dws report outbox list`
    detail  # [deprecated] 已废弃，请改用 `dws report entry get`
    entry  # 日志条目（单条日报操作 — get / stats / submit）
    inbox  # 收件箱（我收到的日报）
    list  # [deprecated] 已废弃，请改用 `dws report inbox list`
    outbox  # 发件箱（我发出的日报）
    sent  # [deprecated] 已废弃，请改用 `dws report outbox list`
    stats  # [deprecated] 已废弃，请改用 `dws report entry stats`
    template  # 日志模版
    add-dimension  # 在末尾追加空行或空列
    append  # 在工作表末尾追加数据
    batch-update  # 批量执行多个写操作（原子事务）
    chart  # 浮动图表管理
    cond-format  # 条件格式管理
    copy  # 复制工作表
    create  # 创建钉钉表格文档
    create-float-image  # 创建浮动图片
    csv-get  # 以 CSV 格式读取工作表数据
    csv-put  # 将 CSV 数据写入表格指定位置（纯值，自动扩容）
    delete-dimension  # 删除指定位置的行或列
    delete-dropdown  # 删除下拉列表
    delete-float-image  # 删除浮动图片
    delete-sheet  # 删除工作表
    export  # 导出表格为 xlsx（异步任务一站式）
    filter  # 全局筛选管理
    filter-view  # 筛选视图管理
    find  # 在工作表中搜索单元格内容
    get-dropdown  # 获取下拉列表配置
    get-float-image  # 获取浮动图片详情
    info  # 获取指定工作表详情
    insert-dimension  # 在指定位置插入行或列
    list  # 获取全部工作表列表
    list-float-images  # 列出工作表所有浮动图片
    media-upload  # 上传附件到表格
    merge-cells  # 合并单元格
    move-dimension  # 移动行或列到指定位置/调整顺序
    new  # 新建工作表
    range  # 数据区域操作
    replace  # 查找替换/批量替换/精确匹配替换/正则替换文本
    set-dropdown  # 设置下拉列表
    template  # 表格模板管理
    unmerge-cells  # 取消合并单元格
    update  # 更新工作表属性
    update-dimension  # 更新指定范围行/列属性（显隐、行高/列宽）
    update-float-image  # 更新浮动图片属性
    write-image  # 上传图片并写入表格单元格
    comment  # 待办评论：新增 / 列表 / 删除
    task  # 创建 / 查询 / 更新 / 删除待办
    member  # 知识库成员管理
    node  # 知识库节点管理
    space  # 知识库管理
    app  # 应用管理
    automation  # 自动化流程管理
    data  # 表单数据管理
    design  # 宜搭设计相关命令
    form  # 表单定义管理
    process  # 流程审批管理
    task  # 待办任务管理
    export  # 导出可迁移认证包
    import  # 导入可迁移认证包
    login  # 登录钉钉（自动刷新 token，必要时扫码）
    logout  # 清除认证信息（默认退出所有组织）
    reset  # 重置认证信息（清除本地 Token，触发重新授权）
    status  # 查看认证状态
    list  # 列出所有可用配置项
    app  # 开放平台应用
    connect  # 建联：把现成机器人接到当前本地 agent（起 Stream，不建号）
    doc  # 开放平台文档搜索
    build  # 将插件 stdio server 编译为原生二进制
    config  # 管理插件配置
    create  # 脚手架生成新插件目录
    dev  # 将本地目录注册为开发态插件
    disable  # 禁用插件
    enable  # 启用插件
    info  # 查看插件详情
    install  # 安装插件
    list  # 列出已安装的插件
    remove  # 卸载已安装的插件
    validate  # 校验 plugin.json
    list  # 列出已登录组织 profile
    switch  # 切换当前组织 profile
    use  # 切换当前组织 profile（兼容 profile switch）
    execute  # 生成面向 Agent 的恢复分析包
    finalize  # 回写恢复闭环结果
    plan  # 基于失败快照生成恢复计划
    install  # 下载并安装技能到悟空
    publish  # 发布技能到企业技能库
    search  # 搜索技能市场
    setup  # 安装 dws 自身 skill 到 Agent 目录
      detail  # 获取经营合约详情
      fields  # 获取经营合约字段列表
      list  # 获取经营合约列表
      update  # 更新经营合约
      list-statistics  # 获取周月报数据跟催列表
      submit-detail  # 获取周月报规则提交详情
      detail  # 获取计分卡详情
      entity-detail  # 获取计分卡实体详情
      update  # 更新计分卡
      detail  # 获取战略解码详情
      list  # 获取战略解码列表
      update  # 更新战略解码
      objectives  # 查询用户目标列表
      rules  # 获取用户的规则周期列表
      disable  # 关闭高级权限总开关（高危）
      enable  # 开启高级权限总开关
      role-create  # 创建自定义角色
      role-delete  # 删除自定义角色（不可逆）
      role-get  # 获取单个角色完整配置
      role-list  # 列出 Base 下所有角色
      role-update  # 增量更新自定义角色配置（patch 语义）
      upload  # 准备附件上传
      copy  # 复制 AI 表格
      create  # 创建 AI 表格
      delete  # 删除 AI 表格
      get  # 获取 AI 表格信息
      get-primary-doc-id  # 获取主键文档ID
      list  # 获取 AI 表格列表
      search  # 搜索 AI 表格
      update  # 更新 AI 表格
      create  # 创建图表
      delete  # 删除图表
      get  # 获取图表信息
      share  # 图表分享管理
      update  # 更新图表
      widgets-example  # 获取图表配置示例
      arrange  # 自动重排仪表盘图表布局
      config-example  # 获取仪表盘配置示例
      create  # 创建仪表盘
      delete  # 删除仪表盘
      get  # 获取仪表盘信息
      share  # 仪表盘分享管理
      update  # 更新仪表盘
      data  # 导出数据
      create  # 创建字段
      delete  # 删除字段
      get  # 获取字段详情
      list  # 获取字段信息（dws aitable field get 的别名）
      search-options  # 搜索单选/多选字段的选项
      update  # 更新字段
      create  # 创建表单视图
      delete  # 删除表单
      field  # 表单字段管理
      get  # 获取单个表单视图详情
      list  # 列出表单视图
      questions  # 表单题目管理（等价于 field create / delete）
      share  # 表单分享管理
      update  # 更新表单配置
      data  # 导入数据
      upload  # 准备导入文件上传
      batch-update  # 批量更新记录（同一 cells 应用到多条 recordId）
      create  # 新增记录
      delete  # 删除行记录
      get  # 按 ID 获取记录（record query --record-ids 的便捷别名，单次最多 100 条）
      history-list  # 查询行记录变更历史
      list  # 获取行记录（dws aitable record query 的别名）
      primary-doc-create  # 为记录创建主键文档
      primary-doc-get  # 查询记录的主键文档
      query  # 获取行记录
      query-empty  # 查询完全没填用户字段的空行
      share-url  # 批量获取记录分享链接
      update  # 更新记录
      upsert  # 批量创建或更新记录（Upsert）
      create  # 创建文件夹
      delete  # 删除文件夹
      list-empty  # 列出空文件夹
      list-nodes  # 列出全部节点
      move-node  # 移动节点
      rename  # 重命名文件夹
      reorder  # 调整文件夹顺序
      create  # 创建数据表
      delete  # 删除数据表
      get  # 获取数据表
      list  # 获取数据表信息（dws aitable table get 的别名）
      update  # 更新数据表
      search  # 搜索模板
      create  # 创建视图
      delete  # 删除视图
      duplicate  # 复制视图
      get  # 获取视图详情
      list  # 获取视图信息（dws aitable view get 的别名）
      lock  # 锁定/解锁视图
      update  # 更新视图
      disable  # 禁用指定工作流（高危）
      enable  # 启用指定工作流
      get  # 获取单个工作流详情
      list  # 列出 Base 下的工作流
      get  # 根据补卡规则主键 ID 查询补卡规则详情
      search  # 查询当前用户可管理的补卡规则列表
      list  # 查询用户审批单（补卡/加班/请假/出差外出）
      templates  # 查询补卡/请假/加班/外出/出差审批提交链接
      record  # 查询打卡流水
      result  # 查询打卡结果
      records  # 查询指定员工的签到记录
      create  # 创建班次
      get  # 根据班次 ID 查询班次详情
      search  # 查询当前用户可管理的所有班次详情
      update  # 更新班次
      get  # 查询全局规则设置（仅管理员）
      save  # 更新保存全局规则设置（仅管理员）
      create  # 创建考勤组
      filtered-get  # 根据考勤组 ID 按需查询成员/打卡地址/蓝牙/Wifi 信息
      get  # 根据考勤组 ID 查询考勤组全量信息
      search  # 查询当前用户可管理的考勤组列表
      update  # 更新考勤组配置（仅修改需要变更的字段）
      update-members  # 更新考勤组成员（添加/删除考勤人员、部门、无需考勤人员）
      get  # 根据加班规则主键 ID 查询加班规则详情
      search  # 查询当前用户可管理的加班规则列表
      get  # 查询个人考勤详情
      columns  # 获取企业考勤字段列表
      query-data  # 根据字段查询考勤数据
      query-leave  # 查询用户假期数据
      get  # 获取指定用户的排班记录
      import  # 导入排班记录到排班制考勤组
      get  # 查询个人规则设置
      save  # 更新保存个人规则设置
      list  # 批量查询员工班次信息
      balance  # 查询指定员工假期余额
      records  # 查询指定员工假期余额变更记录
      save-balance  # 设置员工假期余额
      types  # 查询当前用户假期规则列表
      update-type  # 更新假期规则
      add  # 把我的日历共享给某人
      delete  # 删除日历访问权限
      list  # 查询我的日历共享给了谁
      add  # 添加日程附件
      add  # 添加参会人
      delete  # 移除参会人
      list  # 查看参会人
      get  # 查询指定日历本
      list  # 查询用户的日历列表
      search  # 搜索日历本
      update  # 更新指定日历本
      search  # 查询用户 / 会议室闲忙状态
      create  # 创建日程
      delete  # 删除日程
      get  # 获取日程详情
      list  # 查询日程列表
      list-mine  # 查询归属个人的日程列表
      respond  # 响应日程（接受/拒绝/暂定）
      suggest  # 建议日程时间
      update  # 修改日程
      add  # 预定会议室
      delete  # 移除会议室
      list-groups  # 会议室分组列表
      search  # 搜索会议室 (按名称搜索或按时间段查可用会议室)
      find  # 搜索【全部可用】机器人（含他人/官方，额外返回 openDingTalkId 可发单聊）
      search  # 搜索【我自己创建】的机器人（仅本人创建的，不含他人/官方机器人）
      add-conv  # 将会话移动到指定的自定义分组中
      create  # 创建用户自定义会话分组
      delete  # 删除用户自定义会话分组
      list  # 获取用户自定义会话分组
      list-conversations  # 拉取指定自定义会话分组下的会话
      remove-conv  # 将会话从指定的自定义分组中移出
      rename  # 更新用户自定义会话分组的名称
      cross-org  # 授予跨组织 chat 数据访问权限
      upload  # 上传本地文件或 URL 文件到会话文件空间
      audit-join-validation  # 审批入群验证（通过、拒绝、删除）
      bots  # 查看群内所有机器人
      create  # 创建群（支持内部群/外部群/普通群/话题圈）
      dismiss  # 解散群聊
      get-by-group-id  # 根据群号获取群聊信息
      invite-url  # 获取群邀请链接
      list-all  # 分页拉取我所有群列表
      list-join-validations  # 分页拉取入群验证记录
      list-my-groups  # 拉取我创建/管理的群
      members  # 群成员管理
      quit  # 退出群聊
      rename  # 更新群名称
      set-admin  # 设置 / 取消群管理员
      set-history  # 设置新成员入群可查看历史消息选项
      transfer-owner  # 转让群主
      update-alias  # 设置群备注
      update-icon  # 更新群头像
      update-nick  # 设置用户在群内的群昵称
      update-settings  # 更新群设置
      add  # 添加群身份
      list  # 拉取会话的群身份列表
      query-user  # 查询群成员的群身份
      remove  # 删除群身份
      remove-user  # 移除用户的指定群身份
      set-user  # 设置用户的群身份（覆盖该用户的全部群身份）
      update  # 更新群身份名称
      add-emoji  # 对消息添加 emoji 表情回应
      add-text-emotion  # 对消息添加文字表情回应
      combine-forward  # 合并转发多条消息
      create-text-emotion  # 创建文字表情（获取 emotionId）
      download-media  # 下载消息中的资源（图片/视频/语音等）到本地
      forward  # 转发单条消息（源/目标会话均支持单聊/群聊）
      forward-topic  # 转发话题消息到目标会话
      list  # 拉取会话消息内容
      list-all  # 拉取指定时间范围内当前用户的所有会话消息
      list-by-ids  # 根据消息 ID 批量查询消息
      list-by-sender  # 拉取指定发送者的消息（包含单聊和群聊）
      list-focused  # 拉取特别关注人的消息
      list-mentions  # 拉取 @我 的消息
      list-pin-msg  # 拉取会话中钉住的消息列表
      list-topic-replies  # 拉取群话题回复消息列表
      list-unread-conversations  # 获取未读会话列表
      query-send-status  # 查询消息发送状态
      read-status  # 查询消息的已读/未读状态
      recall  # 撤回用户发送的消息
      recall-by-bot  # 机器人撤回消息（--group 群聊 / 不传为单聊）
      remove-emoji  # 移除消息的 emoji 表情回应
      remove-text-emotion  # 移除消息的文字表情回应
      reply  # 引用回复消息（支持单聊/群聊）
      search  # 按关键词搜索消息
      search-advanced  # 多维度搜索消息
      send  # 以当前用户身份发送消息（--group 群聊 / --user 或 --open-dingtalk-id 单聊）
      send-by-bot  # 机器人发送消息（--group 群聊 / --users 单聊）
      send-by-webhook  # 自定义机器人 Webhook 发送群消息
      send-card  # 创建并推送流式卡片
      set-pin-msg  # 钉住消息（Pin）
      set-top-msg  # 置顶消息
      unset-pin-msg  # 取消钉住消息（Unpin）
      unset-top-msg  # 取消置顶消息
      update-card  # 流式更新卡片内容
      reserve  # 预约会议
      invite  # 邀请指定人入会
      get-info  # 获取部门详情（部门ID、名称、人数）
      list-children  # 查看子部门
      list-members  # 查看部门成员（仅本部门，不含下级）
      search  # 搜索部门
      get  # 根据角色名称查询角色
      list  # 获取企业所有角色列表
      list-members  # 查询角色下的成员
      list-my-followings  # 获取当前用户的特别关注列表
      dismission  # 离职员工查询
      get  # 批量获取用户详情（组织管理信息）
      get-self  # 获取当前用户信息（我是谁 / 本人）
      profile  # 用户档案（花名册）
      search  # 按关键词搜索用户
      search-mobile  # 按手机号搜索用户
      invest  # 对外投资
      shareholder  # 股东信息
      copyright  # 著作权信息
      icp  # ICP备案信息
      patent  # 专利信息
      trademark  # 商标信息
      assist  # 司法协助
      consum  # 限制高消费
      court  # 开庭公告
      dishonest  # 失信被执行
      execute  # 被执行信息
      finalcase  # 终本案件
      litigation  # 涉诉公告
      owetax  # 催缴欠税
      penalty  # 行政处罚
      pledge  # 股权出质
      taxviolation  # 重大税收违法
      verdict  # 裁判文书
      search  # 搜索开放平台文档
      list  # 查询 DING 消息历史
      recall  # 撤回 DING 消息
      recall-personal  # 以用户身份撤回 DING
      receiver-status  # 查看 DING 接收状态
      send  # 发送 DING 消息
      send-by-message  # 消息转 DING（将聊天消息转为 DING 通知）
      send-personal  # 以用户身份发送 DING
      delete  # 删除块元素
      insert  # 插入块元素
      list  # 查询块元素
      update  # 更新块元素
      create  # 创建文档评论
      create-inline  # 创建划词评论
      list  # 查询文档评论列表
      reply  # 回复评论
      get  # 查询导出任务结果（手动兜底）
      create  # 创建文件
      download  # 下载文档附件
      insert  # 上传附件并插入文档
      apply  # 应用文档模板
      list  # 获取文档模板列表
      search  # 搜索文档模板
      list  # 查看文档历史版本列表
      revert  # 回滚文档到指定版本
      save  # 手动保存文档版本快照
      add  # 添加协作者
      list  # 查询协作者列表
      remove  # 移除协作者权限
      update  # 更新协作者权限
      get  # 查询文件公开发布状态
      set  # [危险] 设置文件为互联网公开
      unset  # [危险] 关闭文件互联网公开
      list  # 查看回收站文件列表
      restore  # 还原回收站中的文件
      list  # 分页查询企业账户列表
      create  # 录入银行交易明细
      list  # 分页查询银行交易明细
      query  # 查询银行交易明细
      search  # 搜索收支类别
      save  # 保存主体
      search  # 模糊搜索主体
      update  # 修改主体信息
      get  # 根据名称精确查询客户
      list  # 分页查询客户列表
      save  # 新建客户
      account  # 查询数电账号信息
      batch-draw  # 批量开票（轻量化版）
      batch-draw-query  # 批量开票查询（轻量化版）
      batch-draw-query-saas  # 批量开票查询（SaaS版）
      batch-draw-saas  # 批量开票（SaaS版）
      do-login  # 数电登录认证
      do-login-status  # 查询数电登录状态
      face-qr  # 获取人脸识别二维码
      face-status  # 获取人脸认证状态
      file  # 获取发票版式文件
      get-table  # 获取发票表格配置
      goods-code  # 商品智能赋码
      import-goods  # 导入商品
      issue  # 开具数电发票
      login-page  # 获取数电发票登录页面链接
      search-goods  # 搜索商品
      send-email  # 发送发票邮件（轻量化版）
      send-email-saas  # 发送发票邮件（SaaS版）
      skill-version  # 查询开票 Skill 版本
      sms-code  # 上传数电登录短信验证码
      title  # 智能抬头
      execute  # 执行数据采集（批量）
      execute-general  # 执行数据采集（通用场景）
      query-rule  # 查询采集规则
      save-rule  # 保存采集规则
      try-execute  # 尝试执行数据采集（单条验证）
      add-record  # 添加发票到审批单
      issue  # 开具发票
      issue-result  # 查询开票结果
      list  # 分页查询发票列表
      list-application  # 查询开票申请列表
      recommend-category  # AI 推荐发票收支类别
      upload  # 上传发票
      daily  # 按日查询现金日报
      detail-url  # 获取现金日报明细链接
      account-list  # 查询收款账户列表
      account-url  # 获取收款账户管理页面链接
      cashier-url  # 查询支付收银台链接
      create  # 创建待付款审批单
      list  # 查询待付款列表
      payer-list  # 查询付款账户列表
      form-data  # 根据审批编号查询审批表单信息
      list  # 根据审批模版名查询审批单列表
      create  # 创建付款单
      create-collection  # 基于银行交易明细创建收款单
      search  # 模糊搜索供应商
      entries  # 根据审批单生成会计分录
      generate  # 根据审批单据号生成会计凭证
      list  # 查看我的直播列表
      download  # 下载邮件附件到本地
      list  # 列举邮件附件
      get  # 获取用户的自动回复配置
      batch-delete  # 批量删除邮件联系人
      create  # 创建邮件联系人
      list  # 列举邮件联系人
      update  # 更新邮件联系人
      create  # 创建草稿
      send  # 发送草稿
      update  # 更新草稿
      create  # 创建邮件文件夹
      delete  # 删除邮件文件夹
      list  # 列举邮件文件夹
      update  # 更新邮件文件夹
      list  # 查询可用邮箱地址
      batch-delete  # 批量删除邮件
      batch-move  # 批量移动邮件到指定文件夹
      batch-update  # 批量修改邮件状态（标记已读/未读/添加标签/移除标签）
      forward  # 转发邮件
      get  # 查看邮件完整内容
      list  # 列出文件夹中的邮件
      reply  # 回复邮件
      reply-all  # 回复所有人
      search  # 搜索邮件 (KQL 语法)
      send  # 发送邮件
      verify  # 查询邮件发送状态
      adjust  # 调整收信规则排序
      create  # 创建个人收信规则
      delete  # 删除个人收信规则
      list  # 列出个人收信规则
      update  # 更新个人收信规则
      create  # 创建邮件标签
      delete  # 删除邮件标签
      list  # 列举邮件标签
      update  # 更新邮件标签
      create  # 创建邮件模板
      delete  # 删除邮件模板
      get  # 获取邮件模板详情
      list  # 列举邮件模板
      update  # 更新邮件模板
      batch-trash  # [危险] 批量删除邮件会话
      batch-update  # 批量修改邮件会话状态
      get  # 获取会话详情
      list  # 列出邮件会话
      trash  # [危险] 删除邮件会话
      update  # 修改邮件会话状态
      search  # 搜索邮箱用户
      audio  # 获取听记音频/视频地址
      batch  # 批量查询听记详情
      info  # 获取听记基础信息
      keywords  # 获取听记关键字列表
      summary  # 获取听记 AI 摘要
      todos  # 获取听记中提取的待办事项
      transcription  # 获取听记语音转写原文
      add  # 添加个人热词
      list  # 查询我的热词列表
      all  # 查询我有权限访问的所有听记列表
      mine  # 查询我创建的听记列表
      shared  # 查询他人共享给我的听记列表
      create  # 创建思维导图
      status  # 查询思维导图状态
      add  # 批量添加听记成员并设置权限
      remove  # 批量移除听记成员权限
      pause  # 暂停听记录音
      resume  # 恢复听记录音
      start  # 发起听记（开始录音）
      stop  # 结束听记录音
      replace  # 替换发言人
      summary  # 发言人段落总结
      list  # 查询我的听记标签/分组列表
      query  # 根据标签ID查询听记列表
      summary  # 更新纪要内容
      title  # 修改听记标题
      cancel  # 取消文件上传会话
      complete  # 完成文件上传并创建听记
      create  # 创建文件上传会话
      append-task  # 对审批任务进行加签
      approve  # 同意审批
      detail  # 获取审批实例详情
      ding-info  # 获取审批任务的被催办人 userId（需与 ding message send 串联使用）
      list-cc  # 获取抄送当前用户的审批单列表
      list-executed  # 获取当前用户已经处理过的审批单列表
      list-forms  # 获取当前用户可见的审批表单列表
      list-initiated  # 查询审批模板下已发起的审批记录
      list-pending  # 查询待我处理的审批
      list-submitted  # 获取当前用户已发起的审批单列表
      oa-cc-noticer  # 对审批实例进行抄送
      oa-comments  # 对审批实例添加评论
      records  # 获取审批操作记录
      redirect-task  # 转交审批任务给其他人
      reject  # 拒绝审批
      revert-activities  # 获取审批任务可回退的节点信息（退回前必须调用，获取可回退节点列表）
      revert-task  # 退回审批任务到指定节点（审批人或发起人）
      revoke  # 撤销已发起的审批
      search-forms  # 按关键字模糊搜索当前用户可见的审批表单
      tasks  # 查询待我审批的任务 ID
      get  # 读取单份日报正文（含字段明细 + 钉钉跳转链接）
      stats  # 读取单份日报的已读统计
      submit  # 提交一份新日报（按模版）
      list  # 列出我收到的日报
      list  # 列出我发出的日报
      detail  # [deprecated] 已废弃，请改用 `dws report template get`
      get  # 读取单个日志模版的字段定义
      list  # 获取当前用户可用的日志模版列表
      create  # 创建浮动图表
      delete  # 删除浮动图表
      list  # 获取浮动图表
      update  # 更新浮动图表
      create  # 创建条件格式规则
      delete  # 删除条件格式规则
      list  # 获取条件格式规则
      update  # 更新条件格式规则
      clear-criteria  # 清除单列筛选条件
      create  # 创建全局筛选
      delete  # 删除全局筛选
      get  # 获取全局筛选信息
      sort  # 筛选排序
      update  # 批量更新筛选条件
      create  # 创建筛选视图
      delete  # 删除筛选视图
      delete-criteria  # 删除筛选视图列条件
      get-criteria  # 获取单列筛选条件
      info  # 获取单个筛选视图详情
      list  # 获取所有筛选视图
      list-criteria  # 列出筛选视图所有列条件
      update  # 更新筛选视图属性
      update-criteria  # 更新筛选视图列条件
      batch-clear  # 批量清除多个区域（原子事务）
      batch-set-style  # 按配置文件批量设置样式
      clear  # 清除工作表指定区域
      copy-to  # 复制工作表指定区域到目标位置
      fill  # 自动填充工作表指定区域
      move-to  # 移动工作表指定区域到目标位置
      read  # 读取工作表数据（别名: get）
      set-style  # 设置指定单元格区域的样式
      sort  # 对工作表指定区域排序
      update  # 更新工作表指定区域内容
      apply  # 应用表格模板
      list  # 获取表格模板列表
      search  # 搜索表格模板
      add  # 新增待办评论
      delete  # 删除待办评论
      list  # 查询待办评论列表
      add-attachment  # 上传待办附件
      add-executor  # 添加待办执行人
      add-participant  # 添加待办参与人
      add-reminder  # 添加待办提醒
      create  # 创建待办
      create-sub  # 创建子待办
      delete  # 删除待办
      done  # 修改执行者的待办完成状态
      get  # 待办详情
      list  # 查询待办列表
      list-attachment  # 查询待办任务的附件列表
      list-sub  # 查询子待办列表
      remove-attachment  # 删除待办任务的附件
      remove-executor  # 移除待办执行人
      remove-participant  # 移除待办参与人
      reset-reminder  # 重置待办提醒
      update  # 修改待办任务
      add  # 添加知识库成员
      list  # 查询知识库成员列表
      remove  # 移除知识库成员
      update  # 更新知识库成员权限
      copy  # 复制知识库节点
      create  # 在知识库中创建节点
      delete  # 删除知识库节点
      list  # 列出知识库节点
      move  # 移动知识库节点
      search  # 在知识库中搜索节点
      create  # 创建知识库
      delete  # 删除知识库
      get  # 查看知识库详情
      list  # 列出空间（知识库 / 钉盘空间）
      search  # 搜索知识库
      create  # 创建宜搭应用
      list  # 获取宜搭应用列表
      list-forms  # 获取应用内表单列表
      update-permission  # 更新应用管理员权限
      create  # 创建自动化流程
      update  # 更新自动化流程节点
      create  # 新增表单实例
      detail  # 获取单条记录详情
      export-query  # 查询导出任务结果（返回下载 URL）
      export-start  # 触发异步导出表单数据
      search  # 表单数据条件查询
      update  # 修改表单实例
      form  # 表单设计
      process  # 流程设计
      components  # 获取表单字段定义
      execute  # 执行审批（同意/拒绝）
      records  # 获取流程审批操作记录
      redirect  # 转交审批任务
      running-tasks  # 查询流程运行中的任务节点
      start  # 发起流程表单实例
      list  # 查询当前用户待办任务
      create  # 创建开放平台企业内部应用
      credentials  # 开放平台应用凭证
      delete  # 删除开放平台企业内部应用（不可逆，需 --confirm-name 二次确认）
      disable  # 停用开放平台企业内部应用
      enable  # 启用开放平台企业内部应用
      event  # 开放平台应用事件订阅
      get  # 查询开放平台企业内部应用详情
      list  # 查询开放平台企业内部应用列表
      member  # 开放平台应用成员管理
      permission  # 开放平台应用权限
      robot  # 开放平台应用机器人能力
      security  # 开放平台应用安全设置
      update  # 修改开放平台企业内部应用基础信息
      version  # 开放平台应用版本发布
      webapp  # 开放平台网页应用配置
      list  # 列出本机所有连接器及健康状态（healthy/degraded/down）；--json 供脚本消费
      restart  # 重启连接器守护进程（通过持久化的 unifiedAppId 重新拉取密钥，无需本地存密钥）
      status  # 查看连接器健康状态（healthy/degraded/down，pid、收发活动、日志路径；--json 供外部托管消费）
      stop  # 优雅停止后台连接器守护进程（释放单实例锁与 Stream 连接）
      search  # 搜索开放平台文档
      get  # 读取插件配置项
      list  # 列出插件所有配置项
      set  # 设置插件配置项
      unset  # 删除插件配置项
        get  # 获取图表分享配置
        update  # 更新图表分享配置
        get  # 获取仪表盘分享配置
        update  # 更新仪表盘分享配置
        hide  # 切换表单字段隐藏
        list  # 列出表单字段
        update  # 更新表单字段
        create  # 向表单添加题目（等价于 field create）
        delete  # 从表单删除题目（等价于 field delete，不可逆）
        get  # 获取表单分享配置
        update  # 开启/关闭分享表单
        aggregate  # 获取视图 aggregate 配置
        card  # 获取视图 card 配置
        field-widths  # 获取视图 field-widths 配置
        fill-color-rule  # 获取视图 fill-color-rule 配置
        filter  # 获取视图 filter 配置
        frozen-cols  # 获取视图冻结列数
        group  # 获取视图 group 配置
        lock  # 获取视图锁定状态
        row-height  # 获取视图行高（单元格高度）
        sort  # 获取视图 sort 配置
        timebar  # 获取视图 timebar 配置
        visible-fields  # 获取视图 visible-fields 配置
        aggregate  # 更新视图字段聚合统计（仅 Grid）
        card  # 更新视图 card 配置（Kanban / Gallery）
        field-widths  # 更新视图字段列宽（仅 Grid）
        fill-color-rule  # 更新视图数据高亮规则
        filter  # 更新视图 filter 配置
        frozen-cols  # 更新视图冻结列数
        group  # 更新视图 group 配置
        name  # 重命名视图（= view update --name 的便捷子命令）
        row-height  # 更新视图行高（单元格高度）
        sort  # 更新视图 sort 配置
        timebar  # 更新视图 timebar 配置（仅 Gantt）
        visible-fields  # 更新视图可见字段列表
        add  # 添加群成员
        add-bot  # 将机器人添加到群中
        list-by-ids  # 根据成员 ID 批量查询群成员详情
        remove  # 移除群成员
        remove-bot  # 从群内移除机器人
        search  # 分页获取离职员工列表
        fields  # 查询花名册有权限的字段列表
        get  # 查询员工花名册字段信息（个人档案）
        create  # 触发创建发言人段落总结任务
        get  # 查询发言人段落总结结果
        create  # 新建表单/报表/自定义页面
        get-info  # 读取表单元数据
        get-schema  # 读取表单完整 schema
        preview  # 生成表单预览链接
        update-schema  # 覆盖式更新表单 schema
        update-title  # 修改表单标题
        create-draft  # 基于源流程开新草稿
        get  # 读取流程版本详情
        publish  # 发布流程
        update  # 保存流程定义
        versions  # 流程版本
        get  # 读取开放平台应用凭证
        list  # 查询应用已订阅的事件列表
        subscribe  # 订阅应用事件回调
        unsubscribe  # 取消订阅应用事件
        add  # 添加开放平台应用成员
        list  # 查询开放平台应用成员
        remove  # 移除开放平台应用成员
        add  # 申请开放平台应用权限点
        list  # 查询开放平台应用权限列表
        remove  # 取消开放平台应用权限点
        config  # 创建或更新现有应用的机器人配置（upsert）
        disable  # 停用现有应用的机器人能力
        enable  # 启用现有应用机器人能力（纯启用，无需配置字段）
        get  # 查询现有应用的机器人配置
        result  # 查询机器人异步创建任务结果
        submit  # 异步提交钉钉智能体机器人创建任务（支持失败重试）
        config  # 更新开放平台应用安全配置
        check-approval  # 预检版本发布是否需要审批（不实际发布）
        create  # 基于当前配置创建应用新版本
        get  # 查询指定版本详情
        list  # 分页查询应用版本列表
        publish  # 发布指定版本（含高敏权限需 --confirmed-sensitive）
        status  # 查询版本发布/审批状态
        config  # 配置网页应用能力
        get  # 查询网页应用配置
          list  # 列出流程版本
```
