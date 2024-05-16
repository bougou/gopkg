package snippet

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)

var tmplstr1 = `

// 模板中定义的需要渲染的内容
Hello world!
{{- range . }}
Name: {{ .Name }}	Age: {{ .Age }}
{{- end }}


// 模板中可以定义任意的子模板

// 定义子模板，子模板名字为 title
{{ define "title" }}
>>> title: Just a title
{{ end }}


// 定义子模板 子模板名字为 content
{{ define "content" }}
>>> content start
Hello world!
{{- range . }}
Name: {{ .Name }}	Age: {{ .Age }}
{{- end }}
>>> content end
{{ end }}


// 定义 子模板，子模板名字为 body，body 中又使用到 title, content 子模板
// title 模板中没有需要渲染的变量，所以，template 执行的时候，无需传入数据
// title 模板 is executed with nil data.
// content 中有数据需要渲染，所以 template 执行的时候，需要传入数据，这里特殊变量 $, 代表模板的最顶级作用域"
{{ define "body" }}
>>> body start
{{ template "title" }}
{{ template "content" $ }}
>>> body end
{{ end }}


// top-level 模板中使用 body 子模板，当然也可以直接使用其它的子模板，如 title, content
{{ template "body" . }}
`

func Test_template(t *testing.T) {

	tmpl1, err := template.New("test-tmpl1").Parse(tmplstr1)
	if err != nil {
		panic("parse tmplstr1 failed")
	}

	fmt.Println("tmpl1 name: ", tmpl1.ParseName)
	// tmpl1 name:  test-tmpl1

	type Person struct {
		Name string
		Age  int
	}

	persons := []Person{
		{"Tom", 20},
		{"John", 19},
	}
	tmpl1.Execute(os.Stdout, persons)

	fmt.Println("all defined templates in tmpl1: ")
	for _, tmpl := range tmpl1.Templates() {
		fmt.Printf("template name: %-20s top-level template name: %s\n", tmpl.Name(), tmpl.ParseName)
	}
	// all defined templates in tmpl1:
	// template name: title                top-level template name: test-tmpl1
	// template name: content              top-level template name: test-tmpl1
	// template name: body                 top-level template name: test-tmpl1
	// template name: test-tmpl1           top-level template name: test-tmpl1

	fmt.Println("// tmpl1: top-level test-tmpl1 tempalte")
	tmpl1.ExecuteTemplate(os.Stdout, "test-tmpl1", persons)

	fmt.Println("// tmpl1: subtemplate title")
	tmpl1.ExecuteTemplate(os.Stdout, "title", persons)

	fmt.Println("// tmpl1: subtemplate content")
	tmpl1.ExecuteTemplate(os.Stdout, "content", persons)

	fmt.Println("// tmpl1: subtemplate body")
	tmpl1.ExecuteTemplate(os.Stdout, "body", persons)
}
