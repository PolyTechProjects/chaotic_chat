for dir in $(ls -d */ | cut -f1 -d'/' | grep -iv "proto");
do
    echo "----------------------------------"
    echo "Launch $dir app"
    cd $dir
    for service in $(docker compose config --services | grep -i "_app");
    do
        if docker compose up $service -d; then
            echo "Docker Compose up $service successful"
            echo "$service: " $(docker ps | grep -i "$service" | awk '{print $1}')
        else
            echo "Docker Compose up $service failed"
            break
        fi
    done
    cd ../
    echo "----------------------------------"
done