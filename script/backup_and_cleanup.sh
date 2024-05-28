#!/bin/bash

# 设置备份目录
BACKUP_ROOT="/data/harbor_backups"
BACKUP_DATA_DIR="/data/harbor_backups/artifacts"

# 设置保留天数
RETENTION_DAYS=30

# 从环境变量中获取 harbor host 和认证信息
HARBOR_BASEURL="${HARBOR_BASEURL}"
HARBOR_AUTH="${HARBOR_AUTH}"

# 加载变量
source ~/.bashrc

# 检查环境变量是否设置
if [ -z "$HARBOR_BASEURL" ] || [ -z "$HARBOR_AUTH" ]; then
  echo "Error: HARBOR_BASEURL or HARBOR_AUTH environment variables are not set."
  exit 1
fi

# 根据日期决定执行全量备份或差量备份
DAY_OF_WEEK=$(date +%u)

# 注意脚本 harbor_api_mario 中使用的是相对路径 所以备份时一定要切换到指定路径执行
# 不能使用绝对路径执行 harbor_api_mario 否则会在执行脚本的路径生成备份文件夹
if [ "$DAY_OF_WEEK" -eq 7 ]; then
  # 周日执行全量备份
  echo "Executing full backup..."
  cd "$BACKUP_ROOT" && ./harbor_api_mario --action full_backup
  if [ $? -ne 0 ]; then
    echo "Error: Full backup failed."
    exit 1
  fi
else
  # 周一到周六执行差量备份
  echo "Executing delta backup..."
  cd "$BACKUP_ROOT" && ./harbor_api_mario --action delta_backup
  if [ $? -ne 0 ]; then
    echo "Error: Delta backup failed."
    exit 1
  fi
fi

# 清理指定天数之前的备份文件
echo "Cleaning up old backups..."
find "$BACKUP_DATA_DIR" -type d -mtime +"$RETENTION_DAYS" -exec rm -rf {} \;

# 记录清理操作
echo "Cleanup completed on $(date)" >> "$BACKUP_ROOT"/cleanup.log

echo "Backup and cleanup completed successfully."

#0 2 * * * /data/harbor_backups/backup_and_cleanup.sh >> /data/harbor_backups/backup_and_cleanup.log 2>&1
