package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type command struct {
	Fn        func() error
	Desc      string
	ArgsCount int
	Args      string
}

var gCommands = map[string]command{
	"createPack": command{
		Fn:        cmdCreatePack,
		Desc:      "Create a new mod pack",
		ArgsCount: 3,
		Args:      "<packname> <minecraft version> <forge version>",
	},
	"installPack": command{
		Fn:        cmdInstallPack,
		Desc:      "Install a mod pack from URL",
		ArgsCount: 2,
		Args:      "<packname> <url>",
	},
	"installLocalPack": command{
		Fn:        cmdInstallLocalPack,
		Desc:      "Install specified directory as a pack",
		ArgsCount: 1,
		Args:      "<directory>",
	},
	"info": command{
		Fn:        cmdInfo,
		Desc:      "Show runtime info",
		ArgsCount: 0,
	},
	"registerMod": command{
		Fn:        cmdRegisterMod,
		Desc:      "Register a curseforge mod with an existing pack",
		ArgsCount: 2,
		Args:      "<directory> <url> [<name>]",
	},
	"installMods": command{
		Fn:        cmdInstallMods,
		Desc:      "Install all mods using the manifest",
		ArgsCount: 1,
		Args:      "<directory>",
	},
	"runServer": command{
		Fn:        cmdRunServer,
		Desc:      "Run a Minecraft server using an existing pack",
		ArgsCount: 1,
		Args:      "<directory>",
	},
}

func cmdCreatePack() error {
	name := flag.Arg(1)
	minecraftVsn := flag.Arg(2)
	forgeVsn := flag.Arg(3)

	// Create a new pack directory
	cp, err := NewModPack(name, "")
	if err != nil {
		return err
	}

	// Create the manifest for this new pack
	err = cp.createManifest(name, minecraftVsn, forgeVsn)
	if err != nil {
		return err
	}

	// Create the launcher profile (and install forge if necessary)
	err = cp.createLauncherProfile()
	if err != nil {
		return err
	}

	return nil
}

func cmdInstallPack() error {
	name := flag.Arg(1)
	url := flag.Arg(2)

	// Get ZIP file
	cp, err := NewModPack(name, url)
	if err != nil {
		return err
	}

	// Download the pack
	err = cp.download()
	if err != nil {
		return err
	}

	// Process manifest
	err = cp.processManifest()
	if err != nil {
		return err
	}

	// Create launcher profile
	err = cp.createLauncherProfile()
	if err != nil {
		return err
	}

	// Install mods
	err = cp.installMods()
	if err != nil {
		return err
	}

	// Install overrides
	err = cp.installOverrides()
	if err != nil {
		return err
	}

	return nil
}

func cmdInstallLocalPack() error {
	dir := flag.Arg(1)

	if dir == "." {
		dir, _ = os.Getwd()
	}
	dir, _ = filepath.Abs(dir)

	// Create the mod pack directory (if it doesn't already exist)
	cp, err := OpenModPack(dir)
	if err != nil {
		return err
	}

	// Setup a launcher profile
	err = cp.createLauncherProfile()
	if err != nil {
		return err
	}

	// Install all the mods
	err = cp.installMods()
	if err != nil {
		return err
	}

	return nil
}

func cmdInfo() error {
	fmt.Printf("Env: %+v\n", env())
	return nil
}

func cmdRegisterMod() error {
	dir := flag.Arg(1)
	url := flag.Arg(2)
	name := flag.Arg(3)

	if !strings.Contains(url, "minecraft.curseforge.com") && name == "" {
		return fmt.Errorf("Insufficient arguments")
	}

	cp, err := OpenModPack(dir)
	if err != nil {
		return err
	}

	err = cp.registerMod(flag.Arg(2), flag.Arg(3))
	if err != nil {
		return err
	}

	return nil
}

func cmdInstallMods() error {
	dir := flag.Arg(1)

	cp, err := OpenModPack(dir)
	if err != nil {
		return err
	}

	err = cp.installMods()
	if err != nil {
		return err
	}

	return nil
}

func cmdRunServer() error {
	dir := flag.Arg(1)

	// Open the pack
	cp, err := OpenModPack(dir)
	if err != nil {
		return err
	}

	// Install the server jar, forge and dependencies
	err = cp.installServer()
	if err != nil {
		return err
	}

	return nil
	// Setup the command-line
	// java -jar <forge.jar>
}

func console(f string, args ...interface{}) {
	fmt.Printf(f, args...)
}

func usage() {
	console("usage: mcdex [<options>] <command> [<args>]\n")
	// console(" options:\n")
	// flag.PrintDefaults()
	console(" commands:\n")
	for id, cmd := range gCommands {
		console(" - %s: %s\n", id, cmd.Desc)
	}
}

func usageCmd(name string, cmd command) {
	console("usage: mcdex %s %s\n", name, cmd.Args)
}

func main() {
	// Process command-line args
	flag.Parse()
	if !flag.Parsed() || flag.NArg() < 1 {
		usage()
		os.Exit(-1)
	}

	// Initialize our environment
	err := initEnv()
	if err != nil {
		log.Fatalf("Failed to initialize: %s\n", err)
	}

	commandName := flag.Arg(0)
	command, exists := gCommands[commandName]
	if !exists {
		console("ERROR: unknown command '%s'\n", commandName)
		usage()
		os.Exit(-1)
	}

	// Check that the required number of arguments is present
	if flag.NArg() < command.ArgsCount+1 {
		console("ERROR: insufficient arguments for %s\n", commandName)
		console("usage: mcdex %s %s\n", commandName, command.Args)
		os.Exit(-1)
	}

	err = command.Fn()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}

//mcdex update - download latest mcdex.sqlite
//mcdex forge.install <name> [<vsn>]
//mcdex forge.list

//mcdex init <name> <vsn> <desc>
//mcdex install <modname> [<vsn>]