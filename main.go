package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	cache := make(map[string][]byte)

	conn, err := net.ListenPacket("udp", ":3000")
	if err != nil {
		fmt.Println("Error listening on chosen port:", err)
		return
	} //listen on port 3000 cuz 53 needed sudo and i didn't want to keep typing sudo
	defer conn.Close()

	//default port is 53 though
	fmt.Println("Listening for DNS requests on port 3000...")

	for {
		buffer := make([]byte, 512) //buffer of size 512 bytes

		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}
		log.Println(addr.String())
		log.Println(n)
		log.Println(cache)
		request := buffer[:n] // read upto n bytes, where n is the number of bytes read from buffer
		response, ok := cache[string(request)]

		if !ok {
			// If the request is not cached, forward it to the appropriate DNS server
			dnsServer := net.ParseIP("8.8.8.8")
			dnsAddr := &net.UDPAddr{IP: dnsServer, Port: 53}
			dnsConn, err := net.DialUDP("udp", nil, dnsAddr)
			if err != nil {
				fmt.Println("Error dialing DNS server:", err)
				continue
			}

			_, err = dnsConn.Write(request)
			if err != nil {
				fmt.Println("Error sending DNS request:", err)
				continue
			}

			buffer = make([]byte, 512)
			n, err = dnsConn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading DNS response:", err)
				continue
			}

			// Caching the response
			response = buffer[:n]
			cache[string(request)] = response
		}

		// Sending the response to the client
		_, err = conn.WriteTo(response, addr)
		if err != nil {
			fmt.Println("Error sending DNS response:", err)
			continue
		}
	}
}
