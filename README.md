## Inicialización

Para correr el proyecto ejecutar en el working directory:

```
docker-compose up
```

Este docker levanta 
- la app en el puerto __:8080__, 
- la DB en el __:5432__ y, 
- un Swagger con la especificación de los endpoints en el puerto __:3000__

## Endpoints 

1. Obtener 50 IPs de Argentina (IP, pais y ciudad)
    
    ```
    GET /v1/ips?limit=50&country=Argentina
   ```
2. Obtener toda a información disponible en la base para una determinada dirección IP (la IP debe ser un parámetro)
    
    ```
    GET /v1/ips/{ip}
   ```

3. Obtener el TOP 10 de ISP de Suiza (country code: CH - serían los ISP que más se repiten para este país)

    ```
    GET /v1/ips/isps/top?country=Switzerland
   ```
   
4. Obtener cantidad de IPs por país (país debe ser un parámetro)

     ```
    GET /v1/ips/quantity?country=Argentina
   ```

Traté de tener un diseño orientado a paquetes pensando en la funcionalidad.

Dentro de `cmd/api` se encuentran todos los archivos para inicialización de la API, router, handlers, middlewares y modelos de response.

Dentro de `internal` todos los archivos con la lógica de negocio y bootstrap de la app.
Hay un package que maneja todo lo que la obtención de información de IPs y hay uno que maneja la conversión, este último no es un
servicio sino que es un package con funciones asi como lo es el package `net`

Realicé la obtención de datos del request en los handlers y no en los middlewares ya que estos últimos los suelo utilizar para
validar por ej. seguridad, permisos, etc... Y de obtener en los middlewares información del request habría que pasarla a través del contexto.

Resolví todos los endpoints con queries a la DB por la cantidad de datos que contiene para evitar traerlos a memoria todo el tiempo y
filtrarlos por código. Para mejorar las búsquedas podrían agregarse índices a la tabla, y tener cacheada la respuesta de los requests que más se realizan
o de los requests realizados en los últimos 5 minutos.

Todos los archivos están testeados salvo los destinados a bootstrap de la aplicación y las entidades. 

## Puntos a mejorar

- El endpoint (1) devuelve 50 IPs pero sólo devuelve teniendo en cuenta el campo _ip_from_ no obtiene los conjuntos de IPs
entre _ip_from_ y _ip_to_. Una mejora sería obtener los registros contando cuántas IPs están guardadas en cada registro, porque
puede darse el caso de que para un país existan 50 ips agrupadas en 20 columnas y el endpoint sólo devolverá 20.

- Mejorar el build de las queries para poder tener más parámetros opcionales (por ej. country en esta versión es un parámetro
obligatorio cuando al ser un filtro tendría más sentido que sea opcional) Permitir más parámetros de filtrado.

- Validar query param country y cualquier query param de filtrado que se agregue, el query param `limit` lo restringí a 100 registros de no especificarse.

- Las configuraciones de la aplicación están cargadas en un .env, poco seguro y cada vez que se necesite realizar un cambio de
configuración requiere modificar el archivo dentro de la aplicación y volver a deployarla. 
Modificaría esto en un productivo o de desarrollo para guardar las contraseñas o data sensible como secrets dentro de un servicio como Vault

- Otro punto que va de la mano con el tema configuraciones es la posibilidad de tener distintas inicializaciones dependiendo del
environment en el que se esté corriendo la aplicación, quizas se necesita levantar distintas rutas o no correr ciertos procesos internos.
