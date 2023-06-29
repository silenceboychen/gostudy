package main

import (
	"fmt"
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
	"time"
)

// 解决火狐浏览器断开不会自动重连问题
func main() {
	es := eventsource.New(
		eventsource.DefaultSettings(),
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		})
	defer es.Close()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/events", es)
	go func() {
		for {
			// 设置事件名称为：test-event
			es.SendEventMessage(fmt.Sprintf("send data: %s", time.Now().Format("2006-01-02 15:04:05")), "test-event", "")
			log.Printf("客户端连接数: %d", es.ConsumersCount())
			time.Sleep(2 * time.Second)
		}
	}()

	log.Println("Open URL http://localhost:8080/ in your browser.")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
