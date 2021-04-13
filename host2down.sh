docker-compose -f host2.yaml down -v
docker rm -f $(docker ps -a -q)
docker volume rm $(docker volume ls -q)
docker rmi $(docker images dev-* -q)
