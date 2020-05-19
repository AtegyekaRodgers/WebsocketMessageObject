package websocketMessageObject

import ( 
"fmt"
"log"
"bytes"
"strings"
"encoding/binary"
"encoding/json"
)

	
//=============================================================
type websocketMessageObject {
	BinaryData make([]byte,1024000) //max:1GB 
	filezLength uint32
	jsonnLength uint32
	stringzLength uint32 
	filez []os.File //array of files. functions will keep appending until Encode() is called
	jsonn interface{} //variable to hold json object. A function will set it before Encode() is called
	stringz []string //array of strings. functions will keep appending until Encode() is called
}

func NewWebsocketMessageObject() *websocketMessageObject {
	return &websocketMessageObject{filezLength=0, jsonnLength=0,stringzLength=0}
}

func NewWMO(){  
	return &websocketMessageObject{filezLength=0, jsonnLength=0,stringzLength=0}
} 
//=============================================================
//---------------
type wmoSectionReader struct { 
	readerat io.ReaderAt
}
 
func newWmoSectionReader(readerat io.ReaderAt) *wmoSectionReader { 
	return &wmoSectionReader{readerat: readerat}
} 

func (wsr *wmoSectionReader) ReadAt(dat []byte,offset int64) (int, error) {
	n, err := wsr.readerat.ReadAt(dat,offset)
	if err != nil {
		return n, err
	}   
	return n, nil
}
//----------------
type wmoSectionWriter struct { 
	writerat io.WriterAt
}

func newWmoSectionWriter(writerat io.WriterAt) *wmoSectionWriter { 
	return &wmoSectionWriter{writerat: writerat}
}

//WriteAt writes len(dat) bytes from 'dat' to the underlying data stream (wmo.BinaryData for this case) at offset 'off'.
func (wsw *wmoSectionWriter) ReadAt(dat []byte,offset int64) (int, error) {
	n, err := wsw.writerat.WriteAt(dat,offset)
	if err != nil {
		return n, err
	}   
	return n, nil
}

//----------------- 
type wmoHeader struct {
	FilesHeaderOffset uint32
	FilesHeaderSize uint32
	FilesTotalSize uint32
	JsonOffset uint32
	JsonSize uint32
	StringsOffset uint32
	StringsSize uint32 
} //total header size = 4*7 = 28 bytes = 28*8 = 204 bits = 0.028 KB

type wmoFilesHeader struct {
	NumberOfFiles uint8
	EachFileOffset []uint32
	EachFileSize []uint16
}

type wmoHeaderReader struct { 
	reader io.Reader
}

func newWmoHeaderReader(reader io.Reader) *wmoHeaderReader { 
	return &wmoHeaderReader{reader: reader}
}

func (whr *wmoHeaderReader) Read(hd []byte) (int, error) {
	hbuf := make([]byte, 224)
	n, err := whr.reader.Read(hbuf)
	if err != nil {
		return n, err
	}  
		copy(hd, hbuf) 
		return n, nil
}
  
  func readWmoHeader(dataFromClient []byte) *wmoHeader {
		  //buffer := bytes.NewBuffer(dataFromClient)  
		  /* this is an alternative approach of 
		  receiving data from client in case we 
		  need just a buffer barely. */
	  
		wmoheader := &wmoHeader{}  
		headreader := newWmoHeaderReader(dataFromClient)
		headerBuf := make([]byte, 224) 
		n, err := headreader.Read(headerBuf) 
		err := binary.Read(headerBuf, binary.LittleEndian, &wmoheader)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		} 
		//fmt.Printf("Parsed data:\n%+v\n", wmoheader)
		return wmoheader
	}
	

//------------------Decoders----------
func (wmo *websocketMessageObject) DecodeJson(dataFromClient []byte) string {
	hedr := readWmoHeader(dataFromClient)
	sectionreader := newWmoSectionReader(dataFromClient)
	jsonBuf := make([]byte, hedr.JsonSize) 
	n, err := sectionreader.ReadAt(jsonBuf, hedr.JsonOffset) 
	//convert json binary data into string  
	jsonstr := ''
	if json.Valid(jsonBuf) {
		if err := json.Unmarshal(jsonBuf, &jsonstr); err != nil {
			return err
		}
	}else{  
		fmt.Printf("!! Erro:\n%s\n", err)
	}
	
	return jsonstr 
}

