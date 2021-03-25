if [ $# -eq 0 ]; then
    echo "Usage: $0 [mongo|app]"
    exit 1
fi
echo "ssh into $1"

if [ $1 = "app" ];
then
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec app sh
else
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec mongo bash
fi
