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

full_name=`curl $REPLIT_DB_URL/full_name`

function main() {
  pglet "clean page"
  rowID=`pglet "add row"`
  colID=`pglet "add col to=$rowID"`
  pglet "add text to=$colID value='Enter your name:'"
  pglet "add textbox to=$colID id=fullName value='$full_name'"
  pglet "add button to=$colID id=submit text=Submit"

  events=("submit click welcome")
  pglet_event "${events[@]}"
}

function welcome() {
  # get fullName value
  full_name=`pglet "get fullName value"`

  # save to Repl database
  curl $REPLIT_DB_URL -d "full_name=$full_name"

  # output welcome message
  pglet "clean page"
  pglet "add text value='Hello, $full_name'"
  pglet "add button id=again text=Again"

  events=("again click main")
  pglet_event "${events[@]}" 
}

eval 'main'

echo "Finished"