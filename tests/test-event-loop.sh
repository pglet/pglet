echo "Event pipe: $1"

for i in 1 2 3 4 5
do
    if read line < "$1"; then
        echo $line
    fi
    sleep 1
done

# while true
# do
#     if read line; then
#         echo $line
#     fi
# done < "$1"