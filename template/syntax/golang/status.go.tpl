{{- comment "Status returns an http status code" }}
{{- comment "stride:generate" (key .schema "status") }}
func (x *{{ .schema | camelize }}) Status() int {
	// stride:define body:start
	// NOTE: not implemented
	// stride:define body:end
	return {{ .code }}
}
