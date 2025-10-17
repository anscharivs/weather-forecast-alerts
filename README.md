# Correr en GNU/Linux

## Clonar repo

* `git clone https://github.com/anscharivs/weather-forecast-alerts.git`
* `cd weather-forecast-alerts`

## Copiar .env (EL ENV A USAR CON EL API KEY ESTÁ EN EL DOCUMENTO DE DRIVE)
* `cp .env.example .env`

## Instalar docker

* `sudo apt update`
* `sudo apt install docker.io -y`

## Instalar docker compose

* `sudo curl -L https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose`

## Dar permisos

* `sudo chmod +x /usr/local/bin/docker-compose`

## Verificar instalación

* `docker-compose version`

## Iniciar docker

* `sudo systemctl enable docker`
* `sudo systemctl start docker`

## Hacer build del proyecto

* `sudo /usr/local/bin/docker-compose -f deployments/docker-compose.yml build --no-cache`

## Levantar

* `sudo /usr/local/bin/docker-compose -f deployments/docker-compose.yml up`
