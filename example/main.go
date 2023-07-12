package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coreservice-io/geo_ip_v3/lib"
)

const geo_ip_update_key = ""

func main() {

	//////////
	client, err := lib.NewClient(geo_ip_update_key, "0.0.24", "./example", false, func(log_str string) {
		fmt.Println("log_str:" + log_str)
	}, func(err_log_str string) {
		fmt.Println("err_log_str:" + err_log_str)
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	//initial upgrade
	// upgrade_err := client.Upgrade(true)
	// if upgrade_err != nil {
	// 	log.Fatalln(upgrade_err)
	// 	return
	// }

	log.Println(client.GetInfo("104.233.16.169"))
	log.Println(client.GetInfo("5.78.52.174"))
	log.Println(client.GetInfo("116.227.21.107"))

	log.Println(client.GetInfo("2a7:1c44:39f3:1b::"))

	log.Println(client.GetInfo("601:196:4a00:55c0::"))
	log.Println(client.GetInfo("1000:db8::"))
	log.Println(client.GetInfo("804:b3c::"))

	time.Sleep(30 * time.Second)

	log.Println(client.GetInfo("172.104.160.0"))

	time.Sleep(30 * time.Hour)
}
