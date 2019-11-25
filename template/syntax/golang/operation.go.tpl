reactor := restify.NewReactor(w, r)

var (
   input  = {{.name}}Input{}
   output = {{.name}}Output{}
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
