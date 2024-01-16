package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

const (
	SettingsFolderPathEnvKey   string = "SETTINGS_FOLDER_PATH"
	SettingsFilePrefixEnvKey   string = "SETTINGS_FILE_PREFIX"
	SettingsVariableNameEnvKey string = "SETTINGS_VARIABLE_NAME"
)

type walker struct {
	currentPath         []string
	settingVariableName string
}

func (w *walker) Enter(n js.INode) js.IVisitor {
	// fmt.Println("Enter:", n)
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

	// if computedKey == "AppSettings_isServed" {
	// 	return "false", true
	// }
	// if computedKey == "AppSettings_API_apiVersion_versionMax" {
	// 	return "10", true
	// }
	// if computedKey == "AppSettings_MyKey" {
	// 	return "ModifiedValue", true
	// }
	// if computedKey == "AppSettings_MyArray_[1]" {
	// 	return "ModifiedArrayValue2", true
	// }

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

func main() {
	settingsFolderPath := os.Getenv(SettingsFolderPathEnvKey)
	// settingsFolderPath := "./config"
	if settingsFolderPath == "" {
		panic(SettingsFolderPathEnvKey + " environment variable is not set")
	}

	settingsFilePrefix := os.Getenv(SettingsFilePrefixEnvKey)
	// settingsFilePrefix := "example"
	// TODO : define a default value if not define
	if settingsFilePrefix == "" {
		panic(SettingsFilePrefixEnvKey + " environment variable is not set")
	}

	settingsVariableName := os.Getenv(SettingsVariableNameEnvKey)
	// settingsVariableName := "AppSettings"
	// TODO : define a default value if not define
	if settingsVariableName == "" {
		panic(SettingsVariableNameEnvKey + " environment variable is not set")
	}

	fmt.Println("SettingsFolderPath:", settingsFolderPath)
	fmt.Println("settingsFilePrefix:", settingsFilePrefix)
	fmt.Println("SettingVariableName:", settingsVariableName)

	// TODO : aller chercher le nom du fichier de config dans le index.html
	settingsFilePath, errorDefineFilePath := DefineFilePath(settingsFolderPath, settingsFilePrefix)
	if errorDefineFilePath != nil {
		panic(errorDefineFilePath)
	}

	// Read the JavaScript file
	jsBytes, err := os.ReadFile(settingsFilePath)
	if err != nil {
		panic(err)
	}

	// Parse the JavaScript file
	input := parse.NewInput(bytes.NewReader(jsBytes))
	ast, err := js.Parse(input, js.Options{})
	if err != nil {
		panic(err)
	}

	// Analyse du code javascript et réalisation des modifications si nécessaire
	js.Walk(&walker{settingVariableName: settingsVariableName}, ast)
	fmt.Println("END:", ast.JS())

	// Write the updated JavaScript file
	// TODO : mettre à jour le fichier uniquement si des modifications ont été faite
	// TODO : afficher les modifications apportées
	err = os.WriteFile(settingsFilePath, []byte(ast.JS()), fs.ModePerm)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(settingsFilePath + " updated with success")
	}
}
