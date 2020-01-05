# Hyperledger Fabric 部署

本项目不建议零基础直接上手, 可以先根据官方的 
[fabric-samples/first-network](https://github.com/hyperledger/fabric-samples) 入门

## 已支持功能

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

## 启动测试网络

```bash
git clone https://github.com/yakumioto/hlf-deploy.git && \
    cd hlf-deploy/test-network && \
    ./hlfn.sh up
```

以下是这个示例网络的所有执行步骤.

1. Start test network
2. Create channel
3. Update anchor peer
4. Join channel
5. Install go chaincode (mycc v1.0)
6. Instantiate go chaincode (mycc v1.0)
7. Query go chaincode
8. Invoke go chaincode
9. Add Org3 to mychannel (动态添加组织)
10. Modify Org3's certificate (动态修改组织证书)
11. Remove Org3 from mychannel (动态删除组织)
12. Upgrade from solo consensus to etcdraft consensus (从 solo 升级到 etcdraft)
13. Invoke and query chaincode to ensure that the consensus upgrade is successful (测试升级后的网络)
14. Install java chaincode (mycc v2.0)
15. Upgrade java chaincode (mycc v2.0)
16. Invoke java chaincode
17. Query java chaincode

Here is an example of the output

[![asciicast](https://asciinema.org/a/291386.svg)](https://asciinema.org/a/291386)

### Remove test network

```bash
./hlfn.sh down
```

## 手动部署网络(单节点栗子)

### 启动网络

1. 进入目录

    `cd test-network`

2. 下载二进制程序和Docker镜像（如已存在，可跳过此步骤）
    1. 下载二进制程序 [`hlf-deploy`](https://github.com/yakumioto/hlf-deploy/releases)
    
        `curl -L -o ../bin/hlf-deploy https://github.com/yakumioto/hlf-deploy/releases/download/v0.1.0/hlf-deploy`
        
    2. 拉取动态添加组织的 Docker 镜像 `yakumioto/hlf-tools`
    
        `docker pull yakumioto/hlf-tools:latest`

3. 启动

    `docker-compose up -d`

### 创建 Channel

```bash
../bin/hlf-deploy channel create --configFile config.yaml \
    --channelTxFile channel-artifacts/channel.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1 Org2
```

### 更新 Anchor Peer

```bash
../bin/hlf-deploy channel updateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile channel-artifacts/Org1MSPanchors.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1
    
../bin/hlf-deploy updateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile channel-artifacts/Org2MSPanchors.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org2
```

### 加入 Channel

```bash
../bin/hlf-deploy channel join --configFile config.yaml \
    --channelName mychannel \
    Org1 Org2
```

### 安装 Chaincode

```bash
../bin/hlf-deploy chaincode install --configFile config.yaml \
    --lang golang \
    --goPath chaincode/go \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    Org1 Org2
```
### 实例化 Chaincode

`chaincodePolicy`: 设置需要哪些组织签名 (目前只支持 Member)
`chaincodePolicyNOutOf`: 用来设置多少个组织签名检验成功后返回 true

```bash
../bin/hlf-deploy chaincode instantiate --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --lang golang \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    a 100 b 200
```

### 更新 Chaincode

```bash
hlf-deploy chaincode upgrade --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --lang golang \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v2.0 \
    a 100 b 200
```

### 查询 Chaincode

```bash
../bin/hlf-deploy chaincode query --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodeName mycc \
    query a

../bin/hlf-deploy chaincode query --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodeName mycc \
    query b
```

### 调用 Chaincode

```bash
    ../bin/hlf-deploy chaincode invoke --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --endorsementOrgsName Org1,Org2 \
    --chaincodeName mycc \
    invoke a b 50
```

### 动态添加 Org3 组织

```bash
../bin/hlf-deploy organization join --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/org3.json \
    --orgName Org3MSP \
    Org1 Org2
```

### 动态更新 Org3 组织

```bash
../bin/hlf-deploy organization update --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/modify-org3.json \
    --orgName Org3MSP \
    Org3
```

### 动态删除 Org3 组织

```bash
../bin/hlf-deploy organization delete --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgName Org3MSP \
    Org1 Org2
```

### 从 solo 升级到 etcdraft

参照 [`hlfn.sh`](test-network/hlfn.sh:145) 里的 `soloToRaftConsensus` 方法 

### 动态添加 Orderer

```bash
../bin/hlf-deploy organization join --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/newOrderer.json \
    --orgName OrdererOrg2 \
    --ordererOrg \
    OrdererOrg
```