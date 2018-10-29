package session

import (
	"fmt"
	"os/signal"

	"k8s.io/api/core/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"

	//"k8s.io/client-go/rest"
	"os"
	"path/filepath"

	"strconv"

	"syscall"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const NAME_SPACE string = "default"
const SERVER_AMOUNT int = 1
const MAX_SERVER_AMOUT int = 20
const MAX_PORT int = 31000
const MIN_PORT int = 30000

//MAX_ROOM_IN_POD maxium rooms in pod
const MAX_ROOM_IN_POD = 1

var interrupt chan os.Signal
var PodMod *v1.Pod

//ClusterManager clustermanager's instance
var ClusterManager *clusterManager

//ClusterManager manage kubernete's cluster and gameplayServer
type clusterManager struct {
	client    *kubernetes.Clientset
	pods      map[string]*v1.Pod
	PodAmount int
}

//KubeClientSet create clientset
func (c *clusterManager) KubeClientSet() {
	//config, err := rest.InClusterConfig()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))

	if err != nil {
		log.Warn(err)
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	c.client = clientset
	//log.Debug(config, clientset)

	//clientset.CoreV1().Services("NAME_SPACE")
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	p, _ := os.Getwd()
	f, err := os.Open(filepath.Join(p, "cluster/pod.yaml"))
	if err != nil {
		log.Fatal(ex, "not found:", err)
		return
	}

	yamlde := yaml.NewYAMLOrJSONDecoder(f, 300)
	podInfo := &v1.Pod{}
	if err = yamlde.Decode(podInfo); err != nil {
		log.Debug("yaml parse fail", err)
	}
	//log.Debug("pod info :", podInfo)
	PodMod = podInfo.DeepCopy()
	// id: ''
	// agentToGamePort: ''
	// clientToGamePort: ''
	// clientToAgentPort: ''
	if MIN_PORT+SERVER_AMOUNT*2 > MAX_PORT {
		panic("Usable port is not enough.")
	}
	for i := 0; i < SERVER_AMOUNT; i++ {
		c.CreatePod()
	}
}

//CreatePod scale gameplay server
func (c *clusterManager) CreatePod() {
	if c.PodAmount == MAX_SERVER_AMOUT {
		return
	}
	i := c.PodAmount
	podInfo := PodMod.DeepCopy()
	podInfo.Name = "gameplayserver-" + strconv.Itoa(i)
	podInfo.Labels["type"] = "game"
	podInfo.Labels["id"] = strconv.Itoa(i)
	podInfo.Labels["agentToGamePort"] = strconv.Itoa(MIN_PORT + i*2)
	podInfo.Labels["clientToGamePort"] = strconv.Itoa(MIN_PORT + i*2 + 1)
	_, err := c.client.CoreV1().Pods(NAME_SPACE).Create(podInfo.DeepCopy())
	if err != nil {
		log.Warn(err)
		if errors.IsAlreadyExists(err) {
			policy := metav1.DeletePropagationForeground
			deleteTime := int64(0)
			c.client.CoreV1().Pods(NAME_SPACE).Delete(podInfo.Name, &metav1.DeleteOptions{GracePeriodSeconds: &deleteTime, PropagationPolicy: &policy})

			_, err := c.client.CoreV1().Pods(NAME_SPACE).Create(podInfo.DeepCopy())
			if err != nil {
				log.Fatal("Create Fail", err)
			}
		}
	}
	c.pods[podInfo.Labels["id"]], err = c.client.CoreV1().Pods(NAME_SPACE).Get(podInfo.Name, metav1.GetOptions{})
	if err != nil {
		log.Warn(err)
	}

	nodeName := c.pods[podInfo.Labels["id"]].Spec.NodeName
	node, err := c.client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	extIP := ""
	for _, add := range node.Status.Addresses {
		if add.Type == v1.NodeExternalIP {
			extIP = add.Address
			break
		}
	}

	RoomManager.ConnectGameServer(extIP, podInfo.Labels["clientToGamePort"], podInfo.Labels["agentToGamePort"], podInfo.Labels["id"])
	log.Info("Gameplyer Server created: id: ", podInfo.Labels["id"], " IP: ", extIP, " Client Port: ", podInfo.Labels["clientToGamePort"], " Agent Port: ", podInfo.Labels["agentToGamePort"])
	c.PodAmount++
}

func cleanUp() {
	<-interrupt
	log.Info("Stopping agent")
	log.Info("Stopping gameplayer server")
	for _, pod := range ClusterManager.pods {
		err := ClusterManager.client.CoreV1().Pods(NAME_SPACE).Delete(pod.Name, &metav1.DeleteOptions{})
		if err != nil {
			log.Warn(err)
		}
	}
	os.Exit(1)
}

func init() {
	if os.Getenv("DONT_USE_KUBE") == "true" {
		return
	}
	ClusterManager = &clusterManager{
		pods: make(map[string]*v1.Pod),
	}
	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go cleanUp()
}
