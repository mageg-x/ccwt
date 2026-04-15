export default {
  app: {
    name: 'CCWT'
  },
  nav: {
    toggleSidebar: '切换侧边栏',
    language: '语言',
    voiceInput: '语音输入',
    toggleTheme: '切换主题',
    sessionHistory: '会话历史',
    settings: '设置',
    logout: '退出登录'
  },
  settings: {
    title: '设置',
    tabs: {
      voice: '语音识别',
      proxy: '代理',
      about: '关于'
    },
    voice: {
      enable: '语音识别开关',
      enableDesc: '启用后可用语音输入命令',
      baiduConfig: '百度语音 API',
      baiduConfigTip: '请前往 {link} 创建语音应用，获取 App ID、API Key 和 Secret Key',
      baiduConsole: '百度智能云控制台',
      appId: 'App ID',
      apiKey: 'API Key',
      secretKey: 'Secret Key',
      save: '保存',
      saveSuccess: '保存成功',
      saveFailed: '保存失败',
      loadFailed: '加载设置失败'
    },
    proxy: {
      title: 'SOCKS5 代理服务',
      running: '运行中',
      stopped: '已停止',
      bindConfig: '默认绑定配置',
      ip: 'IP',
      port: '端口',
      save: '保存',
      connectionInfo: '连接信息',
      copy: '复制',
      bindTip: '0.0.0.0 表示绑定所有网卡；127.0.0.1 仅允许本地连接',
      startFailed: '启动失败',
      stopFailed: '停止失败',
      statusFailed: '获取代理状态失败'
    },
    about: {
      version: '版本',
      checkUpdate: '检测更新',
      upToDate: '已是最新版本',
      updateAvailable: '发现新版本 {version}',
      checkFailed: '检测更新失败',
      features: '功能说明',
      webBasedTerminal: '基于 Web 的 Claude Code 终端',
      voiceInput: '支持语音输入命令',
      builtInProxy: '内置 SOCKS5 代理',
      fileManagement: '文件管理和在线编辑'
    }
  },
  common: {
    confirm: '确认',
    cancel: '取消',
    close: '关闭',
    loading: '加载中...',
    error: '错误',
    success: '成功',
    save: '保存',
    saveFailed: '保存失败',
    ok: '确定',
    delete: '删除',
    deleteConfirm: '确定要删除吗？此操作不可恢复。'
  },
  user: {
    admin: '管理员',
    user: '用户',
    adminPanel: '管理面板'
  },
  term: {
    newTab: '新建终端',
    rename: '重命名终端',
    terminalName: '终端名称',
    newTerminal: '新建终端',
    enterFocus: '终端最大化',
    exitFocus: '退出最大化'
  },
  voice: {
    title: '语音输入',
    startRecording: '点击开始录音',
    recording: '正在录音，点击停止...',
    processing: '处理中...',
    recognizing: '正在识别...',
    engine: '当前引擎',
    close: '关闭',
    micPermissionDenied: '无法访问麦克风',
    backendUnavailable: '后端语音不可用',
    backendDisabled: '后端语音识别未启用',
    authFailed: '百度鉴权失败或网络不可达',
    statusFailed: '无法获取后端语音状态',
    recognitionFailed: '语音识别失败'
  },
  fileTree: {
    uploadNoticeTitle: '上传提示',
    uploadResultTitle: '上传结果',
    dirUploadIgnored: '目录拖放上传暂不支持，已忽略目录',
    dirUploadNotSupported: '目录拖放上传暂不支持',
    invalidMoveTarget: '不能移动到该目录',
    moveFailedTitle: '移动失败',
    moveFailedMessage: '移动失败',
    uploadPartialFailed: '上传完成，失败 {count} 个'
  },
  admin: {
    title: '管理面板',
    userManagement: '用户管理',
    loadFailed: '加载用户列表失败',
    confirmRoleChange: '确定将 {username} 的角色改为 {role}？',
    confirmRoleChangeTitle: '确认修改角色',
    roleChangeSuccess: '角色修改成功',
    roleChangeFailed: '修改角色失败',
    confirmDelete: '确定删除用户 {username}？此操作不可恢复。',
    deleteUserTitle: '危险操作确认',
    deleteFailed: '删除失败',
    demoteToUser: '降为用户',
    promoteToAdmin: '升为管理员',
    id: 'ID',
    username: '用户名',
    role: '角色',
    registeredAt: '注册时间',
    actions: '操作'
  },
  contextMenu: {
    newFile: '新建文件',
    newFolder: '新建文件夹',
    rename: '重命名',
    copyPath: '复制路径',
    download: '下载',
    delete: '删除',
    enterFileName: '请输入文件名',
    enterFolderName: '请输入文件夹名',
    enterNewName: '请输入新名称',
    fileNameExample: '例如: notes.md',
    folderNameExample: '例如: docs',
    newName: '新名称',
    create: '创建',
    save: '保存',
    confirmDelete: '确定删除 {name}？',
    confirmDeleteTitle: '确认删除',
    operationFailed: '操作失败',
    createFileFailed: '创建文件失败',
    createFolderFailed: '创建文件夹失败',
    renameFailed: '重命名失败',
    deleteFailed: '删除失败'
  }
}
