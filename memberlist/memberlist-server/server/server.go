package server

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"errors"
	"github.com/hashicorp/memberlist"
	//"context"
	//"os"
	"net"
)


var list *memberlist.Memberlist
var conf *memberlist.Config

func StartMemberlist() {
	// starting up the memberlist
	msgCh := make(chan []byte)

	d := new(MyDelegate)
	d.msgCh = msgCh

	conf = memberlist.DefaultLocalConfig()
	conf.Name          = "node1"
	conf.BindPort      = 7947 // avoid port confliction
	conf.AdvertisePort = conf.BindPort
	// conf.BindAddr      = "3.13.172.219"
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


			if msg.Key == "clientToServerData" {
				// TODO: what to do with the client data? eg. msg.Value
				log.Printf("received msg: key=%s", msg.Key)
			}
		}
	}
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
func retryJoin(ipAddr net.IP, port uint) error {
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
func getNode(ipAddr net.IP, port uint) (*memberlist.Node, error) {
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
func sendToClient(m *MyMessage, ip string, port uint) error {
	// parsing and checking the ip.
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil { return errors.New("invalid IP address") }

	node, err := getNode(ipAddr, port)
	if err != nil {
		return err
	}

	if err := list.SendReliable(node, m.Bytes()); err != nil {
		return err
	}
	return nil
}

// SendClientID
func SendClientID(ID string, ip string, port uint) error {
	idBytes, err := json.Marshal(ID)
	if err != nil {
		return err
	}

	localNode := list.LocalNode()
	// make message
	m := new(MyMessage)
	m.FromAddr = localNode.Addr
	m.FromPort = localNode.Port
	m.Key = "serverToClientID"
	m.Value = idBytes


	if err := sendToClient(m, ip, port); err != nil {
		return err
	}
	return nil
}


// SendClientData
func SendClientData(data []byte, ip string, port uint) error {

	localNode := list.LocalNode()
	// make message
	m := new(MyMessage)
	m.FromAddr = localNode.Addr
	m.FromPort = localNode.Port
	m.Key = "serverToClientData"
	m.Value = data


	if err := sendToClient(m, ip, port); err != nil {
		return err
	}
	return nil
}