var supportedFileTypes = require('./supported-file-types');
var wmoUtils = require('./utils');

// Public API
module.exports = WebsocketMessageObject;
 
//=========================================================================================================
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
        //NOTE: dataFromServer is expected to be an ArrayBuffer. 
        //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function. 
	    var BIGendian = false;
	    var LITTLEendian = true;
	    let dataView = new DataView(dataFromServer);
	    console.log("dataView.byteLength = "+dataView.byteLength);
	    let wmoheader = {
		    FilesHeaderOffset:Number(dataView.getUint32(0,BIGendian)),
		    FilesHeaderSize:Number(dataView.getUint32(4,BIGendian)),
		    FilesTotalSize:Number(dataView.getUint32(8,BIGendian)),
		    JsonOffset:Number(dataView.getUint32(12,BIGendian)),
		    JsonSize:Number(dataView.getUint32(16,BIGendian)),
		    StringsOffset:Number(dataView.getUint32(20,BIGendian)),
		    StringsSize:Number(dataView.getUint32(24,BIGendian))
	    };
	    return wmoheader;  
}

function readFilesHeader(filesBytes, hedrSize){
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
	var BIGendian = false;
	var LITTLEendian = true;
	let wmofHeaderBinary = new Uint8Array(filesBytes,0,(hedrSize+1));
	let dataView = new DataView(wmofHeaderBinary.buffer);
	let noOfFiles = Number(dataView.getUint8(0,BIGendian));
	let filezheader = {
		NumberOfFiles:noOfFiles,
		FilesOffsets:Array.from(new Uint32Array(filesBytes,1,noOfFiles)),
		FilesSizes:Array.from(new Uint32Array(filesBytes,(noOfFiles+1),noOfFiles)),
		FilesTypes:Array.from(new Uint32Array(filesBytes,((2*noOfFiles)+1),noOfFiles))
	};
	return filezheader;
}


