package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	mspprotos "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  `hlf-deploy`,
		Long: `Hyperledger Fabric Deploy`,
	}

	createChannelCmd = &cobra.Command{
		Use:   `createChannel`,
		Short: `Create channel.`,
		Run:   createChannel,
	}

	uptateAnchorPeerCmd = &cobra.Command{
		Use:   `updateAnchorPeer`,
		Short: `Update anchor peer.`,
		Run:   uptateAnchorPeer,
	}

	joinChannelCmd = &cobra.Command{
		Use:   `joinChannel`,
		Short: `Join channel.`,
		Run:   joinChannel,
	}

	installChaincodeCmd = &cobra.Command{
		Use:   `installChaincode`,
		Short: `Install chaincode.`,
		Run:   installChaincode,
	}

	instantiateChaincodeCmd = &cobra.Command{
		Use:   `instantiateChaincode`,
		Short: `Instantiate chaincode.`,
		Run:   instantiateAndUpgradeChaincode,
	}

	upgradeChaincodeCmd = &cobra.Command{
		Use:   `upgradeChaincode`,
		Short: `Upgrade chaincode.`,
		Run:   instantiateAndUpgradeChaincode,
	}

	queryChaincodeCmd = &cobra.Command{
		Use:   `queryChaincode`,
		Short: `Query chaincode.`,
		Run:   queryAdnInvokeChaincode,
	}

	invokeChaincodeCmd = &cobra.Command{
		Use:   `invokeChaincode`,
		Short: `Invoke chaincode.`,
		Run:   queryAdnInvokeChaincode,
	}
)

var (
	fabconfig             string
	channelTX             string
	channelName           string
	ordererOrgName        string
	orgName               string
	anchorPeerTxFile      string
	goPath                string
	chaincodeName         string
	chaincodePath         string
	chaincodeVersion      string
	chaincodePolicy       []string
	chaincodePolicyNOutOf int32
	endorsementOrgsName   []string
)

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(createChannelCmd)
	rootCmd.AddCommand(uptateAnchorPeerCmd)
	rootCmd.AddCommand(joinChannelCmd)
	rootCmd.AddCommand(installChaincodeCmd)
	rootCmd.AddCommand(instantiateChaincodeCmd)
	rootCmd.AddCommand(upgradeChaincodeCmd)
	rootCmd.AddCommand(queryChaincodeCmd)
	rootCmd.AddCommand(invokeChaincodeCmd)

	createChannelCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	createChannelCmd.Flags().StringVar(&channelTX, "channelTxFile", "channel.tx", "Channel TX file path.")
	createChannelCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	createChannelCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")

	uptateAnchorPeerCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	uptateAnchorPeerCmd.Flags().StringVar(&anchorPeerTxFile, "anchorPeerTxFile", "anchor.tx", "Anchor peer TX file.")
	uptateAnchorPeerCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	uptateAnchorPeerCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")

	joinChannelCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	joinChannelCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")

	installChaincodeCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	installChaincodeCmd.Flags().StringVar(&goPath, "goPath", "", "Set the GOPATH env.")
	installChaincodeCmd.Flags().StringVar(&chaincodePath, "chaincodePath", "./", "Chaincode path.")
	installChaincodeCmd.Flags().StringVar(&chaincodeName, "chaincodeName", "chaincode", "Chaincode name.")
	installChaincodeCmd.Flags().StringVar(&chaincodeVersion, "chaincodeVersion", "v0.0.0", "Chaincode version.")

	instantiateChaincodeCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	instantiateChaincodeCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	instantiateChaincodeCmd.Flags().StringVar(&orgName, "orgName", "", "Set the organitztion name.")
	instantiateChaincodeCmd.Flags().StringSliceVar(&chaincodePolicy, "chaincodePolicy", nil, "Set the chaincode policy, separated by ','.")
	instantiateChaincodeCmd.Flags().Int32Var(&chaincodePolicyNOutOf, "chaincodePolicyNOutOf", 1, "Requires N out of the slice of policies to evaluate to true")
	instantiateChaincodeCmd.Flags().StringVar(&chaincodePath, "chaincodePath", "./", "Chaincode path.")
	instantiateChaincodeCmd.Flags().StringVar(&chaincodeName, "chaincodeName", "chaincode", "Chaincode name.")
	instantiateChaincodeCmd.Flags().StringVar(&chaincodeVersion, "chaincodeVersion", "v0.0.0", "Chaincode version.")

	upgradeChaincodeCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	upgradeChaincodeCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	upgradeChaincodeCmd.Flags().StringVar(&orgName, "orgName", "", "Set the organitztion name.")
	upgradeChaincodeCmd.Flags().StringSliceVar(&chaincodePolicy, "chaincodePolicy", nil, "Set the chaincode policy, separated by ','.")
	upgradeChaincodeCmd.Flags().Int32Var(&chaincodePolicyNOutOf, "chaincodePolicyNOutOf", 1, "Requires N out of the slice of policies to evaluate to true")
	upgradeChaincodeCmd.Flags().StringVar(&chaincodePath, "chaincodePath", "./", "Chaincode path.")
	upgradeChaincodeCmd.Flags().StringVar(&chaincodeName, "chaincodeName", "chaincode", "Chaincode name.")
	upgradeChaincodeCmd.Flags().StringVar(&chaincodeVersion, "chaincodeVersion", "v0.0.0", "Chaincode version.")

	queryChaincodeCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	queryChaincodeCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	queryChaincodeCmd.Flags().StringVar(&orgName, "orgName", "", "Set the organitztion name.")
	queryChaincodeCmd.Flags().StringVar(&chaincodeName, "chaincodeName", "chaincode", "Chaincode name.")

	invokeChaincodeCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	invokeChaincodeCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	invokeChaincodeCmd.Flags().StringVar(&orgName, "orgName", "", "Set the organitztion name.")
	invokeChaincodeCmd.Flags().StringSliceVar(&endorsementOrgsName, "endorsementOrgsName", nil, "Set the endorsement organitztions name.")
	invokeChaincodeCmd.Flags().StringVar(&chaincodeName, "chaincodeName", "chaincode", "Chaincode name.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func fabsdkNew() *fabsdk.FabricSDK {
	sdk, err := fabsdk.New(config.FromFile(fabconfig))
	if err != nil {
		log.Fatalln("new fabsdk error:", err)
	}
	return sdk
}

