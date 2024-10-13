package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amjadjibon/raftd/server"
)

var (
	raftDir    string
	raftAddr   string
	raftNodeID string
	grpcAddr   string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Raft server",
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Run(
			cmd.Context(),
			raftDir,
			raftAddr,
			raftNodeID,
			grpcAddr,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().StringVar(&raftDir, "raft-dir", "/tmp/raft", "Raft data directory")
	startCmd.Flags().StringVar(&raftAddr, "raft-addr", "", "Raft bind address")
	startCmd.Flags().StringVar(&raftNodeID, "raft-node-id", "", "Raft node ID")
	startCmd.Flags().StringVar(&grpcAddr, "grpc-addr", ":8080", "gRPC server address")

	_ = startCmd.MarkFlagRequired("raft-addr")
	_ = startCmd.MarkFlagRequired("raft-node-id")
}
