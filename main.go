//main.go
package main

import (
	_ "./cmd/cli"    // Importe le package 'cli' pour que ses init() soient exécutés
	_ "./cmd/server" // Importe le package 'server' pour que ses init() soient exécutés
	"./cmd"
)

func main() {
	// TODO Exécute la commande racine de Cobra.
	cmd.Execute()
}
