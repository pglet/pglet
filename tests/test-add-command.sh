echo "Pipe name: $1"
page_pipe=$1

function pglet() {
    # send command
    echo "$1" > "$page_pipe"

    # read result
    IFS=' ' read result_status result_value < "$page_pipe"
    echo $result_value
}

function pglet_event() {
  # https://askubuntu.com/questions/992439/bash-pass-both-array-and-non-array-parameter-to-function
  arr=("$@")
  IFS=' '
  while read eventTarget eventName eventData < "$page_pipe.events"
  do
    for evt in "${arr[@]}";
    do
      IFS=' ' read -r et en fn <<< "$evt"
      if [[ "$eventTarget" == "$et" && "$eventName" == "$en" ]]; then
        eval "$fn"
        return
      fi      
      #echo "$et - $en - $fn"
    done
  done
}

function main() {
  pglet "clean page"
  rowId=`pglet "add row id=body"`
  colId=`pglet "add col id=form to=$rowId"`
  pglet "add text value='Enter your name:' to=$colId"
  pglet "add textbox id=fullName value='john smith' to=$colId"
  pglet "add button id=submit text=Submit event=btn_event to=$colId"

  events=("body:form:submit click welcome")
  pglet_event "${events[@]}"
}

function welcome() {
  # get fullName value
  full_name=`pglet "get body:form:fullName value"`

  # output welcome message
  pglet "clean page"
  pglet "add text value='Hello, $full_name'"
  pglet "add button id=again text=Again"

  events=("again click main")
  pglet_event "${events[@]}" 
}

eval 'main'

echo "Finished"