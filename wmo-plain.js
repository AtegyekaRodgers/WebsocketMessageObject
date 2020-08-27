
//+++++++++++++++++++++++++++++++++++++++++ wmo start +++++++++++++++++++++++++++++++++++++++++++++++
let wmoHeader = {
	FilesHeaderOffset:0,
	FilesHeaderSize:0,
	FilesTotalSize:0,
	JsonOffset:0,
	JsonSize:0,
	StringsOffset:0,
	StringsSize:0
}; 

function readWmoHeader(dataFromServer){
	var BIGendian = false;
	var LITTLEendian = true;                          //  0                    28   
	let wmoHeaderBinary = Uint32Array(dataFromServer,0, 8); //  |--,--,--,--,--,--,--|~~~ //NOTE:assumption is, dataFromServer came as an arrayBuffer
	let wmoheader = {
		FilesHeaderOffset:Number(wmoHeaderBinary.getUint32(0,BIGendian)),
		FilesHeaderSize:Number(wmoHeaderBinary.getUint32(1,BIGendian)),
		FilesTotalSize:Number(wmoHeaderBinary.getUint32(2,BIGendian)),
		JsonOffset:Number(wmoHeaderBinary.getUint32(3,BIGendian)),
		JsonSize:Number(wmoHeaderBinary.getUint32(4,BIGendian)),
		StringsOffset:Number(wmoHeaderBinary.getUint32(5,BIGendian)),
		StringsSize:Number(wmoHeaderBinary.getUint32(6,BIGendian))
	};
	return wmoheader;
}

function readFilesHeader(filesBytes, hedrSize){
	var BIGendian = false;
	var LITTLEendian = true;
	let wmofHeaderBinary = Uint32Array(filesBytes,0,(hedrSize/4));
	let noOfFiles = Number(wmofHeaderBinary.getUint32(0,BIGendian));
	let filezheader = {
		NumberOfFiles:noOfFiles,
		FilesOffsets:Array.from(Uint32Array(filesBytes,1,noOfFiles)),
		FilesSizes:Array.from(Uint32Array(filesBytes,(noOfFiles+1),noOfFiles))
	};
	return filezheader;
}

class WebsocketMessageObject{
	constructor(objectname){ 
		objectname = objectname || "wmo";
		this.Objectname = objectname;
		this.BinaryData = null //to be defined later as arrayBuffer with appropriate size 
		this._filezLength = 0;
		this._jsonnLength = 0;
		this._stringzLength = 0; 
		this.filez = []; //array of files. functions will keep appending until Build() is called
		this.jsonn = {}; //variable to hold jsob object. A function will set it before Build() is called
		this.stringz = []; //array of strings. functions will keep appending until Build() is called  
	} 
		//--------Encoders--------
	  setFilesSize(size){ this._filezLength=size; }
	  setJsonSize(size){ this._jsonnLength=size; }
	  setStringsSize(size){ this._stringzLength=size; } 
	  
	  addToFilesSize(size){if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._filezLength+=size; }else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToFilesSize()' expects a number ");} }                 
	  addToJsonSize(size){if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._jsonnLength+=size;}else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToJsonSize()' expects a number. ");} }                    
	  addToStringsSize(size){ if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._stringzLength+=size;}else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToStringsSize()' expects a number ");} }  
	   
    AddFile(file) {  
		 if(file.name){
			 this.filez.push(file);  
			 let file_size = file.size; 
			 this.addToFilesSize(file_size);
		 }
	 };
	 
	 AddFileFrom = function(fileInputId) {
		 const file = document.getElementById(fileInputId).files[0]; 
		 if(file.name){
			 this.filez.push(file);  
			 let file_size = file.size; 
			 this.addToFilesSize(file_size);
		 } 
	 } 
	 
	 AddJson = function(myjson) {  
		 this.jsonn = myjson;
		 //determine the size of the json we have just added, then increment the _jsonnLength by that size. 
		 let json_size = JSON.stringify(this.jsonn).length;
		 this.setJsonSize(json_size);
	 };
	 
	 AddString = function(keyy,mystring) {  
			if (typeof(mystring)==='string') {
				//append the string to stringz[] array of this object. 
				var str = keyy+'-'+mystring;
				this.stringz.push(str);
				//determine the size of the string we have just added, then increment the _stringzLength by that size.
				let str_size = str.length;
				this.addToStringsSize(str_size);
			}else{
				this._error(new Error('The function AddString() expects a string'));
				return;
			}  
	 };
	 
