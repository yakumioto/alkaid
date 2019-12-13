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
    docker-compose up -d
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

function addOrganization() {
    ../bin/hlf-deploy addOrgChannel --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/org3.json \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org1 Org2
}

function updateOrganization() {
    ../bin/hlf-deploy updateOrgChannel --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/modify-org3.json \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org3
}

function deleteOrganization() {
    ../bin/hlf-deploy delOrgChannel --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgName Org3MSP \
        --rpcAddress localhost:1234 \
        Org1 Org2
}

function addOrdererOrganization() {
    ../bin/hlf-deploy addOrgChannel --configFile config.yaml \
        --channelName mychannel \
        --ordererOrgName OrdererOrg \
        --orgConfig channel-artifacts/newOrderer.json \
        --orgName OrdererOrg2 \
        --rpcAddress localhost:1234 \
        --ordererOrg \
        OrdererOrg
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
    invokeChaincode
    queryChaincode a
    queryChaincode b
    addOrganization
    addOrdererOrganization
    updateOrganization
    deleteOrganization

elif [[ "${mode}" == "down" ]]; then ## Clear the network
    cleanNetwork
fi