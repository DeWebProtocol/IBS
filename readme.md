## IBS
P2P Information Broadcast Simulator

## Building the source

Building `IBS` requires Go version 1.22 or later. You can install them using your favourite package manager. Once the dependencies are installed, run

```shell
go build
```

## Running `IBS`

### Hardware Requirements

Minimum:

* CPU with 2 cores
* 4GB RAM

Recommended:

* Fast CPU with 4+ cores
* 8GB+ RAM

### Usage
#### Flags

| Flag     | Description                                                                                                            | Default Value |
|----------|------------------------------------------------------------------------------------------------------------------------|---------------|
| net_size | the number of nodes in the network                                                                                     | 1000          |
| broadcast_per_node | the number of broadcast initialized by each node                                                                       | 10            |
| broadcast_interval | the interval between broadcast initialization                                                                          | 5000(μs)      |
| k | (kademlia only) bucket size of kademlia                                                                                | 15            |
| beta | (kademlia only) the broadcast redundancy factor β                                                                      | 1             |
|  table_size    | (flood only) the table size of nodes in flooding based net                                                             | 15            |
|  degree | (flood only) the broadcast degree in flooding based net                                                                | 4             |
| crash_interval | the interval of network disturbance, unit: μs(0.001ms)                                                                 | math.MaxInt   |
| with_ne | using Neighbor Evaluation mechanism                                                                                    | false         |
| malicious | half of nodes in the network will receive messages but do not relay them, instead of just diconnected from the network | false         |

#### Examples
- 500 nodes kademlia broadcasting network with NE mechanism, a new broadcast will be initialized every 2ms, redundancy factor is 3:
```shell
> ibs kademlia latency --net_size 500 --broadcast_interval 2000 --beta 3 --with_ne
```
- 1000 (default) nodes kademlia broadcasting network, a new broadcast will be initialized every 5ms (default), and every 60s about half of nodes (around 500) in the net will be disconnected (rejoin disconnected nodes in the last round of disturbance into the network, clear and rebuild routing tables of them)
```shell
> ibs kademlia coverage --beta 2 --crash_interval 60000000
```