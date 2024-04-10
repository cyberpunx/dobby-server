Dobby es una applicación web para asistir a los moderadores de Hogwarts Rol.

Stack:
* Go
* Echo (https://echo.labstack.com/)
* Turso (https://turso.tech/)
* Google Sheets API
* Go Templ (https://templ.guide/)
* HTMX (https://htmx.org/)


### Requisitos:
* Go 1.22 o superior
* Go Templ (https://templ.guide/)


# Para deployar:
* Se requiere una database en turso (https://turso.tech/)
* Se puede configurar de 3 formas:
  * crear un archivo conf.json
  * crear un archivo .env
  * configurar variables de entorno de sistema

## Mediante archivo conf.json
Crear un archivo conf.json en la raíz del proyecto con el siguiente contenido:
```json
{
  "serverPort": "8080",
  "tursoDbUrl": "libsql://<database>.turso.io",
  "tursoDbToken": "<turso_token>",
  "UnicodeOutput": "true"
}
```

## Mediante archivo .env o variables de entorno
Crear un archivo .env en la raíz del proyecto con el siguiente contenido:
```dotenv
TURSO_DB_URL=libsql://<database>.turso.io
TURSO_DB_TOKEN=<turso_token>
SERVER_PORT=8080
GSHEET_CLIENT_SECRET="contenido del archivo client_secret.json"
GSHEET_TOKEN="contenido del archivo token.json"
```

## Configuración de Google Sheets
* Se requiere un archivo client_secret.json para conectar con el Google Sheet de Moderación (Pedir a Duban)
* Se puede agregar el archivo client_secret.json en la raíz del proyecto
* También se puede agregar el contenido del archivo en una variable de entorno GSHEET_CLIENT_SECRET
* Se requiere un archivo token.json para autenticar con Google Sheets
  * Para generar un token.json, se debe ejecutar el server localmente y seguir el flow de autenticación con Google
* Se puede agregar el archivo token.json en la raíz del proyecto
* También se puede agregar el contenido del archivo en una variable de entorno GSHEET_TOKEN

Al ejecutarse el server se crearán automáticamente las tablas (si no existen) en la database de turso

