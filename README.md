# env2js

## Utilisation / Test
 
### Installer les dépendances
 
```bash
go mod tidy
```
 
### Définir les trois variables d'environments nécessaires
 
```js
// Dossier contenant le fichier de configuration à affecter
export SETTINGS_FOLDER_PATH=../tests
// Nom du fichier avant l'extension, ex : "example.js"
export SETTINGS_FILE_PREFIX=example
// Nom de la clé du fichier à lire
export SETTINGS_VARIABLE_NAME=AppSettings

// Définir les propriété à modifier
export AppSettings_MyKey=TATA

```
 
 ### Lancer le programme
 
```bash
go run .
```