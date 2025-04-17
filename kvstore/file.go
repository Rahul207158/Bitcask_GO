package kvstore

import (
	"encoding/binary"
	"fmt"
	"os"
)
func WriteEntry (filepath string ,entry Entry ) error{
   ///
   file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
  if(err!=nil){
   fmt.Print("ERR in opening file",err)
  }
  defer file.Close()

  err=binary.Write(file,binary.LittleEndian,entry.TimeStamp)
  if(err!=nil){
   fmt.Print("ERR in timestamp",err)
  }

  err=binary.Write(file,binary.LittleEndian,entry.KeySize)
  if(err!=nil){
   fmt.Print("ERR in keysize",err)
  }

  err=binary.Write(file,binary.LittleEndian,entry.ValueSize)
  if(err!=nil){
   fmt.Print("ERR in write valuesize",err)

   _,err=file.Write([]byte(entry.Key))
   if(err!=nil){
      fmt.Print("ERR in valuesize",err)
     }

   _,err=file.Write([]byte(entry.Value))
   if(err!=nil){
      fmt.Print("Error in writing value",err)
   }  
  }
   return nil;
}