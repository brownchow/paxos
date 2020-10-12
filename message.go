package paxos

type msgType int

const (
	Prepare msgType = iota + 1 // send from proposer -> acceptor
	Promise                    // send from acceptor -> proposer
	Propose                    // send from proposer -> acceptor
	Accept                     // send from acceptor -> learner
)

type message struct {
	from   int
	to     int
	typ    msgType
	seq    int
	preSeq int
	val    string
}

func (m *message) getProposeVal() string {
	return m.val
}

func (m *message) getProposeSeq() int {
	return m.seq
}