func (wmo *websocketMessageObject) ReadFilesBytes(dataFromClient []byte) ([]byte,uint32) {
	hedr := readWmoHeader(dataFromClient)
	sectionreader := newWmoSectionReader(dataFromClient)
	filesDataBuf := make([]byte, hedr.FilesTotalSize) 
	n, err := sectionreader.ReadAt(filesDataBuf, hedr.FilesHeaderOffset)
	return filesDataBuf, hedr.FilesHeaderSize
}

func (wmo *websocketMessageObject) DecodeFiles(dataFromClient []byte) []os.File {
	filesBytes, fheaderSize = wmo.ReadFilesBytes(dataFromClient)
	//instantiate wmoFilesHeader
	var wmofheader wmoFilesHeader
	//use section reader to read the files header from filesBytes
	fheadsr := newWmoSectionReader(dataFromClient)
	fheadBuf := make([]byte, fheaderSize) 
	n, err := fheadsr.ReadAt(filesBytes, 0)
	err := binary.Read(fheadBuf, binary.LittleEndian, &wmofheader)
	//--------
	//now use the headers to read files and create an array of os.File objects.
	var allFiles []os.File
	var oneFile os.File 
	fdatasr := newWmoSectionReader(filesBytes)
	i:=0 
	for i<wmofheader.NumberOfFiles {
		oneFileBuf := make([]byte, wmofheader.EachFileSize[i]) 
		//read a section from filesBytes into oneFileBuf
		n, err := fdatasr.ReadAt(oneFileBuf, wmofheader.EachFileOffset[i])
		err := binary.Read(oneFileBuf, binary.LittleEndian, &oneFile)
		//push the 'oneFile' to the array of files - 'allFiles' .
		allFiles.Push(oneFile)
		i++
	}
	
	//now we already have an array of files, which is the return value for this function
	return allFiles
} 

func (wmo *websocketMessageObject) DecodeStringAll(dataFromClient []byte) map([string] string) {
	hedr := readWmoHeader(dataFromClient)
	sectionreader := newWmoSectionReader(dataFromClient)
	strBuf := make([]byte, hedr.StringsSize) 
	n, err := sectionreader.ReadAt(strBuf, hedr.StringsOffset) 
	//convert str binary data into string  
	stringz := string(strBuf) 
	//split the 'stringz' content using dilimitors and create a map of strings with string keys.
	strArray := strings.Fields(stringz) //splitting using spaces: stringz = 'key1-value1 key2-value2 key3-value3 ...'
	var newStrMap map([string] string)
	for _, keyAndVal := range strArray {
		strkey := strings.Split(keyAndVal,'-')[0]
		strval := strings.Split(keyAndVal,'-')[1]
		newStrMap[strkey] = strval
	} 
	return newStrMap 
}

func (wmo *websocketMessageObject) DecodeString(dataFromClient []byte,strkey string) string {
	newStrMap := wmo.DecodeStringAll(dataFromClient)
	return newStrMap[strkey]
}

//------------------Encoders----------

func (wmo *websocketMessageObject) setFilesSize(size int64){ wmo.filezLength=size; }
func (wmo *websocketMessageObject) setJsonSize(size int64){ wmo.jsonnLength=size; }
func (wmo *websocketMessageObject) setStringsSize(size int64){ wmo.stringzLength=size; }
func (wmo *websocketMessageObject) addToFilesSize(size int64){ wmo.filezLength+=size; }                 
func (wmo *websocketMessageObject) addToJsonSize(size int64){ wmo.jsonnLength+=size;}                    
func (wmo *websocketMessageObject) addToStringsSize(size int64){ wmo.stringzLength+=size;}
func (wmo *websocketMessageObject) PrintFilesSize(){ fmt.printf("\nwmo.filezLength: = %d",wmo.filezLength); }
func (wmo *websocketMessageObject) PrintJsonSize(){ fmt.printf("\nwmo.jsonnLength: = %d",wmo.jsonnLength); }
func (wmo *websocketMessageObject) PrintStringsSize(){ fmt.printf("\nwmo.stringzLength: = %d",wmo.stringzLength); }

