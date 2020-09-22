module.exports.typedArrayToString = (typedArr)=>{ 
    var aString = '';
    for (var i in typedArr) {
      aString += String.fromCharCode(typedArr[i]);
    }
    return aString;
}


