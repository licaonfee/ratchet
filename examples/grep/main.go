/*
This program emulates grep
*/

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/licaonfee/ratchet/logger"

	"github.com/licaonfee/ratchet"
	"github.com/licaonfee/ratchet/processors"
)

//not suitable for production code
func exitOnErr(err error) {
	if err != nil {
		log.Printf("E! %v", err)
		os.Exit(1)
	}
}

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 {
		log.Println("use just one argument")
		log.Println("Usage: ")
		log.Printf("%s <regex> \n", os.Args[0])
		os.Exit(2)
	}
	read, err := processors.NewIoReader(processors.WithReader(os.Stdin))
	exitOnErr(err)
	//match all records with a column foo and value "bar"
	filter, err := processors.NewRegexpMatcher(processors.WithPattern(os.Args[1]))
	exitOnErr(err)
	write, err := processors.NewIoWriter(processors.WithWriter(os.Stdout))
	exitOnErr(err)
	pipe := ratchet.NewPipeline(read, filter, write)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	logger.LogLevel = logger.LevelSilent
	perr := pipe.Run()

	select {
	case <-sig:
		log.Println("user interrupt")
		os.Exit(2)
	case err := <-perr:
		exitOnErr(err)
	}

}
