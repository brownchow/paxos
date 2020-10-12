package paxos

import (
	"log"
	"time"
)

type network struct {
	recvQueue map[int]chan message
}

type nodeNetwork struct {
	id  int
	net *network
}

func CreateNetwork(nodes ...int) *network {
	nw := network{recvQueue: make(map[int]chan message, 0)}
	for _, node := range nodes {
		nw.recvQueue[node] = make(chan message, 1024)
	}
	return &nw
}

func (n *network) getNodeNetwork(id int) nodeNetwork {
	return nodeNetwork{id: id, net: n}
}

func (n *network) sendTo(m message) {
	log.Println("send msg from: ", m.from, " send to ", m.to, " val: ", m.val, " type: ", m.typ)
	n.recvQueue[m.to] <- m
}

func (n *network) recvFrom(id int) *message {
	select {
	case retMsg := <-n.recvQueue[id]:
		log.Println("Recv msg from: ", retMsg.from, " send to: ", retMsg.to, " val:", retMsg.val, " type:", retMsg.typ)
		return &retMsg
	case <-time.After(time.Second):
		// log.Println("id: ", id, " don't get message...time out")
		return nil
	}
}

func (n *nodeNetwork) send(m message) {
	n.net.sendTo(m)
}

func (n *nodeNetwork) recv() *message {
	return n.net.recvFrom(n.id)
}
