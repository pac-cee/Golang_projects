package cmd

import (
	"fmt"

	"cloud-cli/internal/cloud"
	"cloud-cli/pkg/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(vmListCmd)
	vmCmd.AddCommand(vmCreateCmd)
	vmCmd.AddCommand(vmDeleteCmd)
	vmCmd.AddCommand(vmStartCmd)
	vmCmd.AddCommand(vmStopCmd)

	// VM create flags
	vmCreateCmd.Flags().StringP("name", "n", "", "VM name")
	vmCreateCmd.Flags().StringP("type", "t", "", "VM type/size")
	vmCreateCmd.Flags().StringP("image", "i", "", "VM image")
	vmCreateCmd.MarkFlagRequired("name")
	vmCreateCmd.MarkFlagRequired("type")
	vmCreateCmd.MarkFlagRequired("image")
}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage virtual machines",
	Long:  `Create, list, start, stop, and delete virtual machines across different cloud providers.`,
}

var vmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		format, _ := cmd.Flags().GetString("output")

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		vms, err := client.ListVMs()
		if err != nil {
			fmt.Printf("Error listing VMs: %v\n", err)
			return
		}

		output.Print(vms, format)
	},
}

var vmCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new virtual machine",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		name, _ := cmd.Flags().GetString("name")
		vmType, _ := cmd.Flags().GetString("type")
		image, _ := cmd.Flags().GetString("image")

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		vm, err := client.CreateVM(cloud.VMConfig{
			Name:  name,
			Type:  vmType,
			Image: image,
		})
		if err != nil {
			fmt.Printf("Error creating VM: %v\n", err)
			return
		}

		fmt.Printf("Created VM: %s\n", vm.ID)
	},
}

var vmDeleteCmd = &cobra.Command{
	Use:   "delete [vm-id]",
	Short: "Delete a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		vmID := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.DeleteVM(vmID); err != nil {
			fmt.Printf("Error deleting VM: %v\n", err)
			return
		}

		fmt.Printf("Deleted VM: %s\n", vmID)
	},
}

var vmStartCmd = &cobra.Command{
	Use:   "start [vm-id]",
	Short: "Start a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		vmID := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.StartVM(vmID); err != nil {
			fmt.Printf("Error starting VM: %v\n", err)
			return
		}

		fmt.Printf("Started VM: %s\n", vmID)
	},
}

var vmStopCmd = &cobra.Command{
	Use:   "stop [vm-id]",
	Short: "Stop a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		vmID := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.StopVM(vmID); err != nil {
			fmt.Printf("Error stopping VM: %v\n", err)
			return
		}

		fmt.Printf("Stopped VM: %s\n", vmID)
	},
}
