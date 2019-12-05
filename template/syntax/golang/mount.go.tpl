{{- comment "Mount mounts the controller to the router" }}
{{- comment "stride:generate" (key .receiver "mount") }}
func (x *{{ .receiver | camelize }}) Mount(r chi.Router) {
	{{- range .operations }}
	r.{{ .Method | titleize }}("{{ .Path }}", x.{{ .Name | camelize }})
	{{- end }}

	// stride:define body:start
	// NOTE: not implemented
	// stride:define body:end
}
