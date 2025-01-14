for dir in $(ls -d */ | cut -f1 -d'/' | grep -iv "proto");
do
    echo "----------------------------------"
    echo "Stop $dir infra"
    cd $dir
    for service in $(docker compose config --services | grep -iv "_app");
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

echo "----------------------------------"
echo "Stop main infra"
for service in $(docker compose config --services | grep -iv "_app");
do
    if docker compose stop $service; then
        echo "Docker Compose stop $service successful"
    else
        echo "Docker Compose stop $service failed"
        break
    fi
done
echo "----------------------------------"

docker ps