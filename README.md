# knas (Knowledge Async)

剪切板自动同步工具 - 在 Mac 本机监听剪切板并同步到 Ubuntu 服务器。

## 功能特点

- 🔔 **实时监听**: 500ms 轮询剪切板变化，CPU 占用极低
- 🔒 **安全传输**: 通过 SSH 协议加密传输，无需在服务器安装接收端
- 🧠 **智能识别**: 自动识别 URL 并抓取网页标题
- 🎯 **智能过滤**: 支持敏感词过滤，避免同步密码等敏感信息
- 📁 **自动归档**: 按日期自动归档到 `YYYY/MM/DD/` 目录结构
- 🔄 **自动重试**: 传输失败自动重试，确保数据不丢失
- 🚀 **后台运行**: 支持 macOS LaunchAgent，开机自启动

## 安装

```bash
# 全局安装
npm install -g knas

# 或从源码安装
git clone https://github.com/yuanguangshan/knas.git
cd knas
npm install
npm run build
npm link
```

## 快速开始

### 1. 初始化配置

```bash
knas init
```

交互式配置：
- SSH 服务器地址
- SSH 端口（默认 22）
- SSH 用户名（默认 root）
- SSH 密钥路径（默认 ~/.ssh/id_rsa）
- 远程存储路径（默认 ~/knas_archive）
- 最小同步长度（默认 100 字符）

### 2. 启动守护进程

```bash
knas start
```

### 3. 查看状态

```bash
knas status
```

### 4. 查看日志

```bash
knas log          # 查看日志
knas log -f       # 实时跟踪日志
```

## 命令

| 命令 | 说明 |
|------|------|
| `knas init` | 初始化配置 |
| `knas start` | 启动守护进程 |
| `knas stop` | 停止守护进程 |
| `knas status` | 查看运行状态 |
| `knas log` | 查看日志 |
| `knas config` | 查看/编辑配置 |
| `knas service install` | 安装为 macOS 系统服务 |

## macOS 系统服务

安装为 LaunchAgent 后，knas 会在登录时自动启动：

```bash
# 安装服务
knas service install

# 加载服务
launchctl load ~/Library/LaunchAgents/com.knas.daemon.plist

# 卸载服务
launchctl unload ~/Library/LaunchAgents/com.knas.daemon.plist
```

## 配置文件

配置文件位于 `~/.knas/config.json`：

```json
{
  "ssh": {
    "host": "your-server.com",
    "port": "22",
    "user": "root",
    "key_path": "~/.ssh/id_rsa",
    "base_path": "~/knas_archive"
  },
  "clipboard": {
    "min_length": 100,
    "poll_interval_ms": 500,
    "exclude_words": ["password", "密码", "token"]
  },
  "sync": {
    "enabled": true,
    "max_retries": 3,
    "retry_delay_ms": 5000
  }
}
```

## 远程文件结构

同步的文件按以下结构存储：

```~/knas_archive/
├── 2026/
│   ├── 04/
│   │   ├── 18/
│   │   │   ├── 133119_knas_已成功运行.md  <-- 语义化后缀
│   │   │   ├── 142545_关于量化交易的思.md
│   │   │   └── ...

```

每个文件包含：

```markdown
---
sync_time: 2026-04-18 14:20:30
source: clipboard
---

[剪切板内容]
```

## 架构

```
┌─────────────────┐
│   Mac 本机      │
│                 │
│  ┌───────────┐  │    SSH    ┌─────────────┐
│  │ 剪切板监听 │  │ ────────> │  Ubuntu     │
│  │   (Go)    │  │            │  服务器     │
│  └───────────┘  │            └─────────────┘
│       ↓         │                  │
│  内容处理       │                  ↓
│  (URL标题抓取)  │          文件存储
│       ↓         │          (按日期归档)
│  ┌───────────┐  │
│  │ SSH 同步  │  │
│  └───────────┘  │
└─────────────────┘
```

## 技术栈

- **剪切板监听**: Go + `github.com/atotto/clipboard`
- **SSH 传输**: Go + `golang.org/x/crypto/ssh`
- **CLI 工具**: Node.js + Commander.js
- **系统服务**: macOS LaunchAgent

## 开发

```bash
# 克隆仓库
git clone https://github.com/yuanguangshan/knas.git
cd knas

# 安装依赖
npm install

# 构建二进制文件
npm run build

# 链接到全局
npm link

# 运行
knas init
knas start
```

## 注意事项

1. **SSH 密钥**: 确保 Mac 的 SSH 公钥已添加到服务器的 `~/.ssh/authorized_keys`
2. **网络连接**: 确保能访问远程服务器
3. **权限**: 确保远程路径有写入权限
4. **敏感信息**: 默认排除包含 "password"、"密码"、"token" 的内容

## 许可证

MIT
