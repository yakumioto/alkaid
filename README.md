# Hyperledger Fabric Deploy

[中文教程](./README-zh.md)

## Supported features

- channel create
- channel updateAnchorPeer
- channel join
- channel update (Update BatchTimeout, BatchSize)
- chaincode install
- chaincode instantiate
- chaincode upgrade
- chaincode invoke
- chaincode query
- organization join (Add organization dynamically, supported system channel)
- organization delete (Delete organization dynamically, supported system channel)
- organization update (Dynamically modify organization certificate, supported system channel)
- channel consensus (Switch consensus algorithm, supported solo, kafka, etcdraft)

## Launch test network

```bash
git clone https://github.com/yakumioto/hlf-deploy.git && \
    cd hlf-deploy/test-network && \
    ./hlfn.sh up
```

The demo includes the following steps.

1. Create test network
2. Create channel
3. Update anchor peer
4. Join channel
5. Install go chaincode
6. Instantiate go chaincode
7. Query go chaincode
8. Invoke go chaincode
9. Add Org3 to mychannel
10. Modify Org3's certificate
11. Remove Org3 from mychannel
12. Upgrade from solo consensus to etcdraft consensus
13. Invoke and query chaincode to ensure that the consensus upgrade is successful
14. Install java chaincode
15. Upgrade java chaincode
16. Invoke java chaincode
17. Query java chaincode

Here is an example of the output

[![asciicast](https://asciinema.org/a/291386.svg)](https://asciinema.org/a/291386)

### Remove test network

```bash
./hlfn.sh down
```

## Deploying the network manually

### Start network

1. `cd test-network`

2. Get binary and docker image(skip if these two already existed).

    1. Download the [`hlf-deploy`](https://github.com/yakumioto/hlf-deploy/releases) binary
    
        `curl -L -o ../bin/hlf-deploy https://github.com/yakumioto/hlf-deploy/releases/download/v0.1.0/hlf-deploy`

    2. Pull docker image `yakumioto/hlf-tools` for dynamically add orgs.
    
        `docker pull yakumioto/hlf-tools:latest`
        
3. `docker-compose up -d`

### Create Channel

```bash
../bin/hlf-deploy channel create --configFile config.yaml \
    --channelTxFile channel-artifacts/channel.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1 Org2
```

### Update Anchor Peer

```bash
../bin/hlf-deploy channel updateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile channel-artifacts/Org1MSPanchors.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1
    
../bin/hlf-deploy channel updateAnchorPeer --configFile config.yaml \
    --anchorPeerTxFile channel-artifacts/Org2MSPanchors.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org2
```

### Join Channel

```bash
../bin/hlf-deploy channel join --configFile config.yaml \
    --channelName mychannel \
    Org1 Org2
```

### Install Chaincode

```bash
../bin/hlf-deploy chaincode install --configFile config.yaml \
    --lang golang \
    --goPath chaincode/go \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    Org1 Org2
```
### Instantiate Chaincode

`chaincodePolicy`: Set which organization signatures are required (currently only members are supported)
`chaincodePolicyNOutOf`: Set how many organization endorsement signatures are checked to return true

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

### Upgrade Chaincode

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

### Query Chaincode

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

### Invoke Chaincode

```bash
    ../bin/hlf-deploy chaincode invoke --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --endorsementOrgsName Org1,Org2 \
    --chaincodeName mycc \
    invoke a b 50
```

### Add Org3 organization dynamically

```bash
../bin/hlf-deploy organization join --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/org3.json \
    --orgName Org3MSP \
    Org1 Org2
```

### Update Org3 organization dynamically

```bash
../bin/hlf-deploy organization update --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/modify-org3.json \
    --orgName Org3MSP \
    Org3
```

### Delete Org3 organization dynamically

```bash
../bin/hlf-deploy organization delete --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgName Org3MSP \
    Org1 Org2
```

### Update solo consensus to edtcraft consensus

check out function `soloToRaftConsensus` in [`hlfn.sh`](test-network/hlfn.sh:145)

### Add new orderer

```bash
../bin/hlf-deploy organization join --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/newOrderer.json \
    --orgName OrdererOrg2 \
    --ordererOrg \
    OrdererOrg
```