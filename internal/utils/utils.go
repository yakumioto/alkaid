package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/yakumioto/hlf-deploy/internal/github.com/hyperledger/fabric/sdkinternal/configtxlator/update"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/util/protolator"
	"github.com/pkg/errors"
)

type Mod string
type ConsensusState string

const (
	ModifiedModAdd Mod = "Add"
	ModifiedModDel Mod = "Del"

	StateNormal      ConsensusState = "STATE_NORMAL"
	StateMaintenance ConsensusState = "STATE_MAINTENANCE"

	ConsensusEtcdRaft = "etcdraft"
)

type ConsensusOpts struct {
	State              string
	Type               string
	OrdererAddress     string
	KafkaBrokerAddress string
}

type RaftOpts struct {
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

type ChannelOpts struct {
	BatchTimeout       string
	BatchSizeAbsolute  string
	BatchSizePreferred string
	BatchSizeMessage   int
}

func (raft *RaftOpts) setDefaultOptions() {
	if raft.ElectionTick == 0 {
		raft.ElectionTick = 10
	}
	if raft.HeartbeatTick == 0 {
		raft.HeartbeatTick = 1
	}
	if raft.MaxInflightBlocks == 0 {
		raft.MaxInflightBlocks = 5
	}
	if raft.SnapshotIntervalSize == "" {
		raft.SnapshotIntervalSize = "20MB"
	}
	if raft.TickInterval == "" {
		raft.TickInterval = "500ms"
	}
}

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

func getProtoMessage(msgName string) proto.Message {
	var msg proto.Message

	switch msgName {
	case "common.Block":
		msg = &common.Block{}
	case "common.Config":
		msg = &common.Config{}
	case "common.Envelope":
		msg = &common.Envelope{}
	case "common.ConfigUpdate":
		msg = &common.ConfigUpdate{}
	default:
		msg = nil
	}
	return msg
}

func protoDecode(msgName string, input []byte) ([]byte, error) {
	var msg proto.Message
	if msg = getProtoMessage(msgName); msg == nil {
		return nil, errors.New("no message type")
	}

	if err := proto.Unmarshal(input, msg); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling")
	}

	output := bytes.NewBuffer(nil)

	err := protolator.DeepMarshalJSON(output, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "error encoding output")
	}

	return output.Bytes(), nil
}

func protoEncode(msgName string, input []byte) ([]byte, error) {
	var msg proto.Message
	if msg = getProtoMessage(msgName); msg == nil {
		return nil, errors.New("no message type")
	}

	intputbuf := bytes.NewBuffer(input)
	err := protolator.DeepUnmarshalJSON(intputbuf, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding input")
	}

	out, err := proto.Marshal(msg)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshaling")
	}

	return out, nil
}

func computeUpdate(channelName string, origin, updated []byte) ([]byte, error) {
	origConf := &common.Config{}
	if err := proto.Unmarshal(origin, origConf); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling original config")
	}

	updtConf := &common.Config{}
	if err := proto.Unmarshal(updated, updtConf); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling updated config")
	}

	cu, err := update.Compute(origConf, updtConf)
	if err != nil {
		return nil, errors.Wrapf(err, "error computing config update")
	}
	cu.ChannelId = channelName

	outBytes, err := proto.Marshal(cu)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshaling computed config update")
	}

	return outBytes, nil
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

	var groups map[string]interface{}
	if sysChannel {
		groups = cfg.(*SystemConfig).ChannelGroup.Groups.Consortiums.Groups.SampleConsortium.Groups

		if ordererOrg {
			groups = cfg.(*SystemConfig).ChannelGroup.Groups.Orderer.Groups
		}
	} else {
		groups = cfg.(*Config).ChannelGroup.Groups.Application.Groups

		if ordererOrg {
			groups = cfg.(*Config).ChannelGroup.Groups.Orderer.Groups
		}
	}

	switch mod {
	case ModifiedModAdd:
		if newOrgConfigBytes != nil {
			if err := json.Unmarshal(newOrgConfigBytes, newOrgConfig); err != nil {
				log.Fatalln(err)
			}
		}

		for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
			groups[orgName] = org
		}
	case ModifiedModDel:
		orgName = string(newOrgConfigBytes)

		delete(groups, orgName)
	}

	modifiedConfigBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal modified cfg json error:", err)
	}

	return modifiedConfigBytes
}

