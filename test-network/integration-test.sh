#!/usr/bin/env bash

function checkExitCode() {
    code=${?}
    if [[ ${code} != "0" ]]; then
        exit ${code}
    fi
}

function createChannel() {
    go run ../*.go channel create --configFile config.yaml \
        --channelTxFile channel-artifacts/channel.tx \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        ${1}

    checkExitCode
    sleep 2s
}

function updateAnchorPeer() {
    go run ../*.go channel updateAnchorPeer --configFile config.yaml \
        --anchorPeerTxFile ${1} \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        ${2}

    checkExitCode
    sleep 2s
}

function joinChannel() {
    go run ../*.go channel join --configFile config.yaml \
        --channelName mychannel \
        ${1}

    checkExitCode
    sleep 2s
}

function installChaincode() {
    go run ../*.go chaincode install --configFile config.yaml \
        --lang ${1} \
        --goPath ${2} \
        --chaincodePath ${3} \
        --chaincodeName mycc \
        --chaincodeVersion ${4} \
        ${5}

    checkExitCode
    sleep 2s
}

function instantiateChaincode() {
    go run ../*.go chaincode instantiate --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodePolicy Org1MSP,Org2MSP \
        --chaincodePolicyNOutOf 2 \
        --lang ${1} \
        --chaincodePath ${2} \
        --chaincodeName mycc \
        --chaincodeVersion ${3} \
        a 100 b 200

    checkExitCode
    sleep 10s
}

function upgradeChaincode() {
    go run ../*.go chaincode upgrade --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodePolicy Org1MSP,Org2MSP \
        --chaincodePolicyNOutOf 2 \
        --lang ${1} \
        --chaincodePath ${2} \
        --chaincodeName mycc \
        --chaincodeVersion ${3} \
        a 100 b 200

    checkExitCode
    sleep 10s
}

function queryChaincode() {
    result=$(go run ../*.go chaincode query --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodeName mycc \
        query ${1} 2>&1)

    checkExitCode

    echo ${result}

    actual=${result:136}

    if [[ ${actual} != ${2} ]]; then
        echo Actual result does not match expected result, expect: ${2}, actual: ${actual}.
        exit -1
    fi

    sleep 2s
}

function invokeChaincode() {
    go run ../*.go chaincode invoke --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --endorsementOrgsName Org1,Org2 \
        --chaincodeName mycc \
        invoke ${1} ${2} ${3}

    checkExitCode
    sleep 2s
}

function addOrganization() {
    go run ../*.go organization join --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/org3.json \
        --orgName Org3MSP \
        Org1 Org2

    checkExitCode
    sleep 2s
}

function updateOrganization() {
    go run ../*.go organization update --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/modify-org3.json \
        --orgName Org3MSP \
        Org3

    checkExitCode
    sleep 2s
}

function deleteOrganization() {
    go run ../*.go organization delete --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgName Org3MSP \
        Org1 Org2

    checkExitCode
    sleep 2s
}

function soloToRaftConsensus() {
    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --state Maintenance

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --state Maintenance


    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --type etcdraft

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --type etcdraft

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer2.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer2.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer3.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer3.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer4.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer4.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer5.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer5.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer2.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer2.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer3.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer3.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer4.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer4.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer5.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer5.example.com:7050

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --state Normal

    checkExitCode

    go run ../*.go channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --state Normal

    checkExitCode

    go run ../*.go channel config \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        ./channel-artifacts/newGenesis.block

    checkExitCode

    docker-compose restart
    docker-compose up -d
}

function upNetwork() {
    docker-compose up -d \
        orderer.example.com \
        peer0.org1.example.com \
        peer1.org1.example.com \
        peer0.org2.example.com \
        peer1.org2.example.com
}

upNetwork

createChannel Org1

updateAnchorPeer channel-artifacts/Org1MSPanchors.tx Org1
updateAnchorPeer channel-artifacts/Org2MSPanchors.tx Org2

joinChannel Org1
joinChannel Org2

installChaincode golang chaincode/go example_02 v1.0 Org1
installChaincode golang chaincode/go example_02 v1.0 Org2

instantiateChaincode golang example_02 v1.0

queryChaincode a 100
queryChaincode b 200
invokeChaincode a b 50
queryChaincode a 50
queryChaincode b 250

addOrganization
updateOrganization
deleteOrganization

soloToRaftConsensus

invokeChaincode b a 50
queryChaincode a 100
queryChaincode b 200

installChaincode java chaincode/go chaincode/java v2.0 Org1
installChaincode java chaincode/go chaincode/java v2.0 Org2
upgradeChaincode java chaincode/java v2.0

invokeChaincode b a 50
queryChaincode a 150
queryChaincode b 150