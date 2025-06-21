package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	host = "localhost"
	port = "8082"
	TYPE = "tcp"
)

type metadata struct {
	name     string // file name denge as part of metadata
	filesize uint64 // ye to pata hi hoga 🫡
	reps     uint32 // KITNE SEGMENTS we are sending means total file size/ chunks size
}

// / file ek object banate hai os library ke ek structure name file ka niche as a parameter paas kiya
func actual_metadata(file *os.File) metadata {
	fileInfo, err := file.Stat() // ye stat tries to get metadata of the file jisko transfer karn hai
	if err != nil {              // AGAR ERROR AA GYA
		log.Fatal(err) // YE BUILT FUNCTION HAI JO ERROR MESSAGE PRINT KAREGA AND PROGRAM GETS STOP
	}
	size := fileInfo.Size()
	// actually fileinfo is one of the return of function stat jo all metadata info of file deta hai
	header := metadata{
		name:     file.Name(),
		filesize: uint64(size),
		reps:     uint32(size/1014) + 1,
	}
	return header
}

func sendfile(path string, conn *net.TCPConn) { // here client ke liye ek socket banaya using inbuilt struct TCPConn of lib net
	file, err := os.OpenFile(path, os.O_RDONLY, 0755) // read only mode BUT WHAT IS THIS 0755 ?
	if err != nil {
		log.Fatal(err)
	}
	header := actual_metadata(file)
	// GO ME HUM LOG MEMORY SPACE KE BARABR BHI SLICE BANA SAKTE CRAZY!!💀
	databuffer := make([]byte, 1014)
	headerbuffer := []byte{1} // 1 BYTE KI SLICE to start header data
	segmentbuffer := []byte{0}
	temp := make([]byte, 4)       // to store the fielsize for another peer
	received := make([]byte, 100) // to send message is file recieved or not
	for i := 0; i < int(header.reps); i++ {
		n, _ := file.ReadAt(databuffer, int64(i*1014)) // readat functions returns n means how many bytes actually read
		// becoz header.repos can be very large GBs ki file ho agr so we'll use
		// ENCODING
		// 🗿HERE n IS LENGTH OF DATA🗿
		if i == 0 {
			// NUMBER OF SEGMENTS KO HEADER FILE ME DAALO
			binary.BigEndian.PutUint32(temp, header.reps)
			headerbuffer = append(headerbuffer, temp...) // temp ki hr ek byte // triple dot is for appending each byte one by one
			// LENGTH OF NAME KO ==>
			binary.BigEndian.PutUint32(temp, uint32(len(header.name)))
			headerbuffer = append(headerbuffer, temp...)
			// name
			headerbuffer = append(headerbuffer, []byte(header.name)...)
			headerbuffer = append(headerbuffer, 0) // so that termination pata chl sake
			_, err := conn.Write(headerbuffer)     // header file ko tcp me send kr diya
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Read(received)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(received))
		}
		// segement number paas karo
		binary.BigEndian.PutUint32(temp, uint32(i))
		segmentbuffer = append(segmentbuffer, temp...)
		// length of data
		binary.BigEndian.PutUint32(temp, uint32(n))
		segmentbuffer = append(segmentbuffer, temp...)
		//
		segmentbuffer = append(segmentbuffer, databuffer...)

		_, err := conn.Write(segmentbuffer)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Read(received)
		fmt.Println(string(received))

		if err != nil {
			log.Fatal(err)
		}
		segmentbuffer = []byte{0}
	}
}
func main() {
	// niche wali line converts above strings of host and port to objectso that go can understand and stores address in
	// in tcpserver
	tcpserver, err := net.ResolveTCPAddr(TYPE, host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	// now trying to connect with server if server is listening we can read / write whatever we want
	conn, err := net.DialTCP(TYPE, nil, tcpserver) // here nil is for os to select any port it wants in local address
	if err != nil {
		log.Fatal(err)
	}
	sendfile("C:\\Users\\myind\\OneDrive\\Pictures\\Screenshots\\Screenshot 2025-06-21 225819.png", conn)
	recieved := make([]byte, 1024) // ek buffer bana liya jo server ke incoming data ko store karega
	_, err = conn.Read(recieved)   //  ye actual me data recieved ko store karta hai
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(recieved))
}
