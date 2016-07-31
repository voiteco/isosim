package spec

import (
	"fmt"
	"log"
	"encoding/hex"
	"github.com/rkbalgi/go/encoding/ebcdic"
)

type Field struct {
	Id               int
	Name             string
	FieldInfo        *FieldInfo
	Position         int;
	fields           []*Field
	//for bitmap only
	fieldsByPosition map[int]*Field
}

func (field *Field) ValueToString(data []byte) string {

	switch(field.FieldInfo.FieldDataEncoding){
	case BCD:
		fallthrough

	case BINARY:{
		return hex.EncodeToString(data);
	}
	case ASCII:{
		return string(data);
	}
	case EBCDIC:{
		return ebcdic.EncodeToString(data)
	}
	default:
		log.Fatal("invalid encoding - ", field.FieldInfo.FieldDataEncoding);

	}
	return "";

}

func (field *Field) ValueFromString(data string) []byte {

	switch(field.FieldInfo.FieldDataEncoding){
	case BCD:
		fallthrough

	case BINARY:{
		str, err := hex.DecodeString(data);
		if err != nil {
			panic(err);
		}
		return str;
	}
	case ASCII:{
		return []byte(data);
	}
	case EBCDIC:{
		return ebcdic.Decode(data)
	}
	default:
		log.Fatal("invalid encoding -", field.FieldInfo.FieldDataEncoding);

	}
	return nil;

}

func (field *Field) HasChildren() bool {
	return len(field.fields) > 0;
}

func (field *Field) Children() []*Field {
	return field.fields;
}

func (field *Field) AddChildField(fieldName string, position int, fieldInfo *FieldInfo) {

	newField := &Field{Name:fieldName, Id:NextId(), Position:position, FieldInfo:fieldInfo}
	field.fields = append(field.fields, newField)

	if (field.FieldInfo.Type == BITMAP) {
		field.fieldsByPosition[position] = newField;
	}
	newField.FieldInfo.Msg = field.FieldInfo.Msg;
}

//Returns properties of the Field as a string
func (field *Field) String() string {

	switch field.FieldInfo.Type{
	case FIXED:{
		return fmt.Sprintf("%-40s - Length: %02d; Encoding: %s", field.Name,
			field.FieldInfo.FieldSize,
			getEncodingName(field.FieldInfo.FieldDataEncoding));

	}
	case BITMAP:{
		return fmt.Sprintf("%-40s - Encoding: %s", field.Name,
			getEncodingName(field.FieldInfo.FieldDataEncoding));
	}
	case VARIABLE:{
		return fmt.Sprintf("%-40s - Length Indicator Size : %02d; Length Indicator Encoding: %s; Encoding: %s",
			field.Name, field.FieldInfo.LengthIndicatorSize,
			getEncodingName(field.FieldInfo.FieldDataEncoding),
			getEncodingName(field.FieldInfo.LengthIndicatorEncoding));

	}
	default:
		log.Fatal("invalid field type -", field.FieldInfo.Type)
	}

	return "";

}
