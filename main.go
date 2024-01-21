package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	// Local packages
	"github.com/fleroy-isagri/env2js/utils"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

// Wrap the library method for the sake of unit-tests that needs to mock the returns.
// GetEnv
// // Testing it without mock would imply to cover the environment variable erasing.
// // Which means that windows & linux would have two different behaviour.
var (
	Getenv      = os.Getenv
	Exit        = os.Exit
	ReadFile    = os.ReadFile
	WriteFile   = os.WriteFile
	HandleError = utils.HandleError
	LogSuccess  = utils.LogSuccess
)

const (
	SettingsFolderPathEnvKey   string = "SETTINGS_FOLDER_PATH"
	SettingsFilePrefixEnvKey   string = "SETTINGS_FILE_PREFIX"
	SettingsVariableNameEnvKey string = "SETTINGS_VARIABLE_NAME"
)

var (
	// variables are set by GoReleaser with this default commandline on build command :
	// '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'
	Commit  string
	Version string
	Date    string
	BuiltBy string
)

type Walker struct {
	CurrentPath         []string
	SettingVariableName string
}

func (w *Walker) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.BindingElement:
		if n.Binding.String() != w.SettingVariableName {
			return nil
		}
	case *js.Property:
		w.CurrentPath = append(w.CurrentPath, n.Name.String())
		if valueExpression, ok := n.Value.(*js.LiteralExpr); ok {
			if newStringValue, ok := GetEnvValue(w.CurrentPath); ok {
				UpdateData(valueExpression, newStringValue)
			}
			w.CurrentPath = w.CurrentPath[:len(w.CurrentPath)-1]
			return nil
		}
	case *js.ArrayExpr:
		for i, item := range n.List {
			if item.Value != nil {
				if valueExpression, ok := item.Value.(*js.LiteralExpr); ok {
					w.CurrentPath = append(w.CurrentPath, "["+fmt.Sprint(i)+"]")
					if newStringValue, ok := GetEnvValue(w.CurrentPath); ok {
						UpdateData(valueExpression, newStringValue)
					}
					w.CurrentPath = w.CurrentPath[:len(w.CurrentPath)-1]
				}
			}
		}
	}
	return w
}

func (w *Walker) Exit(n js.INode) {
	switch n.(type) {
	case *js.Property:
		w.CurrentPath = w.CurrentPath[:len(w.CurrentPath)-1]
	case *js.PropertyName:
	}
}

func GetEnvValue(path []string) (string, bool) {
	// TODO : use settingsVariableName instead
	computedKey := "AppSettings_"
	computedKey += strings.Join(path, "_")

	if envValue := Getenv(computedKey); envValue != "" {
		return envValue, true
	}

	return "", false
}

func GetEnvOrPanic(value string) string {
	env := Getenv(value)
	if env == "" {
		HandleError(errors.New("No environment key for : " + value))
	}

	return env
}

func UpdateData(valueExpression *js.LiteralExpr, newValue string) {
	if valueExpression.TokenType.String() == "String" {
		valueExpression.Data = []byte("'" + newValue + "'")
		return
	}

	// Cas des booléens. Il y a TrueToken et FalseToken pour les valeurs de type booléen retourné par le parser
	if valueExpression.TokenType.String() == "true" || valueExpression.TokenType.String() == "false" {
		valueExpression.Data = []byte(newValue)
		return
	}

	if valueExpression.TokenType.String() == "Decimal" {
		valueExpression.Data = []byte(newValue)
		return
	}
}

func DefineFilePath(settingsFolderPath string, settingsFilePrefix string) (string, error) {
	settingsSearchFilePattern := filepath.Join(settingsFolderPath, settingsFilePrefix) + "*.js"
	fileList, err := filepath.Glob(settingsSearchFilePattern)
	if err != nil {
		return "", err
	}

	if conditions := len(fileList); conditions == 0 {
		return "", fmt.Errorf("No file found with pattern: " + settingsSearchFilePattern)
	}

	settingsFilePath := fileList[0]
	return settingsFilePath, nil
}

