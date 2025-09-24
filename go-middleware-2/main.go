package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /panic", anotherPnaic)
	mux.HandleFunc("GET /debug/", sourceCodeHandler)
	if err := http.ListenAndServe(":8080", devMiddlware(mux)); err != nil {
		log.Fatal("Error while stating the sesrver")
	}
}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	lineNumber, err := strconv.Atoi(r.FormValue("linenumber"))
	if err != nil {
		lineNumber = -1
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	lexer := lexers.Match(path)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	var lineRange [][2]int
	if lineNumber > 0 {
		lineRange = append(lineRange, [2]int{
			lineNumber, lineNumber,
		})
	}
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(true), html.LineNumbersInTable(true), html.HighlightLines(lineRange))

	it, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_ = formatter.Format(w, style, it)
}

func devMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if didRecoverd := recover(); didRecoverd != nil {
				stackTrace := debug.Stack()
				parser(stackTrace)
				if os.Getenv("Env") != "Dev" {

					http.Error(w, "class='backgroundColor:red'> Something went wrong", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Panic haappned at <br/> stack stract:<br/> %s", parser(stackTrace))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func anotherPnaic(w http.ResponseWriter, r *http.Request) {
	panic("another panic")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "<h1>Hello World !!</h1>")
	panic("kya hal")
}
