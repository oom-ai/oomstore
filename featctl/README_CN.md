# featctl

用于管理特征仓库的 CLI 工具

## Usage

**查看特征详情**
```
$ featctl describe --group batch_device --name price
Name:           price
Group:          batch_device
Revision:       20210909
Status:         disabled
Category:       batch
ValueType:      int(11)
Description:    设备价格
RevisionsLimit: 3
CreateTime:     2021-09-10T15:20:43Z
ModifyTime:     2021-09-13T18:58:34Z
```

**点查询特征值**

```
$ featctl query -h
query feature values

Usage:
  featctl query [flags]

Examples:

1. featctl query -g user_info -n gender,city -k 1,2,3 
2. featctl query -g user_info -n gender,'user name' -k 1,2,3 


Flags:
  -g, --group string      feature group
  -h, --help              help for query
  -k, --key strings       entity keys
  -n, --name strings      feature names
  -r, --revision string   revision

$ featctl query --group batch_180d_userinfo -k 24031290,24036534,24039010 -n sex,city
entity_key,sex,city
24031290,0,上海
24036534,1,泰安
24039010,2,
```

**导出特征组**

- 将特征组 `batch_180d_userinfo` 的全部特征下载到 `/tmp/featctl/users.csv` ：

```sh
featctl export --group batch_180d_userinfo --output-file users.csv
```

**导入特征组**

- 将 `device.csv` 导入特征组 `batch_device` ：

1. 准备样例数据和 schema 模板
```sh
#!/usr/bin/env bash
set -euo pipefail

mkdir -p /tmp/featctl && cd /tmp/featctl

cat <<-EOF > device.csv
entity_key,brand,model,price
a9f0d6af575bb7e427fde2dcc81adbed,小米,MIX3,3999
134d9facd06ff355bf53846c0407d4f4,华为,P40,5299
0c66da7c680c4c44f33cb34881f1b104,苹果,IPHONE11,4999
EOF

cat <<-EOF > schema-template.sql
CREATE TABLE {{TABLE_NAME}} (
    entity_key VARCHAR(32) COMMENT '设备ID' PRIMARY KEY,
    brand      VARCHAR(16) COMMENT '设备厂商',
    model      VARCHAR(32) COMMENT '设备型号',
    price      INT         COMMENT '设备价格'
) SHARD_ROW_ID_BITS = 4 PRE_SPLIT_REGIONS = 3;
EOF
```

2. 执行导入
```sh
featctl import \
    --group batch_device \
    --revision 20210909 \
    --schema-template schema-template.sql \
    --input-file device.csv \
    --has-header \
    --separator "," \
    --description "test import version 6"
```

**创建特征配置**
```sh
featctl create feature \
    --name price \
    --group batch_device \
    --category batch \
    --revision 20210909 \
    --revisions-limit 3 \
    --description "设备价格"
```

**修改特征配置**
```sh
# 启用特征并将其版本指定为 20210909
featctl set \
    --name price \
    --group batch_device \
    --revision 20210909 \
    --status "enabled"
```

## Config

featctl 默认读取 `$XDG_CONFIG_HOME/featctl/config.yaml` 作为配置文件（可通过 `--config` 手动指定）：

```yaml
host: 127.0.0.1
port: 4000
user: <user>
pass: <password>
```

如果没有提供配置文件，也可以在执行命令时手动指定以上参数。

## Dependency

- [tidb-toolkit-v4.0.13](https://download.pingcap.org/tidb-toolkit-v4.0.13-linux-amd64.tar.gz):
下载解压后将其中的 `bin` 目录添加到环境变量 `PATH` 中：

```sh
export PATH="$PATH:/path/to/tidb-toolkit/bin"
```

生产环境建议采用容器的方式运行，docker 镜像已经打包好了依赖：

```sh
docker run --rm aiinfra/featctl:latest sh -c 'featctl ...'
```

如果需要和宿主机交换文件，可以将宿主机目录挂载到容器的 `/work` 目录。

## Development

Clone 项目之后先安装 [`pre-commit`](https://pre-commit.com/):

```sh
pip install pre-commit
pre-commit install
```

- 本地编译

```
make build
```

- 构建镜像

```
make image
```

- 推送镜像

```
make image-push
```

## TODO

- [x] `featctl export`
- [x] `featctl import`
- [x] `featctl create feature`
- [x] `featctl set`
- [x] `featctl describe`
- [ ] `featctl get revisions`
- [ ] `featctl get features`
