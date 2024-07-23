package nacos_viper_remote

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"io"
)

//var nacosOptions = &Option{}

func SetOptions(option *Option) {
	manager, _ := NewNacosConfigManager(option)
	viper.SupportedRemoteProviders = []string{"nacos"}
	viper.RemoteConfig = &remoteConfigProvider{ConfigManager: manager}
}

type remoteConfigProvider struct {
	ConfigManager *nacosConfigManager
}

func (rc *remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cmt, err := rc.getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	var b []byte
	switch cm := cmt.(type) {
	case viperConfigManager:
		b, err = cm.Get(rp.Path())
	}
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc *remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return rc.Get(rp)
}

func (rc *remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cmt, err := rc.getConfigManager(rp)
	if err != nil {
		return nil, nil
	}

	switch cm := cmt.(type) {
	case viperConfigManager:
		quit := make(chan bool)
		viperResponseCh := cm.Watch("dataId", quit)
		return viperResponseCh, quit
	}

	return nil, nil
}

func (rc *remoteConfigProvider) getConfigManager(rp viper.RemoteProvider) (interface{}, error) {
	if rp.Provider() == "nacos" {
		return rc.ConfigManager, nil
	} else {
		return nil, errors.New("The Nacos configuration manager is not supported!")
	}
}
