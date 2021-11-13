package main

import (
	"encoding/json"
	. "fmt"
	"reflect"
)

/**
结构体标签

`json:",omitempty"`
`json:"-"`   当前字段总被忽略
`json:"-,"`  当前字段序列化后的key为 -
`json:",string"`  JSON-encoded string
其他例子参考  encode_test.go
*/
type jsonobject struct {
	A  int    `json:"aAlia"`       // 结构体标签，json序列化后的key就是这里指定的key
	B  string `json:"b,omitempty"` // json序列化时若B为空则不会被序列化
	C  bool
	AA []int `json:"aa"`
	D  *int
	// 字段首字母小写不会被序列化
	privateField string
	MM           map[string]interface{} `label:"xx" id:"22"`
	nestobject
}

func (o jsonobject) Done() {
	Println("method done")
}

type nestobject struct {
	Name string
}

/**
type StructField struct {
    Name string          // 字段名
    PkgPath string       // 字段路径
    Type      Type       // 字段反射类型对象
    Tag       StructTag  // 字段的结构体标签
    Offset    uintptr    // 字段在结构体中的相对偏移
    Index     []int      // Type.FieldByIndex中的返回的索引值
    Anonymous bool       // 是否为匿名字段
}
*/
func main() {

	var d = 5
	o := jsonobject{
		A: 1,
		//B:            "ss",
		C:            true,
		AA:           []int{2, 4, 1},
		D:            &d,
		privateField: "nothing",
		MM: map[string]interface{}{
			"k": 1,
		},
		nestobject: nestobject{
			Name: "kk",
		},
	}
	Println(o)

	oValue := reflect.ValueOf(o) // <==> reflect.ValueOf(&o).Elem()
	Println("oValue", oValue, oValue.CanAddr())

	dField := oValue.FieldByName("D")
	Println("dField", dField, dField.Type(), dField.CanInterface(), dField.CanSet(), dField.Elem().CanSet(), dField.CanAddr())
	if dField.Elem().CanSet() {
		dField.Elem().SetInt(55)
	}
	Println("aField SetInt", dField.Elem())

	if content, err := json.Marshal(o); err == nil {
		s := string(content)
		Println(s)
	}

	typeO := reflect.TypeOf(o)
	// 类型 与 种类 ： jsonobject struct
	// 类型指：系统原生数据类型，如 int、string、bool、float32 等类型，以及使用 type 关键字定义的类型
	// 种类指：在reflect中定义的对象归属的品种
	Println(typeO.Name(), typeO.Kind())

	var p = &o
	typeP := reflect.TypeOf(p)
	Println(typeP.Name(), typeP.Kind())

	// 获取这个指针指向的元素类型，这个获取过程被称为取元素，等效于对指针类型变量做了一个*操作，等价于 typeO
	typeOP := typeP.Elem()
	Println("typeO == typeOP", typeO == typeOP)
	// jsonobject struct
	Println(typeOP.Name(), typeOP.Kind())

	// 获取全部字段
	Println(typeOP.NumField())
	for i := 0; i < typeOP.NumField(); i++ {
		Println("field", typeOP.Field(i))
	}
	Println(typeOP.FieldByName("privateField"))
	// 模糊搜索
	sfield, _ := typeOP.FieldByNameFunc(func(fieldName string) bool {
		return "privateField" == fieldName
	})
	Println("FieldByNameFunc", sfield)
	Println("FieldByNameFunc", sfield.Name)

	// get tag
	bField, _ := typeOP.FieldByName("MM")
	bTag := bField.Tag
	Println(bTag)
	Println(bTag.Get("label"))
	if tag, ok := bTag.Lookup("id"); ok {
		Println("tag id", tag)
	}

	// 多层成员访问时，根据 []int 提供的每个结构体的字段索引，返回字段的信息，没有找到时返回零值
	Println(typeOP.FieldByIndex([]int{7}))    // nestobject
	Println(typeOP.FieldByIndex([]int{7, 0})) // nestobject.name

	Println(typeOP.NumMethod())
	for i := 0; i < typeOP.NumMethod(); i++ {
		Println("method", typeOP.Method(i))
	}
	if m, found := typeOP.MethodByName("Done"); found {
		Println("MethodByName", m)
	}

}
