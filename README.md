# Hyperledger Fabric Deploy

[中文教程](./README-zh.md)

## Supported features

- createChannel
- updateAnchorPeer
- joinChannel
- installChaincode
- instantiateChaincode
- upgradeChaincode
- invokeChaincode
- queryChaincode
- addOrgChannel (Add organization dynamically, supported system channel)
- deleteOrgChannel (Delete organization dynamically, supported system channel)
- changeOrgCertificate (Dynamically modify organization certificate, supported system channel)

## Launch test network

```bash
git clone https://github.com/yakumioto/hlf-deploy.git && \
    cd hlf-deploy/test-network && \
    ./hlfn.sh up
```

Here is an example of the output

```bash
./hlfn.sh up

# output
Creating network "test-network_byfn" with the default driver
Creating volume "test-network_orderer.example.com" with default driver
Creating volume "test-network_peer0.org1.example.com" with default driver
Creating volume "test-network_peer0.org2.example.com" with default driver
Creating volume "test-network_peer1.org1.example.com" with default driver
Creating volume "test-network_peer1.org2.example.com" with default driver
Creating peer0.org2.example.com ... done
Creating peer1.org2.example.com ... done
Creating peer0.org1.example.com ... done
Creating peer1.org1.example.com ... done
Creating orderer.example.com    ... done
2019/12/10 17:55:18 create channel txID: ec92b11651b4457cb10cceb4f67670ddc90f938c0f329fbce6116be6b1a50602
2019/12/10 17:55:18 Org1 update anchor peer txID: abd5de7ea4cc72728a60fa9437da766c8df990ebfb47de49d95353a210dcc81f
2019/12/10 17:55:18 Org2 update anchor peer txID: 2e4680f3c0409387cd499f4d521555f638612714de37355b9cfa6d87b17741fc
2019/12/10 17:55:18 Org1 join channel successful
2019/12/10 17:55:18 Org2 join channel successful
2019/12/10 17:55:18 Org1 install chaincode successful
2019/12/10 17:55:18 Org2 install chaincode successful
2019/12/10 17:55:24 Org1 instantiate chaincode txID: 969805f856bf98ac6aa9fc1afa2b52ff2343f5711c51677031d73e090c634beb args: [a 100 b 200]
2019/12/10 17:55:34 Org1 query chaincode txID: 356450b292f3636d6cbe913c1e2ffe7286c17e30819160170d49a093916cb669 args: [query a] result: 100
2019/12/10 17:55:34 Org1 query chaincode txID: fef70c04942c28308d30df3e62dccac1f7eb2244491c34ecb7b6b2a3af8134f4 args: [query b] result: 200
2019/12/10 17:55:36 Org1 invoke chaincode txID: 3fce7db1f779d9d8187f26cd66431aa88807ed32ca7996dc1c94c37c140e4296 args: [invoke a b 50]
2019/12/10 17:55:36 Org1 query chaincode txID: 09323d03d505a09918c9585d25c5dce2456702bda167b1f127ea80d5930979bb args: [query a] result: 50
2019/12/10 17:55:37 Org1 query chaincode txID: 01e939dc3aef2fc32100dba1cfc6c446448fcb103fbe84f20f2cf645b405235c args: [query b] result: 250
2019/12/10 17:55:37 save Org3MSP to mychannel txID: b9146c47a22470d31d9456abdfb7cf106384f019a5c1f2b0ad8a6fe172d736c9
2019/12/10 17:55:37 save Org3MSP to mychannel txID: 3d31bbdb4ad19d42be364ac1ae2326ee0a7ebbee799ae9751b840813e587b977
2019/12/10 17:55:37 delete Org3MSP to mychannel txID: 4ea217bf2d2ded1304e9df2b9dc3ef6764d4d4e72aab033e99880184834889e1
```

```bash
./hlfn.sh down

# output
Stopping orderer.example.com    ... done
Stopping peer1.org2.example.com ... done
Stopping peer1.org1.example.com ... done
Stopping peer0.org1.example.com ... done
Stopping peer0.org2.example.com ... done
Removing orderer.example.com    ... done
Removing peer1.org2.example.com ... done
Removing peer1.org1.example.com ... done
Removing peer0.org1.example.com ... done
Removing peer0.org2.example.com ... done
Removing network test-network_byfn
Deleted Volumes:
test-network_orderer.example.com
test-network_peer0.org1.example.com
test-network_peer1.org2.example.com
test-network_peer0.org2.example.com
test-network_peer1.org1.example.com

Total reclaimed space: 1.475MB
```

## Deploying the network manually

First download the `hlf-deploy` binary and `yakumioto/hlf-tools:latest` mirror

`curl -L -o ../bin/hlf-deploy https://github.com/yakumioto/hlf-deploy/releases/download/v0.1.0/hlf-deploy`

`docker pull yakumioto/hlf-tools:latest`

### Start network

Go to the test-network directory

`docker-compose up -d`

### Create Channel

```bash
../bin/hlf-deploy createChannel --configFile config.yaml \
    --channelTxFile channel-artifacts/channel.tx \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    Org1 Org2
```

### Update Anchor Peer

```bash
../bin/hlf-deploy updateAnchorPeer --configFile config.yaml \
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

### Join Channel

```bash
../bin/hlf-deploy joinChannel --configFile config.yaml \
    --channelName mychannel \
    Org1 Org2
```

### Install Chaincode

```bash
../bin/hlf-deploy installChaincode --configFile config.yaml \
    --goPath chaincode \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    Org1 Org2
```
### Instantiate Chaincode

`chaincodePolicy`: Set which organization signatures are required (currently only members are supported)
`chaincodePolicyNOutOf`: Set how many organization endorsement signatures are checked to return true

```bash
../bin/hlf-deploy instantiateChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v1.0 \
    a 100 b 200
```

### Update Chaincode

```bash
hlf-deploy upgradeChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodePolicy Org1MSP,Org2MSP \
    --chaincodePolicyNOutOf 2 \
    --chaincodePath example_02 \
    --chaincodeName mycc \
    --chaincodeVersion v2.0 \
    a 100 b 200
```

### Query Chaincode

```bash
../bin/hlf-deploy queryChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodeName mycc \
    query a

../bin/hlf-deploy queryChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --chaincodeName mycc \
    query b
```

### Invoke Chaincode

```bash
    ../bin/hlf-deploy invokeChaincode --configFile config.yaml \
    --channelName mychannel \
    --orgName Org1 \
    --endorsementOrgsName Org1,Org2 \
    --chaincodeName mycc \
    invoke a b 50
```

### Add Org3 organization dynamically

```bash
../bin/hlf-deploy addOrgChannel --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/org3.json \
    --orgName Org3MSP \
    --rpcAddress localhost:1234 \
    Org1 Org2
```

### Update Org3 organization dynamically

```bash
../bin/hlf-deploy updateOrgChannel --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/modify-org3.json \
    --orgName Org3MSP \
    --rpcAddress localhost:1234 \
    Org3
```

### Delete Org3 organization dynamically

```bash
../bin/hlf-deploy delOrgChannel --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgName Org3MSP \
    --rpcAddress localhost:1234 \
    Org1 Org2
```

### Add OrdererOrg2 organization dynamically

```bash
../bin/hlf-deploy addOrgChannel --configFile config.yaml \
    --channelName mychannel \
    --ordererOrgName OrdererOrg \
    --orgConfig channel-artifacts/newOrderer.json \
    --orgName OrdererOrg2 \
    --rpcAddress localhost:1234 \
    --ordererOrg \
    OrdererOrg
```
