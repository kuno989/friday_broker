package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kuno989/friday_connect/connect/schema"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *Server) vmRequest(key, sha256, fileType, action, method string) {
	client := &http.Client{}
	client.Timeout = time.Second * 20
	res := schema.RequestMalwareScan{
		ObjectKey: key,
		Sha256:    sha256,
		FileType: fileType,
	}
	buff, err := json.Marshal(res)
	if err != nil {
		logrus.Errorf("Failed to json marshall object: %v ", err)
	}
	uri := s.Config.AgentURI + s.Config.AgentPort + "/api/" + action
	body := bytes.NewBuffer(buff)

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		logrus.Fatalf("failed to request %v", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatalf("failed to request %v", err)
	}
	defer resp.Body.Close()

	if action == "start"{
		b, err := ioutil.ReadAll(resp.Body)
		var resp = schema.ResponsePid{}
		err = json.Unmarshal(b, &resp)
		if err != nil {
			logrus.Fatalf("failed to read response %v", err)
		}
		pid := fmt.Sprintf("%v",resp.Pid)
		fmt.Printf("탐지된 프로세스:%s\npid : %s",resp.MalwareName, pid)
	}
}