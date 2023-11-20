package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// Creo un nuovo FileSet per poter
	fset := token.NewFileSet()

	// "MACRO" per il file .go
	const path_to_source = "/home/duilio/Desktop/containerd/cmd/containerd-shim-runc-v2/runc/container.go"

	// Prendo il file da analizzare e mutare
	source, err := os.Open(path_to_source)
	if err != nil {
		panic("Errore apertura file")
	} // Alla fine di tutto chiudo il file
	defer source.Close()

	file, err := parser.ParseFile(fset, path_to_source, source, 0)
	if err != nil {
		panic("Errore parsing")
	}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:

			// Cerco la funzione Checkpoint(...)
			if x.Name.Name == "Checkpoint" && len(x.Body.List) > 1 {

				// Controllo il blocco di return, in modo da rimodellare l'inizializzazione di lista
				/*
						&process.CheckpointConfig{
						Path:                     r.Path,
						Exit:                     opts.Exit,
						AllowOpenTCP:             opts.OpenTcp,
						AllowExternalUnixSockets: opts.ExternalUnixSockets,
						AllowTerminal:            opts.Terminal,
						FileLocks:                opts.FileLocks,
						EmptyNamespaces:          opts.EmptyNamespaces,
						WorkDir:                  opts.WorkPath,
					}
				*/
				for _, stmt := range x.Body.List {
					if result, ok := stmt.(*ast.ReturnStmt); ok {

						// Cerco l'assegnazione dei parametri alla configurazione del Checkpoint e li elimino
						if ident, ok := result.Results[0].(*ast.CallExpr); ok && ident.Fun.(*ast.SelectorExpr).Sel.Name == "Checkpoint" {
							result.Results[0].(*ast.CallExpr).Args[1].(*ast.UnaryExpr).X.(*ast.CompositeLit).Elts = nil
						}
					}
				}
			}
		}
		return true
	})

	// Stampo l'albero modificato
	log.Println("AST modificato:")
	printer.Fprint(os.Stdout, fset, file)

	// Inietto il fault generato
	injectedFile, err := os.Create(path_to_source)
	if err != nil {
		panic("Errore riscrittura file")
	}
	defer injectedFile.Close()
	err = printer.Fprint(injectedFile, fset, file)
	if err != nil {
		panic("Errore modifica fileSet")
	}
	log.Println("File riscritto, eseguo il test")

}
