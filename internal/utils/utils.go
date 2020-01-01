package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"strconv"
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
type ConsensusState string

const (
	ModifiedModAdd Mod = "Add"
	ModifiedModDel Mod = "Del"

	StateNormal      ConsensusState = "STATE_NORMAL"
	StateMaintenance ConsensusState = "STATE_MAINTENANCE"
)

var (
	client *rpc.Client
)

func GetConsensusState(status string) ConsensusState {
	switch status {
	case "Normal":
		return StateNormal
	case "Maintenance":
		return StateMaintenance
	}
	return ""
}

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
		return nil, fmt.Errorf("get orgs target peers error: %s", err)
	}

	networkConfig := fab.NetworkConfig{}
	if err := lookup.New(configBackend).UnmarshalKey("organizations", &networkConfig.Organizations); err != nil {
		return nil, fmt.Errorf("lookup unmarshal key error: %s", err)
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
			log.Fatalln("dialing rpc error:", err)
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

	err = json.Unmarshal(blockBytes, block)
	if err != nil {
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

func GetModifiedConfig(configBytes, newOrgConfigBytes []byte, mod Mod, ordererOrg, sysChannel bool) []byte {
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
			if ordererOrg {
				for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
					cfg.(*SystemConfig).ChannelGroup.Groups.Orderer.Groups[orgName] = org
				}
				break
			}
			for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
				cfg.(*SystemConfig).ChannelGroup.Groups.Consortiums.Groups.SampleConsortium.Groups[orgName] = org
			}
		} else {
			if ordererOrg {
				for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
					cfg.(*Config).ChannelGroup.Groups.Orderer.Groups[orgName] = org
				}
				break
			}
			for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
				cfg.(*Config).ChannelGroup.Groups.Application.Groups[orgName] = org
			}
		}
	case ModifiedModDel:
		if sysChannel {
			if ordererOrg {
				delete(cfg.(*SystemConfig).ChannelGroup.Groups.Orderer.Groups, orgName)
				break
			}
			delete(cfg.(*SystemConfig).ChannelGroup.Groups.Consortiums.Groups.SampleConsortium.Groups, orgName)
		} else {
			if ordererOrg {
				delete(cfg.(*Config).ChannelGroup.Groups.Orderer.Groups, orgName)
			}
			delete(cfg.(*Config).ChannelGroup.Groups.Application.Groups, orgName)
		}
	}

	modifiedConfigBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal modified cfg json error:", err)
	}

	return modifiedConfigBytes
}

func GetChannelParamsModifiedConfig(configBytes []byte,
	batchTimeout, batchSizeAbsolute, batchSizePreferred string, batchSizeMessage int, sysChannel bool) []byte {

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

	var values map[string]interface{}
	if sysChannel {
		values = cfg.(*SystemConfig).ChannelGroup.Groups.Orderer.Values
	} else {
		values = cfg.(*Config).ChannelGroup.Groups.Orderer.Values
	}

	batchTimoutValueMap := getMap(getMap(values, "BatchTimeout"), "value")
	batchSizeValueMap := getMap(getMap(values, "BatchTimeout"), "value")

	if batchTimeout != "" {
		batchTimoutValueMap["timeout"] = batchTimeout
	}

	if batchSizeAbsolute != "" {
		batchSizeValueMap["absolute_max_bytes"] = convertStorageUnit(batchSizeAbsolute)
	}

	if batchSizeMessage != 0 {
		batchSizeValueMap["max_message_count"] = batchSizeMessage
	}

	if batchSizePreferred != "" {
		batchSizeValueMap["preferred_max_bytes"] = convertStorageUnit(batchSizePreferred)
	}

	modifiedConfigBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal modified cfg json error:", err)
	}

	return modifiedConfigBytes
}

type ConsensusOptions struct {
	State              string
	Type               string
	OrdererAddress     string
	KafkaBrokerAddress string
}

type EtcdRaftOptions struct {
	ElectionTick         int
	HeartbeatTick        int
	MaxInflightBlocks    int
	SnapshotIntervalSize string
	TickInterval         string
	Host                 string
	Port                 int
	ClientTLSCertPath    string
	ServerTLSCertPath    string
}

