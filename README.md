# project-kosmos-server
Server-side for https://github.com/seren1tyy/projectkosmos

## HOW TO USE

1. Install MariaDB 12+
2. Create users 'kosmos@localhost' and 'kosmos@127.0.0.1'
3. Alter it with 'mysql_native_password'
4. Сreate database 'kosmos_db' and grant all privileges to both of 'kosmos' users
5. In MariaDB, execute "source path/to/file" one by one, but be sure to execute schema.sql first
6. build the server "cmd/server/main.go"
7. finally run the server and the client

## WHAT IMPLEMENTED
* Simple auth and login system
* Session manager
* Character manager
* Item manager
* Universe manager
* Stub for universe
* Simple XOR encryption

## TODO
* Translate all comments to english
* Rework UI (client)
* Rework item manager if needed
* Add SSL encryption
* Market
* Simple corporation system
* Simple character card (server-side)
