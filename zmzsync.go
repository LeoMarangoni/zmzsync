package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var host1 = flag.String("host1", "", "Source zimbra server.")
var host2 = flag.String("host2", "", "Destination zimbra server.")
var port1 = flag.String("port1", "7071", " zimbra Admin Port on host1.")
var port2 = flag.String("port2", "7071", " zimbra Admin Port on host2.")
var user1 = flag.String("user1", "", "User to migrate from host1.")
var user2 = flag.String("user2", "", "User to migrate to host2.")
var password1 = flag.String("password1", "", "Password to authenticate in host1.")
var password2 = flag.String("password2", "", "Password to authenticate in host2.")
var authuser1 = flag.String("authuser1", "", "Admin user to authenticate in host1. (Optional)")
var authuser2 = flag.String("authuser2", "", "Admin user to authenticate in host2. (Optional)")
var workdir = flag.String("workdir", "/tmp", "Directory to store exported tgz files.")

var typesHelp = fmt.Sprint(
	"Which data type to migrate. default is to migrate ALL\n" +
		"\tValid Types are:\n" +
		"\t\tmessage\n" +
		"\t\tconversation\n" +
		"\t\tcontact\n" +
		"\t\tappointment\n" +
		"\t\ttask\n" +
		"\tExample, migrating only calendar and contacts: -types contact,appointment\n",
)
var types = flag.String("types", "", typesHelp)
var start = flag.String("datestart", "", "Migrate from that date forward. Date format YYYY-MM-DD")
var end = flag.String("dateend", "", "Migrate from that date back. Date format YYYY-MM-DD")

func validateFlags() {
	flag.Parse()
	if *host1 == "" {
		log.Println("host1 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set host1 and try again")

	}
	if *host2 == "" {
		log.Println("host2 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set host2 and try again")

	}
	if *user1 == "" {
		log.Println("user1 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set user1 and try again")

	}
	if *user2 == "" {
		log.Println("user2 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set user2 and try again")

	}
	if *password1 == "" {
		log.Println("password1 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set password1 and try again")

	}
	if *password2 == "" {
		log.Println("password2 is required")
		flag.PrintDefaults()
		log.Fatal("Exitting... set password2 and try again")

	}
}

func main() {
	//Init parser
	validateFlags()
	//Log Download status and save cartridge position to use in status counter
	log.SetOutput(os.Stdout)
	log.Print("Downloading " + *user1 + "\033[s")
	err := downloadMailbox(*workdir+"/"+*user1, *authuser1, *user1, *password1, *host1, *port1, *types, *start, *end)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Finished downloading " + *user1)
	//Log Upload status and save cartridge position to use in status counter
	log.Println("Uploading " + *user1 + " to " + *user2 + "\033[s")
	err = uploadMailbox(*workdir+"/"+*user1, *authuser2, *user2, *password2, *host2, *port2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Finished Uploading " + *user1 + " to " + *user2)
}