function WebsocketMessageObject(objectname) {
	objectname = objectname || "wmo";
	if (!(this instanceof WebsocketMessageObject)) {
		return new WebsocketMessageObject();
	}
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
WebsocketMessageObject.prototype.setFilesSize = function(size){ this._filezLength=size; }
WebsocketMessageObject.prototype.setJsonSize = function(size){ this._jsonnLength=size; }
WebsocketMessageObject.prototype.setStringsSize = function(size){ this._stringzLength=size; }
WebsocketMessageObject.prototype.addToFilesSize = function(size){if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._filezLength+=size; }else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToFilesSize()' expects a number ");} }                 
WebsocketMessageObject.prototype.addToJsonSize = function(size){if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._jsonnLength+=size;}else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToJsonSize()' expects a number. ");} }                    
WebsocketMessageObject.prototype.addToStringsSize = function(size){if(!isNaN(size)){ if(typeof(size)==="string"){size=Number(size);} this._stringzLength+=size;}else{ console.log("'"+size+"'"+" is not a number. The function: '"+this.Objectname+".addToStringsSize()' expects a number ");} }  

WebsocketMessageObject.prototype.AddFile = function(file) {
	if(file.name){
		this.filez.push(file);
		let file_size = file.size;
		if(supportedFileTypes.findFileTypeIndex(file.type)!=99){
		    this.addToFilesSize(file_size);
		}
	}
};

WebsocketMessageObject.prototype.AddFileFrom = function(fileInputId) {
	const file = document.getElementById(fileInputId).files[0];
	if(file.name){
		this.filez.push(file);
		let file_size = file.size;
		if(supportedFileTypes.findFileTypeIndex(file.type)!=99){
		    this.addToFilesSize(file_size);
		}
	}
};

WebsocketMessageObject.prototype.AddFiles = function(filesArray) {

};

WebsocketMessageObject.prototype.AddJson = function(myjson) {
	this.jsonn = myjson;
	let json_size = JSON.stringify(this.jsonn).length;
	this.setJsonSize(json_size);
};

WebsocketMessageObject.prototype.AddString = function(keyy,mystring) {
	if (typeof(mystring)==='string') {
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

WebsocketMessageObject.prototype.AddStringFrom = function(keyy, textInputId) {
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
 
WebsocketMessageObject.prototype.Encode = function() {
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
	
	this.AddString("ends","rightpadding");
	var filez_start_point; 
	var json_start_point;
	var stringz_start_point; 
   var file_offsets = [];
   var file_sizes = [];
   var file_types = [];
   var total_size_of_files = 0;
	  
	  var wmo_offset_track = 31;  //28 bytes (4*7) for the wmo header,1 for endianness, 2 for 'dont care' bytes, 1 byte for holding number of files
	  var files_ofst_track = 0; 
	  var numberoffilez = supportedFileTypes.countSupportedFilesOnly(this.filez);
	  var size_of_files_header=1+(numberoffilez*(4+4+4)); //1 byte holds No.of files, 4 bytes(size of uint32) for each file_sizes[element], 4 bytes for each file_offsets[element] and 4 bytes for each file_types[element]
	  
	  total_size_of_files += size_of_files_header;
	  wmo_offset_track =31 + size_of_files_header;  //increment the '*_track' by 'size_of_files_header'
	  files_ofst_track = 0 + size_of_files_header;
	  
	  //ready to create the BinaryData
	  this.BinaryData = new ArrayBuffer(31+size_of_files_header+this._filezLength+this._jsonnLength+this._stringzLength+8);  //8 is some extra just in case we need it.
	  console.log("total size = "+(31+size_of_files_header+this._filezLength+this._jsonnLength+this._stringzLength+8));
	  let mainHeaderView = new Uint32Array(this.BinaryData, 0, 7);  //Uint32Array(buffer, offset, size); where 'size' is the number of items with the specific size eg 32 bits in this case.
	  let endiannessView = new Uint8Array(this.BinaryData, 28, 1);
	  let dontCareView = new Uint8Array(this.BinaryData, 29, 2);
	  let noOfFilesView = new Uint8Array(this.BinaryData, 31, 1);
	  let fileOffsetsView = new Uint32Array(this.BinaryData, 32, numberoffilez);
	  let fileSizesView = new Uint32Array(this.BinaryData, (32+(1*(4*numberoffilez))), numberoffilez);
	  let fileTypesView = new Uint32Array(this.BinaryData, (32+(2*(4*numberoffilez))), numberoffilez); 
	  let filesDataView = new Uint8Array(this.BinaryData, (32+(3*(4*numberoffilez))), this._filezLength);
		let filesDataStart = (32+(3*(4*numberoffilez)));
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
	  for (var i in this.filez) {
			//transform a file into arraybuffer
			const freadr = new FileReader();
			freadr.readAsArrayBuffer(this.filez[i]);
			//create a view of the arraybuffer
			let fileBytesView = new Uint8Array(freadr.result);
			//determine the bytelength of the typed array 
			var fiLength = fileBytesView.byteLength; 
			var fiType = supportedFileTypes.findFileTypeIndex(this.filez[i].type);
		  if(fiType != 99){ //99 is for unsupported file types.
		    file_offsets.push(files_ofst_track);
		    file_sizes.push(fiLength); 
		    file_types.push(fiType);
		    filesDataView.set(fileBytesView,fdataOffsett); 
		    total_size_of_files += fiLength;
		    files_ofst_track += fiLength;
		    wmo_offset_track += fiLength;
		    fdataOffsett+=fiLength;
		  }
		  
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
  noOfFilesView.set([numberoffilez],0);
  //now write sizes and offsets of each file previously written to this.BinaryData
  for (var i in file_sizes) {
	  let thisFilesize = file_sizes[i];
	  let thisfdataOffset = file_offsets[i]; 
	  let thisfdataType = file_types[i]; 
	  fileSizesView.set([thisFilesize],i); 
	  fileOffsetsView.set([thisfdataOffset],i);
	  fileTypesView.set([thisfdataType],i);
  } 
	
  //write all wmo header attributes at their known file_offsets  
  mainHeaderView.set([filez_start_point],0); 
  mainHeaderView.set([size_of_files_header],1); 
  mainHeaderView.set([total_size_of_files],2); 
  mainHeaderView.set([json_start_point],3); 
  mainHeaderView.set([json_size],4); 
  mainHeaderView.set([stringz_start_point],5); 
  mainHeaderView.set([all_strings_size],6); 
  endiannessView.set([(isLittleEndian?6:112)],0); 
} 

WebsocketMessageObject.prototype.toString = function () {
	return '[object WebsocketMessageObject]';
};


//--------Decoders-------- 
WebsocketMessageObject.prototype.DecodeJson = (dataFromServer) => { 
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
	let wmoheader = readWmoHeader(dataFromServer);
	let jsonvieww = new Uint8Array(dataFromServer,wmoheader.JsonOffset, wmoheader.JsonSize); 
	let jsonstr = wmoUtils.typedArrayToString(jsonvieww);        
	console.log("DecodeJson: jsonstr = "+jsonstr);
	let jsonObject = JSON.parse(jsonstr);
	return jsonObject;
}

WebsocketMessageObject.prototype.ReadFilesBytes = (dataFromServer) => {
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
	let hedr = readWmoHeader(dataFromServer);
	let filesDataArr = new Uint8Array(dataFromServer,hedr.FilesHeaderOffset, hedr.FilesTotalSize); 
	return filesDataArr.buffer;
}

WebsocketMessageObject.prototype.DecodeFiles = (dataFromServer, returnDataUrls) => {
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
    let types = supportedFileTypes.filetypes;
    let preferDataUrls = returnDataUrls || true; //if false is passed, file objects will be returned instead of data URLs
	var hedr = readWmoHeader(dataFromServer); 
	var filesBytes = this.ReadFilesBytes(dataFromServer);
	let files_hedr = readFilesHeader(filesBytes, hedr.FilesHeaderSize);
	let filesWithKeys = [];
	for(var i in files_hedr.FilesOffsets){
	    let oneFileTypedArr = new Uint8Array(filesBytes,files_hedr.FilesOffsets[i], files_hedr.FilesSizes[i]);
	    //From the typedarray 'oneFileView', decode a single file and acquire its File boject or image URL 
	    /*
	     * The File constructor (as well as the Blob constructor) takes an array of parts. 
	     * A part doesn't have to be a DOMString. It can also be a Blob, File, or a typed array. 
	     * You can easily build a File out of a Blob like this:
	     * let file = new File([blob], "filename");
	     * //---
	     */
	     var fyleDataType = types[FilesTypes[i]];
	     let generateKey = (fyltype)=>{
	        var offset = files_hedr.FilesOffsets[i];
	        var sizze = files_hedr.FilesSizes[i];
	        if(!Date.now){ Date.now = function(){return new Date().getTime(); }}
	        var currentTimestamp = new Date.now(); 
	        return fyltype+offset+sizze+currentTimestamp;
	     }
	     
	     var blob = new Blob( [ oneFileTypedArr ], { type: "image/png" } );
	     var filename = "file"+files_hedr.FilesOffsets[i]+files_hedr.FilesSizes[i];
	     var file = new File([blob], filename, {type:"image/png", lastModified:new Date()});
	     var urlCreator = window.URL || window.webkitURL;
	     var dataUrl = urlCreator.createObjectURL( blob );dataUrl
	     //generate a unique key for this file 
	     let newKey = generateKey(fyleDataType);
	     console.log("wmo.DecodeFiles generated a unique file key = "+newKey);
	     //add the file data into the reuturned array with its generated key specified (associative array).
	     if(preferDataUrls){
	         filesWithKeys[newKey] = dataUrl;
	     }else{
	         filesWithKeys[newKey] = file;
	     }
	     /* 
	     * let imageUrl = dataUrl;
	     * var img = document.querySelector( "#photo" ); 
	     * img.src = imageUrl;
	     * urlCreator.revokeObjectURL(); 
	     */ 
	 }
	 return filesWithKeys;
}

WebsocketMessageObject.prototype.DecodeStringAll = (dataFromServer) => {
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
	let wmoheader = readWmoHeader(dataFromServer);
	let stringsview = new Uint8Array(dataFromServer,wmoheader.StringsOffset, wmoheader.StringsSize); 
	let stringz = wmoUtils.typedArrayToString(stringsview); 
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

WebsocketMessageObject.prototype.DecodeString = (dataFromServer,strkey) => {
	//NOTE: dataFromServer is expected to be an ArrayBuffer. 
    //If not, first convert the dataFromServer into ArrayBuffer, then pass to this function.
	let newStrMap = this.DecodeStringAll(dataFromServer);
	return newStrMap[strkey];
}

//export default WebsocketMessageObject;
//================================================================================================================ 

/* 
 * 
 FileReader 
 METHODS:
 -------------------
 
 FileReader.abort()
 	Aborts the read operation. Upon return, the readyState will be DONE .
 
 FileReader.readAsArrayBuffer()
 	Starts reading the contents of the specified Blob , once finished, the result attribute
 	contains an ArrayBuffer representing the file's data.
 
 FileReader.readAsBinaryString()
 	Starts reading the contents of the specified Blob , once finished, the result attribute
 	contains the raw binary data from the file as a string.
 
 FileReader.readAsDataURL()
 	Starts reading the contents of the specified Blob , once finished, the result attribute
 	contains a data: URL representing the file's data.
 
 FileReader.readAsText()
 	Starts reading the contents of the specified Blob , once finished, the result attribute
 	contains the contents of the file as a text string. An optional encoding name can be
 	specified.
 	
 	
 EVENT HANDLERS
 --------------
 	
 	FileReader.onabort
 		A handler for the abort event. This event is triggered each time the reading operation is
 		aborted.
 		
 	FileReader.onerror
 		A handler for the error event. This event is triggered each time the reading operation
 		encounter an error.
 	
 	FileReader.onload
 		A handler for the load event. This event is triggered each time the reading operation is
 		successfully completed.
 	
 	FileReader.onloadstart
 		A handler for the loadstart event. This event is triggered each time the reading is
 		starting.
 	
 	FileReader.onloadend
 		A handler for the loadend event. This event is triggered each time the reading operation is
 		completed (either in success or failure).
 	
 	FileReader.onprogress
 	A handler for the progress event. This event is triggered while reading a Blob content.
 
 */