	 AddStringFrom = function(keyy, textInputId) { 
		 const strng = document.getElementById(textInputId).value;  
		 if (typeof(strng)==='string') {
			 //append the string to stringz[] array of this object. 
			 var str = keyy+'-'+strng;
			 this.stringz.push(str);
			 //determine the size of the string we have just added, then increment the _stringzLength by that size.
			 let str_size = str.length;
			 this.addToStringsSize(str_size);
		 }else{
			 this._error(new Error('The function AddStringFrom() expects a string'));
			 return;
		 }  	
	 };
	 /*
	 AddTextFrom = function(keyy, elementId) { 
		 const innerTextt = document.getElementById(elementId).textContent || "testString";  
		 if (typeof(elementId)==='string') {  
			 var str = keyy+'-'+innerTextt;
			 this.stringz.push(str); 
			 let str_size = str.length;
			 this.addToStringsSize(str_size);
		 }else{
			 this._error(new Error('The function AddTextFrom() expects a string as id of an html element'));
			 return;
		 }  	
	 }; */
	 AddNumber = function(keyy,mynum) {  
		 if (typeof(Number(mynum))==='integer') { 
			 var numstr = keyy+'-'+mynum;
			 this.numberz.push(numstr);
			 //determine the size of the string we have just added, then increment the _stringzLength by that size.
			 let numstr_size = numstr.length;
			 this.addToStringsSize(numstr_size);
		 }else{
			 this._error(new Error('The function AddNumber() expects a number'));
			 return;
		 }  
	 };    
	 
