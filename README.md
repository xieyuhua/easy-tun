## What does easy-tun do?
To realize communication between different LAN machines.
实现不同网络机器之间的通信。组网为局域网

# serve

```
[root@VM-16-5-centos server]# ./server -h
Usage of ./server:
  -ip string
    	server address (default ":8006")

```

# client
```
[root@VM client]# ./client -h
Usage of ./client:
  -dev string
    	local tun device name (default "gtun")
  -ip string
    	子网掩码是 255.255.255.0、 10.10.10.1/24 (default "10.10.10.1/24")
  -ser string
    	server address (default "47.105.115.26:8006")
[root@VM-16-5-centos client]# 

```


# 路由设置
```
[root@Centos sss]# ip addr add 10.10.10.1/24 dev gtun
[root@Centos sss]# ip link set gtun up
[root@Centos sss]# route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.10.10.0      0.0.0.0         255.255.255.0   U     0      0        0 gtun


route print
netsh interface ip add address "gtun" 10.10.10.1 255.255.255.0
```


# todo
p2p打洞

## Guide
[基于TUN/TAP实现简单VPN](https://blog.csdn.net/qq_63445283/article/details/123779498)