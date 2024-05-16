package main

import (
	"fmt"

	"github.com/ceph/go-ceph/rados"
	"github.com/kr/pretty"
)

func main() {

	fmt.Println("rados conn")
	conn, err := rados.NewConn()
	if err != nil {
		panic(err)
	}
	fmt.Println("rados conn created")
	if err := conn.ReadDefaultConfigFile(); err != nil {
		panic(err)
	}

	fmt.Println("read ceph.conf")
	if err := conn.ReadConfigFile("ceph.conf"); err != nil {
		panic(err)
	}

	v, err := conn.GetConfigOption("mon_host")
	if err != nil {
		panic(v)
	}
	fmt.Println(v)

	fmt.Println("connect")
	if err := conn.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("get")
	stat, err := conn.GetClusterStats()
	if err != nil {
		panic(err)
	}
	pretty.Println(stat)

	poolId, err := conn.GetPoolByName("myfs-data")
	if err != nil {
		panic(err)
	}

	fmt.Println(poolId)
}
