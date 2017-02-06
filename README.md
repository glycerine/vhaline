# vhaline

high availability from a chain of processes

# summary of vhaline properties:

The 'ha' part of the name is for high-availability.

`vhaline` is inspired by chain-replication,
but is simpler and doesn't depend on a paxos system.

We are cloud and firewall friendly since we 
only need go one direction through the firewall.

We don't offer the same levels of proof or
consistency guarantees that chain-replication
does. This is just a best-effort primary-backup
chain. It is susceptible to split brain
situations on network partition, which could result in two
processes writing at once. If you need the
protection of a quorum, then vhaline is not
for you. In my use, I can tolerate having
more than one writer under a degraded network.
AP systems like Cassandra and Dynamo make
similar (configurable) choices.

# architecture

Like a linked list, each node has two pointers
or remote addresses. The
upstream pointer points to our parent, and
the downstream pointer point to our child.
Either or both pointers can be nil.

The most upstream node is the root, and the
root is the only one that does writes.

Regular state checkpoints are passed down
the chain every 10-30 seconds to the middle.

The middle then passes checkpoints down to
its child (the tail), if present.

clients pointing at servers: root (started) first,
is always a server without a parent:
~~~
     root <--- middle <--- tail
~~~

dataflow of checkpoints in the line:
~~~
     root ---> middle ---> tail
~~~

In vhaline, there is only ever
one parent and one child at each node, so the
graph is always a straight line chain of nodes.

Although active only in a line, nodes will be
informed of the whole chain structure on a
regular basis in order to enable working around
a crashed process.

I. Electing ourselves root in the chain:

* a) If we have no parent, then we are root. We write.

* b) If we detect parent failure, then we take over
     as root, and write.
     
* c) We regularly check for parent (and child) failure.
     This is done with pings. We typically require
     at least 3 failed pings before the TTL expires
     and we elect ourselves writer.

* d) If the upstream parent is configured (not nil),
     then we are a middle or last node. As middle or last
     nodes, we listen for checkpoints
     from upstream, persistent them, rotate them, and
     copy them to our child (if we have one). We
     don't write, but we do standby to write if our
     parent fails.

II. On child failure:

* a) Child failures should be detected. Once detected and cleared,
     we should allow another, different node to subscribe as our new child.
     
* b) If we have a child, we replicate checkpoints to them.

III. Misc. notes:

Ideally we should have a means of dedup-ing the checkpoints so we can
recognize that we've already gotten a checkpoint and
we don't propagate it downstream. If we only
propagate things that are new to us, that is much
more efficient/saves on bandwidth. The blake2b
function is already available on the Frames will suffice
if we find this critical. For now it is left undone.


testing
-------------

In the shell, do `ulimit -n 5000` first to raise the file limit. Otherwise `go test` may run out of files on the Test001 stress test for failure detection. Particularly on OSX.

administrative
--------------------

Copyright (c) 2017, Jason E. Aten.

license: MIT
