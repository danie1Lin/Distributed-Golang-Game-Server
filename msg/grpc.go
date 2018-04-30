package msg

import (
	p "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func CallAllClient(entityTypeName string, id uuid.UUID, f string, args ...p.Message) {
	for _, a := range args {
		v := reflect.ValueOf(a)
		t := v.Type()
		anything, _ := ptypes.MarshalAny(a)
		x := &Pos{}
		AnyDecode(anything, x)
		log.Debug(a, v, t, anything, x)
	}
}

func AnyEecode(in p.Message) *any.Any {
	anything, _ := ptypes.MarshalAny(in)
	return anything
}

func AnyDecode(anything *any.Any, out p.Message) {
	err := ptypes.UnmarshalAny(anything, out)
	if err != nil {
		log.Info(err)
	}
}

func RegistMessage() {

}

type MessageManager struct {
	a map[reflect.Type]p.Message
	b map[p.Message]reflect.Type
}
