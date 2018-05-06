package main

import (
	"text/template"
	"bytes"
	"fmt"
)

type Inventory struct {
	Material string
	Count    uint
	Aa       A
}

type A struct {
	name string
}

func main() {
	/*buf := bytes.NewBuffer([]byte{})
	temp := &Inventory{
		Material: "wool",
		Count:    17,
		Aa: A{
			"sadas",
		},
	}
	tmpl1, err := template.ParseGlob("E:\\NetworkRepository\\src\\GoLibApp\\template\\temp.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl1.ExecuteTemplate(buf, "temp.tmpl", temp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v",buf.String())*/
	temp := &Inventory{
		Material: "wool",
		Count:    17,
		Aa: A{
			"sadas",
		},
	}

	tmpl := template.New("").Parse("")

}