func GetConfigFileLocationValue() (string, string, string) {
	settingsFolderPath := GetEnvOrPanic(SettingsFolderPathEnvKey)

	settingsFilePrefix := GetEnvOrPanic(SettingsFilePrefixEnvKey)

	settingsVariableName := GetEnvOrPanic(SettingsVariableNameEnvKey)

	LogSuccess("✓ "+SettingsFolderPathEnvKey+": ", settingsFolderPath)
	LogSuccess("✓ "+SettingsFilePrefixEnvKey+": ", settingsFilePrefix)
	LogSuccess("✓ "+SettingsVariableNameEnvKey+": ", settingsVariableName)

	return settingsFolderPath, settingsFilePrefix, settingsVariableName
}

// https://eli.thegreenplace.net/2020/testing-flag-parsing-in-go-programs/
type CommandLineConfig struct {
	Version bool

	// args are the positional (non-flag) command-line arguments.
	Args []string
}

// parseFlags parses the command-line arguments provided to the program.
// Typically os.Args[0] is provided as 'progname' and os.args[1:] as 'args'.
// Returns the Config in case parsing succeeded, or an error. In any case, the
// output of the flag.Parse is returned in output.
// A special case is usage requests with -h or -help: then the error
// flag.ErrHelp is returned and output will contain the usage message.
func ParseFlags(progname string, args []string) (config *CommandLineConfig, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf CommandLineConfig
	// -version / --version
	flags.BoolVar(&conf.Version, "version", false, "Display version and exit")

	err = flags.Parse(args)
	// When triggered by the "go test" or "ginkgo" command the args starts with "-test.-timeout=..." or "-ginkgo..."
	isTestCommand := strings.Contains(strings.Join(flags.Args(), ""), "-test") || strings.Contains(strings.Join(flags.Args(), ""), "-ginkgo")

	// Prompt error that exit if there is a parse error and that we are out of test context
	if err != nil && !isTestCommand {
		return nil, buf.String(), err
	}
	conf.Args = flags.Args()

	if conf.Version {
		buf.WriteString(fmt.Sprintf("version : %s\n", Version))
		buf.WriteString(fmt.Sprintf("commit  : %s\n", Commit))
		buf.WriteString(fmt.Sprintf("date    : %s\n", Date))
		buf.WriteString(fmt.Sprintf("builtBy : %s\n", BuiltBy))
	}

	return &conf, buf.String(), nil
}

func LogFlags(config *CommandLineConfig, output string, err error) {
	if err == flag.ErrHelp {
		HandleError(errors.New(output))
		Exit(2)
	} else if err != nil {
		HandleError(errors.New(output))
		Exit(1)
	}

	if config.Version {
		Exit(0)
	}
}

func WriteInConfigFile(settingsFilePath string, settingsVariableName string) {
	// Read the JavaScript file
	jsBytes, err := ReadFile(settingsFilePath)
	HandleError(err)

	// Parse the JavaScript file
	input := parse.NewInput(bytes.NewReader(jsBytes))
	ast, err := js.Parse(input, js.Options{})
	HandleError(err)

	// Analyse du code javascript et réalisation des modifications si nécessaire
	js.Walk(&Walker{SettingVariableName: settingsVariableName}, ast)

	// Write the updated JavaScript file
	// TODO : mettre à jour le fichier uniquement si des modifications ont été faite
	// TODO : afficher les modifications apportées
	var buffer bytes.Buffer
	ast.JS(&buffer)
	err = WriteFile(settingsFilePath, buffer.Bytes(), fs.ModePerm)
	HandleError(err)

	LogSuccess("🎉 Successfuly updated : ", settingsFilePath+" 🎉")
}

func Init() {
	config, output, err := ParseFlags(os.Args[0], os.Args[1:])
	LogFlags(config, output, err)

	settingsFolderPath, settingsFilePrefix, settingsVariableName := GetConfigFileLocationValue()
	settingsFilePath, errorDefineFilePath := DefineFilePath(settingsFolderPath, settingsFilePrefix)
	HandleError(errorDefineFilePath)

	WriteInConfigFile(settingsFilePath, settingsVariableName)
}

// Because of the lowercase letter not being accessible in the main_test package,
// we give to main responsability as little as possible to cover most of the code
func main() {
	Init()
}
