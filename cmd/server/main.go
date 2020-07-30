package main

import (
	"github.com/jonyhy96/kt-connect/pkg/apiserver/cluster"
	"github.com/jonyhy96/kt-connect/pkg/apiserver/common"
	"github.com/jonyhy96/kt-connect/pkg/apiserver/server"
	"github.com/jonyhy96/kt-connect/pkg/apiserver/util"
)

func main() {

	client, config, err := util.GetKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	watcher, err := cluster.Construct(client, config)
	if err != nil {
		panic(err.Error())
	}

	context := common.Context{
		Cluster: watcher,
	}

	err = server.Init(context)
	if err != nil {
		panic(err.Error())
	}

}
