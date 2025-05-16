# Booking API

API REST para gestión de reservas, desarrollada en Go usando Fiber, GORM y JWT.

---

## Características

- Registro y autenticación de usuarios con JWT.
- Crear y listar reservas con validaciones:
  - Validación de horarios (entre 09:00 y 18:00).
  - No solapamiento de reservas.
  - Duración mínima de 1 hora.
- Control de errores con respuestas claras.
- Arquitectura modular con controladores, servicios y repositorios.
- Tests unitarios con mocks usando Testify.

---

## Tecnologías

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [JWT](https://github.com/golang-jwt/jwt)
- [Testify](https://github.com/stretchr/testify)

---


## Instalación

1. Clona el repositorio

```bash
git clone https://github.com/tuusuario/booking-api.git
cd booking-api
```

2.Configura las variables de entorno (ejemplo en .env):

Puedes crearte un link en Neon.tech y reemplazar el valor del DB_URL
link de [neon tech](https://neon.tech/)
link de [correoFake](https://temp-mail.org/es/)
Si te lo creas la cuenta y quieres el link debes:
1. Elegir el proyecto creado
2. Ve a Overview en el menu izquierdo 
3. Luego dentro dale en el boton 'Connect' y te aparecera un panel con la info de la BD que creaste
4. Busca 'Connection String' copia y pega la URL en tu .env.
```
DB_URL=postgresql://user:password@host:port/dbname?sslmode=require
JWT_SECRET=tu_secreto_jwt
```

3.Instala las dependencias:
```
go mod download
```
4. Inicia la aplicación:
```
go run main.go
```


## 📥 Endpoints

. Registro y login para obtener token JWT.
. Endpoints para crear y listar reservas (requieren autenticación).


1. URL (POST) - register:
```
http://localhost:3000/register
```
Body (JSON):
```
{
    "name": "Test",
    "email": "test@example.com",
    "password": "123456"
}
```
Respuesta esperada (status 200):
```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```


2. URL (POST) - login:
```
http://localhost:3000/login
```
Body (JSON):
```
{
  "email": "test@example.com",
  "password": "123456"
}
```
Respuesta esperada (status 200):
```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```


3. URL (POST) - reservations :
```
http://localhost:3000/reservations
```
Body (JSON):
```
{
    "date": "2025-05-23",
    "start_time": "15:00",
    "end_time": "16:00"
}
```
Respuesta esperada (status 201):
```
{
    "ID": 7,
    "CreatedAt": "2025-05-16T00:53:05.836307-05:00",
    "UpdatedAt": "2025-05-16T00:53:05.836307-05:00",
    "DeletedAt": null,
    "user_id": 1,
    "date": "2025-05-23",
    "start_time": "16:00",
    "end_time": "17:00"
}
```


4. URL (GET) - reservations?date :
```
http://localhost:3000/reservations?date=2025-05-23
```
Respuesta esperada (status 201):
```
[
    {
        "ID": 7,
        "CreatedAt": "2025-05-16T00:53:05.836307-05:00",
        "UpdatedAt": "2025-05-16T00:53:05.836307-05:00",
        "DeletedAt": null,
        "user_id": 1,
        "date": "2025-05-23",
        "start_time": "16:00",
        "end_time": "17:00"
    }
]
```

Link de la coleccion usada en postman [CollectionPostman](https://drive.google.com/file/d/1kjpVtM97l_cvcW8C4NvPFvHesAMIgX0Q/view?usp=sharing)



## Testing
Para ejecutar los tests:
```
go test ./...
```

## Estructura del proyecto

booking-api/
├── controllers/ # Controladores HTTP
├── models/ # Modelos y estructuras de datos
├── repositories/ # Acceso a datos
├── services/ # Lógica de negocio
├── middlewares/ # Middlewares personalizados
├── tests/ # Tests unitarios y de integración
├── utils/ # Utilidades generales
├── config/ # Configuración y variables de entorno
├── database/ # Conexión y migraciones a la base de datos
├── main.go # Archivo principal de la aplicación
├── .env # Variables de entorno
├── go.mod # Módulo y dependencias
└── go.sum # Checksum de dependencias
