$ErrorActionPreference = "Stop"

$pipe = $null
$pipeReader = $null
$pipeWriter = $null
$eventPipe = $null
$eventPipeReader = $null

function pglet_event {
    $line = $eventPipeReader.ReadLine()
    Write-Host "Event: $line"
    if ($line -match "(?<target>[^\s]+)\s(?<name>[^\s]+)(\s(?<data>.+))*") {
        return @{
            Target = $Matches["target"]
            Name = $Matches["name"]
            Data = $Matches["data"]
        }
    } else {
        throw "Invalid event data: $line"
    }
}

function pglet_send {
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

    #Write-Host "Result: $result"

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

try {
    Write-Host "dddd"
    $res = (pglet page)

    if ($res -match "(?<pipeName>[^\s]+)\s(?<url>[^\s]+)") {
        $pipeName = $Matches["pipeName"]
        $pageUrl = $Matches["url"]
    } else {
        throw "Invalid event data: $res"
    }

    Write-Host "Page URL: $pageUrl"

    $pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
    $pipe.Connect(5000)
    $pipeReader = new-object System.IO.StreamReader($pipe)
    $pipeWriter = new-object System.IO.StreamWriter($pipe)
    $pipeWriter.AutoFlush = $true
    
    $eventPipe = new-object System.IO.Pipes.NamedPipeClientStream("$pipeName.events")
    $eventPipe.Connect(5000)
    $eventPipeReader = new-object System.IO.StreamReader($eventPipe)
    
    pglet_send "clean page"
    #pglet_send "remove body"
    $rowId = pglet_send "add row id=body
      aaa=bbb"
    $colId = pglet_send "add col id=form to=$rowId"
    pglet_send "add text value='Enter your name:' to=$colId"
    pglet_send "add textbox id=fullName value='someone' to=$colId"
    pglet_send "add button id=submit text=Submit event=btn_event to=$colId"
    
    pglet_send "set body:form:fullName value='John Smith'"
    
    while($true) {
        pglet_event
        $fullName = pglet_send "get body:form:fullName value"
        Write-Host "Full name: $fullName"
    }
} catch {
    Write-Host "$_"
} finally {
    $pipe.Close()
    $eventPipe.Close()
}