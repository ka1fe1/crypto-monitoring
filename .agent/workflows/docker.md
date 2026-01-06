---
description: build & push docker image
---

按照如下步骤 build & push docker image

1. 使用当前北京时间格式年月日时分 (如 `2512171818`) 作为 tag 来构建 docker image，如:

```bash
  docker build -t crypto-monitoring:{tag} .
```

2. 打 tag

```bash
  docker tag crypto-monitoring:{tag} tataka1takes2/crypto-monitoring:{tag}
```

3. 推送 tag
```bash
  docker push tataka1takes2/crypto-monitoring:{tag}
```