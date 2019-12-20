#!/usr/bin/env bash


function printHelp() {
  echo "Usage: "
  echo "  hlfn.sh <mode>"
  echo "    <mode> - one of 'up', 'down'"
  echo "      - 'up' - bring up the network with docker-compose up"
  echo "      - 'down' - clear the network with docker-compose down"
  echo
  echo "Taking all defaults:"
  echo "	hlfn.sh up"
  echo "	hlfn.sh down"
}

if [[ ! -f "../bin/hlf-deploy" ]]; then
    mkdir -p ../bin
    curl -L -o ../bin/hlf-deploy https://github.com/yakumioto/hlf-deploy/releases/download/v0.1.0/hlf-deploy
    chmod +x ../bin/hlf-deploy
fi

if [[ -z "$(docker images -q yakumioto/hlf-tools:latest)" ]]; then
    echo "docker pull yakumioto/hlf-tools:latest"
    docker pull yakumioto/hlf-tools:latest
fi

function upNetwork() {
    docker-compose up -d hlf-tools \
        orderer.example.com \
        peer0.org1.example.com \
        peer1.org1.example.com \
        peer0.org2.example.com \
        peer1.org2.example.com
}

function createChannel() {
    ../bin/hlf-deploy channel create --configFile config.yaml \
        --channelTxFile channel-artifacts/channel.tx \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        Org1 Org2
}

function updateAnchorPeer() {
    ../bin/hlf-deploy channel updateAnchorPeer --configFile config.yaml \
        --anchorPeerTxFile ${1} \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        ${2}
}

function joinChannel() {
    ../bin/hlf-deploy channel join --configFile config.yaml \
        --channelName mychannel \
        Org1 Org2
}

function installChaincode() {
    ../bin/hlf-deploy chaincode install --configFile config.yaml \
        --goPath chaincode \
        --chaincodePath example_02 \
        --chaincodeName mycc \
        --chaincodeVersion v1.0 \
        Org1 Org2
}

function instantiateChaincode() {
    ../bin/hlf-deploy chaincode instantiate --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodePolicy Org1MSP,Org2MSP \
        --chaincodePolicyNOutOf 2 \
        --chaincodePath example_02 \
        --chaincodeName mycc \
        --chaincodeVersion v1.0 \
        a 100 b 200
}

function queryChaincode() {
    ../bin/hlf-deploy chaincode query --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodeName mycc \
        query ${1}
}

function invokeChaincode() {
    ../bin/hlf-deploy chaincode invoke --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --endorsementOrgsName Org1,Org2 \
        --chaincodeName mycc \
        invoke ${1} ${2} ${3}
}

function addOrganization() {
    ../bin/hlf-deploy organization join --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/org3.json \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org1 Org2
}

function updateOrganization() {
    ../bin/hlf-deploy organization update --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/modify-org3.json \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org3
}

function deleteOrganization() {
    ../bin/hlf-deploy organization delete --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org1 Org2
}

function addOrdererOrganization() {
    ../bin/hlf-deploy organization join --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/newOrderer.json \
        --orgName OrdererOrg2 \
        --rpcAddress localhost:1234 \
        --ordererOrg \
        OrdererOrg
}

function soloToRaftConsensus() {
    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --state Maintenance

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --state Maintenance

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --type etcdraft

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --type etcdraft

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer2.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer2.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer3.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer3.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer4.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer4.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer4.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --host orderer5.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --host orderer5.example.com \
        --port 7050 \
        --clientTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt \
        --serverTLSCertPath ./crypto-config/ordererOrganizations/example.com/orderers/orderer5.example.com/tls/server.crt

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer2.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer2.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer3.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer3.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer4.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer4.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer5.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --ordererAddress orderer5.example.com:7050

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --state Normal

    ../bin/hlf-deploy channel consensus \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        --state Normal

    go run ../*.go channel config \
        --configFile ./config.yaml \
        --channelName byfn-sys-channel \
        --sysChannel \
        --ordererOrgName OrdererOrg \
        ./channel-artifacts/newGenesis.block

    docker-compose restart
    docker-compose up -d
}

function cleanNetwork() {
    docker-compose down

    sleep 5s

    docker volume prune -f
}

if [[ ! -f "../bin/hlf-deploy" ]]; then
    echo "hlf-deploy tool not found. exiting"
    exit 1
fi

if [[ "${1}" = "-m" ]]; then
  shift
fi
mode=${1}
shift

if [[ "${mode}" == "up" ]]; then
    :
elif [[ "${mode}" == "down" ]]; then
    :
else
    printHelp
    exit 1
fi

while getopts "h?:" opt; do
  case ${opt} in
  h | \?)
    printHelp
    exit 0
    ;;
  esac
done

if [[ "${mode}" == "up" ]]; then
    upNetwork
    createChannel
    updateAnchorPeer channel-artifacts/Org1MSPanchors.tx Org1
    updateAnchorPeer channel-artifacts/Org2MSPanchors.tx Org2
    joinChannel
    installChaincode
    instantiateChaincode
    sleep 10s
    queryChaincode a
    queryChaincode b
    invokeChaincode a b 50
    queryChaincode a
    queryChaincode b
    addOrganization
    updateOrganization
    deleteOrganization
    soloToRaftConsensus
    sleep 20s
    invokeChaincode b a 50
    queryChaincode a
    queryChaincode b
elif [[ "${mode}" == "down" ]]; then ## Clear the network
    cleanNetwork
fi