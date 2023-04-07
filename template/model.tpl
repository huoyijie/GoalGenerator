// auto generated by GoalGenerator@v0.1.0
package {{.Package}}
{{if .Imports}}
import ({{range .Imports}}
    {{.}}{{end}}
)
{{end}}
type {{.Name}} struct { {{if .Lazy}}
    goalgenerator.Lazy{{end}}{{if .EmbeddingBase}}
    goalgenerator.Base{{end}}
{{range .Fields}}
    {{.Name}} {{.Type}} `{{.Tag}}`{{end}}
}
