package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

// RegisterConsulParam stores the message to be sent to consul
type RegisterConsulParam struct {
	Key   string
	Value string
}

// AddDomainName add domain name to the redis
func AddDomainName(domainName string, client *redis.Client) error {
	return client.Set(domainName, "", time.Until(time.Now().Add(time.Minute*3600))).Err()
}

// RegisterToConsul register domain name to consul
func (r *RegisterConsulParam) RegisterToConsul(url string) error {
	client := &http.Client{}
	body := []byte(r.Value)
	req, err := http.NewRequest("PUT", url+r.Key, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	// TODO: check response to see whether it is successful
	defer res.Body.Close()
	return nil
}

// GetConsulNodeInformation return node uuid
func GetConsulNodeInformation(url string) (string, error) {
	type ret struct {
		Config struct {
			NodeID string
		}
	}
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	j := &ret{}
	err = json.Unmarshal(body, j)
	if err != nil {
		return "", err
	}
	if len(j.Config.NodeID) == 0 {
		return "", errors.New("can not get uuid from consul")
	}
	return j.Config.NodeID, nil
}

// IsUUIDExists judge if the uuid exists in redis
func IsUUIDExists(name string, client *redis.Client) (bool, error) {
	_, err := client.Get(name).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
