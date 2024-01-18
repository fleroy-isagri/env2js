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

const (
	SettingsFolderPathEnvKey   string = "SETTINGS_FOLDER_PATH"
	SettingsFilePrefixEnvKey   string = "SETTINGS_FILE_PREFIX"
	SettingsVariableNameEnvKey string = "SETTINGS_VARIABLE_NAME"
)

var (
	// variables are set by GoReleaser with this default commandline on build command :
	// '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'
	commit  string
	version string
	date    string
	builtBy string
)

type walker struct {
	currentPath         []string
	settingVariableName string
}

func (w *walker) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.BindingElement:
		if n.Binding.String() != w.settingVariableName {
			return nil
		}
	case *js.Property:
		w.currentPath = append(w.currentPath, n.Name.String())
		if valueExpression, ok := n.Value.(*js.LiteralExpr); ok {
			if newStringValue, ok := GetEnvValue(w.currentPath); ok {
				UpdateData(valueExpression, newStringValue)
			}
			w.currentPath = w.currentPath[:len(w.currentPath)-1]
			return nil
		}
	case *js.ArrayExpr:
		for i, item := range n.List {
			if item.Value != nil {
				if valueExpression, ok := item.Value.(*js.LiteralExpr); ok {
					w.currentPath = append(w.currentPath, "["+fmt.Sprint(i)+"]")
					if newStringValue, ok := GetEnvValue(w.currentPath); ok {
						UpdateData(valueExpression, newStringValue)
					}
					w.currentPath = w.currentPath[:len(w.currentPath)-1]
				}
			}
		}
	}
	return w
}

func (w *walker) Exit(n js.INode) {
	// fmt.Println("Exit:", n)
	switch n.(type) {
	case *js.Property:
		w.currentPath = w.currentPath[:len(w.currentPath)-1]
	case *js.PropertyName:
	}
}

func GetEnvValue(path []string) (string, bool) {
	// TODO : use settingsVariableName instead
	computedKey := "AppSettings_"
	computedKey += strings.Join(path, "_")

	if envValue := os.Getenv(computedKey); envValue != "" {
		return envValue, true
	}

	return "", false
}

func UpdateData(valueExpression *js.LiteralExpr, newValue string) {
	if newValue == "" {
		return
	}

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
	settingsFolderPath, err := GetEnvOrError(SettingsFolderPathEnvKey)
	utils.HandleError(err)

	settingsFilePrefix, err := GetEnvOrError(SettingsFilePrefixEnvKey)
	utils.HandleError(err)

	settingsVariableName, err := GetEnvOrError(SettingsVariableNameEnvKey)
	utils.HandleError(err)

	utils.LogSuccess("✓ "+SettingsFolderPathEnvKey+": ", settingsFolderPath)
	utils.LogSuccess("✓ "+SettingsFilePrefixEnvKey+": ", settingsFilePrefix)
	utils.LogSuccess("✓ "+SettingsVariableNameEnvKey+": ", settingsVariableName)

	return settingsFolderPath, settingsFilePrefix, settingsVariableName
}

// https://eli.thegreenplace.net/2020/testing-flag-parsing-in-go-programs/
type CommandLineConfig struct {
	version bool

	// args are the positional (non-flag) command-line arguments.
	args []string
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
	flags.BoolVar(&conf.version, "version", false, "Display version and exit")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}
	conf.args = flags.Args()

	if conf.version {
		buf.WriteString(fmt.Sprintf("version : %s\n", version))
		buf.WriteString(fmt.Sprintf("commit  : %s\n", commit))
		buf.WriteString(fmt.Sprintf("date    : %s\n", date))
		buf.WriteString(fmt.Sprintf("builtBy : %s\n", builtBy))
	}

	return &conf, buf.String(), nil
}

func main() {
	config, output, err := ParseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		os.Exit(1)
	}

	if config.version {
		fmt.Println(output)
		os.Exit(0)
	}

	settingsFolderPath, settingsFilePrefix, settingsVariableName := GetConfigFileLocationValue()

	// TODO : aller chercher le nom du fichier de config dans le index.html
	settingsFilePath, errorDefineFilePath := DefineFilePath(settingsFolderPath, settingsFilePrefix)
	if errorDefineFilePath != nil {
		panic(errorDefineFilePath)
	}

	// Read the JavaScript file
	jsBytes, err := os.ReadFile(settingsFilePath)
	utils.HandleError(err)

	// Parse the JavaScript file
	input := parse.NewInput(bytes.NewReader(jsBytes))
	ast, err := js.Parse(input, js.Options{})
	utils.HandleError(err)

	// Analyse du code javascript et réalisation des modifications si nécessaire
	js.Walk(&walker{settingVariableName: settingsVariableName}, ast)

	// Write the updated JavaScript file
	// TODO : mettre à jour le fichier uniquement si des modifications ont été faite
	// TODO : afficher les modifications apportées
	var buffer bytes.Buffer
	ast.JS(&buffer)
	err = os.WriteFile(settingsFilePath, buffer.Bytes(), fs.ModePerm)
	utils.HandleError(err)

	utils.LogSuccess("🎉 Successfuly updated : ", settingsFilePath+" 🎉")
}

func GetEnvOrError(value string) (string, error) {
	env := os.Getenv(value)
	if env == "" {
		return "", errors.New("No environment key for : " + value)
	}
	return env, nil
}
