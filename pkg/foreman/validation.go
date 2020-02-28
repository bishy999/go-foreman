package foreman

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
)

const (
	// CreateUsage message identify what input is expected
	createUsage = `
	##############################################################################################
	#                                                                                            #
	#  Enter the name and size of the new host you would like to create                          #
	#                                                                                            #
	#  Usage:                                                                                    #
	#      ./foreman-client create -name=mytestenv -size=i3.4xlarge -group=1 -profile=2          #
	#                                                                                            #
	##############################################################################################
	`

	// DeleteUsgage message identify what input is expected
	deleteUsgage = `
	########################################################################
	#                                                                      #
	#  Enter the name of the host you would like to delete                 #
	#                                                                      #
	#  Usage:                                                              #
	#      ./foreman-client delete -name=dev99                             #
	#                                                                      #
	########################################################################
	`

	createArg = "create"
	deleteArg = "delete"
)

// CheckUserInput determines whether sufficient credentials and user input has been provided
func (ci *ConnectionInfo) CheckUserInput() (bool, error) {

	var ok bool

	if ci.Username == "" || ci.Password == "" || ci.BaseURL == "" {
		usage()
		msg := "username, password & url should be part of the api call"
		return ok, errors.New(msg)
	}

	err := checkFlags(ci)
	if err != nil {
		ok = false
	} else {
		ok = true
	}

	return ok, err
}

// usage prints both create and delete usage
func usage() {
	log.Print(createUsage)
	log.Print(deleteUsgage)
}

// Logerr error check and logging
func logerr(err error, x ...int) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

// checkFlags ensure user input is correct
func checkFlags(ci *ConnectionInfo) error {

	if len(os.Args) < 2 {
		usage()
		msg := "create or delete sub command is required"
		return errors.New(msg)
	}

	createCommand := flag.NewFlagSet(createArg, flag.ExitOnError)
	createNamePtr := createCommand.String("name", "", "name of instance to create. (Required)")
	createSizePtr := createCommand.String("size", "", "size of instance to create. (Required)")
	createHostGroupPtr := createCommand.Int("group", 0, "hostgroup_id to use. (Required)")
	createHostProfilePtr := createCommand.String("profile", "0", "compute_profile_id to use. (Required)")

	deleteCommand := flag.NewFlagSet(deleteArg, flag.ExitOnError)
	deleteNamePtr := deleteCommand.String("name", "", "name of instance to delete. (Required)")

	switch os.Args[1] {
	case createArg:
		err := createCommand.Parse(os.Args[2:])
		logerr(err)
	case deleteArg:
		err := deleteCommand.Parse(os.Args[2:])
		logerr(err)
		ci.Action = deleteArg
	default:
		usage()
		os.Exit(1)
	}

	if createCommand.Parsed() {
		if *createNamePtr == "" {
			createCommand.PrintDefaults()
			msg := "hostname needs to be provided"
			return errors.New(msg)
		}

		match, _ := regexp.MatchString("^[a-zA-Z0-9.]{1,20}$", *createNamePtr)

		if !match {
			createCommand.PrintDefaults()
			msg := "hostname needs to be alphanumerical and less than 15 characters"
			return errors.New(msg)
		}

		ci.Hostname = *createNamePtr
		if *createSizePtr == "" {
			createCommand.PrintDefaults()
			msg := "size needs to be provided"
			return errors.New(msg)
		}

		ci.Group = *createHostGroupPtr
		if *createHostGroupPtr == 0 {
			createCommand.PrintDefaults()
			msg := "hostgroup_id needs to be provided"
			return errors.New(msg)
		}

		ci.Profile = *createHostProfilePtr
		if *createHostProfilePtr == "" {
			createCommand.PrintDefaults()
			msg := "compute_profile_id needs to be provided"
			return errors.New(msg)
		}
	}
	if deleteCommand.Parsed() {
		if *deleteNamePtr == "" {
			deleteCommand.PrintDefaults()
			msg := "hostname needs to be provided"
			return errors.New(msg)
		}
		ci.Hostname = *deleteNamePtr

	}
	ci.Size = *createSizePtr

	return nil
}
