package main

import (
    "log"
    "time"
    "os"
    "fmt"
    "strings"
    "io"
    "bufio"
    "io/ioutil"
)

const (
    MaxPerLines = 1000
)

var killed bool = false
var serverOk bool = false

func pingLoop(server string) {
    for {
        ok := ping(server) == nil
        if ok != serverOk {
            log.Printf("INFO: Ping server %s, from %v to %v\n", server, serverOk, ok)
            serverOk = ok
        }

        time.Sleep(1 * time.Second)
    }
}

func stop() {
    killed = true
}

func start(conf *Conf) {
    for {
        for _, src := range conf.SrcList {
            key := src.Key
            for _, logPath := range src.LogList {
                if killed {
                    return
                }
                processLog(conf, key, logPath)
            }
        }
        time.Sleep(time.Duration(conf.Interval) * time.Millisecond)
    }
}

func processLog(conf *Conf, key, file string) {
    if !serverOk {
        return
    }

    log.Printf("INFO: Begin to process log file `%s` ...\n", file)

    logFile, err := os.Open(file)
    if err != nil {
        log.Printf("ERROR: Can not open log `%s`\n", file)
        return
    }
    defer logFile.Close()

    historyFile := fmt.Sprintf("%s/%s", conf.LogDir, strings.Replace(file, "/", ".", -1))
    bytes, err := ioutil.ReadFile(historyFile)
    if err != nil {
        bytes = make([]byte, 0)
        log.Printf("INFO: No history file of log `%s`\n", file)
    }

    last := string(bytes)
    rb := bufio.NewReader(logFile)
    cnt := 0
    cur := ""
    reqs := make([]string, 0)
    rlogs := make([]string, 0)

    for {
        line, err := rb.ReadString('\n')
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("ERROR: Fail to read line of log `%s`\n", file)
            return
        }

        if line <= last {
            continue
        }

        if cnt >= MaxPerLines {
            break
        }

        items := strings.SplitN(line, " ", 4)
        switch items[2] {
        case "request":
            reqs = append(reqs, items[3])
        case "request-log":
            rlogs = append(rlogs, items[3])
        default:
            log.Printf("ERROR: Un-supported event `%s` in log `%s`, skipped\n", items[2], file)
            continue
        }

        cur = line
        cnt++
    }

    if cnt == 0 {
        log.Printf("INFO: No more new messages\n")
        return
    }

    data := new(Data)
    data.Reqs = reqs
    data.Rlogs = rlogs

    if err := send(conf.Dest, key, data); err != nil {
        log.Printf("ERROR: Fail to send data, with err %s\n", err.Error())
        return
    }

    log.Printf("INFO: Successfully to send %d messages\n", cnt)

    if err := ioutil.WriteFile(historyFile, []byte(cur), 0644); err != nil {
        log.Printf("ERROR: Fail to write to history file, with err %s\n", err.Error())
        return
    }
}
