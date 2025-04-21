package kvstore

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)
func WriteEntry (filepath string ,entry Entry ) (int64,error){
   ///
   file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
  if(err!=nil){
   fmt.Print("ERR in opening file",err)
  }
  defer file.Close()
   offset,err:=file.Seek(0,io.SeekEnd)    // return the offset 
   if(err!=nil){
      fmt.Print("ERR in seeking the offset before writing ",err)
     }
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
  }
   _,err=file.Write([]byte(entry.Key))
   if(err!=nil){
      fmt.Print("ERR in valuesize",err)
     }

   _,err=file.Write([]byte(entry.Value))
   if(err!=nil){
      fmt.Print("Error in writing value",err)
   }  
   return offset,nil;
}
func ReadEntry(filepath string, offset int64) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Err in opening file", err)
		return "", err
	}
	defer file.Close()

	// Skip timestamp (8 bytes)
	_, err = file.Seek(offset+8, io.SeekStart)
	if err != nil {
		fmt.Println("Err in seeking file", err)
		return "", err
	}

	// Read keySize and valueSize (4 + 4 = 8 bytes)
	sizeBuf := make([]byte, 8)
	if _, err := file.Read(sizeBuf); err != nil {
		fmt.Println("Err in reading size buffer", err)
		return "", err
	}
	keySize := int32(binary.LittleEndian.Uint32(sizeBuf[0:4]))
	valueSize := int32(binary.LittleEndian.Uint32(sizeBuf[4:8]))

	fmt.Println("key size is", keySize, "Value Size is", valueSize)

	// Skip the key
	_, err = file.Seek(int64(keySize), io.SeekCurrent)
	if err != nil {
		fmt.Println("Error in skipping key", err)
		return "", err
	}

	// Read the value
	valueBuf := make([]byte, valueSize)
	n, err := file.Read(valueBuf)
   fmt.Printf("Read %d bytes, expected %d\n", n, valueSize)
  // fmt.Printf("Raw value bytes: %v\n", valueBuf)

	if err != nil {
		fmt.Println("Error in reading value", err)
		return "", err
	}
	if int32(n) != valueSize {
		fmt.Println("Warning: read less than value size")
	}

	return string(valueBuf), nil
}
