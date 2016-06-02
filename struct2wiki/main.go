package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

// StructField 结构体
type StructField struct {
	PkgName string //所在包名
	Name    string
	Fields  []*Field
	Comment string
}

// Field 字段
type Field struct {
	Name    string
	Type    string
	Comment string
}

// Parse 传入一个go文件 读取里边的struct类型和字段
func Parse(filename string) (structs []*StructField, err error) {

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}

	decls := astf.Decls
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		strut := new(StructField)
		strut.PkgName = astf.Name.Name
		strut.Name = typeSpec.Name.Name
		strut.Comment = strings.Trim(typeSpec.Comment.Text(), "\n")

		structType, ok := typeSpec.Type.(*ast.StructType)
		if ok {
			for _, field := range structType.Fields.List {
				f := new(Field)
				if len(field.Names) > 0 {
					f.Name = field.Names[0].Name
				}
				f.Comment = strings.Trim(field.Comment.Text(), "\n")
				if f.Name == "ID" || f.Name == "UpdatedAt" || f.Name == "CreatedAt" {
					// skip thid field
				} else {

					var buf bytes.Buffer
					printer.Fprint(&buf, fset, field.Type)
					f.Type = buf.String()

					strut.Fields = append(strut.Fields, f)
				}

			}
		}
		structs = append(structs, strut)
	}

	return

}

var path = flag.String("f", "", "文件路劲")

func main() {

	flag.Parse()

	if *path == "" {
		log.Fatal("missing file path param")
	}
	structs, err := Parse(*path)
	if err != nil {
		log.Fatal(err)
	}

	structType := structs[0]
	tpl := "|_.Name|_.Required|_.Type|_.Sample|_.Description | \n"
	for _, field := range structType.Fields {
		fstr := fmt.Sprintf("|%s| YES | %s|  |%s| \n",
			strings.ToLower(field.Name[:1])+field.Name[1:],
			strings.ToUpper(field.Type[:1])+field.Type[1:],
			field.Comment,
		)
		tpl += fstr
	}
	fmt.Println(tpl)
}
