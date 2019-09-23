package zcli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/ztype"
	"strings"
)

func parse(outHelp bool) {
	if Version != "" {
		flagVersion = SetVar("version", getLangs("version")).Bool()
	}
	parseCommand(outHelp)
	parseSubcommand()
}

func parseRequiredFlags(fs *flag.FlagSet, requiredFlags RequiredFlags) (err error) {
	requiredFlagsLen := len(requiredFlags)
	if requiredFlagsLen > 0 {
		flagMap := zarray.New(requiredFlagsLen)
		for _, flagName := range requiredFlags {
			flagMap.Push(flagName)
		}
		fs.Visit(func(f *flag.Flag) {
			_, _ = flagMap.RemoveValue(f.Name)
		})
		flagMapLen := flagMap.Length()
		if flagMapLen > 0 && !*flagHelp {
			arr := make([]string, flagMapLen)
			for i := 0; i < flagMapLen; i++ {
				value, _ := flagMap.Get(i)
				arr[i] = "-" + ztype.ToString(value)
			}
			err = errors.New(fmt.Sprintf("required flags: %s", strings.Join(arr, ", ")))
		}
	}
	return
}

func Parse() {
	flag.Parse()
}

func parseCommand(outHelp bool) {
	Parse()
	if *flagVersion {
		showVersionNum()
		osExit(0)
		return
	}
	if len(cmds) < 1 {
		return
	}
	flag.Usage = usage
	requiredErr := parseRequiredFlags(flag.CommandLine, requiredFlags)
	if requiredErr != nil {
		if len(flag.Args()) > 0 {
			Error(requiredErr.Error())
		} else {
			Help()
		}
	}

	if flag.NArg() < 1 {
		if outHelp {
			Help()
		}
		return
	}
}

func parseSubcommand() {
	name := flag.Arg(0)
	if cont, ok := cmds[name]; ok {
		matchingCmd = cont
		FirstParameter += " " + name
		fsArgs := flag.Args()[1:]
		fs := flag.NewFlagSet(name, flag.ExitOnError)
		flag.CommandLine = fs
		subcommand := &Subcommand{
			Name:        cont.name,
			Desc:        cont.desc,
			Supplement:  cont.Supplement,
			Parameter:   FirstParameter,
			CommandLine: fs,
		}
		cont.command.Flags(subcommand)
		fs.SetOutput(&errWrite{})
		fs.Usage = func() {
			Log.Printf("usage of %s:\n", subcommand.Parameter)
			Log.Printf("\n  %s", subcommand.Desc)
			if subcommand.Supplement != "" {
				Log.Printf("\n%s", subcommand.Supplement)
			}
			ShowFlags(fs)
			ShowRequired(fs, cont.requiredFlags)
		}
		_ = fs.Parse(fsArgs)
		args = fs.Args()
		argsIsHelp(args)
		flagMap := zarray.New(len(cont.requiredFlags))
		for _, flagName := range cont.requiredFlags {
			flagMap.Push(flagName)
		}
		fs.Visit(func(f *flag.Flag) {
			_, _ = flagMap.RemoveValue(f.Name)
		})
		flagMapLen := flagMap.Length()
		if flagMapLen > 0 && !*flagHelp {
			arr := make([]string, flagMapLen)
			for i := 0; i < flagMapLen; i++ {
				value, _ := flagMap.Get(i)
				arr[i] = "-" + ztype.ToString(value)
			}
			Error("required flags: %s", strings.Join(arr, ", "))
		}
	} else if name != "" {
		unknownCommandFn(name)
		osExit(1)
	}
	return
}