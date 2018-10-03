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
	"github.com/mlab-lattice/lattice/pkg/latticectl/deploys"
	"github.com/mlab-lattice/lattice/pkg/util/cli"
	"github.com/mlab-lattice/lattice/pkg/util/cli/color"
	"github.com/mlab-lattice/lattice/pkg/util/cli/printer"

	"k8s.io/apimachinery/pkg/util/wait"
)

// Deploys returns a *cli.Command to list a system's deploys with subcommands to interact
// with individual deploys.
func Deploys() *cli.Command {
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
				return WatchDeploys(ctx.Client, ctx.System, format)
			}

			return PrintDeploys(ctx.Client, ctx.System, format, os.Stdout)
		},
		Subcommands: map[string]*cli.Command{
			"status": deploys.Status(),
		},
	}

	return cmd.Command()
}

// PrintBuilds prints the system's deploys to the supplied writer.
func PrintDeploys(client client.Interface, system v1.SystemID, format printer.Format, w io.Writer) error {
	deploys, err := client.V1().Systems().Deploys(system).List()
	if err != nil {
		return err
	}

	switch format {
	case printer.FormatTable:
		t := deploysTable(w)
		r := deploysTableRows(deploys)
		t.AppendRows(r)
		t.Print()

	case printer.FormatJSON:
		j := printer.NewJSON(w)
		j.Print(deploys)

	default:
		return fmt.Errorf("unexpected format %v", format)
	}

	return nil
}

// WatchDeploys watches the system's deploys, updating output based on changes.
// When passed in printer.Table as f, the table uses some ANSI escapes to overwrite some of the terminal buffer,
// so it always writes to stdout and does not accept an io.Writer.
func WatchDeploys(client client.Interface, system v1.SystemID, format printer.Format) error {
	// Poll the API for the systems and send it to the channel
	deploys := make(chan []v1.Deploy)

	go wait.PollImmediateInfinite(
		5*time.Second,
		func() (bool, error) {
			d, err := client.V1().Systems().Deploys(system).List()
			if err != nil {
				return false, err
			}

			deploys <- d
			return false, nil
		},
	)

	var handle func([]v1.Deploy)
	switch format {
	case printer.FormatTable:
		t := deploysTable(os.Stdout)
		handle = func(deploys []v1.Deploy) {
			r := deploysTableRows(deploys)
			t.Overwrite(r)
		}

	case printer.FormatJSON:
		j := printer.NewJSON(os.Stdout)
		handle = func(deploys []v1.Deploy) {
			j.Print(deploys)
		}

	default:
		return fmt.Errorf("unexpected format %v", format)
	}

	for d := range deploys {
		handle(d)
	}

	return nil
}

func deploysTable(w io.Writer) *printer.Table {
	return printer.NewTable(w, []string{"ID", "TARGET", "STATE", "STARTED", "COMPLETED"})
}

func deploysTableRows(deploys []v1.Deploy) [][]string {
	var rows [][]string
	for _, deploy := range deploys {
		stateColor := color.WarningString
		switch deploy.Status.State {
		case v1.DeployStateSucceeded:
			stateColor = color.SuccessString

		case v1.DeployStateFailed:
			stateColor = color.FailureString
		}

		target := "-"
		switch {
		case deploy.Build != nil:
			target = fmt.Sprintf("build %v", *deploy.Build)

		case deploy.Path != nil:
			target = fmt.Sprintf("path %v", deploy.Path.String())

		case deploy.Version != nil:
			target = fmt.Sprintf("version %v", *deploy.Version)
		}

		started := "-"
		if deploy.Status.StartTimestamp != nil {
			started = deploy.Status.StartTimestamp.Local().Format(time.RFC1123)
		}

		completed := "-"
		if deploy.Status.CompletionTimestamp != nil {
			completed = deploy.Status.CompletionTimestamp.Local().Format(time.RFC1123)
		}

		rows = append(rows, []string{
			color.IDString(string(deploy.ID)),
			target,
			stateColor(string(deploy.Status.State)),
			started,
			completed,
		})
	}

	// sort the rows by start timestamp
	startedIdx := 3
	sort.Slice(
		rows,
		func(i, j int) bool {
			ts1, ts2 := rows[i][startedIdx], rows[j][startedIdx]
			if ts1 == "-" {
				return true
			}

			if ts2 == "-" {
				return false
			}

			t1, err := time.Parse(time.RFC1123, ts1)
			if err != nil {
				panic(err)
			}

			t2, err := time.Parse(time.RFC1123, ts2)
			if err != nil {
				panic(err)
			}
			return t1.After(t2)
		},
	)

	return rows
}
