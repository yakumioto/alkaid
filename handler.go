package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gogo/protobuf/proto"
	mspprotos "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/spf13/cobra"
	"github.com/yakumioto/hlf-deploy/utils"
)

func createChannel(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no organitztion")
	}

	sdk := utils.SDKNew(fabconfig)

	ordererCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(ordererOrgName))
	resMgmt, err := resmgmt.New(ordererCtx)
	if err != nil {
		log.Fatalln("resmgmt new error: ", err)
	}

	signingIdentities := make([]msp.SigningIdentity, 0)

	for _, orgName := range args {
		mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
		if err != nil {
			log.Fatalf("%s msp new error: %s", orgName, err)
		}
		identity, err := mspClient.GetSigningIdentity("Admin")
		if err != nil {
			log.Fatalf("%s get signing identity error: %s", orgName, err)
		}
		signingIdentities = append(signingIdentities, identity)
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfigPath: channelTX,
		SigningIdentities: signingIdentities,
	}

	txID, err := resMgmt.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalln("create channel error:", err)
	}

	log.Println("create channel txID:", string(txID.TransactionID))
}

func uptateAnchorPeer(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatalln("organitztion != 1")
	}

	sdk := utils.SDKNew(fabconfig)

	orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(args[0]))
	resMgmt, err := resmgmt.New(orgCtx)
	if err != nil {
		log.Fatalf("%s new resmgmt error: %s", args[0], err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(args[0]))
	if err != nil {
		log.Fatalf("%s msp new error: %s", args[0], err)
	}

	identity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		log.Fatalf("%s get signing identity error: %s", args[0], err)
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfigPath: anchorPeerTxFile,
		SigningIdentities: []msp.SigningIdentity{identity},
	}

	txID, err := resMgmt.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("%s update anchor peer error: %s", args[0], err)
	}

	log.Printf("%s update anchor peer txID: %s", args[0], string(txID.TransactionID))
}

func joinChannel(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no organitztion")
	}

	sdk := utils.SDKNew(fabconfig)

	for _, orgName := range args {
		orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
		resMgmt, err := resmgmt.New(orgCtx)
		if err != nil {
			log.Fatalf("%s new resmgmt error: %s", orgName, err)
		}

		if err := resMgmt.JoinChannel(channelName, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
			log.Fatalf("%s join channel error: %s", orgName, err)
		}

		log.Printf("%s join channel successful", orgName)
	}
}

func installChaincode(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no organitztion")
	}

	sdk := utils.SDKNew(fabconfig)

	ccpkg, err := gopackager.NewCCPackage(chaincodePath, goPath)
	if err != nil {
		log.Fatalln("new chaincode package error:", err)
	}

	for _, orgName := range args {
		orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
		resMgmt, err := resmgmt.New(orgCtx)
		if err != nil {
			log.Fatalf("%s new resmgmt error: %s", orgName, err)
		}

		if _, err := resMgmt.InstallCC(resmgmt.InstallCCRequest{
			Name:    chaincodeName,
			Path:    chaincodePath,
			Version: chaincodeVersion,
			Package: ccpkg,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
			log.Fatalf("%s install chaincode error: %s", orgName, err)
		}

		log.Printf("%s install chaincode successful", orgName)
	}
}