func createChannel(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no organitztion")
	}

	sdk := fabsdkNew()

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

	sdk := fabsdkNew()

	orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(args[0]))
	resMgmt, err := resmgmt.New(orgCtx)
	if err != nil {
		log.Fatalf("%s org new resmgmt error: %s", args[0], err)
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

	sdk := fabsdkNew()

	for _, orgName := range args {
		orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
		resMgmt, err := resmgmt.New(orgCtx)
		if err != nil {
			log.Fatalf("%s org new resmgmt error: %s", orgName, err)
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

	sdk := fabsdkNew()

	ccpkg, err := gopackager.NewCCPackage(chaincodePath, goPath)
	if err != nil {
		log.Fatalln("new chaincode package error:", err)
	}

	for _, orgName := range args {
		orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
		resMgmt, err := resmgmt.New(orgCtx)
		if err != nil {
			log.Fatalf("%s org new resmgmt error: %s", orgName, err)
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
	sdk := fabsdkNew()

	ccPolicy := cauthdsl.SignedByNOutOfGivenRole(chaincodePolicyNOutOf, mspprotos.MSPRole_MEMBER, chaincodePolicy)

	orgCtx := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
	resMgmt, err := resmgmt.New(orgCtx)
	if err != nil {
		log.Fatalf("%s org new resmgmt error: %s", orgName, err)
	}

	ccArgs := make([][]byte, 0)
	ccArgs = append(ccArgs, []byte("init"))
	for _, arg := range args {
		ccArgs = append(ccArgs, []byte(arg))
	}

	switch cmd.Use {
	case "instantiateChaincode":
		resp, err := resMgmt.InstantiateCC(channelName, resmgmt.InstantiateCCRequest{
			Name:    chaincodeName,
			Path:    chaincodePath,
			Version: chaincodeVersion,
			Args:    ccArgs,
			Policy:  ccPolicy,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			log.Fatalf("%s org instantiate chaincode error: %s", orgName, err)
		}

		log.Printf("%s org instantiate chaincode txID: %s", orgName, resp.TransactionID)

	case "upgradeChaincode":
		res, err := resMgmt.UpgradeCC(channelName, resmgmt.UpgradeCCRequest{
			Name:    chaincodeName,
			Path:    chaincodePath,
			Version: chaincodeVersion,
			Args:    ccArgs,
			Policy:  ccPolicy,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			log.Fatalf("%s org instantiate chaincode error: %s", orgName, err)
		}

		log.Printf("%s org instantiate chaincode txID: %s", orgName, res.TransactionID)
	}
}

func queryAdnInvokeChaincode(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("no args")
	}

	sdk := fabsdkNew()

	channelCtx := sdk.ChannelContext(channelName, fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(channelCtx)
	if err != nil {
		log.Fatalf("%s org new channel error: %s", orgName, err)
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
			log.Fatalf("%s org invoke error: %s", orgName, err)
		}

		log.Printf("%s org query chaincode txID: %s result: %s", orgName, res.TransactionID, string(res.Payload))

	case "invokeChaincode":
		peers, err := getOrgsTargetPeers(sdk, endorsementOrgsName)
		if err != nil {
			log.Fatalf("get target peers error: %s", err)
		}

		res, err := channelClient.Execute(channel.Request{
			ChaincodeID: chaincodeName,
			Fcn:         args[0],
			Args:        ccArgs,
		}, channel.WithTargetEndpoints(peers...), channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			log.Fatalf("%s org invoke error: %s", orgName, err)
		}

		log.Printf("%s org invoke chaincode txID: %s", orgName, res.TransactionID)
	}
}

func getOrgsTargetPeers(sdk *fabsdk.FabricSDK, orgsName []string) ([]string, error) {
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
