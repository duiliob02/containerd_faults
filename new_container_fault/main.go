package main

import (
	"go/ast"
	"go/parser"
	//"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// Creo un nuovo FileSet per poter 
	fset := token.NewFileSet()

	const path_to_source = "/home/parallels/Desktop/Progetto_Tesi/containerd/runtime/v2/runc/container.go"

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
			if x.Name.Name == "NewContainer" {
				log.Println("Trovata funzione: ", x.Name.Name)
			}
		}
		return true
	})


}