{{- comment (camelize .function) "handles endpoint" (uppercase .method) .path }}
{{- comment .summary }}
{{- comment .description }}
{{- comment .deprecated }}
{{- comment "stride:generate" (key .receiver .function) }}
func (x *{{ .receiver | camelize }}) {{ .function | camelize }}(w http.ResponseWriter, r *http.Request) {
	reactor := restify.NewReactor(w, r)

	var (
		input  = &{{ .function | camelize }}Input{}
		output = &{{ .function | camelize }}Output{}
	)

	if err := reactor.Bind(input); err != nil {
		reactor.Render(err)
		return
	}

	// stride:define body:start
	// NOTE: not implemented
	// stride:define body:end

	if err := reactor.Render(output); err != nil {
		reactor.Render(err)
	}
}
