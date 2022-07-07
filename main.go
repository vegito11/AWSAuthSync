package main

import (
	"flag"
	"fmt"
	"time"

	utils "github.com/vegito11/AWSAuthSync/pkg/Utils"
	crInf "github.com/vegito11/AWSAuthSync/pkg/client/informers/externalversions"
	"github.com/vegito11/AWSAuthSync/pkg/controller"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "C:\\Users\\Vegito\\.kube\\config", "location to your kubeconfig file")
	crClientset, kubeClient, err := utils.GetKubeConn(kubeconfig)

	if err != nil {
		fmt.Printf("Error %s, creating crClientset \n ", err.Error())
	}

	ch := make(chan struct{})

	CrInfF := crInf.NewSharedInformerFactory(crClientset, 10*time.Minute)
	if err != nil {
		fmt.Printf("Getting informer factory %s \n", err.Error())
	}

	c := controller.NewController(kubeClient, crClientset, CrInfF.Vegito11().V1beta().AWSAuthMaps())
	CrInfF.Start(ch)
	c.Run(ch)
}
