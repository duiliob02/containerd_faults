package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"strings"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// Creo un nuovo FileSet per poter 
	fset := token.NewFileSet()

	// "MACRO" per il file .go
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

			// Cerco la funzione NewContainer(...)
			if x.Name.Name == "NewContainer" && len(x.Body.List) > 1 {

				// Controllo tutti i blocchi if
				for _, stmt := range x.Body.List {
					if blocco_if, ok := stmt.(*ast.IfStmt); ok {

						// Cerco il blocco con condizione "r.Options.GetValue() != nil"
						var condition strings.Builder
						printer.Fprint(&condition, fset, blocco_if.Cond)

						// La cerco hardcoded anche se non è la soluzione più efficiente e sicura
						// Per il caso di studio, conoscendo il codice usato da containerd lo possofare
						if condition.String() == "r.Options.GetValue() != nil" {

							// Essendo che sono già all'interno di un nodo, non posso
							// chiamare c.Replace(...), quindi modifico direttamente la condizione dell'if
							log.Println("Trovata condizione da mutare: ", condition.String())
							blocco_if.Cond.(*ast.BinaryExpr).Op = token.EQL
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


}