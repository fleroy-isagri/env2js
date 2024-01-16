
# env2js

[![GitHub Release](https://img.shields.io/github/v/release/fleroy-isagri/env2js)](https://github.com/fleroy-isagri/env2js/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/fleroy-isagri/env2js.svg)](https://pkg.go.dev/github.com/fleroy-isagri/env2js)
[![go.mod](https://img.shields.io/github/go-mod/go-version/fleroy-isagri/env2js)](go.mod)
[![LICENSE](https://img.shields.io/github/license/fleroy-isagri/env2js)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/fleroy-isagri/env2js/build.yml?branch=main)](https://github.com/fleroy-isagri/env2js/actions?query=workflow%3Abuild+branch%3Amain)

env2js is a tool that will transfer environment variable to a Javascript configuration file.

This is especially useful within virtual machines environments like Docker, where you want to modify the configuration based on the container environment variables.

## How it works :

Given the **config.js** file :
```js
const AppSettings = {
  isServed: true,
  API: {
    apiRoot: "url/server/app",
  },
}
```

and the **environment variables** :
`AppSettings_API_apiRoot="custom/url/app"`


After env2js usage, **config.js** equals :
```js
const AppSettings = {
  isServed: true,
  API: {
    apiRoot: "custom/url/app",
  },
}
```

## Start developing :

**1) Install the dependencies :**

```bash
go mod tidy
```

> :information_source: More informations about : [go mod tidy](https://go.dev/ref/mod#go-mod-tidy)

**2) Define the three required environment variables :**

**SETTINGS_FOLDER_PATH** : Folder that includes the configuration files

**SETTINGS_FILE_PREFIX** : File name without the extension, eg : "example.js"

**SETTINGS_VARIABLE_NAME** : Key name to read inside the file


`export SETTINGS_FOLDER_PATH=/path/to/my/config`

`export SETTINGS_FILE_PREFIX=example`

`export SETTINGS_VARIABLE_NAME=AppSettings`


**3) Run the program :**

```bash
go run .
```
> :information_source: More informations about : [go run](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program)


## Environment variables format

- **String** : `AppSettings_myValue="MyValue"`
- **Int** : `AppSettings_myInt=10`
- **Boolean** : `AppSettings_myBool=true`
- **Array index** : `AppSettings_MyArray_[0]="MyValue"`
- **Nested value** : `AppSettings_MyObject_MyValue="MyValue"`