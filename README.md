# Projeto integrador 2020 [![Go Report Card](https://goreportcard.com/badge/github.com/MrNinso/ProjetoUnivesp2020)](https://goreportcard.com/report/github.com/MrNinso/ProjetoUnivesp2020)

## Build
**Dependencies**

- go >= 1.14
- docker (client in $PATH) >= 19.03.11 
- yarn >= 1.21.1
- openssl >= 1.6.1
- make >= 4.2.1

````shell script
git clone https://github.com/MrNinso/ProjetoUnivesp2020.git
cd ProjetoUnivesp2020
make build
````

if you have Vagrant u can deploy just do:

 ````shell script
 make vagrant-build
 ````
after the deploy just acess `` https://<your server ip>:1443/app `` and use user `` admin@admin.com `` with `` admin `` password
