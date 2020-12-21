# AiuzuBotDiscord
AiuzuBot - Discord

Bot diseñado especificamente para administrar el servidor Virigamers.

Hay 3 niveles de permisos en este bot: Owner, Admin and Mod
Un argumento entre {} es obligatorio
Un argumento entre [] es opcional

Commands by level:

Owner:
reloadConfig: reload modifications to the roles in the configuration file
setStatus {status}: Actualiza el estatus del bot y borra el mensaje original
syncTodos: Revisa todos los usuarios del servidor y agrega a base de datos a los que no esten registrados, operacion pesada

Admin:

Mod:
detallesFull {nombre o id}: mostrar todos los detalles del usuario, excepto el desglose de las sanciones
detalles {nombre o id}: mostrar los detalles basicos de un usuario
detalleSanciones {nombre o id}: mostrar el detalle de las sanciones del usuario
ultimatum {userID} [razon]: Se pasa al usuario con ese ID a ultimatum, se actualiza en DB, se le quitan todos los roles y se le asigna solo el rol de Ultimatum
primerAviso {userID} [razon]: Si tiene derecho a primer aviso se aplica y notifica, si no lo tiene se notifica
sancion {id} [razon]: 
sancionFuerte {id} [razon]: 
Comandos Say personalizados, se configuran en config.json

Todos:
say {msg}: El bot dice lo que el comando le indique y borra el mensaje original