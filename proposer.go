package paxos

import "log"

type proposer struct {
	id         int
	seq        int
	proposeNum int
	proposeVal string
	acceptors  map[int]message
	nt         nodeNetwork
}

func NewProposer(id int, val string, nt nodeNetwork, acceptors ...int) *proposer {
	pro := proposer{id: id, proposeVal: val, seq: 0, nt: nt}
	pro.acceptors = make(map[int]message, len(acceptors))
	log.Println("proposer has ", len(acceptors), " acceptors, val: ", pro.proposeVal)
	for _, acceptor := range acceptors {
		pro.acceptors[acceptor] = message{}
	}
	return &pro
}

// detail process for Proposor
func (p *proposer) run() {
	log.Println("Proposer start run... val:", p.proposeVal)
	//stage1: Propopsor send prepare mesage to acceptor to reach accept from majority
	for !p.majorityReached() {
		log.Println("[Proposor:Prepare]")
		outMsgs := p.prepare()
		log.Println("[Proposer: prepare ", len(outMsgs), "msg")
		for _, msg := range outMsgs {
			p.nt.send(msg)
			log.Println("[Proposer: send ", msg)
		}
		log.Println("[Proposor prepare recv...")
		m := p.nt.recv()
		if m == nil {
			log.Println("[Proposor: no msg")
			continue
		}
		log.Println("[Proposor recv:", m)
		switch m.typ {
		case Promise:
			log.Println(" Proposer recv a promise from ", m.from)
			p.checkRecvPromise(*m)
		default:
			panic("Unsupport message.")
		}
	}
	log.Println("[Proposor: Propose")

	// stage2: proposor send promose value to get acceptor to learn
	log.Println("Proposor psopose seq:", p.getProposeNum(), " value:", p.proposeVal)
	proposeMsgs := p.propose()
	for _, msg := range proposeMsgs {
		p.nt.send(msg)
	}
}

// stage 1:
// prepare will prepare message to send to majority of acceptors
// according to spec, we only send our prepare msg to the "majority" not all acceptors
func (p *proposer) prepare() []message {
	p.seq++

	sendMsgCount := 0
	var msgList []message
	log.Println("proposer: prepare major msg:", len(p.acceptors))
	for acepId := range p.acceptors {
		msg := message{from: p.id, to: acepId, typ: Prepare, seq: p.getProposeNum(), val: p.proposeVal}
		msgList = append(msgList, msg)
		sendMsgCount++
		if sendMsgCount > p.majority() {
			break
		}
	}
	return msgList
}

// after receipt the promise from acceptor and reach majority
// Proposor will propose value to thos acceptors and let them know the consusence already ready
func (p *proposer) propose() []message {
	sendMsgCount := 0
	var msgList []message
	log.Println("proposer: propose msg: ", len(p.acceptors))
	for acepId, acepMsg := range p.acceptors {
		log.Println("check promise id: ", acepMsg.getProposeSeq(), p.getProposeNum())
		if acepMsg.getProposeSeq() == p.getProposeNum() {
			msg := message{from: p.id, to: acepId, typ: Propose, seq: p.getProposeNum()}
			msg.val = p.proposeVal
			log.Println("Propose val:", msg.val)
			msgList = append(msgList, msg)
		}
		sendMsgCount++
		if sendMsgCount > p.majority() {
			break
		}
	}
	log.Println("proposer propose msg list: ", msgList)
	return msgList
}

// checkRecvPromise 检查收到的promise消息
func (p *proposer) checkRecvPromise(promise message) {
	previousPromise := p.acceptors[promise.from]
	log.Println("prevMsg: ", previousPromise, " promiseMsg", promise)
	log.Println(previousPromise.getProposeSeq(), promise.getProposeSeq())
	if previousPromise.getProposeSeq() < promise.getProposeSeq() {
		log.Println("Proposor: ", p.id, " get new promise: ", promise)
		p.acceptors[promise.from] = promise
		if promise.getProposeSeq() > p.getProposeNum() {
			p.proposeNum = promise.getProposeSeq()
			p.proposeVal = promise.getProposeVal()
		}
	}
}

// majorityReached 提案者的提案是否获得半数通过
func (p *proposer) majorityReached() bool {
	return p.getRecvPromiseCount() > p.majority()
}

// majority 获取当前系统中半数的 acceptor
func (p *proposer) majority() int {
	return len(p.acceptors)/2 + 1
}

// getRecvPromiseCount 提案者获取的允诺数
func (p *proposer) getRecvPromiseCount() int {
	recvCount := 0
	for _, acepMsg := range p.acceptors {
		log.Println(" proposor has total ", len(p.acceptors), " acceptors ", acepMsg, " current Num: ", p.getProposeNum(), " msgNum: ", acepMsg.getProposeSeq())
		if acepMsg.getProposeSeq() == p.getProposeNum() {
			log.Println("recv ++", recvCount)
			recvCount++
		}
	}
	log.Println("Current proposer recv promise count=", recvCount)
	return recvCount
}

// getProposeNum 拿到提案者的提案数
func (p *proposer) getProposeNum() int {
	p.proposeNum = p.seq<<4 | p.id
	return p.proposeNum
}
