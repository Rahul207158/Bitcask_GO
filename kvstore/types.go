package kvstore

type Entry struct{
	TimeStamp  int64
	KeySize int32
	ValueSize int32
	Key string 
	Value string

}
var KeyDir map[string]int64
type RequestPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
