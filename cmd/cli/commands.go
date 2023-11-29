package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var message string

var boot = &cobra.Command{
	Use:   "boot",
	Short: "Boots the node with the defined flags. (Defaults to \"<Refresh: True><Port: 3003>\")",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Flags: %s\n", message)
	},
}

var refresh = &cobra.Command{
	Use:   "refresh",
	Short: "Updates node to latest version",
	Long:  "Updates the node to the latest approved version. (Updates to Helix after v1.0.0 are made purely through the evm, unless they are updates relating to efficiency or speed.)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Node Booted")
	},
}

func CommandInit() {
	rootCmd.AddCommand(refresh)

	rootCmd.AddCommand(boot)

	boot.Flags().StringVarP(&message, "flags", "f", "<FlagName: Parameter>...", "Custom message to print")
}
