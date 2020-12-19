res=`pglet page page2`

IFS=' ' read -r page_pipe page_url <<< "$res"

echo "Pipe name: $page_pipe"
#page_pipe="$1"

function pglet_send() {
  # send command
  echo "$1" > "$page_pipe"

  # read result
  local firstLine="true"
  local result_value=""
  IFS=''
  while read -r line; do
    if [[ $firstLine == "true" ]]; then
      IFS=' ' read -r result_status result_value <<< "$line"
      firstLine="false"
    else
      result_value="$line"
    fi
    echo "$result_value"
  done <"$page_pipe"
}

function pglet_event() {
  # https://askubuntu.com/questions/992439/bash-pass-both-array-and-non-array-parameter-to-function
  arr=("$@")
  IFS=' '
  while true
  do
    read eventTarget eventName eventData < "$page_pipe.events"
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
  pglet_send "clean page"
  rowID=`pglet_send "add row"`
  colID=`pglet_send "add col to=$rowID"`
  pglet_send "add text to=$colID value='Enter your name:'"
  pglet_send "add textbox to=$colID id=fullName value='$full_name' multiline=true"
  pglet_send "add button to=$colID id=submit text=Submit"

  events=("submit click welcome")
  pglet_event "${events[@]}"
}

function welcome() {
  # get fullName value
  echo "before get name"
  full_name=`pglet_send "get fullName value"`
  full_name="${full_name//$'\n'/\\n}"
  echo "full name: $full_name"

  # output welcome message
  pglet_send "clean page"
  pglet_send "add text value='Hello, $full_name'"
  pglet_send "add button id=again text=Again"

  events=("again click main")
  pglet_event "${events[@]}" 
}

eval 'main'

echo "Finished"