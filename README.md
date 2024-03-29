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

En la siguiente imagen se presenta la lógica de la máquina de estado correspondiente al patrón.


<img src="https://github.com/marcospostemsky/APIGoCircuitBreaker/blob/master/res/img/Maquina%20de%20estado%20Circuit%20Breaker.png" width="550" >



### Implementación de los timeouts

Se implementan dos timeout de diferentes maneras, en el caso del encargado de limitar
el tiempo de la requests de resultados se realizar mediante channels y select-case. Para el timeout del ping (en el estado
Half-Open), se utiliza el provisto por el http.client (que podría extenderse a toda la implementación).

## Mock de comprobación
Para comprobar el correcto funcionamiento de la implementación se utiliza un [MOCK](https://github.com/marcospostemsky/MockSites) que remplaza la API de MELI.
Se configuran diferentes fallas y tiempos de respuesta para verificar el correcto funcionamiento de la APIGo. 
Además, se le suma al mock la capacidad de responder a un ping (como lo hace [MELI](https://api.mercadolibre.com/ping)).

