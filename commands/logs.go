package commands

import (
	"bufio"
	"encoding/json"
	"kool-dev/kool/core/builder"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// KoolLogsFlags holds the flags for the logs command
type KoolLogsFlags struct {
	Tail   int
	Follow bool
}

type logEntryJSON struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

var execLogsCmd = exec.Command

// KoolLogs holds handlers and functions to implement the logs command logic
type KoolLogs struct {
	DefaultKoolService
	Flags *KoolLogsFlags

	list builder.Command
	logs builder.Command
}

func AddKoolLogs(root *cobra.Command) {
	var (
		logs    = NewKoolLogs()
		logsCmd = NewLogsCommand(logs)
	)

	root.AddCommand(logsCmd)
}

// NewKoolLogs creates a new handler for logs logic
func NewKoolLogs() *KoolLogs {
	return &KoolLogs{
		*newDefaultKoolService(),
		&KoolLogsFlags{25, false},
		builder.NewCommand("docker", "compose", "ps", "-aq"),
		builder.NewCommand("docker", "compose", "logs"),
	}
}

// Execute runs the logs logic with incoming arguments.
func (l *KoolLogs) Execute(args []string) (err error) {
	var services string

	if services, err = l.Shell().Exec(l.list, args...); err != nil {
		return
	}

	if services = strings.TrimSpace(services); services == "" {
		l.Shell().Warning("There are no containers")
		return
	}

	if l.Flags.Tail == 0 {
		l.logs.AppendArgs("--tail", "all")
	} else {
		l.logs.AppendArgs("--tail", strconv.Itoa(l.Flags.Tail))
	}

	if l.Flags.Follow {
		l.logs.AppendArgs("--follow")
	}

	if l.Shell().IsJSONOutput() {
		return l.printLogsJSON(args...)
	}

	err = l.Shell().Interactive(l.logs, args...)
	return
}

// NewLogsCommand initializes new kool logs command
func NewLogsCommand(logs *KoolLogs) (logsCmd *cobra.Command) {
	logsCmd = &cobra.Command{
		Use:   "logs [OPTIONS] [SERVICE...]",
		Short: "Display log output from running service containers",
		Long: `Display log output from all running service containers,
or one or more specified [SERVICE...] containers. Add a '-f' option to the
the command to follow the log output (i.e. 'kool logs -f [SERVICE...]').`,
		RunE: DefaultCommandRunFunction(logs),

		DisableFlagsInUseLine: true,
	}

	logsCmd.Flags().IntVarP(&logs.Flags.Tail, "tail", "t", 25, "Number of lines to show from the end of the logs for each container. A value equal to 0 will show all lines.")
	logsCmd.Flags().BoolVarP(&logs.Flags.Follow, "follow", "f", false, "Follow log output.")
	return
}

func (l *KoolLogs) printLogsJSON(args ...string) (err error) {
	if l.Flags.Follow {
		return l.streamLogsJSON(args...)
	}

	var output string
	if output, err = l.Shell().Exec(l.logs, args...); err != nil {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		entry := parseLogLine(line)
		var payload []byte
		if payload, err = json.Marshal(entry); err != nil {
			return
		}
		l.Shell().Println(string(payload))
	}
	return
}

func (l *KoolLogs) streamLogsJSON(args ...string) (err error) {
	cmdArgs := l.logs.Args()
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, args...)
	}
	cmd := execLogsCmd(l.logs.Cmd(), cmdArgs...)
	cmd.Env = os.Environ()
	cmd.Stderr = l.Shell().ErrStream()

	stdout, e := cmd.StdoutPipe()
	if e != nil {
		err = e
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entry := parseLogLine(line)
		if payload, e := json.Marshal(entry); e == nil {
			l.Shell().Println(string(payload))
		}
	}

	err = cmd.Wait()
	return
}

// parseLogLine parses a docker-compose log line into a logEntryJSON.
// Docker compose log format: "service_name | message" (with optional padding).
// If the line doesn't match, service is empty and message is the full line.
func parseLogLine(line string) logEntryJSON {
	if idx := strings.Index(line, "|"); idx >= 0 {
		service := strings.TrimSpace(line[:idx])
		message := strings.TrimSpace(line[idx+1:])
		return logEntryJSON{Service: service, Message: message}
	}
	return logEntryJSON{Service: "", Message: line}
}
