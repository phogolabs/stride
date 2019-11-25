{{- comment (camelize .operation) "handles endpoint" (uppercase .method) .path }}
{{- comment .summary }}
{{- comment .description }}
{{- comment .deprecated }}
{{- comment "stride:generate" (key .controller .operation) }}
func (x *{{ .controller | camelize }}) {{ .operation | camelize }}(w http.ResponseWriter, r *http.Request) {
	reactor := restify.NewReactor(w, r)

	var (
		input  = &{{ .operation | camelize }}Input{}
		output = &{{ .operation | camelize }}Output{}
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
