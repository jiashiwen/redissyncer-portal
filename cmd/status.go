package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"redissyncer-portal/commons"
	"strconv"
)

func NewStatusCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "status",
		Short: "show server status",
		Run:   statusCommandFunc,
	}

	return sc
}

func statusCommandFunc(cmd *cobra.Command, args []string) {
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
	cmd.Println(dir + "/" + os.Args[0])
	cmd.Println("pidfilepath: " + dir + "/pid")
	cmd.Println("pid:" + strconv.Itoa(pidMap["pid"]))
	if commons.CheckPid(pidMap["pid"]) {
		cmd.Println("server is running")
		return
	}

	cmd.Println("server is stopped")
}
