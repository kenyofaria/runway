docker compose down
docker build -t back-end ./back-end
docker build -t front-end ./front-end
docker compose up --build --force-recreate --remove-orphans

