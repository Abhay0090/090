package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	stop        = 0
	count       = 0
	bit         = 0
	connections int
	stopMutex   sync.Mutex // Mutex for synchronizing access to 'stop'
)

func main() {
	fmt.Println("|--------------------------------------|")
	fmt.Println("|   Golang : UDP Stress Test Tool      |")
	fmt.Println("|          C0d3d By Lee0n123           |")
	fmt.Println("|--------------------------------------|")

	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s host port connections seconds timeout(second)\r\n", os.Args[0])
		os.Exit(1)
	}

	connections, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("connections should be an integer")
		return
	}
	times, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("seconds should be an integer")
		return
	}
	timeout, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Println("timeout should be an integer")
		return
	}
	addr := os.Args[1] + ":" + os.Args[2]

	var wg sync.WaitGroup

	for i := 0; i < connections; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			conn, err := net.DialTimeout("udp", addr, time.Duration(timeout)*time.Second)
			if err != nil {
				fmt.Println("Error connecting to the server:", err)
				return
			}
			defer conn.Close()

			buffer := make([]byte, 1024)
			rand.Seed(time.Now().UnixNano())
			start := time.Now()

			for {
				stopMutex.Lock()
				if stop > 0 {
					stopMutex.Unlock()
					break
				}
				stopMutex.Unlock()

				conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				rand.Read(buffer)
				_, err := conn.Write(buffer)
				if err != nil {
					fmt.Println("Error sending UDP packet:", err)
					continue
				}
				count++
				bit += 1024
				time.Sleep(time.Millisecond * 100)
			}
			elapsed := time.Since(start)
			fmt.Printf("UDP Stress Test Complete: %d packets sent in %.2f seconds\n", count, elapsed.Seconds())
		}(&wg)
	}

	time.Sleep(time.Second * time.Duration(times))
	stopMutex.Lock()
	stop++
	stopMutex.Unlock()
	wg.Wait()

	fmt.Printf("Total Sent: %d packets\n", count)
	fmt.Printf("Throughput: %.2f Mb/s\n", float64(bit)/1024/1024/float64(times))
	fmt.Printf("Packets per Second: %.2f packets/s\n", float64(count)/float64(times))
}
