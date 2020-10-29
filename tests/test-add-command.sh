echo "Pipe name: $1"
pipe_name=$1

function pglet() {
    # send command
    echo "$1" > "$pipe_name"

    # read result
    IFS=' ' read result_status result_value < "$pipe_name"
    echo $result_value
}

rowId=`pglet "add row id=body"`
colId=`pglet "add col id=form to=$rowId"`
pglet "add text value='Enter your name:' to=$colId"
pglet "add textbox id=fullName value='john smith' to=$colId"
pglet "add button id=submit text=Submit event=btn_event to=$colId"

IFS=' '
while read eventTarget eventName eventData < "$pipe_name.events"
do
  echo "$eventTarget $eventName $eventData"
done