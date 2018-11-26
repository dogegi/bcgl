package main

import (
	"./block"
	"./cli"
)

func main() {
	bc := block.NewBlockchain()
	defer bc.DB.Close()

	cli := cli.CLI{bc}
	cli.Run()
}
