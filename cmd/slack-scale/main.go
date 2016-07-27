package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl"

	"github.com/golang/glog"
	"github.com/nlopes/slack"
)

var (
	inClusterConfig = flag.Bool("in-cluster-config", true, "")
	slackDebug      = flag.Bool("slack-debug", false, "")
	slackToken      = flag.String("slack-token", "xoxb-63745459303-os0yBg2OBC2eENoxXaN6lRsd", "")
)

func main() {
	flag.Parse()

	var cfg *restclient.Config
	var err error
	if *inClusterConfig {
		cfg, err = restclient.InClusterConfig()
	} else {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		cfg, err = kubeConfig.ClientConfig()
	}
	if err != nil {
		glog.Fatalf("unable to initialize kubeconfig: %v", err)
	}
	client, err := unversioned.New(cfg)
	if err != nil {
		glog.Fatalf("unable to initialize client: %v", err)
	}

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api := slack.New(*slackToken)
	api.SetDebug(*slackDebug)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	glog.Infof("starting loop")
	for m := range rtm.IncomingEvents {
		if m.Type != "message" {
			continue
		}
		message := m.Data.(*slack.MessageEvent)
		parts := strings.Split(message.Text, " ")
		if len(parts) != 3 || parts[0] != "scale" {
			continue
		}
		deployment := parts[1]
		replicas, err := strconv.Atoi(parts[2])
		if err != nil {
			glog.Errorf("atoi")
			continue
		}

		scaler, err := kubectl.ScalerFor(extensions.Kind("Deployment"), client)
		if err != nil {
			glog.Errorf("foo: %v", err)
			continue
		}

		glog.Infof("scaling %v %v", deployment, replicas)
		err = scaler.ScaleSimple("default", deployment, &kubectl.ScalePrecondition{Size: -1}, uint(replicas))

		if err != nil {
			msgStr := fmt.Sprintf("failed to scale deployment: %v", err)
			api.PostMessage(message.Channel, msgStr, slack.NewPostMessageParameters())
		} else {
			msgStr := fmt.Sprintf("scaled deployment %v to %v replicas", deployment, replicas)
			api.PostMessage(message.Channel, msgStr, slack.NewPostMessageParameters())
		}
		glog.Infof("scaled")
	}
}
