$ErrorActionPreference = "Stop"
$pipeName=$args[0]
$pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
$pipe.Connect(5000)
$pipeReader = new-object System.IO.StreamReader($pipe)
$pipeWriter = new-object System.IO.StreamWriter($pipe)
$pipeWriter.AutoFlush = $true

$eventPipe = new-object System.IO.Pipes.NamedPipeClientStream("$pipeName.events")
$eventPipe.Connect(5000)
$eventPipeReader = new-object System.IO.StreamReader($eventPipe)

function pglet_event {
    return $eventPipeReader.ReadLine()
}

function pglet {
    param (
        $command
    )

    # send command
    $pipeWriter.WriteLine($command)
    $pipeWriter.Flush()

    # parse results
    $OK_RESULT = "ok"
    $ERROR_RESULT = "error"
    
    $result = $pipeReader.ReadLine()
    if ($result -eq $OK_RESULT) {
        return ""
    } elseif ($result.StartsWith("$OK_RESULT ")) {
        return $result.Substring($OK_RESULT.Length + 1)
    } elseif ($result.StartsWith("$ERROR_RESULT ")) {
        throw $result.Substring($ERROR_RESULT.Length + 1)
    } else {
        throw "Unexpected result: $result"
    }
}

$rowId = pglet "add row id=body"
$colId = pglet "add col id=form to=$rowId"
pglet "add text value='Enter your name:' to=$colId"
pglet "add textbox id=fullName value='john smith' to=$colId"
pglet "add button id=submit text=Submit event=btn_event to=$colId"

$pipe.Close()

while($true) {
    pglet_event
}
