# 🛡️ YouTube Toxicity Analyzer — Backend

> Microservicio en Go para análisis automático de toxicidad en comentarios de YouTube mediante DistilBERT.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat-square)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)
![Supabase](https://img.shields.io/badge/Supabase-PostgreSQL-3ECF8E?style=flat-square&logo=supabase&logoColor=white)
![DistilBERT](https://img.shields.io/badge/DistilBERT-HuggingFace-FFD21E?style=flat-square&logo=huggingface&logoColor=black)

---

## 📋 Tabla de Contenidos

- [Descripción](#-descripción)
- [Arquitectura](#-arquitectura)
- [Stack Tecnológico](#-stack-tecnológico)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Requisitos Previos](#-requisitos-previos)
- [Instalación y Configuración](#-instalación-y-configuración)
- [Variables de Entorno](#-variables-de-entorno)
- [API Reference](#-api-reference)
- [Base de Datos](#-base-de-datos)
- [Docker](#-docker)
- [Modelo de IA](#-modelo-de-ia)

---

## 📖 Descripción

Backend de un sistema de análisis inteligente de comentarios de YouTube basado en **Inteligencia Artificial** y **Procesamiento de Lenguaje Natural**.

El sistema permite introducir la URL de cualquier vídeo de YouTube y obtiene automáticamente:
- 📊 Métricas de toxicidad global del vídeo
- 💬 Análisis comentario a comentario
- 🏥 Health Score de la comunidad
- ⚠️ Alertas de comentarios críticos
- 💾 Persistencia de resultados en base de datos

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTE                              │
│              React Frontend (localhost:5173)                 │
└─────────────────────────┬───────────────────────────────────┘
                          │ POST /api/v1/analyze
                          │ { "video_url": "https://..." }
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    BACKEND (Go + Gin)                        │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │   YouTube    │  │  AI Service  │  │  Database Client │  │
│  │   Service    │  │  (DistilBERT)│  │  (Supabase/PG)   │  │
│  └──────┬───────┘  └──────┬───────┘  └────────┬─────────┘  │
└─────────┼─────────────────┼───────────────────┼────────────┘
          │                 │                   │
          ▼                 ▼                   ▼
┌──────────────┐  ┌──────────────────┐  ┌─────────────────┐
│  YouTube     │  │  FastAPI         │  │  Supabase       │
│  Data API v3 │  │  + DistilBERT    │  │  PostgreSQL     │
│  (Google)    │  │  (HuggingFace)   │  │                 │
└──────────────┘  └──────────────────┘  └─────────────────┘
```

### Flujo de Datos

```
1. Frontend envía URL del vídeo
        ↓
2. Go extrae el VideoID de la URL
        ↓
3. YouTube Data API v3 devuelve comentarios
        ↓
4. Cada comentario se envía al modelo DistilBERT
        ↓
5. DistilBERT clasifica: tóxico / no tóxico + confianza
        ↓
6. Go calcula métricas globales (promedio toxicidad, health score)
        ↓
7. Resultados se persisten en Supabase (padre → hijos)
        ↓
8. Response completa se devuelve al Frontend
```

---

## 🛠️ Stack Tecnológico

| Tecnología | Versión | Uso |
|-----------|---------|-----|
| **Go** | 1.26 | Lenguaje principal del backend |
| **Gin** | 1.12.0 | Framework HTTP web |
| **Google YouTube API** | v3 | Extracción de comentarios |
| **Supabase** | — | Base de datos PostgreSQL en la nube |
| **lib/pq** | 1.12.3 | Driver PostgreSQL para Go |
| **godotenv** | 1.5.1 | Gestión de variables de entorno |
| **Docker** | — | Containerización (multi-stage build) |
| **DistilBERT** | — | Modelo NLP en HuggingFace (servido vía FastAPI) |

---

## 📁 Estructura del Proyecto

```
Project_9_Backend_Team2/
│
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada — servidor Gin y endpoint principal
│
├── internal/
│   ├── ai/
│   │   └── client.go            # Cliente HTTP para comunicación con DistilBERT
│   ├── database/
│   │   └── supabase.go          # Cliente PostgreSQL — operaciones CRUD en Supabase
│   ├── dtos/
│   │   └── request.go           # DTOs: estructuras de entrada/salida de la API
│   └── youtube/
│       └── client.go            # Cliente oficial YouTube Data API v3
│
├── .dockerignore                # Archivos excluidos del contexto Docker
├── .env.example                 # Plantilla de variables de entorno
├── .gitignore                   # .env excluido del control de versiones
├── Dockerfile                   # Multi-stage build (builder + alpine)
├── go.mod                       # Módulo Go y dependencias
├── go.sum                       # Checksums de dependencias
└── README.md                    # Este archivo
```

---

## ✅ Requisitos Previos

- [Go 1.26+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) (opcional, para contenedor)
- Cuenta en [Google Cloud Console](https://console.cloud.google.com/) con **YouTube Data API v3** habilitada
- Cuenta en [Supabase](https://supabase.com/) con las tablas creadas
- API de IA (FastAPI + DistilBERT) corriendo localmente o en servidor

---

## 🚀 Instalación y Configuración

### 1. Clonar el repositorio

```bash
git clone https://github.com/Bootcamp-IA-P6/Project_9_Backend_Team2.git
cd Project_9_Backend_Team2
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Configurar variables de entorno

```bash
cp .env.example .env
```

Edita el archivo `.env` con tus credenciales (ver sección [Variables de Entorno](#-variables-de-entorno)).

### 4. Ejecutar en local

```bash
go run ./cmd/api
```

El servidor arrancará en `http://localhost:8080`

---

## 🔐 Variables de Entorno

Copia `.env.example` como `.env` y rellena los valores:

```env
# Puerto del servidor (por defecto: 8080)
PORT=8080

# API Key de YouTube Data API v3
# Obtener en: https://console.cloud.google.com/apis/credentials
YOUTUBE_API_KEY=AIza...

# URL de la API de IA (FastAPI + DistilBERT)
# Ejemplo local: http://localhost:8000/predict
# Ejemplo producción: https://tu-api.railway.app/predict
AI_MODEL_URL=http://localhost:8000/predict

# Cadena de conexión a Supabase (PostgreSQL)
# Formato: postgres://user:password@host:port/database?sslmode=require
DATABASE_URL=postgres://postgres.xxx:password@aws-0-eu-central-1.pooler.supabase.com:5432/postgres
```

> ⚠️ **Nunca** subas el archivo `.env` al repositorio. Ya está incluido en `.gitignore`.

---

## 📡 API Reference

### POST `/api/v1/analyze`

Analiza los comentarios de un vídeo de YouTube.

**Request Body:**
```json
{
  "video_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}
```

**Response exitosa (200):**
```json
{
  "video_id": "dQw4w9WgXcQ",
  "total_comments": 10,
  "promedio_toxicidad": 30.0,
  "alertas_criticas": 3,
  "health_score": 70.0,
  "detalle": [
    {
      "autor": "Usuario123",
      "texto": "This is a great video!",
      "evaluacion": {
        "texto_analizado": "This is a great video!",
        "etiqueta_modelo": "LABEL_0",
        "score_confianza": 0.9823,
        "es_toxico": false
      }
    },
    {
      "autor": "Usuario456",
      "texto": "You are an idiot",
      "evaluacion": {
        "texto_analizado": "You are an idiot",
        "etiqueta_modelo": "LABEL_1",
        "score_confianza": 0.9145,
        "es_toxico": true
      }
    }
  ]
}
```

**Response de error — comentarios desactivados (400):**
```json
{
  "error": "Este vídeo tiene los comentarios desactivados o no existe."
}
```

### Campos de la respuesta

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `video_id` | string | ID del vídeo de YouTube |
| `total_comments` | int | Total de comentarios analizados |
| `promedio_toxicidad` | float | Porcentaje de comentarios tóxicos (0-100) |
| `alertas_criticas` | int | Número de comentarios clasificados como tóxicos |
| `health_score` | float | Salud de la comunidad (100 - promedio_toxicidad) |
| `detalle` | array | Array con análisis individual de cada comentario |

### Campos de evaluación por comentario

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `texto_analizado` | string | Texto original del comentario |
| `etiqueta_modelo` | string | Etiqueta raw del modelo (LABEL_0 / LABEL_1) |
| `score_confianza` | float | Confianza de la predicción (0.0 - 1.0) |
| `es_toxico` | bool | `true` si el comentario es tóxico |

---

## 🗄️ Base de Datos

El sistema utiliza **Supabase (PostgreSQL)** con dos tablas relacionadas:

### Tabla `VideoSummary` (Padre)

Almacena el resumen de análisis por vídeo.

```sql
CREATE TABLE "VideoSummary" (
    id                  BIGSERIAL PRIMARY KEY,
    video_id            TEXT NOT NULL,
    total_comments      INT,
    promedio_toxicidad  FLOAT,
    alertas_criticas    INT,
    health_score        FLOAT,
    created_at          TIMESTAMP DEFAULT NOW()
);
```

### Tabla `ToxicFilterAIPredictions` (Hijo)

Almacena las predicciones individuales por comentario.

```sql
CREATE TABLE "ToxicFilterAIPredictions" (
    id          BIGSERIAL PRIMARY KEY,
    comment_id  TEXT,
    video_id    TEXT NOT NULL,
    text        TEXT,
    is_toxic    BOOLEAN,
    created_at  TIMESTAMP DEFAULT NOW()
);
```

### Relación

```
VideoSummary (1) ──── (N) ToxicFilterAIPredictions
     video_id                    video_id
```

> ⚠️ El backend inserta primero el **padre** (VideoSummary) y después los **hijos** (ToxicFilterAIPredictions) para mantener la integridad referencial.

---

## 🐳 Docker

### Build de la imagen

```bash
docker build -t youtube-toxicity-backend .
```

### Ejecutar el contenedor

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e YOUTUBE_API_KEY=tu_api_key \
  -e AI_MODEL_URL=http://host.docker.internal:8000/predict \
  -e DATABASE_URL=postgres://... \
  youtube-toxicity-backend
```

### Docker Compose (con la API de IA)

```yaml
version: '3.8'
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - YOUTUBE_API_KEY=${YOUTUBE_API_KEY}
      - AI_MODEL_URL=http://ai-service:8000/predict
      - DATABASE_URL=${DATABASE_URL}
    depends_on:
      - ai-service

  ai-service:
    image: tu-imagen-fastapi
    ports:
      - "8000:8000"
```

### Multi-stage Build

El `Dockerfile` usa **multi-stage build** para optimizar el tamaño de la imagen:

```
Etapa 1 (builder): golang:1.26-alpine
  → Compila el binario estático (CGO_ENABLED=0)
  → Resultado: binario "main" ~10MB

Etapa 2 (producción): alpine:latest
  → Solo el binario + certificados CA
  → Imagen final: ~15MB (vs ~300MB con Go completo)
```

---

## 🧠 Modelo de IA

El backend se comunica con una **API FastAPI** que sirve el modelo **DistilBERT** fine-tuned para detección de toxicidad.

### Modelo utilizado

- **Base:** `distilbert-base-uncased`
- **Fine-tuned:** Dataset YouToxic English (1,000 comentarios)
- **Disponible en:** HuggingFace Hub — `PROJECT_9_NLP_Team2/distilbert-toxicity`

### Métricas del modelo

| Métrica | Valor |
|---------|-------|
| F1-Score (Test) | 77.85% |
| Recall (Test) | 86.57% |
| ROC-AUC | 85.29% |
| Overfitting F1 | 3.48% ✅ |

### Contrato de la API de IA

**Request que envía Go:**
```json
{
  "texto": "You are an idiot"
}
```

**Response que espera Go:**
```json
{
  "texto_analizado": "You are an idiot",
  "etiqueta_modelo": "LABEL_1",
  "score_confianza": 0.9145,
  "es_toxico": true
}
```

### Descargar el modelo localmente

```bash
python scripts/download_model_huggingface.py \
    --repo_name PROJECT_9_NLP_Team2/distilbert-toxicity \
    --output_path ./modelo_descargado
```

---

## 🤝 Contribución

1. Crea una rama desde `main`: `git checkout -b feat/nueva-funcionalidad`
2. Haz commits descriptivos: `git commit -m "feat(api): add pagination to analyze endpoint"`
3. Abre un Pull Request con descripción detallada

---

## 👥 Equipo

Este proyecto fue desarrollado por el **Grupo 1 - Bootcamp IA P6**:

| Nombre | Rol |
| :--- | :--- |
| **Gema Yébenes** | AI Product Owner |
| **Naizabeth Bermudez** | AI Scrum Master |
| **Maryori Cruz** |AI Full Stack Developer |
| **Raúl Machaca** | AI Full Stack Developer |



---

*Construido con ❤️ usando Go, DistilBERT y YouTube Data API v3*