	 Encode = function() {
			/*
		              
			0                    28   31               files-content         json                            strings
			|--,--,--,--,--,--,--|-|~~|--,--,...,--,...|-----,-------,---,...|-------------------------------|-------|
			wmo-header                files-header
			                          0                                      #
			                          |<-----------------files-------------->|
			*/ 
			
			/* You can use the following function to determine the endianness of a platform. 
			const BIG_ENDIAN = Symbol('BIG_ENDIAN');
			const LITTLE_ENDIAN = Symbol('LITTLE_ENDIAN');
			function getPlatformEndianness() {  
				let arr32 = Uint32Array.of(0x12345678);
				let arr8 = new Uint8Array(arr32.buffer);
				switch ((arr8[0]*0x1000000) + (arr8[1]*0x10000) + (arr8[2]*0x100) + (arr8[3])) {
					case 0x12345678:
						return BIG_ENDIAN;
					case 0x78563412:
						return LITTLE_ENDIAN;
					default:
						throw new Error('Unknown endianness');
				 }
			} 
			//----------
				//convert arraybuffer to blob
				var array = new Uint8Array([0x04, 0x06, 0x07, 0x08]); 
				var blob = new Blob([array]);
			//----------
			 */
			var BIGendian = false;
			var LITTLEendian = true;
			
			//alternative way to check the endianness of this machine
			var isLittleEndian = (function() {
				var buffer = new ArrayBuffer(2);
				new DataView(buffer).setInt16(0, 256, true); //true -> littleEndian
				// Int16Array uses the platform's endianness.
				return new Int16Array(buffer)[0] === 256;
			})();
			var logmessage = isLittleEndian?"The endianness of this machine is : LittleEndian":"The endianness of this machine is : bigEndian";
			console.log(logmessage);
			
			this.AddString("endstrings","rightpadding");
			var filez_start_point; 
			var json_start_point;
			var stringz_start_point; 
		   var file_offsets = [];
		   var file_sizes = [];
		   var total_size_of_files = 0;
			  
			  var wmo_offset_track = 31;  //28 bytes (4*7) for the wmo header,1 for emdianness, 2 for 'dont care' bytes, 1 byte for holding number of files
			  var files_ofst_track = 0;   //8 bits == 1 byte for holding number of files, starting from the end of wmo header 
			  var numberoffilez = this.filez.length;
			  var size_of_files_header=1+(numberoffilez*(4+4)); //1 byte holds No.of files, 4 bytes(size of uint32) for each file_sizes[element] and 4 bytes for each file_offsets[element]
			  
			  total_size_of_files += size_of_files_header;
			  wmo_offset_track =31 + size_of_files_header;  //increment the '*_track' by 'size_of_files_header'
			  files_ofst_track = 0 + size_of_files_header;
			  
			  //ready to create the BinaryData
			  this.BinaryData = new ArrayBuffer(31+size_of_files_header+this._filezLength+this._jsonnLength+this._stringzLength+8);  //8 is some extra just in case we need it.
			  console.log("total size = "+(31+size_of_files_header+this._filezLength+this._jsonnLength+this._stringzLength+4));
			  let mainHeaderView = new Uint32Array(this.BinaryData, 0, 7);  //Uint32Array(buffer, offset, size); where 'size' is the number of items with the specific size eg 32 bits in this case.
			  let endiannessView = new Uint8Array(this.BinaryData, 28, 1);
			  let dontCareView = new Uint8Array(this.BinaryData, 29, 2);
			  let noOfFilesView = new Uint8Array(this.BinaryData, 31, 1);
			  let fileOffsetsView = new Uint32Array(this.BinaryData, 32, numberoffilez);
			  let fileSizesView = new Uint32Array(this.BinaryData, ((4*numberoffilez)+32), numberoffilez);
				  let newOffstt = (32+(4*numberoffilez))+(4*numberoffilez);
			  let filesDataView = new Uint8Array(this.BinaryData, newOffstt, this._filezLength);
				let filesDataStart = newOffstt;
				wmo_offset_track=filesDataStart;
			  let jsonDataView = new Uint8Array(this.BinaryData,(filesDataStart+this._filezLength), this._jsonnLength);
			  console.log("jsonDataView.byteLength = "+jsonDataView.byteLength);
			  let jsonDataStart = (filesDataStart+this._filezLength);
			  let stringsDataView = new Uint8Array(this.BinaryData,(jsonDataStart+this._jsonnLength), this._stringzLength);
			  console.log("strings offset = "+(jsonDataStart+this._jsonnLength)+" - size="+(this._stringzLength));
			   
			   //loop through files to get info about each of them,
				//and write each into wmo.BinaryData.
			  wmo_offset_track=filesDataStart;
			  var fdataOffsett=0;
			  for (var i=0; i<numberoffilez; i++) {
					//transform a file into arraybuffer
					const freadr = new FileReader();
					freadr.readAsArrayBuffer(this.filez[i]);
					//create a view of the arraybuffer
					let fileBytesView = new Uint8Array(freadr.result);
					//determine the bytelength of the typed array 
					var fiLength = fileBytesView.byteLength; 
				  file_offsets.push(files_ofst_track);
				  file_sizes.push(fiLength);
				  filesDataView.set(fileBytesView,fdataOffsett); 
				  total_size_of_files += fiLength;
				  files_ofst_track += fiLength;
				  wmo_offset_track += fiLength;
				  fdataOffsett+=fiLength;
			  }
			  
			 //read everything that is in 'this.jsonn', stringfy & make it binary, 
			 //then add to this.BinaryData at appropriate offset.
			 wmo_offset_track=jsonDataStart;
			 json_start_point = wmo_offset_track; 
			 let jsonnstr = JSON.stringify(this.jsonn); 
			 let json_size = jsonnstr.length;
			 for (var jchr in jsonnstr){
				 var ascii = jsonnstr.charCodeAt(jchr);
				 jsonDataView[jchr]=ascii;
			 }
			 wmo_offset_track += json_size;
			 //wmo_offset_track+=this._jsonnLength; //alternative to line above
			 
			 //read everything that is in 'wmo.stringz', make it binary, then 
			 //add to this.stringz at appropriate offset.
			let tmpStr = "";
			for (var strIndex in this.stringz) {
				var oneStr = this.stringz[strIndex];
				tmpStr = tmpStr+" "+oneStr; 
			}
			var strAsciiValues = [];
			for (var chr in tmpStr){
				var ascii = tmpStr.charCodeAt(chr); 
				stringsDataView[chr]=ascii;
			}
			let all_strings_size = tmpStr.length+2;
			stringz_start_point = wmo_offset_track;
			wmo_offset_track += all_strings_size;
			
			//==writing headers==
			//now write files header info to the this.BinaryData
		  filez_start_point = 31; //28+1+2=31, offset: beginning of 31st byte 
		  //noOfFilesView.setUint8(0, numberoffilez, isLittleEndian);
		  noOfFilesView.set([numberoffilez],0);
		  //now write sizes and offsets of each file previously written to this.BinaryData
		  for (var i in file_sizes) {
			  let thisFilesize = file_sizes[i];
			  let thisfdataOffset = file_offsets[i];
			  //fileSizesView.setUint32(fheadoffset, thisFilesize, isLittleEndian);
			  fileSizesView.set([thisFilesize],i);
			  //fileOffsetsView.setUint32(fheadoffset,thisfdataOffset, isLittleEndian);
			  fileOffsetsView.set([thisfdataOffset],i);
		  } 
			
		  //write all wmo header attributes at their known file_offsets 
		  //mainHeaderView.setUint32(0, filez_start_point, isLittleEndian); 
		  console.log("header: filez_start_point = ",filez_start_point);
		  mainHeaderView.set([filez_start_point],0);
		  //mainHeaderView.setUint32(1, size_of_files_header, isLittleEndian); 
		  console.log("header: size_of_files_header = ",size_of_files_header);
		  mainHeaderView.set([size_of_files_header],1);
		  //mainHeaderView.setUint32(2, total_size_of_files, isLittleEndian);
		  console.log("header: total_size_of_files = ",total_size_of_files);
		  mainHeaderView.set([total_size_of_files],2);
		  //mainHeaderView.setUint32(3, json_start_point, isLittleEndian); 
		  console.log("header: json_start_point = ",json_start_point);
		  mainHeaderView.set([json_start_point],3);
		  //mainHeaderView.setUint32(4, json_size, isLittleEndian);
		  console.log("header: json_size = ",json_size);
		  mainHeaderView.set([json_size],4);
		  //mainHeaderView.setUint32(5, stringz_start_point, isLittleEndian); 
		  console.log("header: stringz_start_point = ",stringz_start_point);
		  mainHeaderView.set([stringz_start_point],5);
		  //mainHeaderView.setUint32(6, all_strings_size, isLittleEndian);
		  console.log("header: all_strings_size = ",all_strings_size);
		  mainHeaderView.set([all_strings_size],6);
		  //endiannessView.setUint8(0,isLittleEndian?6:112,isLittleEndian)
		  console.log("header: (isLittleEndian?6:112) = ",(isLittleEndian?6:112));
		  endiannessView.set([(isLittleEndian?6:112)],0);
		  //littleEndian: 00000110 = 6 , bigEndian: 01110000 = 112
		}
	   
