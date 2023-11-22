#!/bin/sh

echo "running"

is_up=false
is_down=false

while true; do
  case "$1" in
    -a|--app)
      application_name="$2"
      shift 2;;
    up)
      echo "up"
      is_up=true
      shift ;;
    down)
      echo "down"
      is_down=true
      shift ;;
    --|*)
      break;;
  esac
done

echo $is_up
echo $is_down

# shellcheck disable=SC2004
# shellcheck disable=SC2039
if (($is_up && $is_down)) || ((!$is_up && !$is_down)) ; then
  echo "invalid deploy type"
  exit 1
fi

if $is_up ; then
  if [ "$application_name" = "" ]; then
    echo "running all applications"
    docker-compose -f ./deployment/docker-compose.yml up -d
  else
    echo "running application $application_name"
    docker-compose -f ./deployment/docker-compose.yml up -d "$application_name"
  fi
fi

if $is_down ; then
  if [ "$application_name" = "" ]; then
      echo "running all application"
      docker-compose -f ./deployment/docker-compose.yml down
    else
      echo "running application $application_name"
      docker-compose -f ./deployment/docker-compose.yml down "$application_name"
  fi
fi

echo "done"