package latticectl

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/mlab-lattice/lattice/pkg/api/client"
	"github.com/mlab-lattice/lattice/pkg/api/v1"
	"github.com/mlab-lattice/lattice/pkg/latticectl/command"
	"github.com/mlab-lattice/lattice/pkg/latticectl/services"
	"github.com/mlab-lattice/lattice/pkg/util/cli"
	"github.com/mlab-lattice/lattice/pkg/util/cli/color"
	"github.com/mlab-lattice/lattice/pkg/util/cli/printer"

	"k8s.io/apimachinery/pkg/util/wait"
)

func Services() *cli.Command {
	var (
		output string
		watch  bool
	)

	cmd := command.SystemCommand{
		Flags: map[string]cli.Flag{
			command.OutputFlagName: command.OutputFlag(
				&output,
				[]printer.Format{
					printer.FormatJSON,
					printer.FormatTable,
				},
				printer.FormatTable,
			),
			command.WatchFlagName: command.WatchFlag(&watch),
		},
		Run: func(ctx *command.SystemCommandContext, args []string, flags cli.Flags) error {
			format := printer.Format(output)

			if watch {
				WatchServices(ctx.Client, ctx.System, os.Stdout, format)
				return nil
			}

			return PrintServices(ctx.Client, ctx.System, os.Stdout, format)
		},
		Subcommands: map[string]*cli.Command{
			"logs":   services.Logs(),
			"status": services.Status(),
		},
	}

	return cmd.Command()
}

func PrintServices(client client.Interface, id v1.SystemID, w io.Writer, f printer.Format) error {
	services, err := client.V1().Systems().Services(id).List()
	if err != nil {
		return err
	}

	switch f {
	case printer.FormatTable:
		services, err := client.V1().Systems().Services(id).List()
		if err != nil {
			return err
		}

		t := servicesTable(w)
		r := servicesTableRows(services)
		t.AppendRows(r)
		t.Print()

	case printer.FormatJSON:
		j := printer.NewJSON(w)
		j.Print(services)

	default:
		return fmt.Errorf("unexpected format %v", f)
	}

	return nil
}

func WatchServices(client client.Interface, id v1.SystemID, w io.Writer, f printer.Format) {
	services := make(chan []v1.Service)

	// Poll the API for the builds and send it to the channel
	go wait.PollImmediateInfinite(
		5*time.Second,
		func() (bool, error) {
			s, err := client.V1().Systems().Services(id).List()
			if err != nil {
				// TODO: handle errors
				return false, nil
				//return false, err
			}

			services <- s
			return false, nil
		},
	)

	var handle func(services []v1.Service)
	switch f {
	case printer.FormatTable:
		t := servicesTable(w)
		handle = func(services []v1.Service) {
			r := servicesTableRows(services)
			t.Overwrite(r)
		}

	case printer.FormatJSON:
		j := printer.NewJSON(w)
		handle = func(services []v1.Service) {
			j.Print(services)
		}

	default:
		panic(fmt.Sprintf("unexpected format %v", f))
	}

	for s := range services {
		handle(s)
	}
}

func servicesTable(w io.Writer) *printer.Table {
	return printer.NewTable(w, []string{"PATH", "STATE", "AVAILABLE", "UPDATED", "STALE", "TERMINATING"})
}

func servicesTableRows(services []v1.Service) [][]string {
	var rows [][]string
	for _, service := range services {
		var stateColor color.Formatter
		switch service.Status.State {
		case v1.ServiceStateStable:
			stateColor = color.SuccessString
		case v1.ServiceStateFailed:
			stateColor = color.FailureString
		default:
			stateColor = color.WarningString
		}

		var addresses []string
		for port, address := range service.Status.Ports {
			addresses = append(addresses, fmt.Sprintf("%v: %v", port, address))
		}

		rows = append(rows, []string{
			color.IDString(service.Path.String()),
			stateColor(string(service.Status.State)),
			fmt.Sprintf("%d", service.Status.AvailableInstances),
			fmt.Sprintf("%d", service.Status.UpdatedInstances),
			fmt.Sprintf("%d", service.Status.StaleInstances),
			fmt.Sprintf("%d", service.Status.TerminatingInstances),
		})

	}

	// sort the rows by service ID
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })

	return rows
}
