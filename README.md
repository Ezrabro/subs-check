<h1 align="center">🚀 Subscription Detection & Conversion Tool</h1>

<p align="center">
	<a href="https://github.com/Ezrabro/subs-check/releases"><img src="https://img.shields.io/github/v/release/Ezrabro/subs-check?style=flat-square&include_prereleases&label=version" /></a>
	<a href="https://github.com/Ezrabro/subs-check/releases"><img src="https://img.shields.io/github/downloads/Ezrabro/subs-check/total.svg?style=flat-square" /></a>
  <a href="https://hub.docker.com/r/beck8/subs-check/tags"><img src="https://img.shields.io/docker/pulls/beck8/subs-check" /></a>
	<a href="https://github.com/Ezrabro/subs-check/issues"><img src="https://img.shields.io/github/issues-raw/Ezrabro/subs-check.svg?style=flat-square&label=issues" /></a>
	<a href="https://github.com/Ezrabro/subs-check/graphs/contributors"><img src="https://img.shields.io/github/contributors/Ezrabro/subs-check?style=flat-square" /></a>
	<a href="https://github.com/Ezrabro/subs-check/blob/master/LICENSE"><img src="https://img.shields.io/github/license/Ezrabro/subs-check?style=flat-square" /></a>
</p>

---

> **✨ Fixed logic, simplified operations, added features, memory saving, one-click start without config**

