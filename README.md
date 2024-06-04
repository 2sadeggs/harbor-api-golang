# Harbor Backup Tool

Harbor Backup Tool 是一个用于备份和管理 Harbor 仓库的工具，支持全量备份和差量备份。

## 功能
- `ping`：检查 Harbor 是否可用。
- `health`：检查 Harbor API 健康状态。
- `statistics`：获取 Harbor 统计信息。
- `projects`：获取所有项目。
- `repositories`：获取所有仓库。
- `artifacts`：获取所有制品。
- `uris`：获取所有 URI 列表。
- `pull`：下载制品。
- `save`：下载并保存制品。
- `full_backup`：全量备份。
- `delta_backup`：差量备份，并生成差异列表清单。
  
## 环境变量设置

```bash
# 设置环境变量 HARBOR_BASEURL
export HARBOR_BASEURL="your_harbor_base_url"

# 设置环境变量 HARBOR_AUTH
export HARBOR_AUTH="your_harbor_auth"
