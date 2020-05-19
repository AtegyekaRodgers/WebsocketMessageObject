/*
 
 Attribute Location Hexadecimal value Meaning                                *
 -----------------------------------------------------------------------------------------------
 Format 0x0000 -> 0x0003 4D 50 51 1B MPQ\x1b (format name)
 UserDataMaxSize 0x0004 -> 0x0007 00 02 00 00 512
 HeaderOffset 0x0008 -> 0x0011 00 04 00 00 1024
 UserDataSize 0x0012 -> 0x0015 3C 00 00 00 60
 	0x0016 -> 0x0020 05 08 00 02 2C ?
 Starcraft2 0x0021 -> 0x0042 53 74 61 [...] "Starcraft II Replay 11" in binary
 
 */



package main

import ( 
"fmt"
"log"
"os"
"bytes"
"encoding/binary"
)

type Header struct {
	UserDataMaxSize uint32
	HeaderOffset uint32
	UserDataSize uint32
	_ [5]byte
	Starcraft2 [22]byte
}

func main() {
// ...
header := Header{}

/*
 func readNextBytes(file *os.File, number int) []byte {
 	bytes := make([]byte, number)
 	_, err := file.Read(bytes)
 	if err != nil {
 		log.Fatal(err)
 	}
 	return bytes
 }

 func main() {
 path := "data/replay.SC2Replay"  //'.SC2Replay' si the format of the file we are decoding here. ie, it could also be .png, .pdf, etc.
 		file, err := os.Open(path)
 		if err != nil {
	 		log.Fatal("Error while opening file", err)
		}
		defer file.Close()
		fmt.Printf("%s opened\n", path)
		
	// ...
	
		formatName := readNextBytes(file, 4)
			
		fmt.Printf("Parsed format: %s\n", formatName)
		if string(formatName) != "MPQ\x1b" {
			log.Fatal("Provided replay file is not in correct format.")
		}
}

*/
//data := readNextBytes(file, 39) // 4 + (4 * uint32 (3) ) + (5 * byte (1) ) + (22 * byte (1) ) = 43
//NOTE: The above is used in case the data we want to decode is a file, located on the server. 

//TODO: read binary data from a websocket of a conneted client.

buffer := bytes.NewBuffer(data)  //'data' comes from client already encoded into binary, here we are going to decode it.

err = binary.Read(buffer, binary.LittleEndian, &header)
if err != nil {
	log.Fatal("binary.Read failed", err)
}

fmt.Printf("Parsed data:\n%+v\n", header)


}