> **⚠️ Note:** Features updated frequently, check latest [config file](https://github.com/Ezrabro/subs-check/blob/master/config/config.example.yaml)for latest features.  
> **⚠️ Note:** To view feature updates, refer to [config commit history](https://github.com/Ezrabro/subs-check/commits/master/config/config.example.yaml), with change notes on new features/logic updates

## 📸 Preview


![preview](./doc/images/preview.png)  
![result](./doc/images/results.png)  
![admin](./doc/images/admin.png)
| | |
|---|---|
| ![tgram](./doc/images/tgram.png) | ![dingtalk](./doc/images/dingtalk.png)  |

## ✨ Features

- **🔗 Subscription Merging**
- **🔍 Node Availability Check**
- **🗑️ Node Deduplication**
- **⏱️ Node Speed Test**
- **🎬 Streaming Platform Unlock Detection**
- **✏️ Node Renaming**
- **🔄 Any Format Subscription Conversion**
- **🔔 100+ Notification Channels Supported**
- **🌐 Built-in Sub-Store**
- **🖥️ Web Control Panel**
- **⏰ Supports Crontab Expressions**
- **🖥️ Multi-Platform Support**

## 🛠️ Deployment & Usage 
> First run creates default config file in current directory.

### 🚀 One-Click Install (Linux)

```bash
# Default Install
bash <(curl -fsSL https://raw.githubusercontent.com/Ezrabro/subs-check/master/install.sh)

# Using wget
bash <(wget -qO- https://raw.githubusercontent.com/Ezrabro/subs-check/master/install.sh)

# If GitHub is unreachable, use proxy
bash <(curl -fsSL https://ghfast.top/https://raw.githubusercontent.com/Ezrabro/subs-check/master/install.sh) https://ghfast.top/

# Alpine or environments without bash
wget -qO /tmp/install.sh https://raw.githubusercontent.com/Ezrabro/subs-check/master/install.sh && sh /tmp/install.sh && rm -f /tmp/install.sh
```

<details>
  <summary>Script Description</summary>

Install script performs automatically:
1. Detect system architecture（x86_64 / aarch64 / armv7 / i386）
2. Download latest version from GitHub Releases
3. Install to `/opt/subs-check` directory
4. Configure systemd service
5. Interactive choice: auto-start on boot
6. Interactive choice: start immediately

**Service Management:**
```bash
systemctl start subs-check    # Start
systemctl stop subs-check     # Stop
systemctl restart subs-check  # Restart
systemctl status subs-check   # Status
journalctl -u subs-check -f   # Logs
```

**Uninstall Method:**
```bash
systemctl stop subs-check
systemctl disable subs-check
rm -rf /opt/subs-check /etc/systemd/system/subs-check.service
systemctl daemon-reload
```

</details>

### 🪜 Proxy Settings (Optional)
<details>
  <summary>Show Details</summary>

If non-GitHub subscription download is slow，Can use HTTP_PROXY HTTPS_PROXY env variables to speed up; this variable does not affect node test speed
```bash
# HTTP Proxy Example
export HTTP_PROXY=http://username:password@192.168.1.1:7890
export HTTPS_PROXY=http://username:password@192.168.1.1:7890

# SOCKS5 Proxy Example
export HTTP_PROXY=socks5://username:password@192.168.1.1:7890
export HTTPS_PROXY=socks5://username:password@192.168.1.1:7890

# SOCKS5H Proxy Example
export HTTP_PROXY=socks5h://username:password@192.168.1.1:7890
export HTTPS_PROXY=socks5h://username:password@192.168.1.1:7890
```
To speed up GitHub links，Can use public github proxies，Or use worker.js from custom speed-test address below
```
# GitHub Proxy for fetching subscriptions, must include trailing /
# github-proxy: "https://ghfast.top/"
github-proxy: "https://custom-domain/raw/"
```

</details>

### 🌐 Custom Speed Test Address (Optional)
<details>
  <summary>Show Details</summary>

> **⚠️ Note:** Avoid Speedtest or Cloudflare download links — some nodes block speed test sites.

1. Deploy [worker.js](./doc/cloudflare/worker.js) to Cloudflare Workers.
2. Bind a custom domain (to avoid being blocked by nodes).
3. In config file set `speed-test-url` to your Workers address:

```yaml
# 100MB
speed-test-url: https://custom-domain/speedtest?bytes=104857600
# 1GB
speed-test-url: https://custom-domain/speedtest?bytes=1073741824
```

</details>

### 🐳 Docker Run

> **⚠️ Note:**  
> - Use `--memory="500m"` to limit memory.  
> - Can set Web Control Panel API Key via env `API_KEY`.

```bash
# Basic Run
docker run -d \
  --name subs-check \
  -p 8299:8299 \
  -p 8199:8199 \
  -v ./config:/app/config \
  -v ./output:/app/output \
  --restart always \
  ghcr.io/Ezrabro/subs-check:latest

# Run with Proxy
docker run -d \
  --name subs-check \
  -p 8299:8299 \
  -p 8199:8199 \
  -e HTTP_PROXY=http://192.168.1.1:7890 \
  -e HTTPS_PROXY=http://192.168.1.1:7890 \
  -v ./config:/app/config \
  -v ./output:/app/output \
  --restart always \
  ghcr.io/Ezrabro/subs-check:latest
```

### 📜 Docker Compose

```yaml
version: "3"
services:
  subs-check:
    image: ghcr.io/Ezrabro/subs-check:latest
    container_name: subs-check
    volumes:
      - ./config:/app/config
      - ./output:/app/output
    ports:
      - "8299:8299"
      - "8199:8199"
    environment:
      - TZ=Asia/Shanghai
      # - HTTP_PROXY=http://192.168.1.1:7890
      # - HTTPS_PROXY=http://192.168.1.1:7890
      # - API_KEY=subs-check
    restart: always
    network_mode: bridge
```
### 📦 Binary File Run

Download [Releases](https://github.com/Ezrabro/subs-check/releases) appropriate version, extract and run directly.

### 🖥️ Source Code Run

```bash
go run . -f ./config/config.yaml
```

## 🔔 Notification Channel Config (Optional)
<details>
  <summary>Show Details</summary>

> **📦 Supports 100+ notification channels**，通过 [Apprise](https://github.com/caronc/apprise) 发送通知。

### 🌐 Vercel Deployment

1. Click [**here**](https://vercel.com/new/clone?repository-url=https://github.com/Ezrabro/apprise_vercel)部署 Apprise。
2. After deployment get API link, e.g. `https://testapprise-beck8s-projects.vercel.app/notify`.
3. Recommend setting custom domain `diydomain.com` for Vercel project (access may be restricted in China).

### 🐳 Docker Deployment

> **⚠️ Note:** arm/v7 not supported.

```bash
# Basic Run
docker run --name apprise -p 8000:8000 --restart always -d caronc/apprise:latest

# Run with Proxy
docker run --name apprise \
  -p 8000:8000 \
  -e HTTP_PROXY=http://192.168.1.1:7890 \
  -e HTTPS_PROXY=http://192.168.1.1:7890 \
  --restart always \
  -d caronc/apprise:latest
```

### 📝 Configure Notifications in Config File

```yaml
# Fill in the built apprise API server address
# https://notify.xxxx.us.kg/notify
apprise-api-server: "https://diydomain.com/notify"
# Fill in notification targets
# Supports 100+ notification channels, detailed format at https://github.com/caronc/apprise
recipient-url: 
  # Telegram format:tgram://{bot_token}/{chat_id}
  # - tgram://xxxxxx/-1002149239223
  # DingTalk format:dingtalk://{Secret}@{ApiKey}
  # - dingtalk://xxxxxx@xxxxxxx
# Custom Notification Title
notify-title: "🔔 Node Status Update"
```
</details>

## 💾 Save Method Config

> **⚠️ Note:** When selecting save method, change `save-method` config.

- **本地保存**：保存到 `./output` 文件夹。
- **R2**：保存到 Cloudflare R2 [配置方法](./doc/r2.md)。
- **Gist**：保存到 GitHub Gist [配置方法](./doc/gist.md)。
- **WebDAV**：保存到 WebDAV 服务器 [配置方法](./doc/webdav.md)。
- **S3**：保存到 S3 对象存储。

## 📲 Subscription Usage Method

> **💡 提示：** Built-in Sub-Store，可生成多种订阅格式；高级玩家可DIY很多功能

**🚀 Universal Subscription**
```bash
# Universal Subscription
http://127.0.0.1:8299/download/sub

# URI Subscription
http://127.0.0.1:8299/download/sub?target=URI

# Mihomo/ClashMeta
http://127.0.0.1:8299/download/sub?target=ClashMeta

# Clash
http://127.0.0.1:8299/download/sub?target=Clash

# V2Ray
http://127.0.0.1:8299/download/sub?target=V2Ray

# ShadowRocket
http://127.0.0.1:8299/download/sub?target=ShadowRocket

# Quantumult
http://127.0.0.1:8299/download/sub?target=QX

# Sing-Box
http://127.0.0.1:8299/download/sub?target=sing-box

# Surge
http://127.0.0.1:8299/download/sub?target=Surge

# Surfboard
http://127.0.0.1:8299/download/sub?target=Surfboard
```

**🚀 Mihomo/Clash Subscription (with rules):**
> 默认使用 `https://raw.githubusercontent.com/Ezrabro/override-hub/refs/heads/main/yaml/ACL4SSR_Online_Full.yaml` overwrite  
Can change `mihomo-overwrite-url` in config.
```bash
http://127.0.0.1:8299/api/file/mihomo
```

## 🌐 Built-in Port Description
> subs-check saves 3 files to output after testing; all files in output served by port 8199

| 服务地址                        | 格式说明                | 来源说明|
|-------------------------------|-------------------|----|
| `http://127.0.0.1:8199/sub/all.yaml`   | Clash format nodes | Generated directly by subs-check|
| `http://127.0.0.1:8199/sub/mihomo.yaml`| Mihomo/Clash subscription with split rules | Converted/downloaded by sub-store above|
| `http://127.0.0.1:8199/sub/base64.txt` | Base64 format subscription | Converted/downloaded by sub-store above|
| `http://127.0.0.1:8199/export/surge` 等 | Surge / Loon / Quantumult X / Shadowrocket / Stash / Surfboard / Egern / Clash.Meta / Clash / sing-box / URI Subscription |在 Web 控制面板生成后提供|

> Results & Export: After check completes, click 'Speed Test' card in panel to enter `/admin/results`，Can filter/sort nodes (protocol, server, SNI, TLS/UDP, stream unlock, speed) and generate subscription links per client.
> - Only formats previously generated in panel (requires API key) accessible via `/export/<format>`; public access won't trigger sub-store conversion
> - 生成过的格式每轮检测完成后自动更新，程序Restart后继续有效；在「导出订阅」中停用后链接失效
> - Requires sub-store enabled (`sub-store-port`)
> - Result snapshots and subscriptions saved in `cache/` beside config (Docker mount `/app/config` also persists); files contain node credentials — don't expose directory

## 🗺️ Architecture Diagram
<details>
  <summary>Show Details</summary>

```mermaid
graph TD
    A[Subscription Link] -->|获取Subscription Link| B[subs-check]
    subgraph subs-check 处理流程
        B -->|转成 YAML 格式| B1[Node Deduplication]
        B1 -->|Remove Redundant Nodes| B2[Liveness Check]
        B2 -->|Node Available| B3[Streaming + Rename]
        B2 -->|Node Unavailable| X[Discard]
        B3 -->|Filter Passed| B4[Speed Test]
        B3 -->|Filter Not Passed| X[Discard]
        B4 -->|Speed Test达标| B5[Generate all.yaml]
        B4 -->|Speed Test不达标| X[Discard]
    end
    B5 -->|Save to output directory| C[output 目录]
    B5 -->|Upload all.yaml| D[sub-store]
    C -->|Save to various locations| H1[R2/Gist/WebDAV/S3]
    H1 -->|Storage Complete| H2[Send Notification Message]
    D -->|Provide Subscription Conversion Service| E[sub-store Conversion Service]
    subgraph sub-store 独立功能
        E -->|Generate Config File| E1[mihomo.yaml, base64.txt]
        E -->|Other Format Conversion| E2[Clash, V2Ray, ShadowRocket 等]
        E -->|Subscription Sharing| E3[分享Subscription Link]
    end
    E1 -->|Save to output directory| C
    C -->|File Service| F[8199 端口: /sub]
    B -->|Web Management| G[8199 端口: /admin]
``` 

</details>

## 🙏 Credits
[cmliu](https://github.com/cmliu)、[Sub-Store](https://github.com/sub-store-org/Sub-Store)、[bestruirui](https://github.com/bestruirui/BestSub)、[1password](https://1password.com/)、[ipinfo.io](https://ipinfo.io/)

## ⭐ Star History

[![Stargazers over time](https://starchart.cc/Ezrabro/subs-check.svg?variant=adaptive)](https://starchart.cc/Ezrabro/subs-check)

## ⚖️ Disclaimer

This tool is for learning and research use only; users assume all risks and must comply with relevant laws.
