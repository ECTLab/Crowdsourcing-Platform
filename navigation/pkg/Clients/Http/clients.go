package Http

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
    "io"
    "net/http"
)

func Get(address string, client *http.Client) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithError(fmt.Errorf("%s", err)).WithField("address", address).Errorf("recovered from panic when calling httpGet: %s", err)
		}
	}()
	resp, err := client.Get(address)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)
    err = resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func Post(address string, data interface{}, client *http.Client) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithError(fmt.Errorf("%s", err)).WithField("address", address).WithField("data", data).Errorf("recovered from panic while calling httpPost: %s", err)
		}
	}()
	bytesData, err := json.Marshal(data)
	if err != nil {
		log.Warn("failed to json marshal the data for vroom")
		return nil, err
	}
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}