# WSYNC TOOL 

## Introduction

Multi role project synchronization tool:
- sender： the machine where the project source is located
- accepter：the machine where the project destination is located
- puller：the machine where the project destination is located

## How to use?

```bash
# source machine ip: 192.168.75.10, destination machine ip: 192.168.75.11
cd wsync
./build amd64
scp wsync-amd64 config.yaml root@192.168.75.10:/root
scp wsync-amd64 config.yaml root@192.168.75.11:/root

# in source machine
cd /root
ssh-copy-id root@192.168.75.11
vim config.yaml
./wsync

# in destination machine
cd /root
ssh-copy-id root@192.168.75.10
vim config.yaml
./wsync
```

In the configuration file, the role represents the character of the machine where the wsync tool is located, which can be sender, accepter, or puller.

If the machine role is sender, it is configured as the source machine hosting the project to be synchronized.

If the sender endpoint is configured with accepter related information, the sender will perform the project synchronization operation, and the destination machine (where the project is synchronized to) does not need to execute any actions.

If the sender endpoint is configured with puller related information, the destination machine must deploy the wsync tool with the puller role. The puller will execute the synchronization based on the pullmethod field:
- period: Synchronization is performed periodically.
- web: Synchronization is triggered immediately on the puller side when changes occur to the project on the sender.

## config.yaml example
```yaml
# source machine
#========================================================
#
# sender/accepter/puller
role: "sender"
#
#========================================================
#
sender:
  dir: "/root/your_project_name/"
  user: "root"
  ip: "192.168.75.10"
#
#========================================================
#
accepter:
  dir: ""
  user: ""
  ip: ""
#
#========================================================
#
puller:
  addr: "192.168.75.11:19876"
  https: false
# peried or web
  pullmethod: "web"
# unit: second
  pullperiod: 5
  dir: "/root/your_project_name/"
  user: "root"
  ip: "192.168.75.11"
```
```yaml
# destination machine
#========================================================
#
# sender/accepter/puller
role: "puller"
#
#========================================================
#
sender:
  dir: "/root/your_project_name/"
  user: "root"
  ip: "192.168.75.10"
#
#========================================================
#
accepter:
  dir: ""
  user: ""
  ip: ""
#
#========================================================
#
puller:
  addr: "192.168.75.11:19876"
  https: false
# peried or web
  pullmethod: "web"
# unit: second
  pullperiod: 5
  dir: "/root/your_project_name/"
  user: "root"
  ip: "192.168.75.11"
```

## Contributors

Welcome to participate in the joint development of this project