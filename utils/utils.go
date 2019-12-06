package utils

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	client *rpc.Client
)

func SDKNew(fabconfig string) *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(config.FromFile(fabconfig))
	if err != nil {
		log.Fatalln("new fabsdk error:", err)
	}
	return sdk
}

func GetOrgsTargetPeers(sdk *fabsdk.FabricSDK, orgsName []string) ([]string, error) {
	configBackend, err := sdk.Config()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get orgs target peers error: %s", err))
	}

	networkConfig := fab.NetworkConfig{}
	if err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations); err != nil {
		return nil, errors.New(fmt.Sprintf("lookup unmarshal key error: %s", err))
	}

	var peers []string
	for _, org := range orgsName {
		orgConfig, ok := networkConfig.Organizations[strings.ToLower(org)]
		if !ok {
			continue
		}
		peers = append(peers, orgConfig.Peers...)
	}

	return peers, nil
}

func InitRPCClient(address string) {
	var err error

	client, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatalln("dialling rpc error:", err)
	}
}

func ProtoDecode(msgName string, input []byte) ([]byte, error) {
	return protoEncodeAndDecode("Proto.Decode", msgName, input)
}

func ProtoEncode(msgName string, input []byte) ([]byte, error) {
	return protoEncodeAndDecode("Proto.Encode", msgName, input)
}

func protoEncodeAndDecode(typ, msgName string, input []byte) ([]byte, error) {
	var reply []byte

	if err := client.Call(typ, struct {
		MsgName string
		Input   []byte
	}{
		msgName,
		input,
	}, &reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func ComputeUpdate(channelName string, origin, updated []byte) ([]byte, error) {
	var reply []byte

	if err := client.Call("Compute.Update", struct {
		ChannelName string
		Origin      []byte
		Updated     []byte
	}{
		channelName,
		origin,
		updated,
	}, &reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func GetStdConfigBytes(sysChannel bool, mspID string, configBytes []byte) []byte {
	var format string

	if sysChannel {
		format = `{"channel_group":{"groups":{"Consortiums":{"groups":{"SampleConsortium":{"groups":{"%s":%s}}}}}}}`
	} else {
		format = `{"channel_group":{"groups":{"Application":{"groups":{"%s":%s}}}}}`
	}

	return []byte(fmt.Sprintf(format, mspID, string(configBytes)))
}

func GetStdUpdateEnvelopBytes(channelName string, updateEnvelopBytes []byte) []byte {
	format := `{"payload":{"header":{"channel_header":{"channel_id":"%s", "type":2}},"data":{"config_update":%s}}}`
	return []byte(fmt.Sprintf(format, channelName, string(updateEnvelopBytes)))
}
