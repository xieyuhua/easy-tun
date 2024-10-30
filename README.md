## What does easy-tun do?
To realize communication between different LAN machines.
实现不同网络机器之间的通信。组网为局域网

Rust https://github.com/vnt-dev/vnt   

# serve

```
[root@VM-16-5-centos server]# ./server -h
Usage of ./server:
  -ip string
    	server address (default ":8006")

```

# client
```
[root@VM client]# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o wintun-windows.exe
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
# 点对网设置

两个设备组网后，如果要访问对方设备内网下的其他设备(或者把对方的设备当跳板访问其他IP)，就需要用到点对网的配置了   
![image](https://github.com/user-attachments/assets/b35d8f89-fd90-4d25-b68d-4d31b213cb95)

- 前提条件：两个局域网均有运行了client的设备(图中的A1和C1)，并且跳板节点能访问对方的IP( 图中C1能访问C2，也就是说C2的防火墙对C1开放)
- 原理：图中的跳板节点是C1，A1通过C1访问C2，跳板节点也可以是局域网C的路由器
- 配置参数1：跳板节点设备C1增加参数 -o 0.0.0.0/0 ，该参数表示允许转发所有目标的流量，这个参数是处于内网安全考虑的，也可以是-o 192.168.2.12/32，表示只能转发到C2
- 配置参数2：设备A1增加参数 -i 192.168.2.0/24,10.26.0.3，该参数表示拦截目标是192.168.2.0/24的流量，转发到10.26.0.3节点，也就是图中的跳板节点C1
- 配置上述参数后，即可在A1上使用C2的内网IP(192.168.2.12)直接访问C2，可以尝试在A1上ping 192.168.2.12，会发现能访问了( 如果C2开启ping)

```
 使用系统的IP转发，这会进一步提高点对网的性能，在低性能设备上这很有用，
 在跳板节点上如下操作，执行如下命令
 Windows：
  #设置nat,名字可以自己取，网段是gtun的网段
  New-NetNat -Name vntnat -InternalIPInterfaceAddressPrefix 10.26.0.0/24
  #查看设置
  Get-NetNat
 Linux：
  # 开启ip转发
  sudo sysctl -w net.ipv4.ip_forward=1
  # 开启nat转发  表示来源10.26.0.0/24的数据通过nat映射后再从gtun以外的其他网卡发出去
  sudo iptables -t nat -A POSTROUTING ! -o gtun -s 10.26.0.0/24 -j MASQUERADE
  # 或者这样  表示来源10.26.0.0/24的数据通过nat映射后再从eth0网卡发出去
  sudo iptables -t nat -A POSTROUTING  -o eth0 -s 10.26.0.0/24 -j MASQUERADE
  # 查看设置
  iptables -vnL -t nat
 MacOS：
  # 开启ip转发
  sudo sysctl -w net.ipv4.ip_forward=1
  # 配置NAT转发规则
  # 在/etc/pf.conf文件中添加以下规则,en0是出口网卡，10.26.0.0/24是来源网段
  nat on en0 from 10.26.0.0/24 to any -> (en0)
  # 加载规则
  sudo pfctl -f /etc/pf.conf -e
```

# 网对网设置
我们把client运行到双方局域网的默认网关路由器中(不在路由器中也可以，无非就是再加路由，这里对此不作讨论)
![image](https://github.com/user-attachments/assets/b63072ed-8e73-4eb1-bdb5-21a5869ef2b5)

- 路由器配置点对网：路由器A增加参数-i 192.168.2.0/24,10.26.0.3，路由器C增加参数-o 0.0.0.0/0 ，该参数表示将192.168.2.0/24目标的转发到10.26.0.3，也就是路由器C
- 根据点对网那一节可知，这样配置就能让路由器A能访问路由器C下的所有内网IP，当然我们的目标是让局域网A下的所有设备都能访问到局域网C下的设备
- 原理：由于局域网A的默认网关是路由器A，也就是说在设备A1访问C2的内网IP(192.168.2.12)，数据会发到路由器A的网卡( 假设是eth0)，而路由器A是能访问C2的，所以再添加转发规则，将数据转发到gtun网卡，即可连通整个链路
- 命令操作 ：在路由器A上执行sudo iptables -t nat -A POSTROUTING -o client的虚拟网卡(一般为gtun) -s 192.168.1.0/24 -j MASQUERADE 该参数表示将数据发送到gtun
- 配置完之后，你会发现设备A1能通过路由器A、路由器C访问到设备C2了，如果要C2访问A1或者A2，则再反过来进行上述配置，就能双方互通了

# todo
p2p打洞

## Guide
[基于TUN/TAP实现简单VPN](https://blog.csdn.net/qq_63445283/article/details/123779498)

