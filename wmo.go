package websocketMessageObject

import (
"fmt"
"os"
"io"
"log"
"bytes" 
"strings"
"encoding/binary"
//"encoding/json"
//"github.com/h2non/filetype"
)
 
//=============================================================
type WebsocketMessageObject struct {
	BinaryData []byte 
	filezLength uint32
	jsonnLength uint32
	stringzLength uint32
	filez []*os.File   //array of files. functions will keep appending until Encode() is called
	Fdata []byte
	jsonn interface{} //variable to hold json object. A function will set it before Encode() is called
	stringz []string  //array of strings. functions will keep appending until Encode() is called
}

func NewWebsocketMessageObject() WebsocketMessageObject {
	return WebsocketMessageObject{filezLength:0, jsonnLength:0, stringzLength:0 }
}

func NewWMO() WebsocketMessageObject {
	return WebsocketMessageObject{filezLength:0, jsonnLength:0, stringzLength:0 }
}
//=============================================================
//----------------
type wmoSectionReader struct {
	reader *bytes.Reader
}

func newWmoSectionReader(dataFromClient []byte) *wmoSectionReader {
	bw := new(bytes.Buffer)
	_ = binary.Write(bw,binary.BigEndian,dataFromClient)
	breader := bytes.NewReader(bw.Bytes())
	return &wmoSectionReader{reader: breader}
}

func (wsr *wmoSectionReader) ReadAt(dat []byte,offset int64) (int, error) {
	n, err := wsr.reader.ReadAt(dat,offset)
	if err != nil {
		return n, err
	}
	return n, nil
}
//----------------
type wmoSectionWriter struct { 
	writerat io.WriterAt
}

func newWmoSectionWriter(writerat io.WriterAt) wmoSectionWriter { 
	return wmoSectionWriter{writerat: writerat}
}

//WriteAt writes len(dat) bytes from 'dat' to the underlying data stream (eg tempBuf) at offset 'off'.
func (wsw *wmoSectionWriter) WriteAt(dat []byte,offset int64) (int, error) {
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
	Endianness uint8
} //total header size = (4*7)+1 = 29 bytes = 29*8 = 240 bits = 0.29 KB

type wmoFilesHeader struct {
	NumberOfFiles uint8
	EachFileOffset []uint32
	EachFileSize []uint32
	EachFileType []uint32
}
//----
type wmoHeaderReader struct {
	sreader io.SectionReader
} 
func newWmoHeaderReader(dataFromClient []byte) wmoHeaderReader {
	//create a new Buffer
	var tempheaderbuf bytes.Buffer
	//define a write func for the above writer
	_ = binary.Write(&tempheaderbuf,binary.LittleEndian,dataFromClient)
	reader := bytes.NewReader(tempheaderbuf.Bytes())
	sr := io.NewSectionReader(reader, 0, 29)
	return wmoHeaderReader{sreader: *sr}
} 
func (whr *wmoHeaderReader) Read(hd []byte) (int, error) {
	n, err := whr.sreader.Read(hd)
	if err != nil {
		return n, err
	}
	return n, nil
}
func (whr *wmoHeaderReader) ReadAt(hd []byte,offst int64) (int, error) {
	n, err := whr.sreader.ReadAt(hd,offst)
	if err != nil {
		return n, err
	}
		return n, nil
}

  func readWmoHeader(dataFromClient []byte) wmoHeader {
		var wmoheader wmoHeader
		headerBuf := bytes.NewBuffer(dataFromClient[:29])
		headerBufReader := bytes.NewReader(headerBuf.Bytes())
		err := binary.Read(headerBufReader, binary.LittleEndian, &wmoheader)
		if err != nil {
			log.Fatal("readWmoHeader: binary.Read failed; ", err)
		}
		return wmoheader
	}
	

//------------------Decoders----------
func (wmo *WebsocketMessageObject) DecodeJson(dataFromClient []byte) string {
	hedr := readWmoHeader(dataFromClient)
	jsonReader := bytes.NewBuffer(dataFromClient[hedr.JsonOffset:(hedr.JsonOffset+hedr.JsonSize)])
	jsonstr := string(jsonReader.Bytes())
	return jsonstr 
}

func (wmo *WebsocketMessageObject) ReadFilesBytes(dataFromClient []byte) ([]byte,uint32) {
	hedr := readWmoHeader(dataFromClient)
	fReader := bytes.NewBuffer(dataFromClient[hedr.FilesHeaderOffset:(hedr.FilesHeaderOffset+hedr.FilesTotalSize)])
	fBuf := make([]byte, hedr.FilesTotalSize)
	_ = binary.Read(fReader,binary.BigEndian,&fBuf)
	return fBuf, hedr.FilesHeaderSize
}

