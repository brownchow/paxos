package paxos

import "log"

type acceptor struct {
	id         int
	learners   []int
	acceptMsg  message
	promiseMsg message
	nt         nodeNetwork
}

// create a acceptor will also assign learning IDs int acceptor
// Acceptor: will response request from proposor, promise the first and largest seq number propose
//           After proposer reach the majority promise, Acceptor will pass the proposal value to learner to confirm and choose
func NewAcceptor(id int, nt nodeNetwork, learners ...int) acceptor {
	newAcceptor := acceptor{id: id, nt: nt}
	newAcceptor.learners = learners
	return newAcceptor
}

// Accept process detail logic
func (a *acceptor) run() {
	for {
		// log.Println("acceptor:", a.id, "wait to recv msg")
		m := a.nt.recv()
		if m == nil {
			continue
		}

		// log.Println("acceptor:", a.id, "recv message ", *m)
		switch m.typ {
		case Prepare:
			promiseMsg := a.recvPrepare(*m)
			a.nt.send(*promiseMsg)
			continue
		case Propose:
			accepted := a.recvPropose(*m)
			if accepted {
				for _, lId := range a.learners {
					m.from = a.id
					m.to = lId
					m.typ = Accept
					a.nt.send(*m)
				}
			}
		default:
			log.Fatal("Unsupport message in acceptor ID:", a.id)
		}
	}
}

// After acceptor receive prepare message
// It will check prepare number and return acceptor if it is bigest one
func (a *acceptor) recvPrepare(prepare message) *message {
	if a.promiseMsg.getProposeSeq() >= prepare.getProposeSeq() {
		log.Println("ID:", a.id, " Already accept bigger one")
		return nil
	}
	log.Println("ID:", a.id, " Promise")
	prepare.to = prepare.from
	prepare.from = a.id
	prepare.typ = Promise
	a.acceptMsg = prepare
	return &prepare
}

// Recv Propose only check if acceptor already accept bigger propose before
// Otherwise, will just forward this message out and change its type to "Accept" to learning later
func (a *acceptor) recvPropose(proposeMsg message) bool {
	// already accept message is identical with previous promise message
	log.Println("accept: check propose.", a.acceptMsg.getProposeSeq(), proposeMsg.getProposeSeq())
	if a.acceptMsg.getProposeSeq() > proposeMsg.getProposeSeq() || a.acceptMsg.getProposeSeq() < proposeMsg.getProposeSeq() {
		log.Println("ID:", a.id, " acceptor not take propose: ", proposeMsg.val)
		return false
	}
	log.Println("ID:", a.id, "Accept")
	return true
}
