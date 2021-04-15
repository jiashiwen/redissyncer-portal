package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func NewStopCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "stop ",
		Short: "stop server",
		Run:   stopCommandFunc,
	}

	return sc
}

func stopCommandFunc(cmd *cobra.Command, args []string) {
	pidMap := make(map[string]int)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	content, err := ioutil.ReadFile(dir + "/pid")
	if err != nil {
		cmd.PrintErr("error: ", err)
		return
	}

	if err := yaml.Unmarshal(content, pidMap); err != nil {
		cmd.PrintErr("error: ", err)
		return
	}

	cmd.Println("pid:" + strconv.Itoa(pidMap["pid"]))
	cmd.Println("stopping ...")
	if err := syscall.Kill(pidMap["pid"], syscall.SIGTERM); err != nil {
		cmd.PrintErr("error: ", err)
	}

	cmd.Println("server stopped successes")
}
