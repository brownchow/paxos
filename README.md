# Go 实现 paxos 共识算法

paper:

https://www.microsoft.com/en-us/research/uploads/prod/2016/12/paxos-simple-Copy.pdf

https://github.com/oldratlee/translations/tree/master/paxos-made-simple



CAP 理论：

只能保证 CA 或者 CP 



## 一致性模型：

### 1、弱一致性

   最终一致性 

- DNS
- Gossip( Cassandra的通信协议)

### 2、强一致性

- 同步
- Paxos
- Raft(multi-paxos)
- ZAB(multi-paxos)

### Paxos算法

Paxos 是一个共识算法，系统的最终一致性，不仅需要达成共识，还会取决于 Client 的行为。

Paxos 作者: LaTeX 作者

#### 1、角色介绍

Client: 系统外部角色，请求发起者，**像民众**

Proposer: 接收Client 请求，向集群提出提议（propose），并在冲突发生时，起到冲突调节的作用。**像议员**，替民众提出提案

Acceptor(Voter): 提议投票和接收者，只有在形成法定人数（Quorum，一般即为majority多数派）时，提议才会被接收，**像国会**。 

Learner：提议接受者，backup，备份，对集群一致性没什么影响，**像记录员**。

#### 2、步骤、阶段

Phase1a: prepare

   Proposer 提出一个提案，编号为N，此N大于之前这个proposer提出的提案编号，请求Accepter 的 quorum 接受（超过半数accepter）

Phase1b: promise

  如果N大于之前此acceptor之前接受的任何提案编号则接受，否则拒绝

phase2a:Accept

  如果达到了多数派，proposer 会发出accept请求，此请求包含提案编号N，以及提案内容

phase2b:Accepted

如果此Accepter 在此期间没有收到任何编号大于N的提案 ，则接受此提案内容，否则忽略



#### 基本流程



Basic Paxos 的问题：难实现，效率低（2轮RPC）,活锁

Multi Paxos: 新概念，Leader，唯一的proposer，所有请求都需要经过此Leader



### raft 算法



### ZAB 算法

与raft类似



### etcd

文档更全，比 zookeeper 更友好一点



docker 模拟 etcd 集群

## 参考：

https://www.bilibili.com/video/BV1TW411M7Fx

https://raft.github.io/  动画演示







## TODO:

1、raft 实现

2、zab 实现















































