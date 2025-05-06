# Stock Data Processor API

![Go Version](https://img.shields.io/badge/Go-1.20+-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Build Status](https://github.com/jesusignaciob/stock-api/actions/workflows/ci.yml/badge.svg)](https://github.com/jesusignaciob/stock-api/actions)

API para procesamiento de datos bursátiles con arquitectura hexagonal, paginación por cursor y almacenamiento por lotes.

## ✨ Características Destacadas

| Icono | Función | Descripción |
|-------|---------|-------------|
| ✅ | Arquitectura | **Hexagonal** bien estructurada para mantenibilidad |
| 📈 | Procesamiento | **Alta eficiencia** en manejo de datos de stocks |
| 🔄 | Paginación | **Por cursor** basada en ticker para navegación óptima |
| 🚀 | Almacenamiento | **Por lotes** optimizado para máxima velocidad |
| 💹 | API RESTFul | **API RESTful minimalista** con endpoints para obtener datos datos bursátiles |
| ⚡ | Performance | Gracias a GO...
| 📁 | Comando CLI | Para carga masiva de datos desde un API Externo
<!-- - 🐳 **Dockerizado** para fácil despliegue -->

## 📁 Estructura del Proyecto
```plaintext
/home/jbecerra/Test/Truora/stock-api/
├── cmd/                    # Archivos principales para iniciar la aplicación
│   └── main.go             # Punto de entrada de la aplicación
├── infrastructure/         # Código interno de la aplicación (lógica de negocio)
│   ├── domain/             # Entidades y lógica de dominio
│   ├── handler/            # Casos de uso de la aplicación
│   └── core/               # Implementaciones específicas (DB, APIs, etc.)
├── migrations/             # Migraciones
├── config/                 # Archivos de configuración
│   └── config.go           # Configuración principal
<!-- ├── scripts/                # Scripts útiles para desarrollo y despliegue -->
├── test/                   # Pruebas unitarias y de integración
<!-- ├── Dockerfile              # Dockerfile para construir la imagen de la aplicación -->
<!-- ├── docker-compose.yaml     # Configuración de Docker Compose -->
├── go.mod                  # Archivo de dependencias de Go
└── Makefile                # Makefile
```


## ⚙️ Makefile
```
Makefile for stock-api

Usage:
  make [target]

Targets:
  all            Build the application
  run            Run the application
  run-data       Run the fech data of the application
  build          Build the application
  test           Run tests
  clean          Clean build artifacts
  fmt            Format code
  lint           Lint code
  deps           Install dependencies
  analyze        Analyze code
  format         Format code
  fix            Fix lint issues
  migrate-up     Run database migrations up
  migrate-down   Run database migrations down
  help           Show this help message

Environment Variables:
  DB_HOST        Database host
  DB_PORT        Database port
  DB_USER        Database user
  DB_PASSWORD    Database password
  DB_NAME        Database name
  DB_SSLMODE     Database SSL mode

Note: Make sure to set the environment variables before running the make commands.

For more information, visit the project repository.

Happy coding!
```

## 🚀 API Endpoints  
**RESTful** y **high-performance** para gestión de stocks  

| Método | Ruta               | Descripción                          |  
|--------|--------------------|--------------------------------------|  
| GET    | `/api/stocks`      | Obtiene stocks (📈 cursor-paginated) |


## 📜 Licencia  
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  

Este proyecto está licenciado bajo **MIT License** - Ver el archivo [LICENSE](https://github.com/jesusignaciob/stock-api/blob/master/LICENSE) para términos completos.  