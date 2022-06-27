package commands

import (
	"errors"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	stockCmd = &cobra.Command{
		Use:   "stock",
		Short: "Query stock",
	}
	listStockCmd = &cobra.Command{
		Use:   "list",
		Short: "List stock",
		RunE:  listStock,
	}
)

func init() {
	stockCmd.AddCommand(listStockCmd)
	listStockCmd.Flags().String("type", "gpu", "Instance type (gpu or cpu)")
	listStockCmd.Flags().Bool("all", false, "Include out-of-stock instances")

	rootCmd.AddCommand(stockCmd)
}

func listStock(cmd *cobra.Command, args []string) error {
	instanceType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	switch instanceType {
	case "gpu":
		res, err := client.ListGpuStock()
		if err != nil {
			return nil
		}

		if !res.Success {
			return errors.New(res.Error)
		}

		t.AppendHeader(table.Row{"GPU", "Region", "Available Now", "Available Reserve"})

		for gpuModel, regionStock := range res.Stock {
			for region, stock := range regionStock {
				if (stock.AvailableNow > 0 || stock.AvailableReserve > 0) || all {
					t.AppendRow(table.Row{gpuModel, region, stock.AvailableNow, stock.AvailableReserve})
				}
			}
		}

	case "cpu":
		res, err := client.ListCpuStock()
		if err != nil {
			return nil
		}

		if !res.Success {
			return errors.New(res.Error)
		}

		t.AppendHeader(table.Row{"CPU Model", "Region", "Available Now"})

		for cpuModel, regionStock := range res.Stock {
			for region, stock := range regionStock {
				if stock.AvailableNow != "None" || all {
					t.AppendRow(table.Row{cpuModel, region, stock.AvailableNow})
				}
			}
		}

	default:
		return errors.New("unknown instance type")
	}

	t.Render()

	return nil
}
