for dir in $(ls -d */ | cut -f1 -d'/' | grep -iv "proto");
do
    echo "----------------------------------"
    echo "Stop $dir services"
    cd $dir
    for service in $(docker compose config --services | grep -i "_app");
    do
        if docker compose stop $service; then
            echo "Docker Compose stop $service successful"
        else
            echo "Docker Compose stop $service failed"
            break
        fi
    done
    cd ../
    echo "----------------------------------"
done

docker ps