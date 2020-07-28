package cmd

import (
	"fmt"

	cfg "github.com/libonomy/libonomy-light/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config = cfg.DefaultConfig()
)

// AddCommands adds cobra commands to the app.
func AddCommands(cmd *cobra.Command) {

	/** ======================== BaseConfig Flags ========================== **/

	cmd.PersistentFlags().StringVarP(&config.BaseConfig.ConfigFile,
		"config", "c", config.BaseConfig.ConfigFile, "Set Load configuration from file")
	cmd.PersistentFlags().StringVarP(&config.BaseConfig.DataDir, "data-folder", "d",
		config.BaseConfig.DataDir, "Specify data directory for libonomy")
	cmd.PersistentFlags().BoolVar(&config.TestMode, "test-mode",
		config.TestMode, "Initialize testing features")

	/** ======================== P2P Flags ========================== **/

	cmd.PersistentFlags().IntVar(&config.P2P.TCPPort, "tcp-port",
		config.P2P.TCPPort, "TCP Port to listen on")
	cmd.PersistentFlags().BoolVar(&config.P2P.AcquirePort, "acquire-port",
		config.P2P.AcquirePort, "Should the node attempt to forward the port to this machine on a NAT?")
	cmd.PersistentFlags().DurationVar(&config.P2P.DialTimeout, "dial-timeout",
		config.P2P.DialTimeout, "Network dial timeout duration")
	cmd.PersistentFlags().DurationVar(&config.P2P.ConnKeepAlive, "conn-keepalive",
		config.P2P.ConnKeepAlive, "Network connection keep alive")
	cmd.PersistentFlags().Int8Var(&config.P2P.NetworkID, "network-id",
		config.P2P.NetworkID, "NetworkID to run on (0 - mainnet, 1 - testnet)")
	cmd.PersistentFlags().DurationVar(&config.P2P.ResponseTimeout, "response-timeout",
		config.P2P.ResponseTimeout, "Timeout for waiting on resposne message")
	cmd.PersistentFlags().DurationVar(&config.P2P.SessionTimeout, "session-timeout",
		config.P2P.SessionTimeout, "Timeout for waiting on session message")
	cmd.PersistentFlags().StringVar(&config.P2P.NodeID, "node-id",
		config.P2P.NodeID, "Load node data by id (pub key) from local store")
	cmd.PersistentFlags().IntVar(&config.P2P.BufferSize, "buffer-size",
		config.P2P.BufferSize, "Size of the messages handler's buffer")
	cmd.PersistentFlags().IntVar(&config.P2P.MaxPendingConnections, "max-pending-connections",
		config.P2P.MaxPendingConnections, "The maximum number of pending connections")
	cmd.PersistentFlags().IntVar(&config.P2P.OutboundPeersTarget, "outbound-target",
		config.P2P.OutboundPeersTarget, "The outbound peer target we're trying to connect")
	cmd.PersistentFlags().IntVar(&config.P2P.MaxInboundPeers, "max-inbound",
		config.P2P.MaxInboundPeers, "The maximum number of inbound peers ")
	cmd.PersistentFlags().BoolVar(&config.P2P.SwarmConfig.Gossip, "gossip",
		config.P2P.SwarmConfig.Gossip, "should we start a gossiping node?")
	cmd.PersistentFlags().BoolVar(&config.P2P.SwarmConfig.Bootstrap, "bootstrap",
		config.P2P.SwarmConfig.Bootstrap, "Bootstrap the swarm")
	cmd.PersistentFlags().IntVar(&config.P2P.SwarmConfig.RoutingTableBucketSize, "bucketsize",
		config.P2P.SwarmConfig.RoutingTableBucketSize, "The rounding table bucket size")
	cmd.PersistentFlags().IntVar(&config.P2P.SwarmConfig.RoutingTableAlpha, "alpha",
		config.P2P.SwarmConfig.RoutingTableAlpha, "The rounding table Alpha")
	cmd.PersistentFlags().IntVar(&config.P2P.SwarmConfig.RandomConnections, "randcon",
		config.P2P.SwarmConfig.RoutingTableAlpha, "Number of random connections")
	cmd.PersistentFlags().StringSliceVar(&config.P2P.SwarmConfig.BootstrapNodes, "bootnodes",
		config.P2P.SwarmConfig.BootstrapNodes, "Number of random connections")
	// cmd.PersistentFlags().DurationVar(&config.TIME.MaxAllowedDrift, "max-allowed-time-drift",
	// 	config.TIME.MaxAllowedDrift, "When to close the app until user resolves time sync problems")
	cmd.PersistentFlags().StringVar(&config.P2P.SwarmConfig.PeersFile, "peers-file",
		config.P2P.SwarmConfig.PeersFile, "addrbook peers file. located under data-dir/<publickey>/<peer-file> not loaded or saved if empty string is given.")
	// cmd.PersistentFlags().IntVar(&config.TIME.NtpQueries, "ntp-queries",
	// 	config.TIME.NtpQueries, "Number of ntp queries to do")
	// cmd.PersistentFlags().DurationVar(&config.TIME.DefaultTimeoutLatency, "default-timeout-latency",
	// 	config.TIME.DefaultTimeoutLatency, "Default timeout to ntp query")
	// cmd.PersistentFlags().DurationVar(&config.TIME.RefreshNtpInterval, "refresh-ntp-interval",
	// 	config.TIME.RefreshNtpInterval, "Refresh intervals to ntp")
	cmd.PersistentFlags().IntVar(&config.P2P.MsgSizeLimit, "msg-size-limit",
		config.P2P.MsgSizeLimit, "The message size limit in bytes for incoming messages")

	// // Bind Flags to config
	err := viper.BindPFlags(cmd.PersistentFlags())
	if err != nil {
		fmt.Println("an error has occurred while binding flags:", err)
	}

}
