package main

import (
    "os"
    "flag"
    "fmt"
    "log"
    "syscall"
    "os/signal"
    "io/ioutil"

    "gopkg.in/yaml.v2"
    "github.com/VividCortex/godaemon"
)

type Source struct {
    Key     string   `yaml:"key"`
    LogList []string `yaml:"log"`
}

type Conf struct {
    SrcList  []Source `yaml:"source"`
    Dest     string   `yaml:"dest"`
    LogDir   string   `yaml:"log_dir"`
    Interval int      `yaml:"interval"`
}

const (
    DefaultConfig   = "/etc/vstack/vstack.yml"
    DefaultDest     = "api.vchain.com"
    DefaultLogDir   = "/var/log/vstack"
    DefaultInterval = 2000
)

func exist(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

func parseConf() *Conf {
    path := flag.String("c", DefaultConfig, "config file")
    flag.Parse()

    bytes, err := ioutil.ReadFile(*path)
    if err != nil {
        log.Fatalf("ERROR: Failed to read config file `%s`, with err `%s`\n", *path, err.Error())
    }

    conf := new(Conf)
    if err := yaml.Unmarshal(bytes, conf); err != nil {
        log.Fatalf("ERROR: Failed to parse config file `%s`, with err `%s`\n", *path, err.Error())
    }

    if conf.SrcList == nil {
        log.Fatalf("ERROR: Tag `src` can not be empty\n")
    }

    if len(conf.Dest) == 0 {
        conf.Dest = DefaultDest
    }

    if len(conf.LogDir) == 0 {
        conf.LogDir = DefaultLogDir
    }

    if conf.Interval == 0 {
        conf.Interval = DefaultInterval
    }

    if !exist(conf.LogDir) {
        log.Fatalf("ERROR: No such file or directory `%s`\n", conf.LogDir)
    }

    return conf
}

func initLog(path string) {
    file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        os.Exit(1)
    }

    log.SetOutput(file)
}

func registerSignal() {
    c := make(chan os.Signal, 1)
    signal.Notify(c)
    for {
        s := <-c
        log.Printf("INFO: Receive signal %v\n", s)
        if s == syscall.SIGTERM {
            log.Printf("INFO: Stopping process loop ...\n")
            stop()
            break
        }
    }
}

func main() {
    conf := parseConf()
    godaemon.MakeDaemon(&godaemon.DaemonAttr{})
    initLog(fmt.Sprintf("%s/vstack.log", conf.LogDir))
    go registerSignal()
    go pingLoop(conf.Dest)

    log.Printf("INFO: Starting process loop ...")
    start(conf)
}
