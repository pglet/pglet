echo "command1
command2" > $1

while read -r line; do
    echo $line
done <"$1"


echo "command3
command4" > $1

#sleep 2

while read -r line; do
    echo $line
done <"$1"


# while true
# do
#     if read line; then
#         # if [[ "$line" == '---' ]]; then
#         #     break
#         # fi
#         echo $line
#     else
#         echo "none"
#         break
#     fi
# done < "$1.result"


# while :; do
#     read line < "$1.result"
#     # if [[ "$line" == '---' ]]; then
#     #     break
#     # fi
#     echo $line
# done

# while read line < "$1.result"
# do
#     echo $line
# done
