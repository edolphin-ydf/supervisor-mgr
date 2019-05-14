package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/ochinchina/supervisord/xmlrpcclient"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

type options struct {
	//Server func() `short:"s" long:"servers"`
	Config func(string) `short:"c" long:"config" default:"config.yaml" description:"specify the config file"`
}

var o options
var writer *tabwriter.Writer

func init() {
	o.Config = func(s string) {
		loadConfig(s)
	}
}

func init() {
	const padding = 3
	writer = tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
}

var (
	errStartFail                     = errors.New("stop fail")
	errStopFail                      = errors.New("stop fail")
	errServerAndProcessNameAreNeeded = errors.New("server and process name are needed")
)

func getXmlRPCClient(sInfo *ServerInfo) *xmlrpcclient.XmlRPCClient {
	c := xmlrpcclient.NewXmlRPCClient(sInfo.Url, true)
	c.SetTimeout(5 * time.Second)
	c.SetUser(sInfo.UserName)
	c.SetPassword(sInfo.Password)

	return c
}

func printStatus(sInfo *ServerInfo) error {
	c := getXmlRPCClient(sInfo)
	pInfo, err := c.GetAllProcessInfo()
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		_, _ = fmt.Fprintf(writer, "Server:%v :%v\n", sInfo.Name, sInfo.Url)
		_, _ = fmt.Fprintf(writer, "Name\tState\tPid\tStartAt\n")
		for _, info := range pInfo.Value {
			_, _ = fmt.Fprintf(writer, "%v\t%v\t%v\t%v\n", info.Name, info.Statename, info.Pid, time.Unix(int64(info.Start), 0).Format("2006-01-02 15:04:05"))
		}
		_, _ = fmt.Fprintf(writer, "\n")
		_ = writer.Flush()
		return nil
	}
}

func start(sInfo *ServerInfo, processName string) error {
	c := getXmlRPCClient(sInfo)
	r, err := c.ChangeProcessState("start", processName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if r.Value {
			fmt.Println("start Ok!")
			return nil
		} else {
			fmt.Println("start Fail!")
			return errStartFail
		}
	}
}

func stop(sInfo *ServerInfo, processName string) error {
	c := getXmlRPCClient(sInfo)
	r, err := c.ChangeProcessState("stop", processName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if r.Value {
			fmt.Println("stop Ok!")
			return nil
		} else {
			fmt.Println("stop Fail!")
			return errStopFail
		}
	}
}

func executeOnServer(names []string, continueWhenError bool, f func(sInfo *ServerInfo) error) error {
	for _, name := range names {
		for _, info := range cfg.Servers {
			if info.Name != name {
				continue
			}

			if err := f(info); err != nil && !continueWhenError {
				return err
			}
		}
	}
	return nil
}

type cmdStatus struct {
}

func (*cmdStatus) Execute(args []string) error {
	if len(args) > 0 {
		return executeOnServer(args, true, func(sInfo *ServerInfo) error {
			return printStatus(sInfo)
		})
	} else {
		for _, info := range cfg.Servers {
			_ = printStatus(info)
		}
	}
	return nil
}

type cmdStart struct {
}

func (*cmdStart) Execute(args []string) error {
	if len(args) <= 1 {
		fmt.Println("server name and process name needed")
		return errServerAndProcessNameAreNeeded
	}
	return executeOnServer([]string{args[0]}, false, func(sInfo *ServerInfo) error {
		for _, processName := range args[1:] {
			if err := start(sInfo, processName); err != nil {
				return err
			}
		}
		return nil
	})
}

type cmdStop struct {
}

func (*cmdStop) Execute(args []string) error {
	if len(args) <= 1 {
		fmt.Println("server name and process name needed")
		return errServerAndProcessNameAreNeeded
	}
	return executeOnServer([]string{args[0]}, false, func(sInfo *ServerInfo) error {
		for _, processName := range args[1:] {
			if err := stop(sInfo, processName); err != nil {
				return err
			}
		}
		return nil
	})
}

func init() {
	_, _ = parser.AddCommand("status",
		"show process status: supervisor-mgr status [serverName...]",
		"",
		&cmdStatus{})
	_, _ = parser.AddCommand("start",
		"start process on server: supervisor-mgr start serverName processName[processName...]",
		"",
		&cmdStart{})
	_, _ = parser.AddCommand("stop",
		"stop process on server: supervisor-mgr stop serverName processName[processName...]",
		"",
		&cmdStop{})
}

var parser = flags.NewParser(&o, flags.Default & ^flags.PrintErrors)

type ServerInfo struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Servers []*ServerInfo `yaml:"servers"`
}

var cfg Config

func loadConfig(filePath string) {
	d, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(d, &cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//loadConfig("config.yaml")

	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}
}
