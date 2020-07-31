echo "add label id=lbl2" > $1

echo "before read line"

IFS=''
while read -r line; do
    echo $line
done <"$1"


echo "add row id=\"footer\" to=lbl2 controls='
row
  col
row
  col'" > $1

echo "before read line"
sleep 1

while read -r line; do
    echo $line
done <"$1"

#echo "quit" > $1


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
