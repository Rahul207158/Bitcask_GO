package main

import (
	"github.com/Rahul207158/Bitcask_GO/kvstore"
	"fmt"
	"time"
)

func main(){
	entry :=	kvstore.Entry{
		TimeStamp : time.Now().Unix(),
		Key:       "foo",
		Value:     "water",
		KeySize:   int32(len("foo")),
		ValueSize: int32(len("water")),
	}

	err:=kvstore.WriteEntry("data/store_data",entry)
	if(err !=nil){
	fmt.Print("ERR",err)
	}
	val,_:=kvstore.ReadEntry("data/store_data",0)
	fmt.Print("VALUE IS  ",val)
}