func (wmo *WebsocketMessageObject) DecodeFiles(dataFromClient []byte) []*os.File {
	filesBytes, fheaderSize := wmo.ReadFilesBytes(dataFromClient)
	var wmofheader wmoFilesHeader  
	readr := bytes.NewReader(filesBytes[:fheaderSize])
	_ = binary.Read(readr, binary.LittleEndian, &wmofheader) //TODO: identify and use the endianness that was used by client to encode data instead of hardcoding this to 'binary.LittleEndian'
	//-------- 
	var allFiles []*os.File 
	i:=0
	for i<int(wmofheader.NumberOfFiles) {
		var oneFile *os.File
		oneFileBuf := bytes.NewBuffer(filesBytes[wmofheader.EachFileOffset[i]:(wmofheader.EachFileOffset[i]+wmofheader.EachFileSize[i])]) 
		oneFileBufReader := bytes.NewReader(oneFileBuf.Bytes())
		_ = binary.Read(oneFileBufReader, binary.LittleEndian, oneFile)  //TODO: identify and use the endianness that was used by client to encode data instead of hardcoding this to 'binary.LittleEndian'
		allFiles = append(allFiles, oneFile)
		i++
	}
	return allFiles
} 

func (wmo *WebsocketMessageObject) DecodeStringAll(dataFromClient []byte) map[string]string {
	hedr := readWmoHeader(dataFromClient) 
	strReader := bytes.NewBuffer(dataFromClient[hedr.StringsOffset:(hedr.StringsOffset+hedr.StringsSize)]) 
	stringz := string(strReader.Bytes())
	//split the 'stringz' content using dilimitors and create a map of strings with string keys.
	strArray := strings.Fields(stringz) //splitting using spaces: stringz = 'key1-value1 key2-value2 key3-value3 ...'
	newStrMap := map[string]string{}
	for _, keyAndVal := range strArray {
		strkey := strings.Split(keyAndVal,"-")[0]
		strval := strings.Split(keyAndVal,"-")[1]
		newStrMap[strkey] = strval
	}
	return newStrMap 
}

func (wmo *WebsocketMessageObject) DecodeString(dataFromClient []byte,strkey string) string {
	newStrMap := wmo.DecodeStringAll(dataFromClient) 
	return newStrMap[strkey]
}

//------------------Encoders----------

func (wmo *WebsocketMessageObject) setFilesSize(size uint32){ 
	wmo.filezLength=size 
}
func (wmo *WebsocketMessageObject) setJsonSize(size uint32){ 
	wmo.jsonnLength=size
}
func (wmo *WebsocketMessageObject) setStringsSize(size uint32){ 
	wmo.stringzLength=size 
}
func (wmo *WebsocketMessageObject) addToFilesSize(size uint32){ 
	wmo.filezLength+=size 
}                 
func (wmo *WebsocketMessageObject) addToJsonSize(size uint32){ 
	wmo.jsonnLength+=size
}                    
func (wmo *WebsocketMessageObject) addToStringsSize(size uint32){ 
	wmo.stringzLength+=size
}
func (wmo *WebsocketMessageObject) PrintFilesSize(){ 
	fmt.Printf("wmo.filezLength: = %d",wmo.filezLength)
}
func (wmo *WebsocketMessageObject) PrintJsonSize(){ 
	fmt.Printf("wmo.jsonnLength: = %d",wmo.jsonnLength)
}
func (wmo *WebsocketMessageObject) PrintStringsSize(){ 
	fmt.Printf("wmo.stringzLength: = %d",wmo.stringzLength)
}

func (wmo *WebsocketMessageObject) AddFile(file *os.File) {
	wmo.filez = append(wmo.filez, file)
	fileinfo, _ := file.Stat()
	file_size := fileinfo.Size()
	wmo.addToFilesSize(uint32(file_size))
}

func (wmo *WebsocketMessageObject) AddFileFrom(filepathh string) {
	file, err := os.Open(filepathh)
	if err != nil {
		log.Fatal("Error while opening file by wmo.AddFileFrom() function", err)
	}
	defer file.Close()
	wmo.filez = append(wmo.filez, file)
	//-----getting size of the file just appended------
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := fileinfo.Size() //read the size of 'file'
	someExtra := int64(20) 
	bytes := make([]byte, (size+someExtra))
	n, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	} 
	file_size := uint32(n)
	wmo.addToFilesSize(file_size)
} 

