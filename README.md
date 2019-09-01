# Implementación de Circuit Breaker en APIGo

Se realiza la implementación del patrón Circuit Breaker dentro de la APIGo
que mediante el ID de usuario devuelve los datos del usuario, sitio y país al que este pertenece.

Al realizar un GET al endpoint /resultsch/:userid en el puerto :9090, se pueden obtener diversas respuestas:

* Una respuesta correcta con los datos de usuario, sitio y país.
* Errores de cliente por Bad Request.
* Error 503 cuando se abre el interruptor del Circuit Breaker, por las siguientes razones:
    * La API de MELI responde con un error 500 (puede estar caída), 3 veces consecutivas.
    * Tiempo de respuesta de la API de MELI mayor a 3 segundos (sumando las 3 requests), 3 veces consecutivas.
    
## Autor
 * Marcos Postemsky
 

## Lógica del patrón Circuit Breaker

En la siguiente imagen se presenta la lógica de la máquina de estado que respeta el patrón.

![Texto alternativo](/res/img/Maquina de estado Circuit Breaker.png)



### Implementación de los timeouts