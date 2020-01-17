package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	cName = flag.String("name", "duzhenxun", "客户端名称（注意2台客户端名不要用一样的哦）")
	cPort = flag.Int("p", 25251, "客户端要占用本机的端口")
	server = flag.String("s", "39.106.231.36:2525", "服务端地址")
)

func main() {
	flag.Parse()
	sAddr:=strings.Split(*server,":")
	sIP:=sAddr[0]
	sPort,_:=strconv.Atoi(sAddr[1])
	cAddr := &net.UDPAddr{IP: net.IPv4zero, Port: *cPort}
	conn, err := net.DialUDP("udp", cAddr, &net.UDPAddr{IP: net.ParseIP(sIP), Port: sPort})

	if err != nil {
		log.Println(err)
		return
	}
	//向服务端写消息，如果服务端会应答
	s:="Hi~小手25,My name is " + *cName + " ,我想问一下我需要和哪个地址建立连接？"
	_,err = conn.Write([]byte(s))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("与服务端连接并发送消息成功，等待服务器通知另一个客户端的ip地址与端口....")

	data := make([]byte,1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	fmt.Println("服务端发来信息：",n, remoteAddr, err,string(data))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("关闭与服务端的连接",remoteAddr)
	conn.Close()

	//分析服务端发来的另一个客户端信息
	dstAddrArr := strings.Split(string(data[:n]), ":")
	port, _ := strconv.Atoi(dstAddrArr[1])
	dstAddr := &net.UDPAddr{IP: net.ParseIP(dstAddrArr[0]), Port: port}

	log.Printf("本地:%s,对方:%s", cAddr.String(), dstAddr.String())

	//time.Sleep(1 * time.Second)
	//对方地址获取到后与对方进行连接
	conn2, err := net.DialUDP("udp", cAddr, dstAddr)
	if err != nil {
		log.Println("对方进行连接失败！！！")
		log.Println(err)
		return
	}

	//向另一方发送一条udp消息(对方的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方就可以发过来udp消息
	if _, err = conn2.Write([]byte("......")); err != nil {
		log.Println("第一次发送失败",err)
		return
	}
	time.Sleep(2 * time.Second)

	//给对方发数据
	go func() {
		for {
			s := fmt.Sprintf("%s 发来消息 我是：%s", cAddr.String(),*cName)
			n, err = conn2.Write([]byte(s))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//输出对方发来的数据
	for {
		data := make([]byte, 1024)
		n, _, err := conn2.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("--->", string(data[:n]))
		}
	}

}
