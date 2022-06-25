package commands

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	serversCmd = &cobra.Command{
		Use:   "servers",
		Short: "Manage servers",
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List servers",
		RunE:  serverList,
	}
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get server info",
		RunE:  serverInfo,
	}
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a server",
		RunE:  startServer,
	}
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop a server",
		RunE:  stopServer,
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a server",
		RunE:  deleteServer,
	}
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a server",
		RunE:  deployServer,
	}
	manageCmd = &cobra.Command{
		Use:   "manage",
		Short: "Open server management panel in a browser",
		RunE:  manageServer,
	}
)

func init() {
	serversCmd.AddCommand(listCmd)

	serversCmd.AddCommand(infoCmd)
	infoCmd.Flags().String("server", "", "Server Id")
	infoCmd.MarkFlagRequired("server")

	serversCmd.AddCommand(stopCmd)
	stopCmd.Flags().String("server", "", "Server Id")
	stopCmd.MarkFlagRequired("server")

	serversCmd.AddCommand(startCmd)
	startCmd.Flags().String("server", "", "Server Id")
	startCmd.MarkFlagRequired("server")

	serversCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("server", "", "Server Id")
	deleteCmd.MarkFlagRequired("server")

	serversCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("adminUser", "", "Your desired administrator username, e.g. \"user\" or \"tensordock_user\"")
	deployCmd.MarkFlagRequired("adminUser")
	deployCmd.Flags().String("adminPass", "", "Your desired administrator password. Please change it once you access your server")
	deployCmd.MarkFlagRequired("adminPass")
	deployCmd.Flags().String("gpuModel", "", "The GPU model that you would like to provision")
	deployCmd.MarkFlagRequired("gpuModel")
	deployCmd.Flags().String("location", "", "Location")
	deployCmd.MarkFlagRequired("location")
	deployCmd.Flags().String("name", "", "Name of your server in our dashboard")
	deployCmd.MarkFlagRequired("name")
	deployCmd.Flags().String("instanceType", "gpu", "Either \"gpu\" or \"cpu\"")
	deployCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	deployCmd.Flags().Int("vcpus", 1, "Number of vCPUs that you would like")
	deployCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	deployCmd.Flags().String("storageClass", "st1", "io1 or st1, depending on storage class desired")
	deployCmd.Flags().Int("ram", 2, "Number of GB of RAM to be deployed.")
	deployCmd.Flags().String("os", "Ubuntu 18.04 LTS", "Operating system")

	serversCmd.AddCommand(manageCmd)
	manageCmd.Flags().String("server", "", "Server Id")
	manageCmd.MarkFlagRequired("server")

	rootCmd.AddCommand(serversCmd)
}

func serverList(cmd *cobra.Command, args []string) error {
	res, err := client.ListServers()
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New("api call failed")
	}

	for key, elem := range res.Servers {
		fmt.Printf("%v (%v)\n", key, elem.Name)
	}

	return nil
}

func serverInfo(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New("api call failed")
	}

	if len(args) >= 2 {
		prop := args[1]
		val := reflect.ValueOf(res.Server)
		field := reflect.Indirect(val).FieldByName(prop)
		if field == (reflect.Value{}) {
			log.Fatalf("error: property %v not found", prop)
		}

		fmt.Printf("%v", field)
	} else {
		fmt.Printf(`Id: %v
Name: %v
Cost (Charged): %v 
Cost (Hour On): %v
Cost (Minutes On): %v
Cost (Hour Off): %v
Cost (Minutes Off): %v
CPU Model: %v
GPU Count: %v
GPU Model: %v
IP: %v
Location: %v
RAM: %v
Status: %v
Storage: %v
Storage Class: %v
Type: %v
vCPUs: %v`,
			res.Server.Id,
			res.Server.Name,
			res.Server.Cost.Charged,
			res.Server.Cost.HourOn,
			res.Server.Cost.MinutesOn,
			res.Server.Cost.HourOff,
			res.Server.Cost.MinutesOff,
			res.Server.CPUModel,
			res.Server.GPUCount,
			res.Server.GPUModel,
			res.Server.Ip,
			res.Server.Location,
			res.Server.Ram,
			res.Server.Status,
			res.Server.Storage,
			res.Server.StorageClass,
			res.Server.Type,
			res.Server.VCPUs)
	}

	return nil
}

func startServer(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	res, err := client.StartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func stopServer(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	res, err := client.StopServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func deleteServer(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	res, err := client.DeleteServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func deployServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	adminUser, err := flags.GetString("adminUser")
	if err != nil {
		return err
	}

	adminPass, err := flags.GetString("adminPass")
	if err != nil {
		return err
	}

	instanceType, err := flags.GetString("instanceType")
	if err != nil {
		return err
	}

	gpuModel, err := flags.GetString("gpuModel")
	if err != nil {
		return err
	}

	gpuCount, err := flags.GetInt("gpuCount")
	if err != nil {
		return err
	}

	vcpus, err := flags.GetInt("vcpus")
	if err != nil {
		return err
	}

	ram, err := flags.GetInt("ram")
	if err != nil {
		return err
	}

	storage, err := flags.GetInt("storage")
	if err != nil {
		return err
	}

	storageClass, err := flags.GetString("storageClass")
	if err != nil {
		return err
	}

	os, err := flags.GetString("os")
	if err != nil {
		return err
	}

	location, err := flags.GetString("location")
	if err != nil {
		return err
	}

	name, err := flags.GetString("name")
	if err != nil {
		return err
	}

	res, err := client.DeployServer(
		adminUser,
		adminPass,
		instanceType,
		gpuModel,
		gpuCount,
		vcpus,
		ram,
		storage,
		storageClass,
		os,
		location,
		name,
	)

	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	fmt.Println(res.Server.Id)

	return nil
}

func manageServer(cmd *cobra.Command, args []string) error {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	err = browser.OpenURL(res.Server.Links["dashboard"]["href"])
	if err != nil {
		return err
	}

	return nil
}
