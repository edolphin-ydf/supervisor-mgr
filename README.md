
# What is this?

this is a [go-supervisord](https://github.com/ochinchina/supervisord) xml-rpc-client command line tool which can show and control process on local or remote machine

# Install

go get -u github.com/edolphin-ydf/supervisor-mgr

# Usage

## config file

```yaml
servers:
  - name: test  # specify the server name, will be used in command
    url: http://127.0.0.1:9001  # the go-supervisord inet_http_server host port, the schema "http://" is required
    username: your username # the server's inet_http_server.username
    password: your password # the server's inet_http_server.password
```

## command
```
  supervisor-mgr [OPTIONS] <start | status | stop>

Application Options:
  -c, --config= specify the config file (default: config.yaml)

Help Options:
  -h, --help    Show this help message

Available commands:
  start   start process on server: supervisor-mgr start serverName processName[processName...]
  status  show process status: supervisor-mgr status [serverName...]
  stop    stop process on server: supervisor-mgr stop serverName processName[processName...]

```

# Thanks

thanks [ochinchina](https://github.com/ochinchina) rewrite supervisord in go and provide the xml-rpc-client package
