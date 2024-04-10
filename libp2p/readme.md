
# 使用go-libp2p搭建中转服务器（circuit relay server）


## 在公网服务器上执行 relayserver.go

```
[root@VM-16-5-centos p2p_serve]# ./p2p_serve 
2024/04/10 16:10:38 relay1Info ID: 12D3KooWBMgWt2ftRzfvS76JorxajjkjYtm737PmRRfWAdewNqnb 
 Addrs: 
 [/ip4/10.0.16.5/tcp/42394
/ip4/10.0.16.5/udp/41724/quic-v1/webtransport/certhash/uEiCoknJUDEA0ONRCZQDilF9Urq_gKcH8-sPn-wX1unrDTA/certhash/uEiC0eMSC8eL4GCWs-F1Q5mqkYyN2YP6_FywQBK__I4ysvg
/ip4/10.0.16.5/udp/55028/quic-v1
/ip4/127.0.0.1/tcp/42394
/ip4/127.0.0.1/udp/41724/quic-v1/webtransport/certhash/uEiCoknJUDEA0ONRCZQDilF9Urq_gKcH8-sPn-wX1unrDTA/certhash/uEiC0eMSC8eL4GCWs-F1Q5mqkYyN2YP6_FywQBK__I4ysvg
/ip4/127.0.0.1/udp/55028/quic-v1
/ip6/::1/tcp/41118
/ip6/::1/udp/34874/quic-v1/webtransport/certhash/uEiCoknJUDEA0ONRCZQDilF9Urq_gKcH8-sPn-wX1unrDTA/certhash/uEiC0eMSC8eL4GCWs-F1Q5mqkYyN2YP6_FywQBK__I4ysvg
/ip6/::1/udp/52261/quic-v1]
```

## 执行本地的A host

```
[root@VM-16-5-centos p2p_client]# ./p2p_client -ip "47.126.23.45/tcp/42394" -sid "12D3KooWNB7a6vRE9kqqRnJCjLZAWUyEiyU6BCKao4KmnnNEBWUp"

WatingNode is onLine!use " go run .\main.go -d 12D3KooWNytApNFQHPHmNA4YmyJzu4jKafEPUfv4sx65GYdZs9zF " in other CLI
```


## 执行本地的B host

```
[root@VM-16-5-centos p2p_client]# ./p2p_client -ip "47.126.23.45/tcp/42394" -sid 12D3KooWNB7a6vRE9kqqRnJCjLZAWUyEiyU6BCKao4KmnnNEBWUp -d 12D3KooWNytApNFQHPHmNA4YmyJzu4jKafEPUfv4sx65GYdZs9zF
2024/04/10 16:26:15 Yep, that worked!
> 
```
