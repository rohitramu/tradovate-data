#!/bin/bash
set -x -e -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

metabase_version="${1:-latest}"

image_name="metabase/metabase"
container_name="metabase"
container_port="3000"
local_port="3000"

backup_dir="$SCRIPT_DIR/metabase-backup"
timestamp="$(date --utc '+%Y%m%d_%H%M%S_Z')"
backup_file="$(ls -t $backup_dir/metabase_*.db | head -n1)" # Default to latest backup.
container_restore_path="$container_name:/metabase.db/metabase.db.mv.db"

docker pull $image_name:$metabase_version
if [ "$( docker container inspect -f '{{.State.Status}}' $container_name )" = "running" ]; then
    backup_file="$backup_dir/metabase_$timestamp.db"
    docker cp "$container_restore_path" "$backup_file"
    docker container stop $container_name
    docker container rm $container_name
fi

docker run -d -p $local_port:$container_port -v $SCRIPT_DIR/.db/:/db --restart always --name $container_name $image_name
sleep 10
docker cp "$backup_file" "$container_restore_path"
docker container restart $container_name
