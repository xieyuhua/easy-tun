### Introduction
> A p2p tunnel(based on kcptun) which establishs port forwarding channel for both clients.
> Need one server with public ip and one udp port for hole digging. 

假设A开始给B的公网地址发送UDP数据的同时，给服务器S发送一个中继请求，要求B开始给A的公网地址发送UDP信息。A往B的输出信息会导致NAT A打开 一个A的内网地址与与B的外网地址之间的新通讯会话，B往A亦然。一旦新的UDP会话在两个方向都打开之后，客户端A和客户端B就能直接通讯， 而无须再通过引导服务器S了。


### QuickStart
```
Server: ./p2pserver -l ":4000" 
Client_A: ./p2pclient -r "SERVER_IP:4000" -l ":1012" -t "TARGET_IP_A:22"  -k 4321
Client_B: ./p2pclient -r "SERVER_IP:4000" -l ":10022" -t "TARGET_IP_B:80"  -k 4321
```

### Install from source

```
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```

公网服务器实现打洞 p2p成功后设备A和设备B直接通信，根据k来匹配，访问设备A的1012端口就等价于访问设备B的80端口，访问设备B的1022端口就等价于访问设备A的22端口。

#### Usage

```
root@ubuntu:~/p2p_tun# ./p2pclient -h
NAME:
   p2p_tun - client(based on kcptun)

USAGE:
   p2pclient [global options] command [command options] [arguments...]

VERSION:
   SELFBUILD

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --targettcp value, -t value      target server address (default: "127.0.0.1:22")
   --listentcp value, -l value      local listen address (default: ":12948")
   --remoteudp value, -r value      kcp server address (default: "vps:29900")
   --bindudp value, -b value        bind local udp (default: ":29900")
   --key value, -k value            p2p pair key (default: "1234")
   --passwd value                   pre-shared secret between client and server (default: "1234") [$KCPTUN_KEY]
   --crypt value                    aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none (default: "aes")
   --mode value                     profiles: fast3, fast2, fast, normal, manual (default: "fast2")
   --conn value                     set num of UDP connections to server (default: 1)
   --autoexpire value               set auto expiration time(in seconds) for a single UDP connection, 0 to disable (default: 0)
   --mtu value                      set maximum transmission unit for UDP packets (default: 1350)
   --sndwnd value                   set send window size(num of packets) (default: 1024)
   --rcvwnd value                   set receive window size(num of packets) (default: 1024)
   --datashard value, --ds value    set reed-solomon erasure coding - datashard (default: 10)
   --parityshard value, --ps value  set reed-solomon erasure coding - parityshard (default: 3)
   --dscp value                     set DSCP(6bit) (default: 0)
   --nocomp                         disable compression
   --sockbuf value                  per-socket buffer in bytes (default: 4194304)
   --keepalive value                seconds between heartbeats (default: 10)
   --snmplog value                  collect snmp to file, aware of timeformat in golang, like: ./snmp-20060102.log
   --snmpperiod value               snmp collect period, in seconds (default: 60)
   --log value                      specify a log file to output, default goes to stderr
   --quiet                          to suppress the 'stream open/close' messages
   -c value                         config from json file, which will override the command from shell
   --help, -h                       show help
   --version, -v                    print the version

```

