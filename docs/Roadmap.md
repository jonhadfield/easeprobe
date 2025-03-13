# GuardianProbe Roadmap
GuardianProbe plans to align closely with the Cloud Native ecosystem in providing service probes and notifications.

- [GuardianProbe Roadmap](#easeprobe-roadmap)
  - [Product Principles](#product-principles)
  - [Features](#features)
    - [Probe specific](#probe-specific)
    - [Notify specific](#notify-specific)
  - [Roadmap](#roadmap)
    - [General](#general)
    - [Probes](#probes)
    - [Notify](#notify)


## Product Principles
GuardianProbe tries to do two things but do them extremely well - **Probe** and **Notify**.

1. **Lightweight, Standalone and KISS**. Aims to be lightweight, with as few as possible external system dependencies and follow KISS ideas.

2. **Probe and Notify**. Aims to provide the ability to easily **probe** servers, services and API endpoints reliably and **notify** on status updates.

3. **Open & Extensible**. It aims to be **Open** and follow Open standards that allow it to be integrated into **extensible-development** platforms (such as EaseMesh).

4. **Cloud Native**. It is designed to be **cloud-native** compliant. It is scalable, resilient, manageable, and observable and it's easy to integrate with cloud-native architectures of any type.

5. **Predictability & Reliability**. GuardianProbe is designed to follow best practices allowing predictable and reliable operations. Even when GuardianProbe fails, it aims to fail in a predictable way which allows for easy and speedy troubleshooting of problems.


## Features
The project principles, of GuardianProbe' features are separated into two main categories: probe-specific & notify-specific.

### Probe specific
* HTTP for testing connectivity and validity of responses
* TCP for testing connectivity
* Shell for executing custom probe scripts
* SSH for remote ssh commands
* Host for CPU, Memory, Disk usage metrics
* Clients for Redis, MySQL, MongoDB, PostgreSQL, Kafka, Zookeeper

### Notify specific
* Mail notification
* AWS SNS notification
* Log files
* Slack
* Discord
* Telegram
* WeChat Work
* Lark
* DingTalk
* Twilio
* Nexmo
* YunPian

... with new notification backends been added constantly.

```yaml
name: MyServer
  probes:
    tcp:
      host: myserver:11211
    http:
      url: https://myserver.com
```
* [ ] add support for custom metrics and expand thresholds accordingly eg: number of process
```yaml
host:
  servers:
    - name : server
      host: ubuntu@172.20.2.202:22
      key: /path/to/server.pem
      metrics:
        numprocs: "ps axu | wc -l"
        cpu: true
      threshold:
        numprocs: 400 # custom metric
        cpu: 0.80  # cpu usage  80%
```
