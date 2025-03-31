![](anim-logo.svg)

[![](https://shields.io/badge/Docs-English-blue?logo=Read%20The%20Docs)](../README.md)
[![](https://shields.io/badge/Readme-中文-blue?logo=Read%20The%20Docs)](./README_CN.md)
[![codecov](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation/branch/main/graph/badge.svg)](https://codecov.io/gh/alibaba/opentelemetry-go-auto-instrumentation)

该项目为希望利用 OpenTelemetry 的 Golang 应用程序提供了一个自动解决方案。
利用 OpenTelemetry 实现有效可观察性的 Golang 应用程序提供自动解决方案。目标应用程序无需更改代码
在编译时完成。只需在 `go build` 中添加 `otel` 前缀即可开始 :rocket:

# 安装

### 通过 Bash 安装
对于 **Linux 和 MacOS** 用户，运行以下命令即可安装该工具
```bash
$ sudo curl -fsSL https://cdn.jsdelivr.net/gh/alibaba/opentelemetry-go-auto-instrumentation@main/install.sh | sudo bash
```
默认情况下，它将安装在 `/usr/local/bin/otel`中。

### 预编译二进制文件

请从
[Release](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/releases)
页面下载最新的预编译版本。

### 从源代码编译

通过运行以下命令查看源代码并构建工具：

```bash
$ make         # 只构建
$ make install # 构建并安装
```

# 开始

通过运行以下命令检查版本：
```bash
$ otel version
```

通过以下命令配置工具参数：
```bash
$ otel set -verbose                          # 打印详细日志
$ otel set -log=/path/to/file.log            # 设置日志文件路径
$ otel set -debug                            # 启用调试模式
$ otel set -debug -verbose -rule=custom.json # 组合配置参数
$ otel set -disabledefault -rule=custom.json # 禁用默认规则，仅使用自定义规则
$ otel set -rule=custom.json                 # 同时使用默认和自定义规则
$ otel set -rule=a.json,b.json               # 使用默认规则及 a 和 b 自定义规则
```

在 `go build` 中添加 `otel` 前缀，以构建项目：

```bash
$ otel go build
$ otel go build -o app cmd/app
$ otel go build -gcflags="-m" cmd/app
```
工具本身的参数应放在 `go build` 之前：

```bash
$ otel -help # 打印帮助文档
$ otel -debug go build # 启用调试模式
$ otel -verbose go build # 打印详细日志
$ otel -rule=custom.json go build # 使用自定义规则
```

您可以在 [**使用指南**](./usage.md)中找到 `otel` 工具的详细用法。

> [!NOTE] 
> 如果您发现任何编译失败，而 `go build` 却能正常工作，这很可能是一个 bug。
> 请随时在
> [GitHub Issues](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues)
> 提交问题报告以帮助我们改进本项目。

# 示例

您还可以探索 [**这些示例**](../example/) 以获得实践经验。

此外，还有一些 [**文档**](../docs)，您可能会发现它们对了解项目或为项目做出贡献非常有用。

# 支持的库

| 插件名称       | 存储库网址                                      | 最低支持版本           | 最高支持版本     |
|---------------| ---------------------------------------------- |-----------------------|-----------------------|
| database/sql  | https://pkg.go.dev/database/sql                | -                     | -                     |
| echo          | https://github.com/labstack/echo               | v4.0.0                | v4.12.0               |
| elasticsearch | https://github.com/elastic/go-elasticsearch    | v8.4.0                | v8.15.0               |
| fasthttp      | https://github.com/valyala/fasthttp            | v1.45.0               | v1.59.0               |
| fiber         | https://github.com/gofiber/fiber               | v2.43.0               | v2.52.6               |
| gin           | https://github.com/gin-gonic/gin               | v1.7.0                | v1.10.0               |
| go-redis      | https://github.com/redis/go-redis              | v9.0.5                | v9.5.1                |
| go-redis v8   | https://github.com/redis/go-redis              | v8.11.0               | v8.11.5               |
| gomicro       | https://github.com/micro/go-micro              | v5.0.0                | v5.3.0                |
| gorestful     | https://github.com/emicklei/go-restful         | v3.7.0                | v3.12.1               |
| gorm          | https://github.com/go-gorm/gorm                | v1.22.0               | v1.25.9               |
| grpc          | https://google.golang.org/grpc                 | v1.44.0               | v1.71.0               |
| hertz         | https://github.com/cloudwego/hertz             | v0.8.0                | v0.9.2                |
| iris          | https://github.com/kataras/iris                | v12.2.0               | v12.2.11              |
| kitex         | https://github.com/cloudwego/kitex             | v0.5.1                | v0.11.3               |
| kratos        | https://github.com/go-kratos/kratos            | v2.6.3                | v2.8.4                |
| langchaingo   | https://github.com/tmc/langchaingo             | v0.1.13               | v0.1.13               |
| log           | https://pkg.go.dev/log                         | -                     | -                     |
| logrus        | https://github.com/sirupsen/logrus             | v1.5.0                | v1.9.3                |
| mongodb       | https://github.com/mongodb/mongo-go-driver     | v1.11.1               | v1.15.1               |
| mux           | https://github.com/gorilla/mux                 | v1.3.0                | v1.8.1                |
| nacos         | https://github.com/nacos-group/nacos-sdk-go/v2 | v2.0.0                | v2.2.7                |
| net/http      | https://pkg.go.dev/net/http                    | -                     | -                     |
| redigo        | https://github.com/gomodule/redigo             | v1.9.0                | v1.9.2                |
| slog          | https://pkg.go.dev/log/slog                    | -                     | -                     |
| trpc-go       | https://github.com/trpc-group/trpc-go          | v1.0.0                | v1.0.3                |
| zap           | https://github.com/uber-go/zap                 | v1.20.0               | v1.27.0               |
| zerolog       | https://github.com/rs/zerolog                  | v1.10.0               | v1.33.0               |

我们正在逐步开源我们支持的库，非常欢迎您的贡献💖！

> [!IMPORTANT]
> 您期望的框架不在列表中？别担心，您可以轻松地将代码注入到任何官方不支持的框架/库中。
>
> 请参考 [这个文档](./how-to-add-a-new-rule.md) 开始使用。

# 文档

- [如何添加新规则](./how-to-add-a-new-rule.md)
- [如何编写插件测试](./how-to-write-tests-for-plugins.md)
- [兼容性说明](./compatibility.md)
- [实现原理](./how-it-works.md)
- [如何调试](./how-to-debug.md)
- [上下文传播机制](./context-propagation.md)
- [支持的库](./supported-libraries.md)
- [基准测试](../example/benchmark/benchmark.md)
- [OpenTelemetry社区讨论主题](https://github.com/open-telemetry/community/issues/1961)
- [面向OpenTelemetry的Golang应用无侵入插桩技术](https://mp.weixin.qq.com/s/FKCwzRB5Ujhe1stOH2ibXg)

# 社区

我们期待您的反馈和建议。您可以加入我们的 [DingTalk 群组](https://qr.dingtalk.com/action/joingroup?code=v1,k1,GyDX5fUTYnJ0En8MrVbHBYTGUcPXJ/NdsmLODGibd0w=&_dt_no_comment=1&origin=11? )
与我们交流。

<img src="dingtalk.png" height="200">

# 应用案例

以下为部分采用本项目的企业列表，仅供参考。如果您正在使用此项目，请[在此处添加您的公司](https://github.com/alibaba/opentelemetry-go-auto-instrumentation/issues/225)告诉我们您的使用场景，让这个项目变得更好。

- <img src="./alibaba.png" width="80">
- <img src="./aliyun.png" width="100">

# Contributors

<a href="https://github.com/alibaba/opentelemetry-go-auto-instrumentation/graphs/contributors">
  <img alt="contributors" src="https://contrib.rocks/image?repo=alibaba/opentelemetry-go-auto-instrumentation"/>
</a>

# Star History

[![Star History](https://api.star-history.com/svg?repos=alibaba/opentelemetry-go-auto-instrumentation&type=Date)](https://star-history.com/#alibaba/opentelemetry-go-auto-instrumentation&Date)

<p align="right" style="font-size: 14px; color: #555; margin-top: 20px;">
    <a href="#安装" style="text-decoration: none; color: #007bff; font-weight: bold;">
        ↑ 返回顶部 ↑
    </a>
</p>