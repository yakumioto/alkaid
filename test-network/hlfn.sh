#!/usr/bin/env bash


function printHelp() {
  echo "Usage: "
  echo "  hlfn.sh <mode> [-n <nodename name>]"
  echo "    <mode> - one of 'up', 'down'"
  echo "      - 'up' - bring up the network with docker-compose up"
  echo "      - 'down' - clear the network with docker-compose down"
  echo
  echo "Taking all defaults:"
  echo "	hlfn.sh up"
  echo "	hlfn.sh down"
}

if [[ ! -f "../bin/hlf-deploy" ]]; then
    curl -L -O https://github.com/yakumioto/hlf-deploy/releases/download/v0.1.0/hlf-deploy
    chmod +x ../bin/hlf-deploy
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

while getopts "h?:n:" opt; do
  case ${opt} in
  h | \?)
    printHelp
    exit 0
    ;;
  n)
    nodename=${OPTARG}
    ;;
  esac
done

function upNetwork() {
    ORDERER_HOSTNAME=orderer \
        ORDERER_DOMAIN=example.com \
        ORDERER_GENERAL_LOCALMSPID=OrdererMSP \
        FABRIC_LOGGING_SPEC=info \
        NODE_HOSTNAME=${nodename} \
        NETWORK=hlf \
        PORT=7050 \
        NFS_ADDR=127.0.0.1 \
        NFS_PATH=/nfsvolume \
        docker stack up -c ../orderer.yaml orderer

    PEER_HOSTNAME=peer0 \
        PEER_DOMAIN=org1.example.com \
        FABRIC_LOGGING_SPEC=info \
        CORE_PEER_LOCALMSPID=Org1MSP \
        NODE_HOSTNAME=${nodename} \
        NETWORK=hlf \
        PORT=7051 \
        NFS_ADDR=127.0.0.1 \
        NFS_PATH=/nfsvolume \
        docker stack up -c ../peer-leveldb.yaml peer0org1

    PEER_HOSTNAME=peer1 \
        PEER_DOMAIN=org1.example.com \
        FABRIC_LOGGING_SPEC=info \
        CORE_PEER_LOCALMSPID=Org1MSP \
        NODE_HOSTNAME=${nodename} \
        NETWORK=hlf \
        PORT=8051 \
        NFS_ADDR=127.0.0.1 \
        NFS_PATH=/nfsvolume \
        docker stack up -c ../peer-leveldb.yaml peer1org1

    PEER_HOSTNAME=peer0 \
        PEER_DOMAIN=org2.example.com \
        FABRIC_LOGGING_SPEC=info \
        CORE_PEER_LOCALMSPID=Org2MSP \
        NODE_HOSTNAME=${nodename} \
        NETWORK=hlf \
        PORT=9051 \
        NFS_ADDR=127.0.0.1 \
        NFS_PATH=/nfsvolume \
        docker stack up -c ../peer-leveldb.yaml peer0org2

    PEER_HOSTNAME=peer1 \
        PEER_DOMAIN=org2.example.com \
        FABRIC_LOGGING_SPEC=info \
        CORE_PEER_LOCALMSPID=Org2MSP \
        NODE_HOSTNAME=${nodename} \
        NETWORK=hlf \
        PORT=10051 \
        NFS_ADDR=127.0.0.1 \
        NFS_PATH=/nfsvolume \
        docker stack up -c ../peer-leveldb.yaml peer1org2
}

function createChannel() {
    ../bin/hlf-deploy createChannel --configFile config.yaml \
        --channelTxFile channel-artifacts/channel.tx \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        Org1 Org2
}

function updateAnchorPeer() {
    ../bin/hlf-deploy updateAnchorPeer --configFile config.yaml \
        --anchorPeerTxFile ${1} \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        ${2}
}

function joinChannel() {
    ../bin/hlf-deploy joinChannel --configFile config.yaml \
        --channelName mychannel \
        Org1 Org2
}

function installChaincode() {
    ../bin/hlf-deploy installChaincode --configFile config.yaml \
        --goPath chaincode \
        --chaincodePath example_02 \
        --chaincodeName mycc \
        --chaincodeVersion v1.0 \
        Org1 Org2
}

function instantiateChaincode() {
    ../bin/hlf-deploy instantiateChaincode --configFile config.yaml \
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
    ../bin/hlf-deploy queryChaincode --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --chaincodeName mycc \
        query ${1}
}

function invokeChaincode() {
    ../bin/hlf-deploy invokeChaincode --configFile config.yaml \
        --channelName mychannel \
        --orgName Org1 \
        --endorsementOrgsName Org1,Org2 \
        --chaincodeName mycc \
        invoke a b 50
}

function cleanNetwork() {
    docker stack rm orderer
    docker stack rm peer0org1
    docker stack rm peer1org1
    docker stack rm peer0org2
    docker stack rm peer1org2

    sleep 10s

    docker volume prune
}

if [[ ! -f "../bin/hlf-deploy" ]]; then
    echo "hlf-deploy tool not found. exiting"
    exit 1
fi

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
    invokeChaincode
    queryChaincode a
    queryChaincode b
elif [[ "${mode}" == "down" ]]; then ## Clear the network
    cleanNetwork
fi