func GetChannelParamsModifiedConfig(configBytes []byte, channelOpts *ChannelOpts, sysChannel bool) []byte {
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

	if channelOpts.BatchTimeout != "" {
		batchTimoutValueMap["timeout"] = channelOpts.BatchTimeout
	}

	if channelOpts.BatchSizeAbsolute != "" {
		batchSizeValueMap["absolute_max_bytes"] = convertStorageUnit(channelOpts.BatchSizeAbsolute)
	}

	if channelOpts.BatchSizeMessage != 0 {
		batchSizeValueMap["max_message_count"] = channelOpts.BatchSizeMessage
	}

	if channelOpts.BatchSizePreferred != "" {
		batchSizeValueMap["preferred_max_bytes"] = convertStorageUnit(channelOpts.BatchSizePreferred)
	}

	modifiedConfigBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalln("marshal modified cfg json error:", err)
	}

	return modifiedConfigBytes
}

func setEtcdRaftConsenter(valueMap map[string]interface{}, raftOpts *RaftOpts) {
	optionsMap := getMap(getMap(valueMap, "metadata"), "options")

	// If the latest configuration block consensus type is not etcdraft, set the default optsions.
	if valueMap["type"] != ConsensusEtcdRaft {
		raftOpts.setDefaultOptions()
	}

	if raftOpts.ElectionTick != 0 {
		optionsMap["election_tick"] = raftOpts.ElectionTick
	}

	if raftOpts.HeartbeatTick != 0 {
		optionsMap["heartbeat_tick"] = raftOpts.HeartbeatTick
	}

	if raftOpts.MaxInflightBlocks != 0 {
		optionsMap["max_inflight_blocks"] = raftOpts.MaxInflightBlocks
	}

	if raftOpts.SnapshotIntervalSize != "" {
		optionsMap["snapshot_interval_size"] = convertStorageUnit(raftOpts.SnapshotIntervalSize)
	}

	if raftOpts.TickInterval != "" {
		optionsMap["tick_interval"] = raftOpts.TickInterval
	}

	valueMap["type"] = ConsensusEtcdRaft
}

func setRaftAddress(value map[string]interface{}, raftOpts *RaftOpts) {
	metadataMap := getMap(value, "metadata")
	consenters := make([]Consenters, 0)
	if metadataMap["consenters"] != nil {
		data, _ := json.Marshal(metadataMap["consenters"])
		_ = json.Unmarshal(data, &consenters)
	}

	var del bool
	for i, consenter := range consenters {
		if consenter.Host == raftOpts.Host {
			consenters = append(consenters[:i], consenters[i+1:]...)
			del = true
			break
		}
	}

	if !del {
		consenters = append(consenters, Consenters{
			Host:          raftOpts.Host,
			Port:          raftOpts.Port,
			ClientTLSCert: readCert2base64(raftOpts.ClientTLSCertPath),
			ServerTLSCert: readCert2base64(raftOpts.ServerTLSCertPath),
		})
	}

	metadataMap["consenters"] = consenters
}

func setOrdererAddress(configValues map[string]interface{}, ordererAddress string) {
	valueMap := getMap(getMap(configValues, "OrdererAddresses"), "value")
	addresses := valueMap["addresses"].([]interface{})

	var del bool
	for i, address := range addresses {
		if ordererAddress == address.(string) {
			addresses = append(addresses[:i], addresses[i+1:]...)
			del = true
			break
		}
	}

	if !del {
		addresses = append(addresses, ordererAddress)
	}

	valueMap["addresses"] = addresses
}

func setKafkaBroker(ordererValues map[string]interface{}, kafkaAddress string) {
	valueMap := getMap(getMap(ordererValues, "KafkaBrokers"), "value")
	addresses := valueMap["brokers"].([]interface{})
	var del bool
	for i, address := range addresses {
		if kafkaAddress == address.(string) {
			addresses = append(addresses[:i], addresses[i+1:]...)
			del = true
			break
		}
	}

	if !del {
		addresses = append(addresses, kafkaAddress)
	}

	valueMap["brokers"] = addresses
}

func GetConsensusStateModifiedConfig(configBytes []byte, consensus *ConsensusOpts, raftOpts *RaftOpts, sysChannel bool) []byte {
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

	if consensus.Type == ConsensusEtcdRaft {
		setEtcdRaftConsenter(valueMap, raftOpts)
	}

	if raftOpts.Host != "" && raftOpts.Port != 0 && valueMap["type"] == ConsensusEtcdRaft {
		setRaftAddress(valueMap, raftOpts)
	}

	if consensus.OrdererAddress != "" {
		setOrdererAddress(configValues, consensus.OrdererAddress)
	}

	if consensus.KafkaBrokerAddress != "" {
		setKafkaBroker(ordererValues, consensus.KafkaBrokerAddress)
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
