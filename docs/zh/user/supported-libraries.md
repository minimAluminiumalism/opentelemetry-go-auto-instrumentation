# Supported Libraries

## 数据库
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| database/sql   | https://pkg.go.dev/database/sql                             | -               | -            |
| gorm           | https://github.com/go-gorm/gorm                             | v1.22.0         | v1.25.9       |
| sqlx           | https://github.com/jmoiron/sqlx                             | v1.3.0          | v1.4.0        |
| gopg           | https://github.com/go-pg/pg                                 | v10.10.0        | v10.14.0      |
| mongodb        | https://github.com/mongodb/mongo-go-driver                  | v1.11.1         | v1.15.1       |
| elasticsearch  | https://github.com/elastic/go-elasticsearch                 | v8.4.0          | v8.15.0       |

## 缓存
| Library          | Repository Url                                           | Min Version     | Max Version |
|------------------|----------------------------------------------------------|-----------------|-------------|
| redis (go-redis) | https://github.com/redis/go-redis                        | v9.0.5          | v9.5.1      |
| redis v8         | https://github.com/go-redis/redis/v8                     | v8.11.0         | v8.11.5     |
| redigo           | https://github.com/gomodule/redigo                       | v1.9.0          | v1.9.3      |
| rueidis          | https://github.com/redis/rueidis                         | v1.0.30         | -           |

## 消息队列
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| rocketmq        | https://github.com/apache/rocketmq-client-go/v2             | v2.0.0          | -            |
| amqp091         | https://github.com/rabbitmq/amqp091-go                      | v1.10.0         | -            |
| segmentio/kafka-go| https://github.com/segmentio/kafka-go                     | v0.4.0          | -            |

## RPC/通信框架
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| grpc            | https://google.golang.org/grpc                              | v1.44.0         | -            |
| dubbo-go        | https://github.com/apache/dubbo-go                          | v3.3.0          | -            |
| kitex           | https://github.com/cloudwego/kitex                          | v0.5.1          | -            |
| kratos          | https://github.com/go-kratos/kratos                         | v2.6.3          | -            |
| go-micro        | https://github.com/micro/go-micro                           | v5.0.0          | v5.3.0        |
| trpc-go         | https://github.com/trpc-group/trpc-go                       | v1.0.0          | -            |

## HTTP/Web 框架
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| net/http        | https://pkg.go.dev/net/http                                 | -               | -            |
| echo            | https://github.com/labstack/echo                            | v4.0.0          | -            |
| gin             | https://github.com/gin-gonic/gin                            | v1.7.0          | v1.10.1       |
| fiber           | https://github.com/gofiber/fiber                            | v2.43.0         | v2.52.9       |
| fasthttp        | https://github.com/valyala/fasthttp                         | v1.45.0         | v1.65.0       |
| gorilla/mux     | https://github.com/gorilla/mux                              | v1.3.0          | v1.8.1        |
| iris            | https://github.com/kataras/iris                             | v12.2.0         | v12.2.11      |
| hertz           | https://github.com/cloudwego/hertz                          | v0.8.0          | -            |
| go-restful      | https://github.com/emicklei/go-restful                      | v3.7.0          | v3.12.1       |
| gorestful/v3    | https://github.com/emicklei/go-restful/v3                   | v3.7.0          | v3.12.1       |

## 配置/注册中心
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| nacos           | https://github.com/nacos-group/nacos-sdk-go/v2              | v2.0.0          | v2.2.9        |
| k8s client-go   | https://github.com/kubernetes/client-go                     | v0.33.3         | -            |

## 日志
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| log             | https://pkg.go.dev/log                                      | -               | -            |
| zap             | https://github.com/uber-go/zap                              | v1.20.0         | v1.27.0       |
| logrus          | https://github.com/sirupsen/logrus                          | v1.5.0          | v1.9.3        |
| zerolog         | https://github.com/rs/zerolog                               | v1.10.0         | v1.33.0       |
| go-kit/log      | https://github.com/go-kit/log                               | v0.1.0          | v0.2.1        |

## AI/LLM
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| langchaingo     | https://github.com/tmc/langchaingo                          | v0.1.13         | -            |
| ollama          | https://github.com/ollama/ollama                            | v0.3.14         | -            |
| eino            | https://github.com/cloudwego/eino                           | v0.3.51         | -            |

## 限流/熔断
| Library         | Repository Url                                               | Min Version     | Max Version   |
|----------------|-------------------------------------------------------------|-----------------|--------------|
| sentinel        | https://github.com/alibaba/sentinel-golang                  | v1.0.4          | -            |