func instantiateAndUpgradeChaincode(cmd *cobra.Command, args []string) {
	sdk := utils.SDKNew(fabconfig)

	ccPolicy := cauthdsl.SignedByNOutOfGivenRole(chaincodePolicyNOutOf, mspprotos.MSPRole_MEMBER, chaincodePolicy)

	orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
	resMgmt, err := resmgmt.New(orgCtx)
	if err != nil {
		log.Fatalf("%s new resmgmt error: %s", orgName, err)
	}

	ccArgs := make([][]byte, 0)
	ccArgs = append(ccArgs, []byte("init"))
	for _, arg := range args {
		ccArgs = append(ccArgs, []byte(arg))
	}

	switch cmd.Use {
	case "instantiateChaincode":
		res, err := resMgmt.InstantiateCC(channelName, resmgmt.InstantiateCCRequest{
			Name:    chaincodeName,
			Path:    chaincodePath,
			Version: chaincodeVersion,
			Args:    ccArgs,
			Policy:  ccPolicy,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			log.Fatalf("%s instantiate chaincode error: %s", orgName, err)
		}

		log.Printf("%s instantiate chaincode txID: %s args: %s", orgName, res.TransactionID, args)

	case "upgradeChaincode":
		res, err := resMgmt.UpgradeCC(channelName, resmgmt.UpgradeCCRequest{
			Name:    chaincodeName,
			Path:    chaincodePath,
			Version: chaincodeVersion,
			Args:    ccArgs,
			Policy:  ccPolicy,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			log.Fatalf("%s instantiate chaincode error: %s", orgName, err)
		}

		log.Printf("%s instantiate chaincode txID: %s args: %s", orgName, res.TransactionID, args)
	}
}

func queryAdnInvokeChaincode(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no args")
	}

	sdk := utils.SDKNew(fabconfig)

	channelCtx := sdk.ChannelContext(channelName, fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelCtx)
	if err != nil {
		log.Fatalf("%s new channel error: %s", orgName, err)
	}

	ccArgs := make([][]byte, 0)
	for _, arg := range args[1:] {
		ccArgs = append(ccArgs, []byte(arg))
	}

	switch cmd.Use {
	case "queryChaincode":
		res, err := channelClient.Query(channel.Request{
			ChaincodeID: chaincodeName,
			Fcn:         args[0],
			Args:        ccArgs,
		}, channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			log.Fatalf("%s invoke error: %s", orgName, err)
		}

		log.Printf("%s query chaincode txID: %s args: %s result: %s", orgName, res.TransactionID, args, string(res.Payload))

	case "invokeChaincode":
		peers, err := utils.GetOrgsTargetPeers(sdk, endorsementOrgsName)
		if err != nil {
			log.Fatalf("get target peers error: %s", err)
		}

		res, err := channelClient.Execute(channel.Request{
			ChaincodeID: chaincodeName,
			Fcn:         args[0],
			Args:        ccArgs,
		}, channel.WithTargetEndpoints(peers...), channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			log.Fatalf("%s invoke error: %s", orgName, err)
		}

		log.Printf("%s invoke chaincode txID: %s args: %s", orgName, res.TransactionID, args)
	}
}

func addOrgChannel(_ *cobra.Command, args []string) {
	utils.InitRPCClient(rpcAddress)

	sdk := utils.SDKNew(fabconfig)

	ordererCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(ordererOrgName))
	resMgmt, err := resmgmt.New(ordererCtx)
	if err != nil {
		log.Fatalln("resmgmt new error: ", err)
	}

	// get last config.block witch channelName
	blockPB, err := resMgmt.QueryConfigBlockFromOrderer(channelName)
	if err != nil {
		log.Fatalln(err)
	}
	blockPBBytes, err := proto.Marshal(blockPB)
	if err != nil {
		log.Fatalln(err)
	}

	blockBytes, err := utils.ProtoDecode("common.Block", blockPBBytes)
	if err != nil {
		log.Fatalln("proto decode common.Block error:", err)
	}

	block := new(Block)
	if err := json.Unmarshal(blockBytes, block); err != nil {
		log.Fatalln("unmarshal block json error:", err)
	}

	config := block.Data.Data[0].Payload.Data.Config

	configBytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("marshal config json error:", err)
	}

	// get new org config.json
	newOrgFileBytes, err := ioutil.ReadFile(newOrgConfig)
	if err != nil {
		log.Fatalln(err)
	}

	newOrgConfigBytes := utils.GetStdConfigBytes(sysChannel, newOrgMSPID, newOrgFileBytes)

	newOrgConfig := new(Config)
	if err := json.Unmarshal(newOrgConfigBytes, newOrgConfig); err != nil {
		log.Fatalln(err)
	}

	// get modified config.json
	for orgName, org := range newOrgConfig.ChannelGroup.Groups.Application.Groups {
		config.ChannelGroup.Groups.Application.Groups[orgName] = org
	}

	modifiedConfigBytes, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("marshal modified config json error:", err)
	}

	// get config.pb
	configPBBytes, err := utils.ProtoEncode("common.Config", configBytes)
	if err != nil {
		log.Fatalln("proto encode common.Config error:", err)
	}

	// get modified config.pb
	modifiedConfigPBBytes, err := utils.ProtoEncode("common.Config", modifiedConfigBytes)
	if err != nil {
		log.Fatalln("proto encode common.Config error:", err)
	}

	// get update.pb
	updateConfigPBBytes, err := utils.ComputeUpdate(channelName, configPBBytes, modifiedConfigPBBytes)
	if err != nil {
		log.Fatalln("compute update error:", err)
	}

	// get update.json
	updateConfigBytes, err := utils.ProtoDecode("common.ConfigUpdate", updateConfigPBBytes)
	if err != nil {
		log.Fatalln("proto decode common.ConfigUpdate error:", err)
	}
	updateEnvelopeBytes := utils.GetStdUpdateEnvelopBytes(channelName, updateConfigBytes)

	// get update.pb
	updateEnvelopePBBytes, err := utils.ProtoEncode("common.Envelope", updateEnvelopeBytes)
	if err != nil {
		log.Fatalln("proto encode common.Envelope error:", err)
	}

	// sign envelope
	signingIdentities := make([]msp.SigningIdentity, 0)
	for _, orgName := range args {
		mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
		if err != nil {
			log.Fatalf("%s msp new error: %s", orgName, err)
		}
		identity, err := mspClient.GetSigningIdentity("Admin")
		if err != nil {
			log.Fatalf("%s get signing identity error: %s", orgName, err)
		}
		signingIdentities = append(signingIdentities, identity)
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfig:     bytes.NewBuffer(updateEnvelopePBBytes),
		SigningIdentities: signingIdentities,
	}

	txID, err := resMgmt.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		log.Fatalf("add organization to %s error: %s", channelName, err)
	}

	log.Printf("add organization to %s txID: %s", channelName, txID.TransactionID)
}
