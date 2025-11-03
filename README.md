# Strade

API para consultar información de códigos postales de México basada en los datos oficiales del Servicio Postal Mexicano (SEPOMEX).

## Uso de Datos del Servicio Postal Mexicano (SEPOMEX)

**Strade** utiliza información proveniente del *Catálogo Nacional de Códigos Postales*, elaborado y publicado por el **Servicio Postal Mexicano (SEPOMEX)**.  
Dicho catálogo es de acceso público y gratuito, conforme a la información disponible en el portal oficial de SEPOMEX.

La base de datos original puede descargarse directamente desde el sitio oficial:  
[https://www.correosdemexico.gob.mx/](https://www.correosdemexico.gob.mx/)

De acuerdo con el aviso del Servicio Postal Mexicano:  
> “El Catálogo Nacional de Códigos Postales es elaborado por el Servicio Postal Mexicano y se proporciona en forma gratuita, no estando permitida su comercialización, total o parcial.”

En cumplimiento con lo anterior, **Strade** no comercializa los datos contenidos en el catálogo.  
La información proporcionada por este servicio se distribuye únicamente con fines técnicos y de conveniencia, como parte de una interfaz programática (API) que facilita el acceso automatizado a los datos públicos publicados por SEPOMEX.  
Dicha distribución no implica su venta ni su redistribución con fines de lucro.

Cualquier uso de este proyecto que implique redistribución, modificación o monetización de los datos deberá respetar las restricciones establecidas por SEPOMEX.  
El mantenimiento, alojamiento o servicios derivados ofrecidos a través de **Strade** se limitan a proveer infraestructura y funcionalidad adicional, sin alterar la naturaleza gratuita de los datos originales.

## Desarrollo con Docker

El proyecto utiliza Docker Compose para gestionar los servicios necesarios en desarrollo y producción.

### Estructura

```
docker/
├── Dockerfile.api          # Multi-stage build para API
├── Dockerfile.watcher      # Multi-stage build para Watcher
├── docker-compose.yml      # Configuración de producción
├── docker-compose.dev.yml  # Configuración de desarrollo
└── entrypoint.sh          # Script de inicialización y migraciones
```

### Configuración Inicial

Copiar el archivo de ejemplo y configurar variables de entorno:

```bash
cp .env.example .env.dev
```

Editar `.env.dev` con los valores apropiados. Para desarrollo con Docker, asegurar:

```bash
# Database (usar 'db' como host dentro de Docker)
DB_ADDR=postgres://postgres:password@db:5432/strade-db?sslmode=disable
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=strade-db

# Redis (usar 'redis' como host dentro de Docker)
REDIS_ADDR=redis:6379
REDIS_PW=
REDIS_PASSWORD=
```

### Comandos Disponibles

Desarrollo:

```bash
# Iniciar solo DB y Redis (para desarrollo local de Go)
make docker-dev-up

# Iniciar stack completo (DB, Redis, API, Watcher)
make docker-dev-full-up

# Detener servicios
make docker-dev-down

# Ver logs en tiempo real
make docker-dev-logs
```

Producción:

```bash
# Iniciar stack completo
make docker-up

# Detener servicios
make docker-down

# Ver logs
make docker-logs

# Limpiar todo (incluyendo volúmenes)
make docker-clean
```

### Servicios

- **PostgreSQL** - Base de datos principal (puerto 5432 en dev, interno en prod)
- **Redis** - Cache y message broker (puerto 6379 en dev, interno en prod)
- **API** - Servidor HTTP (puerto 8080)
- **Watcher** - Servicio de monitoreo y sincronización (puerto 8081)

## Documentación de la API

La documentación completa de la API está disponible en formato OpenAPI (Swagger) en la siguiente ruta cuando el servidor está en ejecución:

```
http://localhost:8080/docs
```

Para acceder a la documentación, asegúrate de que el servicio de la API esté en ejecución. La documentación incluye todos los endpoints disponibles, sus parámetros, formatos de solicitud y ejemplos de respuestas.

### Migraciones

Las migraciones de base de datos se ejecutan automáticamente al iniciar los contenedores mediante `entrypoint.sh`. Si las migraciones fallan, el contenedor no iniciará.

Para ejecutar migraciones manualmente:

```bash
# Aplicar migraciones pendientes
make migrate

# Revertir última migración
make migrate-down

# Resetear base de datos (drop + up)
make migrate-reset

# Crear nueva migración
make create-migration nombre_migracion
```

### Arquitectura Docker

Desarrollo:
- Restart policy: `no` (facilita debugging)
- Puertos expuestos: DB y Redis accesibles desde host
- Redis sin autenticación
- Resource limits reducidos

Producción:
- Restart policy: `unless-stopped` (alta disponibilidad)
- Puertos internos: DB y Redis solo accesibles dentro de red Docker
- Redis con autenticación requerida
- Resource limits definidos por servicio
- Multi-stage builds para imágenes optimizadas

Todos los servicios se comunican a través de una red Docker aislada llamada `backend`.

### Seguridad

En producción:
- Base de datos y Redis no exponen puertos al host
- Redis requiere autenticación mediante `REDIS_PASSWORD`
- Imágenes multi-stage reducen superficie de ataque
- Health checks configurados para todos los servicios
- Resource limits previenen consumo excesivo de recursos