func GetChannelConsensusStateModifiedConfig(configBytes []byte, consensus ConsensusOptions, raftOptions EtcdRaftOptions,
	sysChannel bool) []byte {
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

	var ordererValues map[string]interface{}
	var configValues map[string]interface{}
	if sysChannel {
		ordererValues = cfg.(*SystemConfig).ChannelGroup.Groups.Orderer.Values
		configValues = cfg.(*SystemConfig).ChannelGroup.Values
	} else {
		ordererValues = cfg.(*Config).ChannelGroup.Groups.Orderer.Values
		configValues = cfg.(*Config).ChannelGroup.Values
	}

	state := GetConsensusState(consensus.State)

	consensusTypeMap := getMap(ordererValues, "ConsensusType")
	valueMap := getMap(consensusTypeMap, "value")

	if state != "" {
		valueMap["state"] = state
	}

	if consensus.Type != "" {
		if consensus.Type == "etcdraft" {

			optionsMap := getMap(getMap(valueMap, "metadata"), "options")

			// If the latest configuration block consensus type is not etcdraft, set the default optsions.
			if valueMap["type"] != "etcdraft" {
				if raftOptions.ElectionTick == 0 {
					raftOptions.ElectionTick = 10
				}
				if raftOptions.HeartbeatTick == 0 {
					raftOptions.HeartbeatTick = 1
				}
				if raftOptions.MaxInflightBlocks == 0 {
					raftOptions.MaxInflightBlocks = 5
				}
				if raftOptions.SnapshotIntervalSize == "" {
					raftOptions.SnapshotIntervalSize = "20MB"
				}
				if raftOptions.TickInterval == "" {
					raftOptions.TickInterval = "500ms"
				}
			}

			if raftOptions.ElectionTick != 0 {
				optionsMap["election_tick"] = raftOptions.ElectionTick
			}

			if raftOptions.HeartbeatTick != 0 {
				optionsMap["heartbeat_tick"] = raftOptions.HeartbeatTick
			}

			if raftOptions.MaxInflightBlocks != 0 {
				optionsMap["max_inflight_blocks"] = raftOptions.MaxInflightBlocks
			}

			if raftOptions.SnapshotIntervalSize != "" {
				optionsMap["snapshot_interval_size"] = convertStorageUnit(raftOptions.SnapshotIntervalSize)
			}

			if raftOptions.TickInterval != "" {
				optionsMap["tick_interval"] = raftOptions.TickInterval
			}
		}

		valueMap["type"] = consensus.Type
	}

	if raftOptions.Host != "" && raftOptions.Port != 0 && valueMap["type"] == "etcdraft" {

		metadataMap := getMap(valueMap, "metadata")
		consenters := make([]Consenters, 0)
		if metadataMap["consenters"] != nil {
			data, _ := json.Marshal(metadataMap["consenters"])
			_ = json.Unmarshal(data, &consenters)
		}

		var del bool
		for i, consenter := range consenters {
			if consenter.Host == raftOptions.Host {
				consenters = append(consenters[:i], consenters[i+1:]...)
				del = true
				break
			}
		}

		if !del {
			consenters = append(consenters, Consenters{
				Host:          raftOptions.Host,
				Port:          raftOptions.Port,
				ClientTLSCert: readCert2base64(raftOptions.ClientTLSCertPath),
				ServerTLSCert: readCert2base64(raftOptions.ServerTLSCertPath),
			})
		}

		metadataMap["consenters"] = consenters
	}

	if consensus.OrdererAddress != "" {
		valueMap := getMap(getMap(configValues, "OrdererAddresses"), "value")
		addresses := valueMap["addresses"].([]interface{})

		var del bool
		for i, address := range addresses {
			if consensus.OrdererAddress == address.(string) {
				addresses = append(addresses[:i], addresses[i+1:]...)
				del = true
				break
			}
		}

		if !del {
			addresses = append(addresses, consensus.OrdererAddress)
		}

		valueMap["addresses"] = addresses
	}

	if consensus.KafkaBrokerAddress != "" {
		valueMap := getMap(getMap(ordererValues, "KafkaBrokers"), "value")
		addresses := valueMap["brokers"].([]interface{})
		var del bool
		for i, address := range addresses {
			if consensus.OrdererAddress == address.(string) {
				addresses = append(addresses[:i], addresses[i+1:]...)
				del = true
				break
			}
		}

		if !del {
			addresses = append(addresses, consensus.KafkaBrokerAddress)
		}

		valueMap["brokers"] = addresses
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

func convertStorageUnit(data string) int64 {
	var KB int64 = 1024
	var MB = 1024 * KB

	num, err := strconv.Atoi(data[:len(data)-2])
	if err != nil {
		log.Fatalln("strconv atoi error:", err)
	}

	data = strings.ToLower(data)
	if strings.Contains(data, "kb") {
		return int64(num) * KB
	}

	if strings.Contains(data, "mb") {
		return int64(num) * MB
	}

	return 0
}

func getMap(data map[string]interface{}, key string) map[string]interface{} {
	if data[key] == nil {
		data[key] = make(map[string]interface{})
	}
	return data[key].(map[string]interface{})
}

func readCert2base64(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read file error: %s", err)
	}

	return base64.StdEncoding.EncodeToString(data)
}
