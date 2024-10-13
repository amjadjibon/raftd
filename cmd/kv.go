package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	raftdv1 "github.com/amjadjibon/raftd/gen/raftd/v1"
)

func NewKvServiceClient(grpcAddr string) (raftdv1.KVServiceClient, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return raftdv1.NewKVServiceClient(conn), nil
}

var kvGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a value by key",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr = cmd.Flag("grpc-addr").Value.String()
		client, err := NewKvServiceClient(grpcAddr)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		key := cmd.Flag("key").Value.String()
		value, err := client.Get(cmd.Context(), &raftdv1.GetRequest{Key: key})
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println(string(value.Value))
	},
}

var kvSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value by key",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr = cmd.Flag("grpc-addr").Value.String()
		client, err := NewKvServiceClient(grpcAddr)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		key := cmd.Flag("key").Value.String()
		value := cmd.Flag("value").Value.String()
		_, err = client.Set(cmd.Context(), &raftdv1.SetRequest{Key: key, Value: []byte(value)})
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("Value set successfully")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a value by key",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr = cmd.Flag("grpc-addr").Value.String()
		client, err := NewKvServiceClient(grpcAddr)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		key := cmd.Flag("key").Value.String()
		_, err = client.Delete(cmd.Context(), &raftdv1.DeleteRequest{Key: key})
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		cmd.Println("Value deleted successfully")
	},
}

func init() {
	kvGetCmd.Flags().String("key", "", "Key to get")
	kvGetCmd.Flags().String("grpc-addr", ":8080", "gRPC server address")

	kvSetCmd.Flags().String("key", "", "Key to set")
	kvSetCmd.Flags().String("value", "", "Value to set")
	kvSetCmd.Flags().String("grpc-addr", ":8080", "gRPC server address")

	deleteCmd.Flags().String("key", "", "Key to delete")
	deleteCmd.Flags().String("grpc-addr", ":8080", "gRPC server address")

	_ = kvGetCmd.MarkFlagRequired("key")
	_ = kvSetCmd.MarkFlagRequired("key")
	_ = kvSetCmd.MarkFlagRequired("value")
	_ = deleteCmd.MarkFlagRequired("key")
	_ = kvGetCmd.MarkFlagRequired("grpc-addr")
	_ = kvSetCmd.MarkFlagRequired("grpc-addr")
	_ = deleteCmd.MarkFlagRequired("grpc-addr")

}
