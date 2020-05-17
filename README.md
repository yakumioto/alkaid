# Alkaid

[![Github Action](https://github.com/yakumioto/alkaid/workflows/alkaid/badge.svg)](https://github.com/yakumioto/alkaid/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yakumioto/alkaid)](https://goreportcard.com/report/github.com/yakumioto/alkaid)
[![codecov](https://codecov.io/gh/yakumioto/alkaid/branch/master/graph/badge.svg)](https://codecov.io/gh/yakumioto/alkaid)

Alikid 是一个基于 Hyperledger Fabric 实现的 BaaS(Blockchan as a Service) 服务

[English](README-en.md)

项目处于开发阶段, 并未正式发布以下功能均为计划支持的功能.

## 支持

- [ ] 集群支持: `docker`, `docker swarm`, `kubernetes`
- [ ] 组织管理: 组织的管理, 创建不同类型的组织, 如: `orderer`, `peer`
- [ ] 证书管理: 自动创建根证书, 并签发用户的 `MSP` 证书
- [ ] 用户管理: 系统将 `MSP` 视为用户, 用户类型: `Admin`, `Client`, `Orderer`, `Peer`
- [ ] 网络管理: 组织可以创建多个网络, 加入多个网络
- [ ] 动态添加: 网络支持动态添加组织功能
- [ ] 共识算法: 支持 `solo`, `etcdraft`
- [ ] 合约管理: 合约的上传, 初始化, 部署, 支持 `go`, `java`, `node`
- [ ] 分布式交互: 支持跨 `BaaS` 的同一网络通信

## 架构

```text
+-------------------------------------------------------------+
|                                                             |
|                      Alkaid Frontend                        |
|                                                             |
+-----+-------------------------------------------------+-----+
      |                                                 |
      |                                                 |
      v                                                 v
+-----+-------------------------------------------------+-----+
|                                                             |
|                      Alkaid Backend                         |
|                                                             |
+-----+-------------------------------------------------+-----+
      |                                                 |
      |                                                 |
      v                                                 v
+-----+-------------------------------------------------+-----+
|                                                             |
|                    Docker / K8S / K3S                       |
|                                                             |
+---------+-------------------+--------------------+----------+
          |                   |                    |
          |                   |                    |
          v                   v                    v
    +-----+-----+       +-----+-----+        +-----+-----+
    |           |       |           |        |           |
    |  Net 001  |       |  Net 002  |        |  Net 003  |
    |           |       |           |        |           |
    +-----+-----+       +-----+-----+        +-----+-----+
          |                   |                    |
          |                   |                    |
          v                   v                    v
 +--------+-------------------+--------------------+----------+
 |                                                            |
 |   Virtual or Physical Machine / Public or Private Cloud    |
 |                                                            |
 +------------------------------------------------------------+

```

## 社群

Telegram: <https://t.me/fab_alkaid>

## hlf-deploy

Alikid 的前身是 [hlf-deploy](https://github.com/yakumioto/alkaid/tree/v0.2.0) 用于快速实现对 Hyperledger Fabric 网络的部署与调整

支持的功能:

- channel create
- channel updateAnchorPeer
- channel join
- channel update (支持 BatchTimeout, BatchSize)
- chaincode install
- chaincode instantiate
- chaincode upgrade
- chaincode invoke
- chaincode query
- organization join (动态添加组织, 支持 system channel)
- organization delete (动态删除组织, 支持 system channel)
- organization update (动态更新组织, 支持 system channel)
- channel consensus (切换 orderer 共识, 支持 solo, kafka, etcdraft)
