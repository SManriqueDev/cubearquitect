# 🏗️ CubeArquitect

**Editor visual para diseñar y desplegar arquitecturas de infraestructura en CubePath**

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![CubePath](https://img.shields.io/badge/CubePath-1E40AF?style=for-the-badge&logo=cloudflare&logoColor=white)](https://cubepath.com/)
[![Fiber](https://img.shields.io/badge/Fiber-00ADD8?style=for-the-badge&logo=&logoColor=white)](https://gofiber.io/)

---

CubeArquitect es una herramienta que permite crear arquitecturas de infraestructura mediante un editor visual de nodos. Cada nodo representa un VPS que se despliega automáticamente en CubePath con código boilerplate preconfigurado según su tipo.

- **Nodos**: App (aplicaciones) y Database (bases de datos)
- **Sistema DAG**: Despliegue por niveles - los nodos sin dependencias entre sí se despliegan en paralelo
- **Inyección de variables**: Las conexiones entre nodos permiten inyectar automáticamente valores (ej: DATABASE_URL de Database → App)

## ✨ Características

- 🎨 **Editor visual drag-and-drop** para diseñar arquitecturas de infraestructura
- 🖥️ **Nodos con boilerplate**: App y Database vienen con código preconfigurado listo para desplegar
- 🔗 **Sistema DAG**: Despliegue por niveles - nodos sin dependencias se despliegan en paralelo
- 💉 **Inyección de variables**: Conecta nodos y las dependencias (ej: DATABASE_URL) se injectan automáticamente
- 🚀 **Despliegue automático** en CubePath desde el editor visual
- 📡 **Logs en tiempo real** del proceso de despliegue
- ⚙️ **Panel de configuración** por tipo de nodo
- 💰 **Calculadora de precios** integrada

## 🧱 Demo

🚀 **Prueba la aplicación:** [http://vps23511.cubepath.net:3001/](http://vps23511.cubepath.net:3001/)

## 🚀 Cómo comenzar

### Prerrequisitos

- Go 1.21+
- Node.js 18+
- Credenciales de CubePath (API Token)

### Instalación

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/cubearquitect.git
cd cubearquitect

# Backend
cd backend
go mod download
go run cmd/api/main.go

# Frontend (en otra terminal)
cd ../frontend
npm install
npm run dev
```

## 🛠️ Tecnologías

| Capa | Tecnología |
|------|------------|
| Backend | Go + Fiber |
| Frontend | React + TypeScript + Vite |
| UI | shadcn/ui + Tailwind CSS |
| Editor visual | React Flow |
| Deployment | **CubePath** |

## 📸 Screenshots

*[PLACEHOLDER: Añadir screenshot del editor visual - arrastrar nodos]*

*[PLACEHOLDER: Añadir screenshot del panel de configuración de nodo]*

*[PLACEHOLDER: Añadir screenshot de logs en tiempo real durante despliegue]*

## 📂 Estructura del proyecto

```
cubearquitect/
├── backend/                   # API en Go + Fiber
│   ├── cmd/api/main.go       # Punto de entrada
│   └── internal/
│       ├── cubepath/         # Cliente de CubePath
│       ├── orchestrator/     # Motor de despliegue
│       │   ├── blueprints_app.go
│       │   ├── blueprints_database.go
│       │   ├── engine.go
│       │   └── dag.go
│       └── service/          # Lógica de negocio
└── frontend/                  # App React
    └── src/
        ├── components/
        │   ├── flow/         # Editor visual
        │   └── nodes/        # Nodos personalizados
        ├── hooks/           # Custom hooks
        ├── services/        # API y servicios
        └── stores/          # Estado global
```

## 🗺️ Roadmap

- [ ] Soporte para más tipos de nodos (Redis, Cache, etc.)
- [ ] Templates predefinidos de arquitecturas
- [ ] Guardar y cargar proyectos
- [ ] Sistema de autenticación de usuarios
- [ ] Dashboard de proyectos y despliegues
- [ ] Historial de despliegues

## ☁️ Despliegue en CubePath

El proyecto está desplegado en **CubePath** utilizando **Dokploy** (Docker Compose):

- **Infraestructura**: Docker Compose en un VPS de CubePath
- **Servicios**: Backend (Go + Fiber) + Frontend (React + Vite) en el mismo VPS
- **Puerto**: 3001 (expuesto via Docker)

### Despliegue con Dokploy

El proyecto utiliza un archivo `docker-compose.yml` para orquestar los servicios:

```yaml
services:
  backend:
    build: ./backend
    environment:
      - CUBE_API_URL=${CUBE_API_URL:-https://api.cubepath.com}
      - PORT=${PORT:-8080}

  frontend:
    build: ./frontend
    ports:
      - "3001:80"  # Puerto 3001 del host -> Puerto 80 del contenedor
```

El backend utiliza el cliente oficial de CubePath para:
1. Crear VPS según la configuración de cada nodo
2. Configurar las conexiones entre nodos
3. Obtener información de estado en tiempo real

## 📝 Requisitos del Hackaton

✅ Proyecto desplegado en **CubePath**  
✅ Repositorio público  
✅ README con descripción, demo y screenshots  
✅ Explicación del uso de CubePath  

---

Hecho con ❤️ para la **Hackaton CubePath 2026**