func (wmo *WebsocketMessageObject) AddJson(myjson interface{}) {
	 wmo.jsonn = myjson
	 	 //jsonnstr, err := json.Marshal(wmo.jsonn) 
	 	 jsonbytes := wmo.jsonn.([]uint8)
	 	 jsonnstr := string(jsonbytes) 
	 json_size := bytes.Count([]byte(jsonnstr),nil)-1
	 wmo.setJsonSize(uint32(json_size))
}

func (wmo *WebsocketMessageObject) AddString(keyy string,mystring string) {
	str := fmt.Sprintf("%s-%s",keyy,mystring)
	wmo.stringz = append(wmo.stringz, str)
	str_size := bytes.Count([]byte(str),nil)-1
	wmo.addToStringsSize(uint32(str_size)) 
}

func Uint32toBinary(uint32num uint32) []byte {
	var buf bytes.Buffer
	mywriter := io.MultiWriter(&buf)
	_ = binary.Write(mywriter,binary.BigEndian,uint32num)
	return buf.Bytes()
}
func Uint8toBinary(uint8num uint8) []byte {
	var buf bytes.Buffer
	mywriter := io.MultiWriter(&buf)
	_ = binary.Write(mywriter,binary.BigEndian,uint8num)
	return buf.Bytes()
}
 
