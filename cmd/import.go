package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/thepoly/shuttletracker"
	"github.com/thepoly/shuttletracker/config"
	"github.com/thepoly/shuttletracker/postgres"
)

// Add is a flag to put the admins command into "add" mode.
// var Add bool

// Remove is a flag to put the admins command into "remove" mode.
// var Remove bool

func init() {
	// adminsCmd.Flags().BoolVar(&Add, "add", false, "add administrator")
	// adminsCmd.Flags().BoolVar(&Remove, "remove", false, "remove administrator")

	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import vehicles, routes, and stops from another Shuttle Tracker instance",
	// Long:  "List, add, or remove Shuttle Tracker administrators by RCS ID.",
	Args: func(cms *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expects exactly one argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.New()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to read configuration.")
			os.Exit(1)
		}

		pg, err := postgres.New(*cfg.Postgres)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to connect to Postgres:", err)
			os.Exit(1)
		}
		var ms shuttletracker.ModelService = pg

		instance := args[0]
		instanceURL, err := url.Parse(instance)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to parse instance URL:", err)
			os.Exit(1)
		}
		fmt.Printf("Importing from \"%s\"...\n", instance)

		c := http.Client{Timeout: time.Second * 10}

		// get stops
		fmt.Printf("Stops... ")
		stopsURL, _ := url.Parse(instanceURL.String())
		stopsURL.Path += "stops"

		resp, err := c.Get(stopsURL.String())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to get stops:", err)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to get stops: unexpected status code")
			os.Exit(1)
		}

		dec := json.NewDecoder(resp.Body)
		stops := []*shuttletracker.Stop{}
		err = dec.Decode(&stops)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to get stops:", err)
			os.Exit(1)
		}

		// map remote stop IDs to newly-inserted local stop IDs
		remoteStopsToLocal := map[int64]int64{}
		for _, stop := range stops {
			remoteID := stop.ID
			err = ms.CreateStop(stop)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Unable to create stop:", err)
				os.Exit(1)
			}
			remoteStopsToLocal[remoteID] = stop.ID
		}

		fmt.Printf("✅\n")

		// get routes
		fmt.Printf("Routes... ")
		routesURL, _ := url.Parse(instanceURL.String())
		routesURL.Path += "routes"

		resp, err = c.Get(routesURL.String())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to get routes:", err)
			os.Exit(1)
		}
		if resp.StatusCode != 200 {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to get routes: unexpected status code %d\n", resp.StatusCode)
			os.Exit(1)
		}

		dec = json.NewDecoder(resp.Body)
		routes := []*shuttletracker.Route{}
		err = dec.Decode(&routes)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Unable to get routes:", err)
			os.Exit(1)
		}

		for _, route := range routes {
			// change all stop IDs to reference the newly-created stops
			for i := range route.StopIDs {
				route.StopIDs[i] = remoteStopsToLocal[route.StopIDs[i]]
			}

			err = ms.CreateRoute(route)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Unable to create route:", err)
				os.Exit(1)
			}
		}

		fmt.Printf("✅\n")
	},
}
