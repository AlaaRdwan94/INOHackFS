# My notes and findings about memberlist server-client

## My verdict overall:
This was supposed to be a server and a client communicating with each other using the memberlist, but after I've finished making it, <br >
I realized and found that it will not work properly because in the client-server model, the client will always be the one who initiate connection to the server and not the opposite.

The memberlist is meant to be used in a machine and make it act like a node in a network where all the nodes have 2 way communication and each node need a direct route to other nodes so they can ping each other.

The problem is that the "client" memberlist doesn't work correctly when runing in private LAN network that is found in almost all household routers, the PC can connect to outside world, but the outside world can't connect to it because the router by default doesn't forward connections received to the public ip (router obtains it from the ISP) to the private IP (PC obtains it from the router).

There's 2 solutions to this problem, but each have their drawbacks:
  1. Make port forwarding rule in the router' settings. <br >
     from WAN (public) Address:Port to LAN (private) Address:Port
  2. Use a tcp proxy server (to use it as a tunnel) or a VPN, there's `ssh` tunneling method and a service like ngrok.

The drawbacks are:
  1. making port forwarding on the router could be a hard task for the common user who doesn't have technical background, Also some router make it even harder to find the port forward settings.
  2. The ssh tunneling and service like ngrok only support TCP forwarding but the memberlist depends on UDP also which make it more difficult to set it up and you'll have to use more tools like `socat`


## Links:

[A stackexcagne answer about public and private IPs](https://networkengineering.stackexchange.com/a/42958)


### Memberlist links:

`Note: Memberlist is a part of Serf which also a part of Consul`<br >
`and all of them are developed and maintained by HashiCorp.`<br >
[The memberlist package page](https://pkg.go.dev/github.com/hashicorp/memberlist)<br >
[Consul's Gossip page](https://www.consul.io/docs/internals/gossip.html)<br >
[SWIM: scalable weakly-consistent infection-style process group membership protocol - white paper](https://ieeexplore.ieee.org/document/1028914)<br >
[NetTransport is a Transport implementation that uses connectionless UDP for packet operations, and ad-hoc TCP connections for stream operations](https://github.com/hashicorp/memberlist/blob/237d410aa2bf83254678ef78dd638480780e54a2/net_transport.go#L40)<br >
[Making Gossip More Robust with Lifeguard - Important article about memberlist](https://www.hashicorp.com/blog/making-gossip-more-robust-with-lifeguard/)<br >
[Making Gossip More Robust with Lifeguard - youtube video](https://youtu.be/u-a7rVJ6jZY)<br >
[Golang: Creating distributed systems using memberlist](https://davidsbond.github.io/2019/04/14/creating-distributed-systems-using-memberlist.html)<br >


### Tcp proxy servers and ssh tunneling:
[An important wikipedia article about the `socks` & `socks 5` protocol](https://en.wikipedia.org/wiki/SOCKS)<br >
[A wikipedia article descriping what is a proxy server](https://en.wikipedia.org/wiki/Proxy_server)<br >
[SSH Port Forwarding for TCP and UDP Packets](https://stackpointer.io/network/ssh-port-forwarding-tcp-udp/)<br >
[A reddit comment that clarifies difference between tcp proxy and true port forwarding](https://www.reddit.com/r/golang/comments/60b9ys/port_forwarding_with_go/df5bd94?utm_source=share&utm_medium=web2x)<br >
[Simple SSH port forward in Golang](https://stackoverflow.com/a/21655505)<br >
[TCP-UDP-Proxy github repo](https://github.com/MengRao/TCP-UDP-Proxy)<br >
[ Ngrok ](https://dashboard.ngrok.com/billing/plan)<br >

### VPN related links:
`OpenVPN is blocked in Egypt!` <br >
[A facebook post about bypassing the blocked OpenVPN. might be useful somehow](https://www.facebook.com/groups/egyptian.geeks/permalink/2697246933648331/)<br >
