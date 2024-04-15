package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"github.com/fatih/color"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	// 	"github.com/songgao/water"
	"github.com/labulakalia/water"
)

var (
	inSer = flag.String("ser", "47.105.115.26:8006", "server address")
	inDev = flag.String("dev", "gtun", "local tun device name")
	ip    = flag.String("ip", "10.10.10.1/24", "子网掩码是 255.255.255.0、 10.10.10.1/24")
)

func main() {
	flag.Parse()

	// 创建tun网卡
	config := water.Config{
		DeviceType: water.TUN, //TAP TUN
	}
	// windows os是config.InterfaceName
	config.Name = *inDev

	ifce, err := water.New(config)
	if err != nil {
		color.Red(err.Error())
		return
	}
	// 连接server，8006
	conn, err := connServer(*inSer)
	if err != nil {
		color.Red(err.Error())
		return
	}

	color.Cyan("Interface tun device Name: %s\n", ifce.Name())
	color.Cyan("server address	%s", *inSer)
	color.Cyan("connect server succeed.")

	// 读取tun网卡，将读取到的数据转发至server端
	go ifceRead(ifce, conn)
	// 接收server端的数据，并将数据写到tun网卡中
	go ifceWrite(ifce, conn)

	//if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
	switch runtime.GOOS {
	case "linux":
		// 添加 IP 地址到接口
		cmdAddIP := exec.Command("ip", "addr", "add", *ip, "dev", *inDev)
		_, errAddIP := cmdAddIP.Output()
		if errAddIP != nil {
			color.Red("Failed to add IP: %v\n", errAddIP)
			return
		}
		color.Cyan("Added IP: %s\n", *ip)
		// 启用接口
		cmdLinkUp := exec.Command("ip", "link", "set", *inDev, "up")
		_, errLinkUp := cmdLinkUp.Output()
		if errLinkUp != nil {
			color.Red("Failed to set link up: %v\n", errLinkUp)
			return
		}
	case "windows":
		//route print
		//netsh interface ip add address "gtun" 10.10.10.1 255.255.255.0
		cmdLinkUp := exec.Command("cmd", "/C", "netsh interface ip add address ", *inDev, *ip)
		_, errLinkUp := cmdLinkUp.CombinedOutput()
		if errLinkUp != nil {
			color.Red("Running generic error: %v\n", errLinkUp)
			return
		}
		color.Cyan("Added IP: %s\n", *ip)
	default:
		color.Red("netsh interface ip add address error: %v\n", runtime.GOOS)
	}

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)
	<-sig
}

// 连接server
func connServer(srv string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", srv)
	if err != nil {
		return nil, err
	}
	return conn, err
}

// 读取tun网卡数据转发到server端
func ifceRead(ifce *water.Interface, conn net.Conn) {
	packet := make([]byte, 2048)
	for {
		// 从tun网卡读取数据
		size, err := ifce.Read(packet)
		if err != nil {
			color.Red(err.Error())
			break
		}
		// 转发到server端
		err = forwardSer(conn, packet[:size])
		if err != nil {
			color.Red(err.Error())
		}
	}
}

// 将server端的数据读取出来写到tun网卡
func ifceWrite(ifce *water.Interface, conn net.Conn) {
	// 定义SplitFunc，解决tcp的粘贴包问题
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// 检查 atEOF 参数和数据包头部的四个字节是否为 0x123456
		if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == 0x123456 {
			// 数据的实际大小
			var size int16
			// 读出数据包中实际数据的大小(大小为 0 ~ 2^16)
			binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &size)
			// 总大小 = 数据的实际长度+魔数+长度标识
			allSize := int(size) + 6
			// 如果总大小小于等于数据包的大小，则不做处理！
			if allSize <= len(data) {
				return allSize, data[:allSize], nil
			}
		}
		return
	}
	// 创建buffer
	buf := bytes.NewBuffer(nil)
	// 定义包，由于标识数据包长度的只有两个字节故数据包最大为 2^16+4(魔数)+2(长度标识)
	packet := make([]byte, 65542)
	for {
		nr, err := conn.Read(packet[0:])
		buf.Write(packet[0:nr])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				color.Red(err.Error())
				break
			}
		}
		scanner := bufio.NewScanner(buf)
		scanner.Split(splitFunc)
		for scanner.Scan() {
			_, err = ifce.Write(scanner.Bytes()[6:])
			if err != nil {
				color.Red(err.Error())
			}
		}
		buf.Reset()
	}
}

// 将tun的数据包写到server端
func forwardSer(srvcon net.Conn, buff []byte) (err error) {
	output := make([]byte, 0)
	magic := make([]byte, 4)
	binary.BigEndian.PutUint32(magic, 0x123456)
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(buff)))

	// magic
	output = append(output, magic...)
	// length
	output = append(output, length...)
	// data
	output = append(output, buff...)

	left := len(output)
	for left > 0 {
		nw, er := srvcon.Write(output)
		if err != nil {
			err = er
		}
		left -= nw
	}
	return err
}