func (wmo *WebsocketMessageObject) Encode() {
	/*
              
	0                    28   31               files-content         json                            strings
	|--,--,--,--,--,--,--|-|~~|--,--,...,--,...|-----,-------,---,...|-------------------------------|-------|
	wmo-header                files-header
	                          0                                      #
	                          |<-----------------files-------------->|
	*/ 
	
	var filez_start_point uint32
	var json_start_point uint32
	var stringz_start_point uint32
	 
	var file_offset_readers []io.Reader
	var file_size_readers []io.Reader
	var file_type_readers []io.Reader
	var file_data_readers []io.Reader
	  
	var total_size_of_files uint32
	total_size_of_files = 0 
	  
	var wmo_offset_track uint32
	var files_ofst_track uint32
	wmo_offset_track = 31  //28 bytes (4*7) for the wmo header,1 for endianness, 2 for 'dont care' bytes, 1 byte for holding number of files
	files_ofst_track = 0  
	//numberoffilez := uint32(len(wmo.filez))
	numberoffilez := uint32(countSupportedFilesOnly(wmo.filez))
	var size_of_files_header uint32
	size_of_files_header=1+(numberoffilez*(4+4+4)) //1 byte holds No.of files, 4 bytes(size of uint32) for each file_sizes[element] and 4 bytes for each file_offsets[element]
	total_size_of_files += size_of_files_header
	wmo_offset_track =31 + size_of_files_header  //increment the '*_track' by 'size_of_files_header'
	files_ofst_track = 0 + size_of_files_header 
	   
	for _, fi := range wmo.filez {
		fileinfo, err := fi.Stat()
		if err != nil {
			log.Fatal(err)
		}
		size := fileinfo.Size() //read the size of 'file' 
		fbytes := make([]byte, size)
		n, err := fi.Read(fbytes)
		if err != nil {
			log.Fatal(err)
		} 
		//create an io reader with data fbytes[:n]
		fbytesreader := bytes.NewReader(fbytes) 
		file_data_readers=append(file_data_readers,fbytesreader)
		fiLength:= uint32(n) 
		size_reader := bytes.NewReader(Uint32toBinary(fiLength)) 
		file_size_readers=append(file_size_readers,size_reader)
		offset_reader := bytes.NewReader(Uint32toBinary(uint32(files_ofst_track)))
		file_offset_readers=append(file_offset_readers,offset_reader) 
		typ := getFileType(fi)
		type_reader := bytes.NewReader(Uint32toBinary(uint32(findFileTypeIndex(typ))))
		file_type_readers=append(file_type_readers,type_reader)
		total_size_of_files += fiLength
		files_ofst_track += fiLength
		wmo_offset_track += fiLength
	}
	var filesDataMultiReader io.Reader
	var fileOffsetsMultiReader io.Reader
	var fileSizesMultiReader io.Reader
	var fileTypesMultiReader io.Reader
	var no_of_files_reader io.Reader
	var allFilesDataReader io.Reader
	if numberoffilez > 0 {
		filesDataMultiReader = io.MultiReader(file_data_readers...)
		fileOffsetsMultiReader=io.MultiReader(file_offset_readers...)
		fileSizesMultiReader = io.MultiReader(file_size_readers...)
		fileTypesMultiReader = io.MultiReader(file_type_readers...)
		no_of_files_reader = bytes.NewReader(Uint8toBinary(uint8(numberoffilez)))
		allFilesDataReader = io.MultiReader(no_of_files_reader,fileOffsetsMultiReader,fileSizesMultiReader,fileTypesMultiReader,filesDataMultiReader)
	}else{
		no_of_files_reader = bytes.NewReader(Uint8toBinary(uint8(numberoffilez))) 
		newreadr := bytes.NewBuffer(wmo.Fdata)
		if len(newreadr.Bytes()) > 0 {
		  allFilesDataReader = newreadr 
		  nofbyteReader := bytes.NewReader(wmo.Fdata[:1]) //byte bits representing number of files
		  var nof uint8
		  _ = binary.Read(nofbyteReader,binary.BigEndian,nof) // decode into a uint8 number
		  size_of_files_header = uint32(uint8(1)+(nof*uint8(4+4+4))) //calculate size of files header using number of files value 
		}else{
		    allFilesDataReader = no_of_files_reader
		} 
		total_size_of_files = uint32(len(wmo.Fdata))
		if total_size_of_files > 1 {
		    wmo_offset_track += (total_size_of_files-1)
		}else{
		    total_size_of_files = 1  //must not be less than this.
		}  
	}
	//read everything that is in 'wmo.jsonn', stringfy & make it binary, also create a reader out of it
 	jsonSlice := wmo.jsonn.([]uint8)
 	jsonnstr := string(jsonSlice)
	jsonbytes := []byte(jsonnstr)
	json_size := uint32(len(jsonbytes)) 
	json_data_reader := bytes.NewReader(jsonbytes) 
	json_start_point = wmo_offset_track
	wmo_offset_track += json_size
	 
	//read everything that is in 'wmo.stringz', make it binary, then 
	//add to this.stringz at appropriate offset. 
	tmpStr := " "
	for _,oneStr := range wmo.stringz { 
		 tmpStr = tmpStr+" "+oneStr 
		 log.Println("encode: tmpStr = ", tmpStr )
	}  
	tmpStrBytes := []byte(tmpStr)
	all_strings_size := uint32(len(tmpStrBytes))
	strings_data_reader := bytes.NewReader(tmpStrBytes) 
	stringz_start_point = wmo_offset_track
	wmo_offset_track += all_strings_size 
		
	//==writing headers== 
	filez_start_point = 31 //28+1+2=31, offset: beginning of 31st byte
	wmo_offset_track = filez_start_point
	 
	//readers for all wmo header values 
	start_of_files_reader := bytes.NewReader(Uint32toBinary(uint32(31))) //31 => filez_start_point  
	size_of_files_header_reader := bytes.NewReader(Uint32toBinary(uint32(size_of_files_header))) 
	total_size_of_files_reader := bytes.NewReader(Uint32toBinary(uint32(total_size_of_files)))  
	json_start_point_reader := bytes.NewReader(Uint32toBinary(uint32(json_start_point)))  
	json_size_reader := bytes.NewReader(Uint32toBinary(uint32(json_size))) 
	stringz_start_point_reader := bytes.NewReader(Uint32toBinary(uint32(stringz_start_point)))  
	all_strings_size_reader := bytes.NewReader(Uint32toBinary(uint32(all_strings_size)))  
	allWmoHeadersReader := io.MultiReader(start_of_files_reader,size_of_files_header_reader,total_size_of_files_reader,json_start_point_reader, json_size_reader, stringz_start_point_reader,all_strings_size_reader)
	 
	endianness := bytes.NewReader(Uint8toBinary(uint8(0)))
	dontcarepart1 := bytes.NewReader(Uint8toBinary(uint8(0))) 
	dontcarepart2 := bytes.NewReader(Uint8toBinary(uint8(0)))
	dontCareBytesReader := io.MultiReader(endianness,dontcarepart1,dontcarepart2)  
	allDataReader := io.MultiReader(allWmoHeadersReader,dontCareBytesReader,allFilesDataReader,json_data_reader,strings_data_reader)
	   
	var tempBuffer bytes.Buffer
	tempBufferReader := &tempBuffer 
	_,_ = io.Copy(tempBufferReader,allDataReader)
	//encode our data into the wmo.BinaryData from tempBufferReader, being able to specify the Endianness. 
	fulDataLength := tempBufferReader.Len()
	fixedSizebuffer := make([]byte, fulDataLength)
	_ = binary.Read(tempBufferReader, binary.BigEndian, &fixedSizebuffer)
	wmo.BinaryData = fixedSizebuffer 
	  
    log.Println("wmo.BinaryData = ",wmo.BinaryData) 
    //TODO: empty all the other wmo variables except wmo.BinaryData. ie, empty: wmo.filez, wmo.stringz, wmo.jsonn, etc
    //...
       
	//DONE !
	//We are now ready to pass wmo.BinaryData to websocket.send() function. 
	//ie, ws.send(wmo.BinaryData) . This sends the encoded data to the client in plain binary form
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

