
while true
do
    if read line; then
        echo $line
    fi
done < "$1.events"