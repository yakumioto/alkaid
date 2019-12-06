# Hyperledger Fabric Deploy

Fabric network deployment of this project is based on docker swarm multi-machine deployment.

[中文教程](./README-zh.md)

## Start Network

### Start Orderer

```bash
ORDERER_HOSTNAME=orderer \
    ORDERER_DOMAIN=example.com \
    ORDERER_GENERAL_LOCALMSPID=OrdererMSP \
    FABRIC_LOGGING_SPEC=debug \
    NODE_HOSTNAME=master \
    NETWORK=hlf \
    PORT=7050 \
    NFS_ADDR=192.168.51.10 \
    NFS_PATH=/nfsvolume \
    docker stack up -c orderer.yaml orderer
```

### Start peer

Use LevelDB

```bash
PEER_HOSTNAME=peer0 \
    PEER_DOMAIN=org1.example.com \
    FABRIC_LOGGING_SPEC=debug \
    CORE_PEER_LOCALMSPID=Org1MSP \
    NODE_HOSTNAME=node1 \
    NETWORK=hlf \
    PORT=7051 \
    NFS_ADDR=192.168.51.10 \
    NFS_PATH=/nfsvolume \
    docker stack up -c peer-leveldb.yaml peer0org1
```

Use CouchDB

```bash
PEER_HOSTNAME=peer0 \
    PEER_DOMAIN=org1.example.com \
    FABRIC_LOGGING_SPEC=debug \
    CORE_PEER_LOCALMSPID=Org1MSP \
    NODE_HOSTNAME=node1 \
    NETWORK=hlf \
    PORT=7051 \
    NFS_ADDR=192.168.51.10 \
    NFS_PATH=/nfsvolume \
    docker stack up -c peer-couchdb.yaml peer0org1
```

### Start CA Server

```bash
PEER_DOMAIN=org1.example.com \
    NODE_HOSTNAME=node1 \
    USERNAME=admin \
    PASSWORD=adminpwd \
    NETWORK=hlf \
    PORT=7054 \
    NFS_ADDR=192.168.51.10 \
    NFS_PATH=/nfsvolume \
    CA_PRIVEATE_KEY=$(cd ${NFS_PATH}/crypto-config/peerOrganizations/${PEER_DOMAIN}/ca && ls *_sk) \
    docker stack up -c ca.yaml peer0org1ca
```

## Deploy Network

### Create Channel

```bash
hlf-deploy createChannel --configFile config.yaml \
    --channelTxFile channel.tx \
    --channelName testchannel \
    --ordererOrgName OrdererOrg \
    Org1 Org2
```

### Update Anchor Peer

```bash
hlf-deploy uptateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile anchor.tx \
    --channelName testchannel \
    --ordererOrgName OrdererOrg \
    Org1
```

### Join Channel

```bash
hlf-deploy joinChannel --configFile config.yaml \
    --channelName testchannel \
    Org1 Org2
```

### Install Chaincode

```bash
hlf-deploy installChaincode --configFile config.yaml \
    --goPath ./chaincode \
    --chaincodePath example02 \
    --chaincodeName example \
    --chaincodeVersion v0.1.0 \
    Org1 Org2
```

### Instantiate Chaincode

```bash
hlf-deploy instantiateChaincode --configFile config.yaml \
    --channelName testchannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 1 \
    --chaincodePath example02 \
    --chaincodeName example \
    --chaincodeVersion v0.0.0 \
    a 100 b 200
```

### Upgrade Chaincode

```bash
hlf-deploy upgradeChaincode --configFile config.yaml \
    --channelName testchannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --chaincodePath example02 \
    --chaincodeName example \
    --chaincodeVersion v0.1.0 \
    a 200 b 100
```

### Query Chaincode

```bash
hlf-deploy queryChaincode --configFile config.yaml \
    --channelName testchannel \
    --orgName Org1 \
    --chaincodeName example \
    query a
```

### Invoke Chaincode

```bash
hlf-deploy invokeChaincode --configFile config.yaml \
    --channelName testchannel \
    --orgName Org1 \
    --endorsementOrgsName Org1,Org2 \
    --chaincodeName example \
    invoke a b 50
```
