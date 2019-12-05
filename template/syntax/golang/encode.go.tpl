{{- comment "MarshalJSON marshals into valid JSON" }}
{{- comment "stride:generate" (key .receiver "MarshalJSON") }}
func (x {{ .receiver | camelize }}) MarshalJSON() ([]byte, error) {
  return json.Marshal(x.Body)
}

{{- comment "MarshalXML marshals into valid XML" }}
{{- comment "stride:generate" (key .receiver "MarshalXML") }}
func (x {{ .receiver | camelize }}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.Body, start)
}
