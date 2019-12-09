# Hyperledger Fabric 部署

本项目不建议零基础直接上手, 可以先根据官方的 
[fabric-samples/first-network](https://github.com/hyperledger/fabric-samples) 入门

本项目基于 `Docker swarm` 进行多机部署

## 已支持功能

- createChannel
- updateAnchorPeer
- joinChannel
- installChaincode
- instantiateChaincode
- upgradeChaincode
- invokeChaincode
- queryChaincode
- addOrgChannel (动态添加组织, 支持 system channel)

## 在实现中

- deleteOrgChannel (动态删除组织)
- changeOrgCertificate (动态修改组织证书)

## 前期准备工作

1. docker swarm 集群 (可单节点)
2. fabric all images
3. nfs server client (用于共享证书等文件)
4. channel-artifacts
5. crypto-config

```text
channel-artifacts
├── channel.tx
├── genesis.block
├── Org1MSPanchors.tx
└── Org2MSPanchors.tx
```

```text
crypto-config
├── ordererOrganizations
│   └── example.com
└── peerOrganizations
    ├── org1.example.com
    └── org2.example.com
```

## 启动测试网络

**请保证上方所有准备工作已经完毕!**

**将 nfs 目录创建在 `/nfsvolume` !**

**以 root 权限运行**

```text
# nfs 配置编写
# 编辑 /etc/exports 
# 添加以下行, 并保存
# 并执行 exportfs -arv 重新加载
# 如实在不懂可 google

/nfsvolume *(ro,sync,no_root_squash)
```

进入 `test-network` 目录

### 启动网络

```bash
./hlfn.sh up -n manjaro

# output
Creating service orderer_orderer
Creating service peer0org1_peer
Creating service peer1org1_peer
Creating service peer0org2_peer
Creating service peer1org2_peer
2019/12/09 18:11:33 create channel txID: 009e848cea0b1731f2c5ff11d08435f9dcb045bfbb472c0cf66e64c7c25af535
2019/12/09 18:11:33 Org1 update anchor peer txID: 0831972c4a4c52a45fade69c82c14ec0c2a9fdc79657e4cb5ff6093d96ecd2c4
2019/12/09 18:11:33 Org2 update anchor peer txID: 1791b42f58e309dfb8517cf049884c1da58c865ff5bb16300610bc2a87f35ff4
2019/12/09 18:11:33 Org1 join channel successful
2019/12/09 18:11:36 Org2 join channel successful
2019/12/09 18:11:36 Org1 install chaincode successful
2019/12/09 18:11:36 Org2 install chaincode successful
2019/12/09 18:11:41 Org1 instantiate chaincode txID: b2dfb3b59ebf732f5c9a765863af40394da349077c554cd1c78f2833059c6bc6 args: [a 100 b 200]
2019/12/09 18:11:51 Org1 query chaincode txID: ce8f61708298abfc5337e362c8be48655cdb7b6e7a504f13ba07751051c4c369 args: [query a] result: 100
2019/12/09 18:11:51 Org1 query chaincode txID: a3dad9a3ffadbb85696c9a7a3b3590bc010433dffdacc9a75f16bf450107414e args: [query b] result: 200
2019/12/09 18:11:53 Org1 invoke chaincode txID: 74d8acd5218a169825b0daea535b99e8dd1b6bba83b96c27e2481747704dd698 args: [invoke a b 50]
2019/12/09 18:11:53 Org1 query chaincode txID: 0f9c2e19de066b915813f4bc676b181fb8d9ba967c4915a7fa8c90e30d25b041 args: [query a] result: 50
2019/12/09 18:11:53 Org1 query chaincode txID: ef9086e500e2f6da399cb5b1a197a175a6987ac8bcb2862e1f9805f56a21fd0c args: [query b] result: 250
```

### 动态添加组织

```bash
./hlfn.sh addOrgChannel

# output
2019/12/09 18:22:13 add Org3MSP to mychannel txID: 3ef3711dea377a5611f036c8b787993779e8610f4e055a3c84b101c3d992aa3d
```

### 停止网络

```bash
./hlfn.sh down

# output
Removing service orderer_orderer
Removing service peer0org1_peer
Removing service peer1org1_peer
Removing service peer0org2_peer
Removing service peer1org2_peer
WARNING! This will remove all local volumes not used by at least one container.
Are you sure you want to continue? [y/N] y
Deleted Volumes:
orderer.example.com.data
peer1.org1.example.com.data
peer0.org2.example.com.data
peer0.org1.example.com.data
peer1.org2.example.com.data

Total reclaimed space: 1.033MB
```

## 启动网络 (单节点为例子)

### 创建集群网络

```bash
docker network create --driver=overlay --attachable hlf
```

### 启动 Orderer

```bash
ORDERER_HOSTNAME=orderer \
    ORDERER_DOMAIN=example.com \
    ORDERER_GENERAL_LOCALMSPID=OrdererMSP \
    FABRIC_LOGGING_SPEC=debug \
    NODE_HOSTNAME=master \
    NETWORK=hlf \
    PORT=7050 \
    NFS_ADDR=127.0.0.1 \
    NFS_PATH=/nfsvolume \
    docker stack up -c orderer.yaml orderer
```

### 启动 peer

使用 LevelDB

```bash
PEER_HOSTNAME=peer0 \
    PEER_DOMAIN=org1.example.com \
    FABRIC_LOGGING_SPEC=debug \
    CORE_PEER_LOCALMSPID=Org1MSP \
    NODE_HOSTNAME=master \
    NETWORK=hlf \
    PORT=7051 \
    NFS_ADDR=127.0.0.1 \
    NFS_PATH=/nfsvolume \
    docker stack up -c peer-leveldb.yaml peer0org1
```

使用 CouchDB

```bash
PEER_HOSTNAME=peer0 \
    PEER_DOMAIN=org1.example.com \
    FABRIC_LOGGING_SPEC=debug \
    CORE_PEER_LOCALMSPID=Org1MSP \
    NODE_HOSTNAME=master \
    NETWORK=hlf \
    PORT=7051 \
    NFS_ADDR=127.0.0.1 \
    NFS_PATH=/nfsvolume \
    docker stack up -c peer-couchdb.yaml peer0org1
```

### 启动 CA Server

```bash
PEER_DOMAIN=org1.example.com \
    NODE_HOSTNAME=master \
    USERNAME=admin \
    PASSWORD=adminpwd \
    NETWORK=hlf \
    PORT=7054 \
    NFS_ADDR=127.0.0.1 \
    NFS_PATH=/nfsvolume \
    CA_PRIVEATE_KEY=$(cd ${NFS_PATH}/crypto-config/peerOrganizations/${PEER_DOMAIN}/ca && ls *_sk) \
    docker stack up -c ca.yaml peer0org1ca
```

## 部署

### 创建 Channel

```bash
hlf-deploy createChannel --configFile config.yaml \
    --channelTxFile channel.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1 Org2
```

### 更新 Anchor Peer

```bash
hlf-deploy uptateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile anchor.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1
```

### 加入 Channel

```bash
hlf-deploy joinChannel --configFile config.yaml \
    --channelName mychannel \
    Org1 Org2
```

### 安装 Chaincode

```bash
hlf-deploy installChaincode --configFile config.yaml \
    --goPath ./chaincode \
    --chaincodePath example_02 \
    --chaincodeName example \
    --chaincodeVersion v1.0 \
    Org1 Org2
```

### 更新 Chaincode

`chaincodePolicy`: 设置需要哪些组织签名 (目前只支持 Member)
`chaincodePolicyNOutOf`: 用来设置多少个组织签名检验成功后返回 true

```bash
hlf-deploy upgradeChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    a 200 b 100
```

### 实例化 Chaincode

```bash
hlf-deploy instantiateChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 1 \
    --chaincodePath example02 \
    --chaincodeName example \
    --chaincodeVersion v0.0.0 \
    a 100 b 200
```

### 查询 Chaincode

```bash
hlf-deploy queryChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodeName example \
    query a
```

### 调用 Chaincode

```bash
hlf-deploy invokeChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --endorsementOrgsName Org1,Org2 \
    --chaincodeName example \
    invoke a b 50
```
