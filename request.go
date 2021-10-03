package connect

import (
	"bytes"
	"encoding/json"
	"github.com/kuno989/friday_connect/connect/schema"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s *Server) vmRequest(key, sha256, fileType, action string) {
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
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		logrus.Fatalf("failed to request %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatalf("failed to request %v", err)
	}
	defer resp.Body.Close()
}