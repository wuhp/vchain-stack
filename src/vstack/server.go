package main

import (
    "fmt"
    "bytes"
    "errors"
    "net/http"
    "encoding/json"
)

type Data struct {
    Reqs  []string `json:"requests"`
    Rlogs []string `json:"request-logs"`
}

func send(server, key string, data *Data) error {
    body, err := json.Marshal(data)
    if err != nil {
        return err
    }

    url := fmt.Sprintf("http://%s/api/v1/data", server)
    header := map[string]string { "Authorization": key }

    return do("POST", url, header, body)
}

func ping(server string) error {
    url := fmt.Sprintf("http://%s/api/v1/ping", server)
    return do("GET", url, nil, nil)
}

func do(method, url string, header map[string]string, body []byte) error {
    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    for k, v := range header {
        req.Header.Set(k, v)
    }

    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode != http.StatusOK {
        return errors.New(fmt.Sprintf("response %s", res.Status))
    }

    return nil
}
