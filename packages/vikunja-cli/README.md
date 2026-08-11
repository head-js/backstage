# Vikunja CLI

`backstage-vikunja` 是 Vikunja HTTP API 的轻量命令行入口。当前只实现 `/version` 系统接口，不包含业务能力。

## 配置

```dotenv
BACKSTAGE_VIKUNJA_URL=http://127.0.0.1:23456
BACKSTAGE_VIKUNJA_TOKEN=your_token_here
```

`BACKSTAGE_VIKUNJA_URL` 只配置 Vikunja 服务根地址，不包含 API 版本路径。Adapter 会统一追加 `/api/v2`。

## 使用

```bash
backstage-vikunja api GET /version
```

`GET /version` 是 `backstage-vikunja` 提供的系统接口。它实际请求 Vikunja API v2 的 `GET /info`，并原样输出完整 Info JSON，不会只提取其中的 `version` 字段。

其他 method 或路径当前均返回 `unsupported path`。

## 开发

```bash
go build ./...
```
