{{- comment "Status returns an http status code" }}
{{- comment "stride:generate" (key .receiver "status") }}
func (x *{{ .receiver | camelize }}) Status() int {
	// stride:define body:start
	// NOTE: not implemented
	// stride:define body:end
	return {{ .code }}
}
