package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/yakumioto/hlf-deploy/internal/utils"

	"github.com/gogo/protobuf/proto"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/cobra"
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

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfigPath: channelTX,
		SigningIdentities: utils.GetSigningIdentities(sdk.Context(), args),
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

func updateOrdererParam(_ *cobra.Command, args []string) {
	utils.InitRPCClient(rpcAddress)
	sdk := utils.SDKNew(fabconfig)

	ordererCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(ordererOrgName))
	resMgmt, err := resmgmt.New(ordererCtx)
	if err != nil {
		log.Fatalln("resmgmt new error: ", err)
	}

	// get newest config
	configBytes := utils.GetNewestConfigWithConfigBlock(resMgmt, channelName, sysChannel)

	// get modified config
	modifiedConfigBytes := utils.GetChannelParamsModifiedConfig(configBytes, batchTimeout, batchSizeAbsolute, batchSizePreferred, batchSizeMessage, sysChannel)

	// get config.pb
	updateEnvelopePBBytes := utils.GetUpdateEnvelopeProtoBytes(configBytes, modifiedConfigBytes, channelName)

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfig:     bytes.NewBuffer(updateEnvelopePBBytes),
		SigningIdentities: utils.GetSigningIdentities(sdk.Context(), args),
	}

	txID, err := resMgmt.SaveChannel(req)
	if err != nil {
		log.Fatalf("update %s channel parameters error: %s", channelName, err)
	}

	log.Printf("update %s channel parameters successfully txID: %s", channelName, txID.TransactionID)
}

func getChannelConfigBlock(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("no output path")
	}

	sdk := utils.SDKNew(fabconfig)

	ordererCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(ordererOrgName))
	resMgmt, err := resmgmt.New(ordererCtx)
	if err != nil {
		log.Fatalln("resmgmt new error: ", err)
	}

	blockPB, err := resMgmt.QueryConfigBlockFromOrderer(channelName)
	if err != nil {
		log.Fatalln("query config block error:", err)
	}
	blockPBBytes, err := proto.Marshal(blockPB)
	if err != nil {
		log.Fatalln("proto marshal error:", err)
	}

	if err := ioutil.WriteFile(args[0], blockPBBytes, 0664); err != nil {
		log.Fatalln("write file error:", err)
	}

	log.Printf("write latest config block to %s", args[0])
}
