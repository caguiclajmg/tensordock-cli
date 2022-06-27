package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	billingCmd = &cobra.Command{
		Use:   "billing",
		Short: "Manage billing",
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := client.GetBillingDetails()
			if err != nil {
				return err
			}

			if !res.Success {
				return errors.New(res.Error)
			}

			fmt.Printf(`Balance: %v
Hourly Spending Rate: %v`,
				res.Balance,
				res.HourlySpendingRate)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(billingCmd)
}
