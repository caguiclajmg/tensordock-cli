package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/caguiclajmg/tensordock-cli/api"
	"github.com/jedib0t/go-pretty/v6/table"
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
		Use:   "info [flags] server_id",
		Short: "Get server info",
		Args:  cobra.ExactArgs(1),
		RunE:  serverInfo,
	}
	startCmd = &cobra.Command{
		Use:     "start [flags] server_id",
		Short:   "Start a server",
		Args:    cobra.ExactArgs(1),
		RunE:    startServer,
		PostRun: logAction("success"),
	}
	stopCmd = &cobra.Command{
		Use:     "stop [flags] server_id",
		Short:   "Stop a server",
		Args:    cobra.ExactArgs(1),
		RunE:    stopServer,
		PostRun: logAction("success"),
	}
	deleteCmd = &cobra.Command{
		Use:     "delete [flags] server_id",
		Short:   "Delete a server",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteServer,
		PostRun: logAction("success"),
	}
	deployCmd = &cobra.Command{
		Use:     "deploy [flags] name admin_user admin_pass",
		Short:   "Deploy a server",
		Args:    cobra.ExactArgs(3),
		RunE:    deployServer,
		PostRun: logAction("success"),
	}
	manageCmd = &cobra.Command{
		Use:   "manage server_id",
		Short: "Open server management panel in a browser",
		Args:  cobra.ExactArgs(1),
		RunE:  manageServer,
	}
	sshCmd = &cobra.Command{
		Use:   "ssh server_id",
		Short: "Launch an SSH sesion with a server",
		Args:  cobra.ExactArgs(1),
		RunE:  sshServer,
	}
	restartCmd = &cobra.Command{
		Use:     "restart [flags] server_id",
		Short:   "Restart a server",
		Args:    cobra.ExactArgs(1),
		RunE:    restartServer,
		PostRun: logAction("success"),
	}
	modifyCmd = &cobra.Command{
		Use:     "modify [flags] server_id",
		Short:   "Modify a server",
		Args:    cobra.ExactArgs(1),
		RunE:    modifyServer,
		PostRun: logAction("success"),
	}
	statusCmd = &cobra.Command{
		Use:   "status server_id",
		Short: "Get server status",
		Args:  cobra.ExactArgs(1),
		RunE:  serverStatus,
	}
)

func init() {
	serversCmd.AddCommand(listCmd)

	serversCmd.AddCommand(infoCmd)

	serversCmd.AddCommand(stopCmd)

	serversCmd.AddCommand(startCmd)

	serversCmd.AddCommand(deleteCmd)

	serversCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("gpuModel", "Quadro_4000", "The GPU model that you would like to provision")
	deployCmd.Flags().String("location", "na-us-chi-1", "Location")
	deployCmd.Flags().String("instanceType", "gpu", "Either \"gpu\" or \"cpu\"")
	deployCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	deployCmd.Flags().String("cpuModel", "Intel_Xeon_v4", "The CPU model that you would like to provision")
	deployCmd.Flags().Int("vcpus", 2, "Number of vCPUs that you would like")
	deployCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	deployCmd.Flags().String("storageClass", "io1", "io1 or st1, depending on storage class desired")
	deployCmd.Flags().Int("ram", 4, "Number of GB of RAM to be deployed.")
	deployCmd.Flags().String("os", "Ubuntu 20.04 LTS", "Operating system")

	serversCmd.AddCommand(manageCmd)

	serversCmd.AddCommand(sshCmd)
	sshCmd.Flags().String("bin", "ssh", "Name of SSH client executable (e.g. ssh, mosh)")
	sshCmd.Flags().String("user", "user", "User account to use for login")
	sshCmd.Flags().String("extraFlags", "", "Extra flags to pass to the SSH client")

	serversCmd.AddCommand(restartCmd)

	serversCmd.AddCommand(modifyCmd)
	modifyCmd.Flags().String("instanceType", "gpu", "Either \"gpu\" or \"cpu\"")
	modifyCmd.Flags().String("gpuModel", "Quadro_4000", "The GPU model that you would like to provision")
	modifyCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	modifyCmd.Flags().String("cpuModel", "Intel_Xeon_v4", "The CPU model that you would like to provision")
	modifyCmd.Flags().Int("vcpus", 2, "Number of vCPUs that you would like")
	modifyCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	modifyCmd.Flags().Int("ram", 4, "Number of GB of RAM to be deployed.")

	serversCmd.AddCommand(statusCmd)

	rootCmd.AddCommand(serversCmd)
}

func serverList(cmd *cobra.Command, args []string) error {
	res, err := client.ListServers()
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Name", "Location", "Status"})
	for _, elem := range res.Servers {
		t.AppendRow(table.Row{elem.Id, elem.Name, elem.Location, elem.Status})
	}
	t.Render()

	return nil
}

