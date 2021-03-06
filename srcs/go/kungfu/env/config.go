package env

import (
	"fmt"
	"os"

	kb "github.com/lsds/KungFu/srcs/go/kungfu/base"
	"github.com/lsds/KungFu/srcs/go/plan"
)

type Config struct {
	ConfigServer string
	Parent       plan.PeerID
	Parents      plan.PeerList
	Self         plan.PeerID
	Strategy     kb.Strategy

	InitClusterVersion string
	InitPeers          plan.PeerList

	Single bool
}

func ParseConfigFromEnv() (*Config, error) {
	if _, ok := os.LookupEnv(SelfSpecEnvKey); !ok {
		return singleEnv(), nil
	}
	self, err := getSelfFromEnv()
	if err != nil {
		return nil, err
	}
	parent, err := getParentFromEnv()
	if err != nil {
		return nil, err
	}
	hostList, err := getHostListFromEnv()
	if err != nil {
		return nil, err
	}
	initPeers, err := getInitPeersFromEnv()
	if err != nil {
		return nil, err
	}
	strategy, err := kb.ParseStrategy(os.Getenv(AllReduceStrategyEnvKey))
	if err != nil {
		return nil, err
	}
	return &Config{
		ConfigServer:       getConfigServerFromEnv(),
		Self:               *self,
		Parent:             *parent,
		Parents:            getParentIDs(hostList, *parent),
		InitPeers:          initPeers,
		Strategy:           *strategy,
		InitClusterVersion: os.Getenv(InitClusterVersionEnvKey),
	}, nil
}

func getParentIDs(hl plan.HostList, parent plan.PeerID) plan.PeerList {
	var ps plan.PeerList
	for _, h := range hl {
		ps = append(ps, plan.PeerID{IPv4: h.IPv4, Port: parent.Port})
	}
	return ps
}

func singleEnv() *Config {
	pl, _ := plan.DefaultHostList.GenPeerList(1, plan.DefaultPortRange)
	self := pl[0]
	return &Config{
		Self:      self,
		InitPeers: plan.PeerList{self},
		Strategy:  kb.DefaultStrategy,
		Single:    true,
	}
}

func getConfigServerFromEnv() string {
	return os.Getenv(ConfigServerEnvKey)
}

func getSelfFromEnv() (*plan.PeerID, error) {
	config, ok := os.LookupEnv(SelfSpecEnvKey)
	if !ok {
		return nil, fmt.Errorf("%s not set", SelfSpecEnvKey)
	}
	return plan.ParsePeerID(config)
}

func getParentFromEnv() (*plan.PeerID, error) {
	val, ok := os.LookupEnv(ParentIDEnvKey)
	if !ok {
		return nil, fmt.Errorf("%s not set", ParentIDEnvKey)
	}
	return plan.ParsePeerID(val)
}

func getInitPeersFromEnv() (plan.PeerList, error) {
	val, ok := os.LookupEnv(PeerListEnvKey)
	if !ok {
		return nil, fmt.Errorf("%s not set", PeerListEnvKey)
	}
	return plan.ParsePeerList(val)
}

func getHostListFromEnv() (plan.HostList, error) {
	val, ok := os.LookupEnv(HostListEnvKey)
	if !ok {
		return nil, fmt.Errorf("%s not set", HostListEnvKey)
	}
	return plan.ParseHostList(val)
}
