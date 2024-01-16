# env2js

[![GitHub Release](https://img.shields.io/github/v/release/fleroy-isagri/env2js)](https://github.com/fleroy-isagri/env2js/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/fleroy-isagri/env2js.svg)](https://pkg.go.dev/github.com/fleroy-isagri/env2js)
[![go.mod](https://img.shields.io/github/go-mod/go-version/fleroy-isagri/env2js)](go.mod)
[![LICENSE](https://img.shields.io/github/license/fleroy-isagri/env2js)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/fleroy-isagri/env2js/build.yml?branch=main)](https://github.com/fleroy-isagri/env2js/actions?query=workflow%3Abuild+branch%3Amain)

Env2js est une utilitaire permettant d'écrire dans un fichier de configuration en Javascript, la valeurs de variables d'environment.

### Exemple :

Compte tenu du fichier **config.js**
```js
const AppSettings = {
  isServed: true,
  API: {
    apiRoot: "url/server/app",
  },
}
```

Et de la **variable d’environnement** :
`AppSettings_API_apiRoot="custom/url/app"`


Après utilisation de env2js, **config.js** vaut :
```js
const AppSettings = {
  isServed: true,
  API: {
    apiRoot: "custom/url/app",
  },
}
```

### Utilisation :

**1) Installer les dépendances :**

```bash
go mod tidy
```

> :information_source: Plus d'infos sur [go mod tidy](https://go.dev/ref/mod#go-mod-tidy)

**2) Définir les trois variables d'environments nécessaires :**

**SETTINGS_FOLDER_PATH** : Dossier contenant le fichier de configuration à affecter

**SETTINGS_FILE_PREFIX** : Nom du fichier avant l'extension, ex : "example.js"

**SETTINGS_VARIABLE_NAME** : Nom de la clé du fichier à lire


`export SETTINGS_FOLDER_PATH=/path/to/my/config`

`export SETTINGS_FILE_PREFIX=example`

`export SETTINGS_VARIABLE_NAME=AppSettings`


**3) Lancer le programme :**

```bash
go run .
```
> :information_source: Plus d'infos sur [go run](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program)