package commands

import (
	"errors"
	"fmt"
	"log"
	"reflect"

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
