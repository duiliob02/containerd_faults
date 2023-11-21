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

			// Cerco la funzione Delete(...)
			if x.Name.Name == "Delete" && len(x.Body.List) > 1 {

				for _, stmt := range x.Body.List {
					if ifStmt, ok := stmt.(*ast.IfStmt); ok {
						ifStmt.Cond = &ast.Ident{Name: "true"}
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

