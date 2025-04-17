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
		Value:     "bar",
		KeySize:   int32(len("foo")),
		ValueSize: int32(len("bar")),
	}

	err:=kvstore.WriteEntry("data/store_data",entry)
	if(err !=nil){
	fmt.Print("ERR",err)
	}

}