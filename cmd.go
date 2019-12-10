package main

import "github.com/spf13/cobra"

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

	addOrgChannelCmd = &cobra.Command{
		Use:   `addOrgChannel`,
		Short: `add organization to channel.`,
		Run:   addAdnUpdateOrgChannel,
	}

	updateOrgChannelCmd = &cobra.Command{
		Use:   `updateOrgChannel`,
		Short: `update organization to channel.`,
		Run:   addAdnUpdateOrgChannel,
	}

	delOrgChannelCmd = &cobra.Command{
		Use:   `delOrgChannel`,
		Short: `delete organization to channel.`,
		Run:   delOrgChannel,
	}
)
