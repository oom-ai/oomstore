# apply 命令文档

apply 的效果是：接受一个资源（entity，group，feature）文件，使用 apply 将数据库中的资源设置成资源文件声明的样子（通过创建和更新）。但对于没有声明的字段，不做修改。

下面是几种资源文件的例子：

```yaml
kind: Entity
name: user
length: 8
description: 'description'
batch-features:
- group: device
  description: a description
  features:
  - name: model
    db-value-type: varchar(16)
    description: 'description'
  - name: price
    db-value-type: int
    description: 'description'
- group: user
  description: a description
  features:
  - name: age
    db-value-type: int
    description: 'description'
  - name: gender
    db-value-type: int
    description: 'description'
---
kind: Entity
name: device
length: 16
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
---
kind: Group
name: device
entity-name: user
category: batch
description: 'description'
features:
- name: model
  db-value-type: varchar(16)
  description: 'description'
- name: price
  db-value-type: int
  description: 'description'
- name: radio
  db-value-type: int
  description: 'description'
---
kind: Feature
name: model
group-name: device
category: batch
db-value-type: varchar(16)
description: 'description'
---
```

问题：
* 每种资源文件需要对应到一个 struct，这个 struct 长什么样子？
* 如何将资源文件正确解析到正确的 struct？
* 拿到 struct 后，如何和数据库中的数据对比，找到需要新增和更新的 action，正确拼装参数，调用相应 api？
