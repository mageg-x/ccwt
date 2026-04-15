export default {
  app: {
    name: 'CCWT'
  },
  nav: {
    toggleSidebar: 'Toggle Sidebar',
    language: 'Language',
    voiceInput: 'Voice Input',
    toggleTheme: 'Toggle Theme',
    sessionHistory: 'Session History',
    settings: 'Settings',
    logout: 'Logout'
  },
  settings: {
    title: 'Settings',
    tabs: {
      voice: 'Voice',
      proxy: 'Proxy',
      about: 'About'
    },
    voice: {
      enable: 'Voice Recognition',
      enableDesc: 'Enable voice input for commands',
      baiduConfig: 'Baidu Voice API',
      baiduConfigTip: 'Go to {link} to create a voice app and get App ID, API Key and Secret Key',
      baiduConsole: 'Baidu Cloud Console',
      appId: 'App ID',
      apiKey: 'API Key',
      secretKey: 'Secret Key',
      save: 'Save',
      saveSuccess: 'Saved successfully',
      saveFailed: 'Save failed',
      loadFailed: 'Failed to load settings'
    },
    proxy: {
      title: 'SOCKS5 Proxy',
      running: 'Running',
      stopped: 'Stopped',
      bindConfig: 'Default Bind',
      ip: 'IP',
      port: 'Port',
      save: 'Save',
      connectionInfo: 'Connection Info',
      copy: 'Copy',
      bindTip: '0.0.0.0 binds all interfaces; 127.0.0.1 for local only',
      startFailed: 'Start failed',
      stopFailed: 'Stop failed',
      statusFailed: 'Failed to get status'
    },
    about: {
      version: 'Version',
      checkUpdate: 'Check for Updates',
      upToDate: 'Already up to date',
      updateAvailable: 'New version {version} available',
      checkFailed: 'Failed to check for updates',
      features: 'Features',
      webBasedTerminal: 'Web-based Claude Code terminal',
      voiceInput: 'Voice command input',
      builtInProxy: 'Built-in SOCKS5 proxy',
      fileManagement: 'File management and online editing'
    }
  },
  common: {
    confirm: 'Confirm',
    cancel: 'Cancel',
    close: 'Close',
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    save: 'Save',
    saveFailed: 'Save failed',
    ok: 'OK',
    delete: 'Delete',
    deleteConfirm: 'Are you sure? This cannot be undone.'
  },
  user: {
    admin: 'Admin',
    user: 'User',
    adminPanel: 'Admin Panel'
  },
  term: {
    newTab: 'New Terminal',
    rename: 'Rename Terminal',
    terminalName: 'Terminal Name',
    newTerminal: 'New Terminal',
    enterFocus: 'Maximize Terminal',
    exitFocus: 'Exit Maximized View'
  },
  voice: {
    title: 'Voice Input',
    startRecording: 'Click to start recording',
    recording: 'Recording, click to stop...',
    processing: 'Processing...',
    recognizing: 'Recognizing...',
    engine: 'Current engine',
    close: 'Close',
    micPermissionDenied: 'Cannot access microphone',
    backendUnavailable: 'Backend voice unavailable',
    backendDisabled: 'Backend voice recognition disabled',
    authFailed: 'Baidu auth failed or network unreachable',
    statusFailed: 'Cannot get backend voice status',
    recognitionFailed: 'Voice recognition failed'
  },
  fileTree: {
    uploadNoticeTitle: 'Upload Notice',
    uploadResultTitle: 'Upload Result',
    dirUploadIgnored: 'Folder drag-and-drop upload is not supported. Folder entries were ignored.',
    dirUploadNotSupported: 'Folder drag-and-drop upload is not supported.',
    invalidMoveTarget: 'Cannot move to this directory.',
    moveFailedTitle: 'Move Failed',
    moveFailedMessage: 'Move failed',
    uploadPartialFailed: 'Upload completed, {count} failed.'
  },
  admin: {
    title: 'Admin Panel',
    userManagement: 'User Management',
    loadFailed: 'Failed to load users',
    confirmRoleChange: 'Are you sure you want to change role of {username} to {role}?',
    confirmRoleChangeTitle: 'Confirm Role Change',
    roleChangeSuccess: 'Role changed successfully',
    roleChangeFailed: 'Failed to change role',
    confirmDelete: 'Are you sure you want to delete user {username}? This cannot be undone.',
    deleteUserTitle: 'Confirm Delete',
    deleteFailed: 'Failed to delete user',
    demoteToUser: 'Demote',
    promoteToAdmin: 'Promote',
    id: 'ID',
    username: 'Username',
    role: 'Role',
    registeredAt: 'Registered At',
    actions: 'Actions'
  },
  contextMenu: {
    newFile: 'New File',
    newFolder: 'New Folder',
    rename: 'Rename',
    copyPath: 'Copy Path',
    download: 'Download',
    delete: 'Delete',
    enterFileName: 'Enter file name',
    enterFolderName: 'Enter folder name',
    enterNewName: 'Enter new name',
    fileNameExample: 'e.g.: notes.md',
    folderNameExample: 'e.g.: docs',
    newName: 'New name',
    create: 'Create',
    save: 'Save',
    confirmDelete: 'Are you sure you want to delete {name}?',
    confirmDeleteTitle: 'Confirm Delete',
    operationFailed: 'Operation failed',
    createFileFailed: 'Failed to create file',
    createFolderFailed: 'Failed to create folder',
    renameFailed: 'Failed to rename',
    deleteFailed: 'Failed to delete'
  }
}
