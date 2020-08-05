package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/memberlist"
	//"context"
	//"os"
	"net"
)


var list *memberlist.Memberlist
var conf *memberlist.Config

func StartMemberlist() {
	// starting up the memberlist

	// fetch the public Dynamic ip of the client machine and use it in the Advertise address for the node.
	ip := fetchIP()

	msgCh := make(chan []byte)

	d := new(MyDelegate)
	d.msgCh = msgCh

	conf = memberlist.DefaultLocalConfig()
	conf.Name          = "node2"
	conf.BindPort      = 7948 // avoid port confliction
	conf.AdvertisePort = conf.BindPort
	//conf.BindAddr      = ""
	conf.AdvertiseAddr = ip
	conf.Delegate      = d
	conf.ProbeTimeout  = 10 * time.Second

	var err error
	list, err = memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	local := list.LocalNode()
	list.Join([]string{
		fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	})


	for {
		select {
		case data := <-d.msgCh:
			msg, ok := ParseMyMessage(data)
			if ok != true {
				continue
			}


			if msg.Key == "serverToClientID" {
				log.Printf("received msg: key=%s", msg.Key)

				var ID string
				_ = json.Unmarshal(msg.Value, ID)

				// get data to send
				// TODO: what is the data to send to the server?
				dataToSend := GetDataToSend()

				//
				// TODO: do you want to retry sending response multiple times?
				_ = RespondDataToServer(dataToSend, msg.FromAddr.To4(), msg.FromPort)

			} else if msg.Key == "serverToClientData" {
				SaveServerData(msg.Value)
			}
		}
	}
}


//
func fetchIP() string {
	url := "https://api.ipify.org?format=text"
	// we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ip := string(ipBytes)
	return ip
}


// getNodeByIP
func getNodeByIP(ipAddr net.IP) *memberlist.Node {

	members := list.Members()
	for _, node := range members {
		if node.Name == conf.Name {
			continue // skip self
		}
		if node.Addr.To4().Equal(ipAddr.To4()) {
			return node
		}
	}
	return nil
}

// func retryJoin
func retryJoin(ipAddr net.IP, port uint16) error {
	ipWithPort := ipAddr.To4().String() + ":" + string(port)

	var retryCount uint8
	var err error
	for retryCount <= 2 {
		if _, err = list.Join([]string{ipWithPort}); err != nil {
			retryCount++
			continue
		} else {
			return nil
		}
	}
	return err
}

// getNode
func getNode(ipAddr net.IP, port uint16) (*memberlist.Node, error) {
	// get node by ip from memberlist, if not found retry  3 times to join it.
	node := getNodeByIP(ipAddr)
	if node == nil {
		if err := retryJoin(ipAddr, port); err != nil {
			return nil, err
		} else {
			node = getNodeByIP(ipAddr)
		}
	}
	return node, nil
}

// SendToClient
func sendToClient(m *MyMessage, ip net.IP, port uint16) error {

	node, err := getNode(ip, port)
	if err != nil {
		return err
	}

	if err := list.SendReliable(node, m.Bytes()); err != nil {
		return err
	}
	return nil
}


func GetDataToSend() []byte {
	// TODO: what is the data to respond to the server?
	return []byte("hello?")
}

// SendClientID
func SaveServerData(data []byte) {
	// TODO: where or how to save the data?
	fmt.Println("saving received data from server", data)
}


// SendClientData
func RespondDataToServer(data []byte, ip net.IP, port uint16) error {

	localNode := list.LocalNode()
	// make message
	m := new(MyMessage)
	m.FromAddr = localNode.Addr
	m.FromPort = localNode.Port
	m.Key = "clientToServerData"
	m.Value = data


	if err := sendToClient(m, ip, port); err != nil {
		return err
	}
	return nil
}