func (wmo *websocketMessageObject) AddFile(file *os.File) {
	wmo.filez = append(wmo.filez, file)
	file_size := len([]byte(file))
	wmo.addToFilesSize(file_size)
}

func (wmo *websocketMessageObject) AddFileFrom(filepathh string) {
	file, err := os.Open(filepathh)
	if err != nil {
		log.Fatal("Error while opening file by wmo.AddFileFrom() function", err)
	}
	defer file.Close()
	wmo.filez = append(wmo.filez, file)
	file_size := len([]byte(file))
	wmo.addToFilesSize(file_size)
} 

func (wmo *websocketMessageObject) AddJson(myjson interface{}) {  
	 wmo.jsonn = myjson
	 	 jsonnstr, err := json.Marshal(wmo.jsonn)
		 if err != nil {
			 fmt.Println("error:", err)
			 return
		 }
	 json_size := len([]byte(jsonnstr))
	 wmo.setJsonSize(json_size)
}

func (wmo *websocketMessageObject) AddString(keyy string,mystring string) { 
	str := keyy+'-'+mystring 
	wmo.stringz = append(wmo.stringz, str) 
	str_size := len([]byte(str))
	wmo.addToStringsSize(str_size) 
}
 
func (wmo *websocketMessageObject) Encode() {
	/*
              
	0                    28    34               files-content         json                            strings
	|--,--,--,--,--,--,--|~~~~~|--,--,...,--,...|-----,-------,---,...|-------------------------------|-------|
	wmo-header                 files-header
	                           0                                      #
	                           |<-----------------files-------------->|
	*/
	
	var filez_start_point uint32 
	var json_start_point uint32 
	var stringz_start_point uint32
	 
	  var file_offsets []uint32
	  var file_sizes []uint32
	  total_size_of_files := 0
	  
	  wmo_offset_track := 34  //28 bytes (4*7) for the wmo header, 5 for 'dont care' bytes, 1 byte for holding number of files
	  files_ofst_track := 0   //8 bits == 1 byte for holding number of files, starting from the end of wmo header 
	  numberoffilez := len(wmo.filez)
	  size_of_files_header:=1+(numberoffilez*(4+4)) //1 byte holds No.of files, 4 bytes(size of uint32) for each file_sizes[element] and 4 bytes for each file_offsets[element]
	  total_size_of_files += size_of_files_header
	  wmo_offset_track =34 + size_of_files_header  //increment the '*_track' by 'size_of_files_header'
	  files_ofst_track = 0 + size_of_files_header
	   
	   wmoswriter := newWmoSectionWriter(wmo.BinaryData) //wmoSectionWriter  
	  
	  //loop through files to get info about each of them, 
		//and write each into wmo.BinaryData.
		for _, fi := range wmo.filez {
		  file_bytes := []byte(fi)
		  fiLength:= len(file_bytes)  
		  file_offsets=append(file_offsets,files_ofst_track)
		  file_sizes = append(file_sizes,fiLength) 
		  //write the 'file_bytes' at offset = wmo_offset_track in wmo.BinaryData
		  wmoswriter.WriteAt(file_bytes, wmo_offset_track) 
		  total_size_of_files += fiLength 
		  files_ofst_track += fiLength 
		  wmo_offset_track += fiLength
	  }
	  
	 //read everything that is in 'this.jsonn', stringfy & make it binary, 
	 //then add to this.BinaryData at appropriate offset. 
	 	 jsonnstr, err := json.Marshal(wmo.jsonn)
		 if err != nil {
			 fmt.Println("error:", err)
			 return
		 }
		 jsonbytes = []byte(jsonnstr)
		 json_size := len(jsonbytes)
		 wmoswriter.WriteAt(jsonbytes, wmo_offset_track)
		 json_start_point = wmo_offset_track
		 wmo_offset_track += json_size
		 
	 //read everything that is in 'wmo.stringz', make it binary, then 
	 //add to this.stringz at appropriate offset.
		tmpStr := ''
		for _,oneStr := range wmo.stringz { 
			 tmpStr = tmpStr+" "+oneStr
		} 
		tmpStrBytes = []byte(tmpStr) 
		all_strings_size := len(tmpStrBytes)
		wmoswriter.WriteAt(tmpStrBytes, wmo_offset_track)
		stringz_start_point = wmo_offset_track
		wmo_offset_track += all_strings_size 
		
		//==writing headers==
		//now write files header info to the wmo.BinaryData
	  filez_start_point = 34 //28+6=34, offset: beginning of 34th byte
	  wmo_offset_track = filez_start_point
	  files_ofst_track = 0
	  wmoswriter.WriteAt(numberoffilez, wmo_offset_track) //numberoffilez is the value writen at offset 34 = filez_start_point
	  files_ofst_track += 1   //now, files_ofst_track = 1
	  wmo_offset_track += 1   //now, wmo_offset_track = 35
	  //now write sizes and offsets of each file previously written wmo.BinaryData
	  i=0
	  for i < numberoffilez {
		  wmoswriter.WriteAt(file_offsets[i], wmo_offset_track)
		  wmoswriter.WriteAt(file_sizes[i], (wmo_offset_track+(numberoffilez*4)))
		  files_ofst_track += 4 
		  wmo_offset_track += 4
		  i++
	  } 
		
		//reset the wmo_offset_track
	  wmo_offset_track=0
		
		//write all wmo header attributes at their known file_offsets
	  wmoswriter.WriteAt(34, wmo_offset_track)  //34 => filez_start_point
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(size_of_files_header, wmo_offset_track)
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(total_size_of_files, wmo_offset_track)
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(json_start_point, wmo_offset_track)
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(json_size, wmo_offset_track)
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(stringz_start_point, wmo_offset_track)
	  wmo_offset_track += 4
	  wmoswriter.WriteAt(all_strings_size, wmo_offset_track)
	  wmo_offset_track += 4
		  
		//TODO: do some house keeping proofs of the offsets and sizes, since some of 
		//the attributes have more than one sources of these values; so we can campare to prove.    

}

