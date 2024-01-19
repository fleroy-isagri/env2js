package main_test

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"

	// Local Module
	. "github.com/fleroy-isagri/env2js"
)

var _ = Describe("Main", func() {
	var mockUtils *MockUtils
	BeforeEach(func() {
		mockUtils = new(MockUtils)
		HandleError = mockUtils.HandleError
		LogSuccess = mockUtils.LogSuccess
	})

	AfterEach(func() {
		Getenv = os.Getenv
		Exit = os.Exit
		ReadFile = os.ReadFile
		WriteFile = os.WriteFile
	})

	Describe("IVisitor - When calling the walk function", func() {
		var mockOs *MockOs
		BeforeEach(func() {
			mockOs = new(MockOs)
		})

		Context("With a string value", func() {
			It("should erase the file key value with the new value", func() {
				// Arrange
				mockOs.On("Getenv", "AppSettings_MyKey").Return("Test1")
				Getenv = mockOs.Getenv
				// Act
				InterpretJSStringAsAst("const AppSettings = {MyKey: 'MyValue1'};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})

		Context("With a boolean value", func() {
			It("should erase the file key value with the new value", func() {
				// Arrange
				mockOs.On("Getenv", "AppSettings_MyKey").Return("true")
				Getenv = mockOs.Getenv
				// Act
				InterpretJSStringAsAst("const AppSettings = {MyKey: false};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})

		Context("With an integer value", func() {
			It("should erase the file key value with the new value", func() {
				// Arrange
				mockOs.On("Getenv", "AppSettings_MyKey").Return("1")
				Getenv = mockOs.Getenv
				// Act
				InterpretJSStringAsAst("const AppSettings = {MyKey: 0};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})

		Context("With an empty value", func() {
			It("should erase the file key value with the new value", func() {
				// Arrange
				mockOs.On("Getenv", "AppSettings_MyKey").Return("")
				Getenv = mockOs.Getenv
				// Act
				InterpretJSStringAsAst("const AppSettings = {MyKey: 'MyValue1'};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})

		Context("With an array value", func() {
			It("should show the current version", func() {
				// Arrange
				mockOs.On("Getenv", "AppSettings_MyArray_[0]").Return("Test1")
				mockOs.On("Getenv", "AppSettings_MyArray_[1]").Return("Test2")
				Getenv = mockOs.Getenv
				// Act
				InterpretJSStringAsAst("const AppSettings = {MyArray: ['MyValue1', 'MyValue2']};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})

		Context("With an invalid binding element value", func() {
			It("should not do anything", func() {
				// Act
				InterpretJSStringAsAst("const WrongBindingElement = {};")
				// Assert
				Expect(HasNotToPanic()).To(Equal(true))
			})
		})
	})

	Describe("DefineFilePath", func() {
		Context("When there is multiple files in the test folder", func() {
			It("should take the fisrt one in the folder", func() {
				want := "tests/example-1.js"
				got, err := DefineFilePath("./tests", "example")
				if err != nil {
					AbortSuite("File reading error")
				}
				Expect(filepath.ToSlash(got)).To(Equal(want))
			})

			It("should target the file with the closest name match in the folder", func() {
				want := "tests/example-2.js"
				got, err := DefineFilePath("./tests", "example-2")
				if err != nil {
					AbortSuite("File reading error")
				}
				Expect(filepath.ToSlash(got)).To(Equal(want))
			})

			It("should return an error if no file was found in the folder", func() {
				_, err := DefineFilePath("./tests", "toto")
				if err == nil {
					AbortSuite("File should not be find")
				}
				Expect(filepath.ToSlash(err.Error())).To(Equal("No file found with pattern: " + filepath.ToSlash("tests/toto*.js")))
			})
		})
	})

	Describe("GetConfigFileLocationValue", func() {
		Context("When one of the environment variable is missing", func() {
			var doesPanic bool
			var mockOs MockOs
			BeforeEach(func() {
				doesPanic = true
				// Create an instance of the mock.
				mockOs := new(MockOs)
				// Set expectations for the Getenv method.
				mockOs.On("Getenv", SettingsFolderPathEnvKey).Return(SettingsFolderPathEnvKey)
				mockOs.On("Getenv", SettingsFilePrefixEnvKey).Return(SettingsFilePrefixEnvKey)
				mockOs.On("Getenv", SettingsVariableNameEnvKey).Return(SettingsVariableNameEnvKey)
			})

			It("should panic regarding settingsFolderPath", func() {
				mockOs.On("Getenv", SettingsFolderPathEnvKey).Return("")
				// Replace the original os.Getenv with the mock implementation.
				Getenv = mockOs.Getenv
				defer ExpectToPanic(doesPanic, "settingsFolderPath panic")
				GetConfigFileLocationValue()
				doesPanic = false
			})
			It("should panic regarding settingsFilePrefix", func() {
				mockOs.On("Getenv", SettingsFilePrefixEnvKey).Return("")
				// Replace the original os.Getenv with the mock implementation.
				Getenv = mockOs.Getenv
				defer ExpectToPanic(doesPanic, "settingsFilePrefix panic")
				GetConfigFileLocationValue()
				doesPanic = false
			})
			It("should panic regarding settingsVariableName", func() {
				mockOs.On("Getenv", SettingsVariableNameEnvKey).Return("")
				// Replace the original os.Getenv with the mock implementation.
				Getenv = mockOs.Getenv
				defer ExpectToPanic(doesPanic, "settingsVariableName panic")
				GetConfigFileLocationValue()
				doesPanic = false
			})
		})
	})

	Describe("ParseFlags", func() {
		Context("The built program is executed with -version flag", func() {
			It("should show the current version", func() {
				// Arrange
				Version = "1.0.0"
				expectedOutput := "version : 1.0.0\ncommit  : \ndate    : \nbuiltBy : \n"

				// Act
				config, output, _ := ParseFlags("prog", []string{"-version"})

				// Assert
				Expect((config)).To(HaveExistingField("Version"))
				Expect((output)).To(Equal(expectedOutput))
			})
		})

		Context("The built program is executed with -help flag", func() {
			It("should display the command list", func() {
				// Arrange
				Version = "1.0.0"
				expectedOutput := "Usage of prog:\n  -version\n    \tDisplay version and exit\n"

				// Act
				config, output, err := ParseFlags("prog", []string{"-help"})

				// Assert
				Expect(config).To(BeNil())
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(flag.ErrHelp))
				Expect(output).To(Equal(expectedOutput))
			})
		})
	})

	Describe("LogVersion", func() {
		BeforeEach(func() {
			mockOs := new(MockOs)
			Exit = mockOs.Exit
		})

		It("should exit when LogVersion function is called with an error", func() {
			Expect(func() { LogFlags(nil, "", errors.New("An error")) }).To(Panic())
		})

		It("should exit when LogVersion function is called with an flag.ErrHelp", func() {
			Expect(func() { LogFlags(nil, "", flag.ErrHelp) }).To(Panic())
		})

		It("should exit when LogVersion function is called with an flag.ErrHelp", func() {
			Expect(func() { LogFlags(&CommandLineConfig{Version: true}, "output", nil) }).To(Panic())
		})
	})

	Describe("WriteInConfigFile", func() {
		It("should exit when LogVersion function is called with an error", func() {
			// Arrange
			mockOs := new(MockOs)
			ReadFile = mockOs.ReadFile
			WriteFile = mockOs.WriteFile

			// Assert
			Expect(func() { WriteInConfigFile("fileName", "variableName") }).NotTo(Panic())
		})
	})
})

////////////// HELPERS //////////////

func InterpretJSStringAsAst(jsString string) {
	// Parse the JavaScript file
	input := parse.NewInputString(jsString)
	ast, _ := js.Parse(input, js.Options{})
	// Analyse du code javascript et réalisation des modifications si nécessaire
	js.Walk(&Walker{SettingVariableName: "AppSettings"}, ast)
}

func ExpectToPanic(doesPanic bool, panicMessage string) {
	recover()
	Expect(doesPanic).To(Equal(true))
}

func HasNotToPanic() bool {
	if r := recover(); r == nil {
		return true
	} else {
		return false
	}
}

// MockOs is a mock implementation of the os interface.
type MockOs struct {
	mock.Mock
}

// Getenv is a mocked implementation of os.Getenv.
func (m *MockOs) Getenv(key string) string {
	args := m.Called(key)
	// args.String(0) will return the Return method parameter as string
	if args.String(0) != "" {
		if envValue := os.Getenv(key); envValue == "" {
			return args.String(0)
		} else {
			return envValue
		}
	} else {
		return ""
	}
}

func (m *MockOs) Exit(code int) {
	panic("Mock Exit panic with code : " + strconv.Itoa(code))
}

func (m *MockOs) ReadFile(name string) ([]byte, error) {
	return []byte("Mocked data"), nil
}

func (m *MockOs) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

// MockOs is a mock implementation of the os interface.
type MockUtils struct {
	mock.Mock
}

// HandleError is a mocked implementation of utils.HandleError.
func (m *MockUtils) HandleError(key error) {}

// LogSuccess is a mocked implementation of utils.LogSuccess.
func (m *MockUtils) LogSuccess(title string, log string) {}
