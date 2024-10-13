package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	raftdv1 "github.com/amjadjibon/raftd/gen/raftd/v1"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the Raft cluster",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr = cmd.Flag("grpc-addr").Value.String()
		client, err := NewRaftServiceClient(grpcAddr)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		status, err := client.Status(cmd.Context(), &raftdv1.StatusRequest{})
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println()
		cmd.Printf("Leader: %s\n", status.Leader)
		cmd.Println("Peers:")
		for _, server := range status.Peers {
			cmd.Printf("  - ID: %s, Address: %s\n", server.Id, server.Address)
		}
	},
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a Raft node to the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr = cmd.Flag("grpc-addr").Value.String()
		client, err := NewRaftServiceClient(grpcAddr)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		joinAddr := cmd.Flag("join-addr").Value.String()
		joinId := cmd.Flag("join-id").Value.String()
		_, err = client.Join(cmd.Context(), &raftdv1.JoinRequest{
			Id:      joinId,
			Address: joinAddr,
		})
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("Node joined successfully")
	},
}

func init() {
	statusCmd.Flags().String("grpc-addr", "", "gRPC server address")
	_ = statusCmd.MarkFlagRequired("grpc-addr")

	joinCmd.Flags().String("grpc-addr", "", "gRPC server address")
	joinCmd.Flags().String("join-addr", "", "Address of the node to join")
	joinCmd.Flags().String("join-id", "", "ID of the node to join")
	_ = joinCmd.MarkFlagRequired("grpc-addr")
	_ = joinCmd.MarkFlagRequired("join-addr")
	_ = joinCmd.MarkFlagRequired("join-id")
}

func NewRaftServiceClient(addr string) (raftdv1.RaftServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	raftClient := raftdv1.NewRaftServiceClient(conn)

	return raftClient, err
}