//===================================================================================================



/*
 * 
 golang how to determine length of an array
 golang how to push to an array:  
	  a := []int{1,2,3}
	  a = append(a, 4)
	  fmt.Println(a)
	  append(a[:3], 5)
	  fmt.Println(a)
	//----
		package main

		import (
		    "fmt"
		)

		var a = make([]int, 7, 8)
		// A slice is a descriptor of an array segment. 
		// It consists of a pointer to the array, the length of the segment, and its capacity (the maximum length of the segment).
		// The length is the number of elements referred to by the slice.
		// The capacity is the number of elements in the underlying array 
		//(beginning at the element referred to by the slice pointer).
		// |-> Refer to: https://blog.golang.org/go-slices-usage-and-internals -> "Slice internals" section

		func Test(slice []int) {
		    // slice receives a copy of slice `a` which point to the same array as slice `a`
		    slice[6] = 10
		    slice = append(slice, 100)
		    // since `slice` capacity is 8 & length is 7, it can add 100 and make the length 8
		    fmt.Println(slice, len(slice), cap(slice), " << Test 1")
		    slice = append(slice, 200)
		    // since `slice` capacity is 8 & length also 8, slice has to make a new slice 
		    // - with double of size with point to new array (see Reference 1 below).
		    // (I'm also confused, why not (n+1)*2=20). But make a new slice of 16 capacity).
		    slice[6] = 13 // make sure, it's a new slice :)
		    fmt.Println(slice, len(slice), cap(slice), " << Test 2")
		}

		func main() {
		    for i := 0; i < 7; i++ {
		        a[i] = i
		    }

		    fmt.Println(a, len(a), cap(a))
		    Test(a)
		    fmt.Println(a, len(a), cap(a))
		    fmt.Println(a[:cap(a)], len(a), cap(a))
		    // fmt.Println(a[:cap(a)+1], len(a), cap(a)) -> this'll not work
		}
//---- 

 */

