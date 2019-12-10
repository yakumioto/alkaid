package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"strings"

	"github.com/gogo/protobuf/proto"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Mod string

const (
	ModifiedModAdd Mod = "Add"
	ModifiedModDel Mod = "Del"
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

func GetSigningIdentities(ctx context.ClientProvider, orgs []string) []msp.SigningIdentity {
	signingIdentities := make([]msp.SigningIdentity, 0)
	for _, orgName := range orgs {
		mspClient, err := mspclient.New(ctx, mspclient.WithOrg(orgName))
		if err != nil {
			log.Fatalf("%s msp new error: %s", orgName, err)
		}
		identity, err := mspClient.GetSigningIdentity("Admin")
		if err != nil {
			log.Fatalf("%s get signing identity error: %s", orgName, err)
		}

		signingIdentities = append(signingIdentities, identity)
	}

	return signingIdentities
}

func InitRPCClient(address string) {
	var err error

	if client == nil {
		client, err = rpc.DialHTTP("tcp", address)
		if err != nil {
			log.Fatalln("dialling rpc error:", err)
		}
	}
}

func protoDecode(msgName string, input []byte) ([]byte, error) {
	return protoEncodeAndDecode("Proto.Decode", msgName, input)
}

func protoEncode(msgName string, input []byte) ([]byte, error) {
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

func computeUpdate(channelName string, origin, updated []byte) ([]byte, error) {
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

func GetStdConfigBytes(mspID string, configBytes []byte) []byte {
	format := `{"channel_group":{"groups":{"Application":{"groups":{"%s":%s}}}}}`
	return []byte(fmt.Sprintf(format, mspID, string(configBytes)))
}

func GetStdUpdateEnvelopBytes(channelName string, updateEnvelopBytes []byte) []byte {
	format := `{"payload":{"header":{"channel_header":{"channel_id":"%s", "type":2}},"data":{"config_update":%s}}}`
	return []byte(fmt.Sprintf(format, channelName, string(updateEnvelopBytes)))
}

func GetNewestConfigWithConfigBlock(resMgmt *resmgmt.Client, channelName string, sysChannel bool) []byte {
	blockPB, err := resMgmt.QueryConfigBlockFromOrderer(channelName)
	if err != nil {
		log.Fatalln(err)
	}
	blockPBBytes, err := proto.Marshal(blockPB)
	if err != nil {
		log.Fatalln(err)
	}

	blockBytes, err := protoDecode("common.Block", blockPBBytes)
	if err != nil {
		log.Fatalln("proto decode common.Block error:", err)
	}

	var block interface{}
	if sysChannel {
		block = new(SystemBlock)
	} else {
		block = new(Block)
	}
	if err := json.Unmarshal(blockBytes, block); err != nil {
		log.Fatalln("unmarshal block json error:", err)
	}

	var cfg interface{}
	if sysChannel {
		cfg = block.(*SystemBlock).Data.Data[0].Payload.Data.Config
	} else {
		cfg = block.(*Block).Data.Data[0].Payload.Data.Config
	}

	configBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal config json error:", err)
	}

	return configBytes
}

func GetNewOrgConfigWithFielePath(filePath, mspID string) []byte {
	newOrgFileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	return GetStdConfigBytes(mspID, newOrgFileBytes)
}

func GetModifiedConfig(configBytes []byte, newOrgConfigBytes []byte, mod Mod, sysChannel bool) []byte {
	var cfg interface{}

	if configBytes != nil {
		if sysChannel {
			cfg = new(SystemConfig)
		} else {
			cfg = new(Config)
		}

		if err := json.Unmarshal(configBytes, cfg); err != nil {
			log.Fatalln(err)
		}
	}

	newOrgConfig := new(Config)
	orgName := ""
	switch mod {
	case ModifiedModAdd:
		if newOrgConfigBytes != nil {
			if err := json.Unmarshal(newOrgConfigBytes, newOrgConfig); err != nil {
				log.Fatalln(err)
			}
		}
	case ModifiedModDel:
		orgName = string(newOrgConfigBytes)
	}

	switch mod {
	case ModifiedModAdd:
		if sysChannel {
			for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
				cfg.(*SystemConfig).ChannelGroup.Groups.Consortiums.Groups.SampleConsortium.Groups[orgName] = org
			}
		} else {
			for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
				cfg.(*Config).ChannelGroup.Groups.Application.Groups[orgName] = org
			}
		}
	case ModifiedModDel:
		if sysChannel {
			delete(cfg.(*SystemConfig).ChannelGroup.Groups.Consortiums.Groups.SampleConsortium.Groups, orgName)
		} else {
			delete(cfg.(*Config).ChannelGroup.Groups.Application.Groups, orgName)
		}
	}

	modifiedConfigBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal modified cfg json error:", err)
	}

	return modifiedConfigBytes
}

func GetUpdateEnvelopeProtoBytes(configBytes, modifiedConfigBytes []byte, channelName string) []byte {
	configPBBytes, err := protoEncode("common.Config", configBytes)
	if err != nil {
		log.Fatalln("proto encode common.Config error:", err)
	}

	// get modified config.pb
	modifiedConfigPBBytes, err := protoEncode("common.Config", modifiedConfigBytes)
	if err != nil {
		log.Fatalln("proto encode common.Config error:", err)
	}

	// get update.pb
	updateConfigPBBytes, err := computeUpdate(channelName, configPBBytes, modifiedConfigPBBytes)
	if err != nil {
		log.Fatalln("compute update error:", err)
	}

	// get update.json
	updateConfigBytes, err := protoDecode("common.ConfigUpdate", updateConfigPBBytes)
	if err != nil {
		log.Fatalln("proto decode common.ConfigUpdate error:", err)
	}
	updateEnvelopeBytes := GetStdUpdateEnvelopBytes(channelName, updateConfigBytes)

	// get update.pb
	updateEnvelopePBBytes, err := protoEncode("common.Envelope", updateEnvelopeBytes)
	if err != nil {
		log.Fatalln("proto encode common.Envelope error:", err)
	}

	return updateEnvelopePBBytes
}