	  toString=function() {
		  return '[object WebsocketMessageObject]';
	  }
	 
		//------------------Decoders----------
		DecodeJson = (dataFromServer) => { 
			let wmoheader = readWmoHeader(dataFromServer);
			let jsonvieww = Uint8Array(dataFromServer,wmoheader.JsonOffset, wmoheader.JsonSize);
			let jsonstr = jsonvieww.join(''); //Converts all elements of the jsonvieww to strings and concatenates them, separated by the specified separator, '' in this case.          
			let jsonObject = JSON.parse(jsonstr);
			return jsonObject;
		}

		ReadFilesBytes = (dataFromServer) => {
			let hedr = readWmoHeader(dataFromServer); 
			let filesDataView = Uint8Array(dataFromServer,hedr.FilesHeaderOffset, hedr.FilesTotalSize); 
			return filesDataView;
		}

		DecodeFiles = (dataFromServer) => {
			let hedr = readWmoHeader(dataFromServer);
			let filesBytes = this.ReadFilesBytes(dataFromServer);
			let files_hedr = readFilesHeader(filesBytes, hedr.FilesHeaderSize);
			//let allFiles = [];
			let imageUrls = [];
			//let oneFileView = Uint8Array(filesBytes,files_hedr.Offsets[i], files_hedr.Sizes[i]);
			//From the typedarray 'oneFileView', decode a single file and acquire its File boject or image URL
			/*
			 The File constructor (as well as the Blob constructor) takes an array of parts. 
			 A part doesn't have to be a DOMString. It can also be a Blob, File, or a typed array. 
			 You can easily build a File out of a Blob like this:
			 let file = new File([blob], "filename");
			 //---
			 var file = new File([blob], "my_image.png",{type:"image/png", lastModified:new Date()})
			 //--- 
			 
			 var blob = new Blob( [ oneFileView ], { type: "image/jpeg" } );
			 var urlCreator = window.URL || window.webkitURL;
			 var imageUrl = urlCreator.createObjectURL( blob );
			 //then push that image Url to the imageUrls array to be returned
			 imageUrls.push(imageUrl);
			 //var img = document.querySelector( "#photo" ); 
			 //img.src = imageUrl;
			 urlCreator.revokeObjectURL(); 
			 
			//return allFiles;
			return imageUrls; */
		} 

		DecodeStringAll = (dataFromServer) => {
			let wmoheader = readWmoHeader(dataFromServer);
			let stringsview = Uint8Array(dataFromServer,wmoheader.StringsOffset, wmoheader.StringsSize);
			let stringz = stringsview.join(''); //Converts all elements of the stringsview to strings and concatenates them, separated by the specified separator, '' in this case.          
			//split the 'stringz' content using dilimitors and create a map of strings with string keys.
			let strArray = stringz.split(" "); //splitting using spaces: stringz = 'key1-value1 key2-value2 key3-value3 ...' 
			var newStrMap = [];
			strArray.forEach(function(key, val){
				if(val.includes("-")){
					let strkey = val.split('-')[0];
					let strval = val.split('-')[1];
					newStrMap[strkey] = strval;
				}
			});  
			return newStrMap 
		}

		DecodeString = (dataFromServer,strkey) => {
			let newStrMap = this.DecodeStringAll(dataFromServer);
			return newStrMap[strkey];
		}
	  
} 
//++++++++++++++++++++++++++++++++++++++++++++++++++ end wmo ++++++++++++++++++++++++++++++++++++++++++++++++





 