func serverInfo(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	props := []map[string]string{
		{"name": "ID", "value": res.Server.Id},
		{"name": "Name", "value": res.Server.Name},
		{"name": "Location", "value": res.Server.Location},
		{"name": "IP", "value": res.Server.Ip},
		{"name": "Charged Cost", "value": fmt.Sprintf("%v", res.Server.Cost.Charged)},
		{"name": "Hour-On Cost", "value": fmt.Sprintf("%v", res.Server.Cost.HourOn)},
		{"name": "Hour-Off Cost", "value": fmt.Sprintf("%v", res.Server.Cost.HourOff)},
		{"name": "Minutes-On", "value": fmt.Sprintf("%v", res.Server.Cost.MinutesOn)},
		{"name": "Minutes-Off", "value": fmt.Sprintf("%v", res.Server.Cost.MinutesOff)},
		{"name": "CPU Model", "value": res.Server.CPUModel},
		{"name": "GPU Count", "value": strconv.Itoa(res.Server.GPUCount)},
		{"name": "GPU Model", "value": res.Server.GPUModel},
		{"name": "RAM", "value": fmt.Sprintf("%vGB", res.Server.Ram)},
		{"name": "Status", "value": res.Server.Status},
		{"name": "Storage", "value": fmt.Sprintf("%vGB", res.Server.Storage)},
		{"name": "Storage Class", "value": res.Server.StorageClass},
		{"name": "Type", "value": res.Server.Type},
		{"name": "vCPUs", "value": strconv.Itoa(res.Server.VCPUs)},
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Property", "Value"})
	for _, elem := range props {
		t.AppendRow(table.Row{elem["name"], elem["value"]})
	}
	t.Render()

	return nil
}

func startServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func stopServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StopServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func deleteServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.DeleteServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func deployServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

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

	cpuModel, err := flags.GetString("cpuModel")
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

	name := args[0]
	adminUser := args[1]
	adminPass := args[2]

	req := api.DeployServerRequest{
		AdminUser:    adminUser,
		AdminPass:    adminPass,
		InstanceType: instanceType,
		VCPUs:        vcpus,
		RAM:          ram,
		Storage:      storage,
		StorageClass: storageClass,
		OS:           os,
		Location:     location,
		Name:         name,
	}

	switch instanceType {
	case "cpu":
		req.CPUModel = cpuModel
	case "gpu":
		req.GPUModel = gpuModel
		req.GPUCount = gpuCount
	default:
		return errors.New("unknown instance type")
	}

	res, err := client.DeployServer(req)

	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	fmt.Println(res.Server.Id)

	return nil
}

func manageServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	err = browser.OpenURL(res.Server.Links["dashboard"]["href"])
	if err != nil {
		return err
	}

	return nil
}

func sshServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	bin, err := flags.GetString("bin")
	if err != nil {
		return err
	}

	user, err := flags.GetString("user")
	if err != nil {
		return err
	}

	extraFlags, err := flags.GetString("extraFlags")
	if err != nil {
		return err
	}

	sshCmd := exec.Command(bin, fmt.Sprintf("%v@%v", user, res.Server.Ip), extraFlags)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return err
	}

	return nil
}

func logAction(message string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) { log.Println(message) }
}

func restartServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.RestartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func modifyServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	serverId := args[0]

	var instanceType *string = nil
	if flags.Changed("instanceType") {
		instanceTypeVal, err := flags.GetString("instanceType")
		if err != nil {
			return err
		}
		instanceType = &instanceTypeVal
	}

	var gpuModel *string = nil
	if flags.Changed("gpuModel") {
		gpuModelVal, err := flags.GetString("gpuModel")
		if err != nil {
			return err
		}
		gpuModel = &gpuModelVal
	}

	var gpuCount *int = nil
	if flags.Changed("gpuCount") {
		gpuCountVal, err := flags.GetInt("gpuCount")
		if err != nil {
			return err
		}
		gpuCount = &gpuCountVal
	}

	var cpuModel *string = nil
	if flags.Changed("cpuModel") {
		cpuModelVal, err := flags.GetString("cpuModel")
		if err != nil {
			return err
		}
		cpuModel = &cpuModelVal
	}

	var vcpus *int = nil
	if flags.Changed("vcpus") {
		vcpusVal, err := flags.GetInt("vcpus")
		if err != nil {
			return err
		}
		vcpus = &vcpusVal
	}

	var ram *int = nil
	if flags.Changed("ram") {
		ramVal, err := flags.GetInt("ram")
		if err != nil {
			return err
		}
		ram = &ramVal
	}

	var storage *int = nil
	if flags.Changed("storage") {
		storageVal, err := flags.GetInt("storage")
		if err != nil {
			return err
		}
		storage = &storageVal
	}

	// based on tests, it seems that the endpoint does not
	// support specifying only parts of the spec to modify
	// (e.g. adjust VCPUs only) and that you need to specify
	// the entirety of the server configuration on every call
	req := api.ModifyServerRequest{
		ServerId:     serverId,
		InstanceType: instanceType,
		VCPUs:        vcpus,
		RAM:          ram,
		Storage:      storage,
	}

	switch *instanceType {
	case "cpu":
		req.CPUModel = cpuModel
	case "gpu":
		req.GPUModel = gpuModel
		req.GPUCount = gpuCount
	default:
		return errors.New("unknown instance type")
	}

	res, err := client.ModifyServer(req)

	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	return nil
}

func serverStatus(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServerStatus(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New(res.Error)
	}

	fmt.Println(res.Status)

	return nil
}
