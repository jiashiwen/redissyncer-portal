package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"redissyncer-portal/core"
	"redissyncer-portal/global"
	"redissyncer-portal/httpserver"
	"redissyncer-portal/httpserver/router"
	"redissyncer-portal/inspection"
	"redissyncer-portal/node"
	"redissyncer-portal/resourceutils"
	"sync"
)

func NewStartCommand() *cobra.Command {
	sc := &cobra.Command{
		Use:   "start [-d ]",
		Short: "start server",
		Run:   startCommandFunc,
	}

	sc.Flags().BoolP("daemon", "d", false, "start as daemon")

	return sc
}

func startCommandFunc(cmd *cobra.Command, args []string) {
	cmd.Println("start ....")
	daemon, err := cmd.Flags().GetBool("daemon")
	if err != nil {
		panic(err)
	}
	cmd.Println(daemon)
	global.RSPViper.Set("daemon", daemon)
	startServer()
}

// 启动服务器
func startServer() {
	global.RSPLog = core.Zap()
	//MainStart(os.Args[1:])

	// -d 后台启动
	if global.RSPViper.GetBool("daemon") {
		cmd, err := background()
		if err != nil {
			fmt.Println(err)
		}

		//根据返回值区分父进程子进程
		if cmd != nil { //父进程
			fmt.Println("PPID: ", os.Getpid(), "; PID:", cmd.Process.Pid, "; Operating parameters: ", os.Args)
			return //父进程退出
		} else { //子进程
			fmt.Println("PID: ", os.Getpid(), "; Operating parameters: ", os.Args)
		}
	}

	global.RSPLog = core.Zap()

	wg := &sync.WaitGroup{}

	//start CursorGC
	resourceutils.NewEtcdCursorGC(10, 300000)
	wg.Add(1)
	go resourceutils.StartCursorGC(wg)

	//start node
	node := node.NewNode()
	if err := node.Registry(); err != nil {
		panic(err)
	}
	wg.Add(1)
	go node.Start(wg)

	//start inspect server
	inspector := inspection.NewInspector()
	wg.Add(1)
	go inspector.Start(wg)

	//启动http server
	//wg.Add(1)
	r := router.RootRouter()
	addr := "0.0.0.0:" + global.RSPViper.GetString("node.nodeport")
	httpserver.StartServer(addr, r)

}

func background() (*exec.Cmd, error) {
	envName := "DAEMON"    //环境变量名称
	envValue := "SUB_PROC" //环境变量值

	val := os.Getenv(envName) //读取环境变量的值,若未设置则为空字符串
	if val == envValue {      //监测到特殊标识, 判断为子进程,不再执行后续代码
		return nil, nil
	}

	/*以下是父进程执行的代码*/

	//因为要设置更多的属性, 这里不使用`exec.Command`方法, 直接初始化`exec.Cmd`结构体
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,      //注意,此处是包含程序名的
		Env:  os.Environ(), //父进程中的所有环境变量
	}

	//为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, envValue))

	//异步启动子进程
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
