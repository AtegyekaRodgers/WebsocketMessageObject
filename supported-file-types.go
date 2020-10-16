package websocketMessageObject

import (
    "fmt"
    "os"
    "strings"
    "mime"
)


type suportedType struct {
    ttype string
    eextension string
    ddescription string
} 

func initSuportedFileTypes() []suportedType {
    var sliceOfTyps []suportedType
    copy(sliceOfTyps,[]suportedType{
    {ttype:"audio/aac",eextension:".aac",ddescription:"AAC audio"},
    {ttype:"application/x-abiword",eextension:".abw",ddescription:"AbiWord document"},
    {ttype:"application/x-freearc",eextension:".arc",ddescription:"Archive document (multiple files embedded)"}, 
    {ttype:"video/x-msvideo",eextension:".avi",ddescription:"AVI: Audio Video Interleave"},
    {ttype:"application/vnd.amazon.ebook",eextension:".azw",ddescription:"Amazon Kindle eBook format"},
    {ttype:"application/octet-stream",eextension:".bin",ddescription:"Any kind of binary data"},
    {ttype:"image/bmp",eextension:".bmp",ddescription:"Windows OS/2 Bitmap Graphics"},
    {ttype:"application/x-bzip",eextension:".bz",ddescription:"BZip archive"},
    {ttype:"application/x-bzip2",eextension:".bz2",ddescription:"BZip2 archive"},
    {ttype:"application/x-csh",eextension:".csh",ddescription:"C-Shell script"},
    {ttype:"text/css",eextension:".css",ddescription:"Cascading Style Sheets (CSS)"},
    {ttype:"text/csv",eextension:".csv",ddescription:"Comma-separated values (CSV)"},
    {ttype:"application/msword",eextension:".doc",ddescription:"Microsoft Word"},
    {ttype:"application/vnd.openxmlformats-officedocument.wordprocessingml.document",eextension:".docx",ddescription:"Microsoft Word (OpenXML)"},
    {ttype:"application/vnd.ms-fontobject",eextension:".eot",ddescription:"MS Embedded OpenType fonts"},
    {ttype:"application/epub+zip",eextension:".epub",ddescription:"Electronic publication (EPUB)"},
    {ttype:"application/gzip",eextension:".gz",ddescription:"GZip Compressed Archive"},
    {ttype:"image/gif",eextension:".gif",ddescription:"Graphics Interchange Format (GIF)"},
    {ttype:"text/html",eextension:".htm",ddescription:"HyperText Markup Language (HTML)"},
    {ttype:"text/html",eextension:".html",ddescription:"HyperText Markup Language (HTML)"},
    {ttype:"image/vnd.microsoft.icon",eextension:".ico",ddescription:"Icon format	"},
    {ttype:"text/calendar",eextension:".ics",ddescription:"iCalendar format"},
    {ttype:"application/java-archive",eextension:".jar",ddescription:"Java Archive (JAR)"},
    {ttype:"image/jpeg",eextension:".jpeg",ddescription:"JPEG images"},
    {ttype:"image/jpeg",eextension:".jpg",ddescription:"JPEG images"},
    {ttype:"text/javascript",eextension:".js",ddescription:"JavaScript"},
    {ttype:"application/json",eextension:".json",ddescription:"JSON format"},
    {ttype:"application/ld+json",eextension:".jsonld",ddescription:"JSON-LD format"},
    {ttype:"audio/midi audio/x-midi",eextension:".mid",ddescription:"Musical Instrument Digital Interface"},
    {ttype:"audio/midi audio/x-midi",eextension:".midi",ddescription:"Musical Instrument Digital Interface (MIDI)"},
    {ttype:"text/javascript",eextension:".mjs",ddescription:"JavaScript module"},
    {ttype:"audio/mpeg",eextension:".mp3",ddescription:"MP3 audio"},
    {ttype:"video/mpeg",eextension:".mpeg",ddescription:"MPEG Video"},
    {ttype:"application/vnd.apple.installer+xml",eextension:".mpkg",ddescription:"Apple Installer Package"},
    {ttype:"application/vnd.oasis.opendocument.presentation",eextension:".odp",ddescription:"OpenDocument presentation document"},
    {ttype:"application/vnd.oasis.opendocument.spreadsheet",eextension:".ods",ddescription:"OpenDocument spreadsheet document"},
    {ttype:"application/vnd.oasis.opendocument.text",eextension:".odt",ddescription:"OpenDocument text document"},
    {ttype:"audio/ogg",eextension:".oga",ddescription:"OGG audio"},
    {ttype:"video/ogg",eextension:".ogv",ddescription:"OGG video"},
    {ttype:"application/ogg",eextension:".ogx",ddescription:"OGG"},
    {ttype:"audio/opus",eextension:".opus",ddescription:"Opus audio"},
    {ttype:"font/otf",eextension:".otf",ddescription:"OpenType font"},
    {ttype:"image/png",eextension:".png",ddescription:"Portable Network Graphics"},
    {ttype:"application/pdf",eextension:".pdf",ddescription:"Adobe Portable Document Format (PDF)"},
    {ttype:"application/x-httpd-php",eextension:".php",ddescription:"Hypertext Preprocessor (Personal Home Page)"},
    {ttype:"application/vnd.ms-powerpoint",eextension:".ppt",ddescription:"Microsoft PowerPoint"},
    {ttype:"application/vnd.openxmlformats-officedocument.presentationml.presentation",eextension:".pptx",ddescription:"Microsoft PowerPoint (OpenXML)"},
    {ttype:"application/vnd.rar",eextension:".rar",ddescription:"RAR archive"},
    {ttype:"application/rtf",eextension:".rtf",ddescription:"Rich Text Format (RTF)"},
    {ttype:"application/x-sh",eextension:".sh",ddescription:"Bourne shell script"},
    {ttype:"image/svg+xml",eextension:".svg",ddescription:"Scalable Vector Graphics (SVG)"},
    {ttype:"application/x-shockwave-flash",eextension:".swf",ddescription:"Small web format (SWF) or Adobe Flash document"},
    {ttype:"application/x-tar",eextension:".tar",ddescription:"Tape Archive (TAR)"},
    {ttype:"image/tiff",eextension:".tif",ddescription:"Tagged Image File Format (TIFF)"},
    {ttype:"image/tiff",eextension:".tiff",ddescription:"Tagged Image File Format (TIFF)"},
    {ttype:"video/mp2t",eextension:".ts",ddescription:"MPEG transport stream"},
    {ttype:"font/ttf",eextension:".ttf",ddescription:"TrueType Font"},
    {ttype:"text/plain",eextension:".txt",ddescription:"Text, (generally ASCII or ISO 8859-n)"},
    {ttype:"application/vnd.visio",eextension:".vsd",ddescription:"Microsoft Visio"},
    {ttype:"audio/wav",eextension:".wav",ddescription:"Waveform Audio Format"},
    {ttype:"audio/webm",eextension:".weba",ddescription:"WEBM audio"},
    {ttype:"video/webm",eextension:".webm",ddescription:"WEBM video"},
    {ttype:"image/webp",eextension:".webp",ddescription:"WEBP image"},
    {ttype:"font/woff",eextension:".woff",ddescription:"Web Open Font Format (WOFF)"},
    {ttype:"font/woff2",eextension:".woff2",ddescription:"Web Open Font Format (WOFF)"},
    {ttype:"application/xhtml+xml",eextension:".xhtml",ddescription:"XHTML"},
    {ttype:"application/vnd.ms-excel",eextension:".xls",ddescription:"Microsoft Excel"},
    {ttype:"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",eextension:".xlsx",ddescription:"Microsoft Excel (OpenXML)"},
    {ttype:"application/xml",eextension:".xml",ddescription:"XML - not readable from casual users (RFC 3023, section 3)"},
    {ttype:"text/xml",eextension:".xml",ddescription:"XML - if readable from casual users (RFC 3023, section 3)"},
    {ttype:"application/vnd.mozilla.xul+xml",eextension:".xul",ddescription:"XUL"},
    {ttype:"application/zip",eextension:".zip",ddescription:"ZIP archive"},
    {ttype:"video/3gpp",eextension:".3gp",ddescription:"3GPP video container"},
    {ttype:"audio/3gpp",eextension:".3gp",ddescription:"3GPP audio container"},
    {ttype:"video/3gpp2",eextension:".3g2",ddescription:"3GPP2 video container"},
    {ttype:"audio/3gpp2",eextension:".3g2",ddescription:"3GPP2 audio container"},
    {ttype:"application/x-7z-compressed",eextension:".7z",ddescription:"7-zip archive"} })
    for i, typp := range sliceOfTyps {
        fmt.Println("type ", i,": ", typp.ttype)
    }
  return sliceOfTyps
 }

func findFileTypeIndex(filetype string) uint {
    filetypes := initSuportedFileTypes()
    for i, ftype := range filetypes {
        if filetype==ftype.ttype {
            return uint(i)
        }
    }
    fmt.Println("Encoding/decoding of unsupported file type: '",filetype,"' skipped.")
    return 99  //this will be used to indicate that file format is unsupported.
}

func getFileType(fi *os.File) string {
    filenme := fi.Name()
    partts := strings.Split(filenme,".")
    ext := partts[(len(partts)-1)]
    typ := mime.TypeByExtension(fmt.Sprintf(".%s",ext))
    return typ
}

func countSupportedFilesOnly(filezArray []*os.File) uint {
    numberOfSupportedFiles := uint(0)
    for _, onefile := range filezArray {
        typ := getFileType(onefile)
        if findFileTypeIndex(typ) != 99 {  
            numberOfSupportedFiles++
        }
    }
    return numberOfSupportedFiles
};




