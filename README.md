# golang-tcp-port-reuse

This repo contains code for reproducing port binding problems in Go and for understanding how socket stuff works in general.


## Problem description

Sometimes your TCP server written in Go fails with error `bind: address already in use`

It could happen in such cases as:
1. Other process already listens on this port on the same address. (non interesting case)
2. Server code explicitly disables SO_REUSEADDR (enabled by default in *nix systems)
and tries to bind on same port and address which have unfinished connection in `TIME_WAIT` state.
3. Client uses random port for itself and make a lot of reconnects so it exhausted all ephemeral ports
(left all potential available port in `TIME_WAIT` state)
4. Client uses the same port to itself

Or sometimes your TCP client fails bind to ephemeral port.

*TODO: describe*


### Scenario #1: Can't bind on the same port by the server w/o `SO_REUSEADDR` option if socket in TIME_WAIT state

1. Server listens on 1215 port, close connection and exit after 2 seconds:
`go run server/main.go -l 127.0.0.1:1215 -s 2`.
2. Client connects to 1215 port, exit without closing connection after 4 seconds:
 `go run client/main.go -e 127.0.0.1:1215 -s 5 -noclose`.
 Lefts server side in `TIME_WAIT` state for 30s (depends on system settings).
3. Check connections statuses: `(netstat -a -n | head -n 2) ; (netstat -a -n | grep 1215)`.
4. Try to bind on port w/o `SO_REUSEADDR` option: `go run server/main.go -l 127.0.0.1:1215 -s 2 -no-reuse-addr`
(should fail if you are fast enough)
5. Try to bind on port with default mode: `go run server/main.go -l 127.0.0.1:1215 -s 2`
   (should works if no other process listens on the same ip:port pair)



## Resources

* SO_REUSEPORT, SO_REUSEADDR related discussion: https://github.com/golang/go/issues/9661
* SO_REUSEPORT/ADDR
  * [How different about the condition of binding](https://medium.com/uckey/the-behaviour-of-so-reuseport-addr-1-2-f8a440a35af6)
  * [How packets forwarded to multiple sockets](https://medium.com/uckey/so-reuseport-addr-2-2-how-packets-forwarded-to-multiple-sockets-ce4b83cd0fd2)
* «why we do not use SO_REUSEADDR on windows»: https://github.com/golang/go/commit/c3733b29d494995859bb6d6241797f67ece4c53d


* [This is strictly a violation of the TCP specification](https://blog.cloudflare.com/this-is-strictly-a-violation-of-the-tcp-specification/)
* [TCP: About FIN_WAIT_2, TIME_WAIT and CLOSE_WAIT](https://benohead.com/tcp-about-fin_wait_2-time_wait-and-close_wait/)
