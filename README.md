# WebsocketMessageObject
WebsocketMessageObject is an API that enables you to package any type of data into a single (binary format) massage so that you can easily send it over websocket protocol. You simply create a wmo instance and keep adding any data into it by calling the appropriate functions provided in the API, then finally call the 'Encode' function which converts everything into a stream of binary data that you can send to your back-end server using the 'websocket.send(theBinaryData)' function. You can add images, videos, documents, json data as well as plain text data. The API is available in a number of programming languages including go (golang) and javascript. NOTE: You need to have two wmo API interfaces; one for the front end (mainly in javascript) and the other for the back-end ( golang, nodejs, php, etc).

## Installation 1

### Let's start with the front-end version, (npm package in javascript). 
The Golang version will be explained after javascript. Please scroll down.

To use Websocket Message Object API, install it as a dependency:

If you use npm:

```bash
$ npm install websocket-message-object
```
Or if you use Yarn:


```bash
$ yarn add websocket-message-object
```

## Usage

 First, you need to import the module into your application as follows:

 In a react app,
```javascript
import {WebsocketMessageObject} from 'websocket-message-object';

```
 Or if you are not using react;
```javascript
const wmo = require("websocket-message-object");

```
 Now that you have a wmo object created, you can add any data to it. Here is how to add an image file.
```javascript
wmo.AddFile(file); 
//NOTE: This assumes you alredy have a file object read from the DOM eg using FileReader
 
//If you want to add the file directly from the html file input field, use this method:
wmo.AddFileFrom("id-of-the-html-file-input"); 

//...
//...
// and many more ... (see 'API' section below)
//...
//...

```
 And now, if you're done with adding data, fire the function that does the encoding for you automatically like so:
```javascript
wmo.Encode();

```
# List of functions available (API)

## Adding data to encode
```javascript
wmo.AddFile(file); 
```
```javascript 
wmo.AddFileFrom("file-input-id");
```
```javascript 
wmo.AddJson({jsonobject});
//eg. wmo.AddJson({name: "Rodgers", nationality: "Ugandan", title: "Software Engineer"});
```
```javascript 
wmo.AddString("key","string you intend to add"); 
// eg wmo.AddString("name", "Rodgers"); 
```
```javascript 
wmo.Encode();
``` 
## Decoding data received from server 
```javascript 
let jsondata = wmo.DecodeJson(messageFromServer);
```
```javascript 
let mystr = wmo.DecodeString(messageFromServer, "key"); 
// eg, let myname = wmo.DecodeString(messageFromServer, "name"); => returns a name such as 'Rodgers'
```
```javascript 
let strMap = wmo.DecodeStringAll(messageFromServer);    
//returns a map of key-vale pairs of strings.
```
```javascript 
let filesArray = wmo.DecodeFiles(messageFromServer);
```
```javascript 
let filesBytes = wmo.ReadFilesBytes(messageFromServer);    
//returns all files data in plain binary format
``` 

# NOTE: 
### 'messageFromServer' MUST be a binary format message encoded using the server side version of wmo. This is currently available only in Go (Golang). I am still implementing wmo for other languages. 

I welcome anyone who want to contribute to this project especially in other languages.
## ___

## Installation 2
### Now the go back-end. 
This is how to import the Go package into your project.

Run this command:
```bash 

$ go get "github.com/AtegyekaRodgers/WebsocketMessageObject"

```
Then import the package into your application
```go
package main

import (
  "fmt"
  wmo "github.com/AtegyekaRodgers/WebsocketMessageObject"
  //...
  //...
)

func main(){
  //...
}

```


