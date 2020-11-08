# AiuzuBotDiscord
AiuzuBot - Discord

Bot diseñado especificamente para administrar el servidor Virigamers.

Hay 3 niveles de permisos en este bot: Owner, Admin and Mod

Commands by level:

Owner:
reloadConfig: reload modifications to the roles in the configuration file
setStatus {status}: Actualiza el estatus del bot y borra el mensaje original
say {msg}: El bot dice lo que el comando le indique y borra el mensaje original
syncTodos: Revisa todos los usuarios del servidor y agrega a base de datos a los que no esten registrados, operacion pesada

Admin:

Mod:
detallesFull {nombre o id}: mostrar todos los detalles del usuario, excepto el desglose de las sanciones
detalles {nombre o id}: mostrar los detalles basicos de un usuario
detalleSanciones {nombre o id}: mostrar el detalle de las sanciones del usuario
