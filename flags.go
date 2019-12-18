package main

import "github.com/spf13/cobra"

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
	rpcAddress            string
	sysChannel            bool
	orgConfig             string
	ordererOrg            bool
	batchTimeout          string
	batchSizeAbsolute     string
	batchSizeMessage      int
	batchSizePreferred    string
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
	rootCmd.AddCommand(addOrgChannelCmd)
	rootCmd.AddCommand(updateOrgChannelCmd)
	rootCmd.AddCommand(delOrgChannelCmd)
	rootCmd.AddCommand(updateChannelParamCmd)

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

	addOrgChannelCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	addOrgChannelCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")
	addOrgChannelCmd.Flags().BoolVar(&sysChannel, "sysChannel", false, "Channel is system channel.")
	addOrgChannelCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	addOrgChannelCmd.Flags().StringVar(&orgConfig, "orgConfig", "org.json", "New organitztion config material in JSON.")
	addOrgChannelCmd.Flags().StringVar(&orgName, "orgName", "", "New organitztion MSP id.")
	addOrgChannelCmd.Flags().StringVar(&rpcAddress, "rpcAddress", "localhost:1234", "hlf-tools Address.")
	addOrgChannelCmd.Flags().BoolVar(&ordererOrg, "ordererOrg", false, "Organization is the orderer organization.")

	updateOrgChannelCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	updateOrgChannelCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")
	updateOrgChannelCmd.Flags().BoolVar(&sysChannel, "sysChannel", false, "Channel is system channel.")
	updateOrgChannelCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	updateOrgChannelCmd.Flags().StringVar(&orgConfig, "orgConfig", "org.json", "New organitztion config material in JSON.")
	updateOrgChannelCmd.Flags().StringVar(&orgName, "orgName", "", "New organitztion MSP id.")
	updateOrgChannelCmd.Flags().StringVar(&rpcAddress, "rpcAddress", "localhost:1234", "hlf-tools Address.")
	updateOrgChannelCmd.Flags().BoolVar(&ordererOrg, "ordererOrg", false, "Organization is the orderer organization")

	delOrgChannelCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	delOrgChannelCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")
	delOrgChannelCmd.Flags().BoolVar(&sysChannel, "sysChannel", false, "Channel is system channel.")
	delOrgChannelCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	delOrgChannelCmd.Flags().StringVar(&orgName, "orgName", "", "New organitztion MSP id.")
	delOrgChannelCmd.Flags().StringVar(&rpcAddress, "rpcAddress", "localhost:1234", "hlf-tools Address.")
	delOrgChannelCmd.Flags().BoolVar(&ordererOrg, "ordererOrg", false, "Organization is the orderer organization")

	updateChannelParamCmd.Flags().StringVar(&fabconfig, "configFile", "config.yaml", "Fabric SDK config file path.")
	updateChannelParamCmd.Flags().StringVar(&ordererOrgName, "ordererOrgName", "OrdererOrg", "Orderer organitztion name.")
	updateChannelParamCmd.Flags().BoolVar(&sysChannel, "sysChannel", false, "Channel is system channel.")
	updateChannelParamCmd.Flags().StringVar(&channelName, "channelName", "mychannel", "Channel name.")
	updateChannelParamCmd.Flags().StringVar(&rpcAddress, "rpcAddress", "localhost:1234", "hlf-tools Address.")
	updateChannelParamCmd.Flags().StringVar(&batchTimeout, "batchTimeout", "2s", "set batch timeout.")
	updateChannelParamCmd.Flags().StringVar(&batchSizeAbsolute, "absoluteMaxBytes", "99MB", "set batch size absolute max bytes.")
	updateChannelParamCmd.Flags().StringVar(&batchSizePreferred, "preferredMaxBytes", "512KB", "set batch size preferred max bytes.")
	updateChannelParamCmd.Flags().IntVar(&batchSizeMessage, "sizeMessageMaxCount", 10, "set batch size max message count.")
}
