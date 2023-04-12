package main

// zero dependency, all of these are built-in packages. Only http package is actually required. The others are for logging purposes.
import (
	"fmt"
	"log"
	"net"
)

func main() {
	cache := make(map[string][]byte) // define the cache as a golang map, this is essentially a key value store

	conn, err := net.ListenPacket("udp", ":3000") //udp port connection that listens for packets
	if err != nil {
		fmt.Println("Error listening on chosen port:", err) // error handling
		return
	} //listen on port 3000 cuz 53 needed sudo and i didn't want to keep typing sudo
	defer conn.Close()

	//default port is 53 though
	fmt.Println("Listening for DNS requests on port 3000...") // logging purposes

	for {
		buffer := make([]byte, 512) // buffer of size 512 bytes where the request is read from

		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err) // error handling
			continue
		}
		log.Println(addr.String())             //logging info
		log.Println(n)                         //logging info
		log.Println(cache)                     //logging the cache
		request := buffer[:n]                  // read upto n bytes, where n is the number of bytes read from buffer
		response, ok := cache[string(request)] // define the response

		if !ok {
			// If the request is not cached, forward it to the appropriate DNS server
			dnsServer := net.ParseIP("8.8.8.8")              // forwarding the request to google's DNS server
			dnsAddr := &net.UDPAddr{IP: dnsServer, Port: 53} // using our default port 53 here since its already set up
			dnsConn, err := net.DialUDP("udp", nil, dnsAddr)
			if err != nil {
				fmt.Println("Error dialing DNS server:", err) // error handling
				continue
			}

			_, err = dnsConn.Write(request) // check if error in sending the request and send the request as well
			if err != nil {
				fmt.Println("Error sending DNS request:", err) // error handling
				continue
			}

			buffer = make([]byte, 512) // define the buffer where the request is read from
			n, err = dnsConn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading DNS response:", err) //error handling
				continue
			}

			// Caching the response
			response = buffer[:n]             // define what the response was
			cache[string(request)] = response // caching
		}

		// Sending the response to the client
		_, err = conn.WriteTo(response, addr)
		if err != nil {
			fmt.Println("Error sending DNS response:", err) // error handling
			continue
		}
	}
}
