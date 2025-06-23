//main.go
package main

import (
	_ "urlshortenerGroupe8/cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "urlshortenerGroupe8/cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
	"urlshortenerGroupe8/cmd"
)

func main() {
	// TODO Exécute la commande racine de Cobra.
	cmd.Execute()
}
