package sql2struct

import (
	"PassageOne/internal/word"
	"fmt"
	"os"
	"text/template"
)

// Define template
const structTpl = `type {{.TableName|ToCamelCase}} struct {
	{{range .Columns}} {{$length := len .Comment}} {{if gt $length 0}} //{{.Comment}} {{else}} // {{.Name}} {{end}}
	{{ $typelen := len .Type}} {{if gt $typelen 0}} {{.Name|ToCamelCase}} {{.Type}} {{.Tag}} {{else}} {{.Name}} {{end}}
	{{end}}

	func (model {{.TableName|ToCamelCase}}) TableName() string {
		return "{{.TableName}}"
	}
}`

type StructTemplate struct {
	structTpl string
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

// Decompositing and transformating for tbColumns
func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, 10)
	for _, column := range tbColumns {
		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", column.ColunmName)
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColunmName,
			Type:    DBTypeTostructType[column.DataType],
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}
	return tplColumns
}

// Parse and Execute template
func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{"ToCamelCase": word.UnderLineToUpperCamelCase}).Parse(t.structTpl))

	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}

	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}
	return nil
}
