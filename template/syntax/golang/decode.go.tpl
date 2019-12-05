{{- comment "UnmarshalJSON unmarshals from valid JSON" }}
{{- comment "stride:generate" (key .receiver "UnmarshalJSON") }}
func (x *{{ .receiver | camelize }}) UnmarshalJSON(data []byte) error {
  x.Body  = &{{ .body | camelize }}{}
  return json.Unmarshal(data, x.Body)
}

{{- comment "UnmarshalXML unmarshals from valid XML" }}
{{- comment "stride:generate" (key .receiver "UnmarshalXML") }}
func (x *{{ .receiver | camelize }}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
  x.Body  = &{{ .body | camelize }}{}
  return d.DecodeElement(x.Body, &start)
}
