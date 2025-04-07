/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/adpg24/devoops/kube"

	"github.com/spf13/cobra"
)

// currentContextCmd represents the currentContext command
var currentContextCmd = &cobra.Command{
	GroupID: "contextGroup",
	Use:     "current-context",
	Aliases: []string{"cc"},
	Short:   "Get current context",
	Long:    `Get the current context as defined in your kube config`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeConfig := kube.NewKubeConfig("")
		fmt.Println("Current context:", kubeConfig.GetCurrentContext())
	},
}

func init() {
	rootCmd.AddCommand(currentContextCmd)
}
