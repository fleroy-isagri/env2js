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