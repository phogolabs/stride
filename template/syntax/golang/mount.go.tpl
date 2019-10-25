{{- comment "Mount mounts the controller to the router" }}
{{- comment "stride:generate" (key .controller "mount") }}
func (x *{{ .controller | camelize }}) Mount(r chi.Router) {
	{{- range .operations }}
	r.{{ .Method | titleize }}("{{ .Path }}", x.{{ .Name | camelize }})
	{{- end }}

	// stride:define body:start
	// NOTE: not implemented
	// stride:define body:end
}
