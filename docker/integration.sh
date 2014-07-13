# Bash script that spins up 'n' number of sshd instances for henchman to login to
# Tested on Mac OSX. In theory, if you set the DOCKER_HOSTNAME it should just work

# For OSX, if you are running boot2docker. Get the ip of the boot2docker vm using `boot2docker ip`.
# Once you have that invoke this script as follows
# $ sh docker/integration.sh <boot2docker_ip>

DOCKER_HOSTNAME=$1
NUM_HOSTS=$2
HENCHMAN_PLAN=$3

if [ -z $DOCKER_HOSTNAME ]
then
    echo "No docker host specified. Defaulting to 127.0.0.1"
    DOCKER_HOSTNAME="127.0.0.1"
fi

if [ -z $NUM_HOSTS ]
then
    echo "Setting number of hosts to 3"
    NUM_HOSTS=3
fi

if [ -z $HENCHMAN_PLAN ]
then
    echo "Setting plan to the sample"
    HENCHMAN_PLAN="samples/plan.yaml"
fi

echo "Starting all the SSHD containers"
HENCHMAN_HOSTS=""
for i in `seq 1 ${NUM_HOSTS}`
do
    ssh_port="${DOCKER_HOSTNAME}:320${i}"
    docker run -d -p $ssh_port:22 --name sshd${i} -t sudharsh/henchman:hosts
    if [[ $HENCHMAN_HOSTS == "" ]]
    then
        HENCHMAN_HOSTS="$ssh_port"
        continue
    fi
    HENCHMAN_HOSTS="$HENCHMAN_HOSTS,$ssh_port"
done
echo "Started ${NUM_HOSTS} containers"

echo "Invoking henchman with the plan: $HENCHMAN_PLAN"
bin/henchman -user root -private-keyfile docker/cmdcentre/insecure_private_key -args "hosts=$HENCHMAN_HOSTS" ${HENCHMAN_PLAN} 2> /dev/null
echo "*******"

rv=$?
echo
echo "Integration tests..."
if [ $rv == 0 ]
then
    echo "SUCCESS!"
else
    echo "FAILURE!"
fi

echo
echo "Spinning down containers"
echo "Takes a while depending on the number of containers"
for i in `seq 1 ${NUM_HOSTS}`
do
    docker stop "sshd${i}" > /dev/null
    docker rm "sshd${i}" > /dev/null
done

exit ${rv}

