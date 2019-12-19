package main

import (
	"log"

	mspprotos "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/spf13/cobra"
	"github.com/yakumioto/hlf-deploy/internal/utils"
)

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
	case "instantiate":
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

	case "upgrade":
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
	case "query":
		res, err := channelClient.Query(channel.Request{
			ChaincodeID: chaincodeName,
			Fcn:         args[0],
			Args:        ccArgs,
		}, channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			log.Fatalf("%s invoke error: %s", orgName, err)
		}

		log.Printf("%s query chaincode txID: %s args: %s result: %s", orgName, res.TransactionID, args, string(res.Payload))

	case "invoke":
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
