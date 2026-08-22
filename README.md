# Seeyouvideo
Sistema de Gestión de Streaming
Aquí tienes el README completo, listo para pegar en README.md:

markdown
# SeeUvideo

Sistema de Gestión de Streaming desarrollado en Go, con persistencia en MySQL.

**Autor:** Ricardo Elihu Piñeiros Vera
**Fecha:** Agosto 2026

## Objetivo

Gestionar el catálogo de contenido y las suscripciones de usuarios de una
plataforma de streaming, controlando el acceso al contenido según el plan
del usuario (básico o premium).

## Funcionalidades principales

- **Registro e inicio de sesión** con validación de datos y contraseñas
  hasheadas con bcrypt.
- **Catálogo de videos** con búsqueda por categoría.
- **Control de acceso por plan**: usuarios básico solo visualizan
  contenido; usuarios premium pueden además agregar videos nuevos.
- **Interfaz de escritorio** construida con Fyne, con tema visual
  personalizado.
- **Página web** que muestra el catálogo con miniaturas de YouTube.
- **API REST** con 8 servicios web que serializan datos en formato JSON.
- **Persistencia en MySQL**, con creación automática de tablas.

## Servicios web (API REST)

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/registro` | Crea una cuenta nueva |
| POST | `/login` | Verifica credenciales |
| GET | `/usuarios/{id}?email=` | Perfil de usuario |
| PUT | `/usuarios/{id}/plan?email=` | Cambia el plan del usuario |
| GET | `/videos` | Lista todos los videos |
| GET | `/videos/{id}` | Video específico |
| POST | `/videos` | Agrega un video nuevo |
| GET | `/videos/categoria/{categoria}` | Filtra videos por categoría |
| GET | `/` | Página web del catálogo |

## Estructura del proyecto
Seeyouvideo/
├── main.go # Interfaz de escritorio (Fyne)
├── src/
│ ├── db/ # Conexión a MySQL
│ ├── usuarios/ # Struct Usuario, repositorio, autenticación
│ ├── videos/ # Struct Video, repositorio
│ └── api/ # Servidor web, 8 servicios REST, página HTML
└── docs/ # Documentación de planeación
## Tecnologías

- **Go 1.22+**
- **Fyne** — interfaz gráfica de escritorio
- **MySQL** — persistencia de datos
- **net/http** (estándar de Go) — servidor web y servicios REST

## Pruebas de software

El proyecto incluye 10 pruebas automatizadas:

- 4 pruebas unitarias (validación del struct `Usuario`)
- 3 pruebas de integración (`ServicioAuth` + repositorio)
- 2 pruebas de aceptación (peticiones HTTP simuladas)
- 1 prueba de concurrencia (20 peticiones simultáneas)

Para ejecutarlas:
```bash
go test ./...
```

## Cómo ejecutar

1. Crear la base de datos en MySQL:
```sql
CREATE DATABASE streaming_db;
```

2. Ajustar las credenciales de MySQL en `main.go` (`db.Config`).

3. Instalar dependencias y ejecutar:
```bash
go mod tidy
go run main.go
```

Esto abre la interfaz de escritorio y levanta el servidor web en
`http://localhost:8081`.
