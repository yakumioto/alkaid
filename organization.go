package main

import (
	"bytes"
	"log"

	"github.com/yakumioto/hlf-deploy/internal/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/cobra"
)

func addAndUpdateOrgChannel(_ *cobra.Command, args []string) {
	utils.InitRPCClient(rpcAddress)
	sdk := utils.SDKNew(fabconfig)

	ordererCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(ordererOrgName))
	resMgmt, err := resmgmt.New(ordererCtx)
	if err != nil {
		log.Fatalln("resmgmt new error: ", err)
	}

	// get newest config
	configBytes := utils.GetNewestConfigWithConfigBlock(resMgmt, channelName, sysChannel)

	// get new organization config
	newOrgConfigBytes := utils.GetNewOrgConfigWithFielePath(orgConfig, orgName)

	// get modified config
	modifiedConfigBytes := utils.GetModifiedConfig(configBytes, newOrgConfigBytes, utils.ModifiedModAdd, ordererOrg, sysChannel)

	// get config.pb
	updateEnvelopePBBytes := utils.GetUpdateEnvelopeProtoBytes(configBytes, modifiedConfigBytes, channelName)

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfig:     bytes.NewBuffer(updateEnvelopePBBytes),
		SigningIdentities: utils.GetSigningIdentities(sdk.Context(), args),
	}

	txID, err := resMgmt.SaveChannel(req)
	if err != nil {
		log.Fatalf("save %s to %s error: %s", orgName, channelName, err)
	}

	log.Printf("save %s to %s txID: %s", orgName, channelName, txID.TransactionID)
}

func delOrgChannel(_ *cobra.Command, args []string) {
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
	modifiedConfigBytes := utils.GetModifiedConfig(configBytes, []byte(orgName), utils.ModifiedModDel, ordererOrg, sysChannel)

	// get config.pb
	updateEnvelopePBBytes := utils.GetUpdateEnvelopeProtoBytes(configBytes, modifiedConfigBytes, channelName)

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelName,
		ChannelConfig:     bytes.NewBuffer(updateEnvelopePBBytes),
		SigningIdentities: utils.GetSigningIdentities(sdk.Context(), args),
	}

	txID, err := resMgmt.SaveChannel(req)
	if err != nil {
		log.Fatalf("delete %s to %s error: %s", orgName, channelName, err)
	}

	log.Printf("delete %s to %s txID: %s", orgName, channelName, txID.TransactionID)
}
