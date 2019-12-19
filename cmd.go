package main

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:  `hlf-deploy`,
		Long: `Hyperledger Fabric Deploy`,
	}

	channelCmd = &cobra.Command{
		Use:   `channel`,
		Short: `Channel subcommand.`,
	}

	chaincodeCmd = &cobra.Command{
		Use:   `chaincode`,
		Short: `Chaincode subcommand.`,
	}

	organiztionCmd = &cobra.Command{
		Use:   `organization`,
		Short: `Organization subcommand.`,
	}

	createChannelCmd = &cobra.Command{
		Use:   `create`,
		Short: `Create channel.`,
		Run:   createChannel,
	}

	uptateAnchorPeerCmd = &cobra.Command{
		Use:   `updateAnchorPeer`,
		Short: `Update anchor peer.`,
		Run:   uptateAnchorPeer,
	}

	joinChannelCmd = &cobra.Command{
		Use:   `join`,
		Short: `Join channel.`,
		Run:   joinChannel,
	}

	installChaincodeCmd = &cobra.Command{
		Use:   `install`,
		Short: `Install chaincode.`,
		Run:   installChaincode,
	}

	instantiateChaincodeCmd = &cobra.Command{
		Use:   `instantiate`,
		Short: `Instantiate chaincode.`,
		Run:   instantiateAndUpgradeChaincode,
	}

	upgradeChaincodeCmd = &cobra.Command{
		Use:   `upgrade`,
		Short: `Upgrade chaincode.`,
		Run:   instantiateAndUpgradeChaincode,
	}

	queryChaincodeCmd = &cobra.Command{
		Use:   `query`,
		Short: `Query chaincode.`,
		Run:   queryAdnInvokeChaincode,
	}

	invokeChaincodeCmd = &cobra.Command{
		Use:   `invoke`,
		Short: `Invoke chaincode.`,
		Run:   queryAdnInvokeChaincode,
	}

	addOrgChannelCmd = &cobra.Command{
		Use:   `join`,
		Short: `add organization to channel.`,
		Run:   addAndUpdateOrgChannel,
	}

	updateOrgChannelCmd = &cobra.Command{
		Use:   `update`,
		Short: `update organization to channel.`,
		Run:   addAndUpdateOrgChannel,
	}

	delOrgChannelCmd = &cobra.Command{
		Use:   `delete`,
		Short: `delete organization to channel.`,
		Run:   delOrgChannel,
	}

	updateChannelParamCmd = &cobra.Command{
		Use:   `update`,
		Short: `update channel params.`,
		Run:   updateOrdererParam,
	}
)
