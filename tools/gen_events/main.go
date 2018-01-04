// Called by `go generate` (from wsevents.go).
package main

import (
	"bytes"
	"flag"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"sort"
	"text/template"
	"unicode"
)

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"toEvent": ToEvent,
}).Parse(`
package dgo2poc

// -- THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT. --
// Instead, edit wsevents.go and/or tools/gen_events/main.go and re-run 'go generate'.

import (
	"sync"
	"context"
	"encoding/json"
)

type wsHandlers struct { {{range .}}
	{{.}} []*func(ctx context.Context, ev *{{.}})
	{{.}}Lock sync.RWMutex
	{{end}}
}

type wsHandler func(hls *wsHandlers) func()

func (hls *wsHandlers) Dispatch(ctx context.Context, t string, data []byte) error {
	switch t { {{range .}}
	case "{{toEvent .}}":
			var ev {{.}}
			if err := json.Unmarshal(data, &ev); err != nil {
				return errors.Wrap(err, t)
			}
			hls.Dispatch{{.}}(ctx, &ev){{end}}
	}
	return nil
}

{{range .}}
func (hls *wsHandlers) Dispatch{{.}}(ctx context.Context, ev *{{.}}) {
	hls.{{.}}Lock.RLock()
	fns := hls.{{.}}
	hls.{{.}}Lock.RUnlock()
	for _, ptr := range fns {
		fn := *ptr
		go fn(ctx, ev)
	}
}
{{end}}

{{range .}}
// Handle a {{.}} event. See WSClient.AddHandler().
func On{{.}}(fn func(ctx context.Context, ev *{{.}})) wsHandler {
	return wsHandler(func(hls *wsHandlers) func() {
		hls.{{.}}Lock.Lock()
		hls.{{.}} = append(hls.{{.}}, &fn)
		hls.{{.}}Lock.Unlock()
		return func() {
			hls.{{.}}Lock.Lock()
			for i, v := range hls.{{.}} {
				if v == &fn {
					hls.{{.}} = append(hls.{{.}}[:i], hls.{{.}}[i+1:]...)
				}
			}
			hls.{{.}}Lock.Unlock()
		}
	})
}
{{end}}
`))

var (
	infile  = flag.String("in", "", "input file")
	outfile = flag.String("out", "", "output file")
)

func main() {
	flag.Parse()
	if *infile == "" || *outfile == "" {
		flag.Usage()
		return
	}

	f, err := parser.ParseFile(token.NewFileSet(), *infile, nil, 0)
	if err != nil {
		panic(err)
	}

	var names []string
	for _, obj := range f.Scope.Objects {
		names = append(names, obj.Name)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, names); err != nil {
		panic(err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(*outfile, out, 0644); err != nil {
		panic(err)
	}
}

func ToEvent(s string) string {
	in := []rune(s)
	out := []rune{}
	for i, r := range in {
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(in[i-1]) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToUpper(r))
	}
	return string(out)
}
