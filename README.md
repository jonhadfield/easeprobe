<h1>GuardianProbe</h1>

GuardianProbe is a simple, standalone, and lightweight tool that can do health/status checking, written in Go.

![](docs/overview.png)

<h2>Table of Contents</h2>

- [1. Introduction](#1-introduction)
  - [1.1 Probe](#11-probe)
  - [1.2 Notification](#12-notification)
  - [1.3 Report \& Metrics](#13-report--metrics)
- [2. Getting Started](#2-getting-started)
  - [2.1 Build](#21-build)
  - [2.2 Configure](#22-configure)
  - [2.3 Run](#23-run)
- [3. Deployment](#3-deployment)
- [4. User Manual](#4-user-manual)
- [5. Benchmark](#5-benchmark)
- [6. Contributing](#6-contributing)
- [7. Community](#7-community)
- [8. License](#8-license)


# 1. Introduction

GuardianProbe is designed to do three kinds of work - **Probe**, **Notify**, and **Report**.

## 1.1 Probe

GuardianProbe supports a variety of methods to perform its probes such as:

- **HTTP**. Checking the HTTP status code, Support mTLS, HTTP Basic Auth, setting Request Header/Body, and XPath response evaluation. ( [HTTP Probe Manual](./docs/Manual.md#12-http) )
- **TCP**. Check whether a TCP connection can be established or not. ( [TCP Probe Manual](./docs/Manual.md#13-tcp) )
- **Ping**. Ping a host to see if it is reachable or not. ( [Ping Probe Manual](./docs/Manual.md#14-ping) )
- **Shell**. Run a Shell command and check the result. ( [Shell Command Probe Manual](./docs/Manual.md#15-shell) )
- **SSH**. Run a remote command via SSH and check the result. Support the bastion/jump server ([SSH Command Probe Manual](./docs/Manual.md#16-ssh))
- **TLS**. Connect to a given port using TLS and (optionally) validate for revoked or expired certificates ( [TLS Probe Manual](./docs/Manual.md#17-tls) )
- **Host**. Run an SSH command on a remote host and check the CPU, Memory, and Disk usage. ( [Host Load Probe Manual](./docs/Manual.md#18-host) )
- **Client**. The following native clients are supported. They all support mTLS and data checking. ( [Native Client Probe Manual](./docs/Manual.md#19-native-client) )
  - **MySQL**. Connect to a MySQL server and run the `SHOW STATUS` SQL.
  - **Redis**. Connect to a Redis server and run the `PING` command.
  - **Memcache**. Connect to a Memcache server and run the `version` command or validate a given key/value pair.
  - **MongoDB**. Connect to a MongoDB server and perform a ping.
  - **Kafka**. Connect to a Kafka server and perform a list of all topics.
  - **PostgreSQL**. Connect to a PostgreSQL server and run `SELECT 1` SQL.
  - **Zookeeper**. Connect to a Zookeeper server and run `get /` command.

## 1.2 Notification

GuardianProbe supports notification delivery to the following:

- **Slack**. Using Slack Webhook for notification delivery
- **Discord**. Using Discord Webhook for notification delivery
- **Telegram**. Using Telegram Bot for notification delivery
- **Teams**. Support the [Microsoft Teams](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using?tabs=cURL#setting-up-a-custom-incoming-webhook) notification delivery
- **Email**. Support email notification delivery to one or more email addresses
- **AWS SNS**. Support the AWS Simple Notification Service
- **WeChat Work**. Support Enterprise WeChat Work notification delivery
- **DingTalk**. Support the DingTalk notification delivery
- **Lark**. Support the Lark(Feishu) notification delivery
- **SMS**. SMS notification delivery with support for multiple SMS service providers
  - [Twilio](https://www.twilio.com/sms)
  - [Vonage(Nexmo)](https://developer.vonage.com/messaging/sms/overview)
  - [YunPian](https://www.yunpian.com/official/document/sms/en/domestic_list?lang=en)
- **Log**. Write the notification into a log file or Syslog.
- **Shell**. Run a shell command to deliver the notification (see [example](resources/scripts/notify/notify.sh))
- **RingCentral**. Using RingCentral Webhook for notification delivery

> **Note**:
>
> 1) The notification is **Edge-Triggered Mode** by default, if you want to config it as **Level-Triggered Mode** with different interval and max notification, please refer to the manual - [Alerting Interval](./docs/Manual.md#112-alerting-interval).
>
> 2) Windows platforms do not support syslog as notification method.

Check the [Notification Manual](./docs/Manual.md#2-notification) to see how to configure it.

## 1.3 Report & Metrics

GuardianProbe supports the following report and metrics:

- **SLA Report Notify**. GuardianProbe would send the daily, weekly, or monthly SLA report using the defined **`notify:`** methods.
- **SLA Live Report**. The GuardianProbe would listen on the `0.0.0.0:8181` port by default. By accessing this service you will be provided with a live SLA report either as HTML at `http://localhost:8181/` or as JSON at `http://localhost:8181/api/v1/sla`
- **SLA Data Persistence**. The SLA data will be persisted in `$CWD/data/data.yaml` by default. You can configure this path by editing the `settings` section of your configuration file.

For more information, please check the [Global Setting Configuration](./docs/Manual.md#73-global-setting-configuration)

- **Prometheus Metrics**. The GuardianProbe would listen on the `8181` port by default. By accessing this service you will be provided with Prometheus metrics at `http://easeprobe:8181/metrics`.

The metrics are prefixed with `easeprobe_` and are documented in [Prometheus Metrics Exporter](./docs/Manual.md#6-prometheus-metrics-exporter)

# 2. Getting Started

You can get started with GuardianProbe, by any of the following methods:
* Download the release for your platform from https://github.com/o2ip/guardianprobe/releases
* Use the available GuardianProbe docker image `docker run -it o2ip/guardianprobe`
* Build `easeprobe` from sources

## 2.1 Build

Compiler `Go 1.21+` (Generics Programming Support), checking the [Go Installation](https://go.dev/doc/install) to see how to install Go on your platform.

Use `make` to build and produce the `easeprobe` binary file. The executable is produced under the `build/bin` directory.

```shell
$ make
```
## 2.2 Configure

Read the [User Manual](./docs/Manual.md) for detailed instructions on how to configure all GuardianProbe parameters.

Create a configuration file (eg. `$CWD/config.yaml`) using the configuration template at [./resources/config.yaml](https://raw.githubusercontent.com/o2ip/guardianprobe/main/resources/config.yaml), which includes the complete list of configuration parameters.

The following simple configuration example can be used to get started:

```YAML
http: # http probes
  - name: GuardianProbe Github
    url: https://github.com/o2ip/guardianprobe
notify:
  log:
    - name: log file # local log file
      file: /var/log/easeprobe.log
settings:
  probe:
    timeout: 30s # the time out for all probes
    interval: 1m # probe every minute for all probes
```

You can check the [GuardianProbe JSON Schema](./docs/Manual.md#81-easeprobe-json-schema) section to use a JSON Scheme file to make your life easier when you edit the configuration file.

## 2.3 Run

You can run the following command to start GuardianProbe once built

```shell
$ build/bin/easeprobe -f config.yaml
```
* `-f` configuration file or URL or path for multiple files which will be automatically merged into one. Can also be achieved by setting the environment variable `PROBE_CONFIG`
* `-d` dry run. Can also be achieved by setting the environment variable `PROBE_DRY`

# 3. Deployment

GuardianProbe can be deployed by Systemd, Docker, Docker-Compose, & Kubernetes.

You can find the details in [Deployment Guide](./docs/Deployment.md)

# 4. User Manual

For detailed instructions and features please refer to the [User Manual](./docs/Manual.md)

# 5. Benchmark

We have performed an extensive benchmark on GuardianProbe. For the benchmark results please refer to - [Benchmark Report](./docs/Benchmark.md)

# 6. Contributing

If you're interested in contributing to the project, please spare a moment to read our [CONTRIBUTING Guide](./docs/CONTRIBUTING.md)

# 7. Community

- Join Slack [Workspace](https://join.slack.com/t/openmegaease/shared_invite/zt-upo7v306-lYPHvVwKnvwlqR0Zl2vveA) for requirements, issues, and development.
- [MegaEase on Twitter](https://twitter.com/megaease)

# 8. License

GuardianProbe is under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
