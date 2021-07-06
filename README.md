# AlertManager-Notifier

接收AlertManager WebHook告警消息并通过指定通道发送。

## 启动参数

```
Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --config.file="./conf/settings.yaml"  
                              AlertManager-Notifier configuration file path. Default is ./conf/settings.yaml
      --web.port=":8080"      Port to listen on for the web interface and API. Default is :8080
      --log.max-backups=5     The maximum number of old log files to retain. Default is 5.
      --log.max-days=30       The maximum number of days to retain old log files. Default is 30
      --log.level=info        Log level, default is info. One of: [debug, info, warn, error]
      --running.mode=release  Running Mode,default is release. One of: [dev, debug, release]
      --version               Show application version.

```

## AlertManager Receiver 配置

``` yaml
receivers:
  - name: system:operator
    webhook_configs:
      - http://127.0.0.1:8080/api/v1/alert
```


## 配置项说明

配置参考 conf/settings.yaml

存储数据库：
- sqlite
- postgres

示例
```yaml
userota: false
static_login: dbb342c8604b24b466a1920002a14858 # 是否从数据库校验用户登陆，密码使用md5编码
database:
  sqlite:
    datapath: 'data/alertmanager_notifier.db'
#  postgres:
#    host: '127.0.0.1'
#    port: '5432'
#    dbname: 'alertmanager_notifier'
#    user: 'postgres'
#    password: 'abc123!'
#    sslmode: 'disable'
```

接收者配置
```yaml
receivers:
  - name: test # 告警通道名称
    receiver_type: uid # 接收名称类型
    shell_config:     # 接收通道配置
      - command: echo
        args:
          - {{ .Receiver }}
          - {{ $labels.instance }}
```

参数模版
```
{{ uuid }} -> 9f59d7bb-6280-48f6-81ee-8e3605061102
{{ uuid32 }} -> b49fad6efd9248b9b56aef631e03a099
{{ transTime $alert.StartsAt "2006-01-02 15:04:05" }}

{{ .Name }}
{{ $alert.<config> }}
{{ $labels.<label_name> }}
```


*In order to learn golang*