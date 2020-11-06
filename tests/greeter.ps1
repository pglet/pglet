$ErrorActionPreference = "Stop"

$pipeName=$args[0]
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
    $pipeName=$args[0]
    $pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
    $pipe.Connect(5000)
    $pipeReader = new-object System.IO.StreamReader($pipe)
    $pipeWriter = new-object System.IO.StreamWriter($pipe)
    $pipeWriter.AutoFlush = $true
    
    $eventPipe = new-object System.IO.Pipes.NamedPipeClientStream("$pipeName.events")
    $eventPipe.Connect(5000)
    $eventPipeReader = new-object System.IO.StreamReader($eventPipe)

    function greet() {
        $fullName = pglet "get body:form:fullName value"
        Write-Host "Full name: $fullName"

        # output welcome message
        pglet "clean page"
        pglet "add text value='Hello, $fullName'"
        pglet "add button id=again text=Again"

        while($true) {
            $e = pglet_event
            if ($e.Target -eq 'again' -and $e.Name -eq 'click') {
                main
                return
            }
        }
    }    
    
    function main() {
        pglet "clean page"
        #pglet "remove body"
        $rowId = pglet "add row id=body
          aaa=bbb"
        $colId = pglet "add col id=form to=$rowId"
        pglet "add text value='Enter your name:' to=$colId"
        pglet "add textbox id=fullName value='someone' to=$colId"
        pglet "add button id=submit text=Submit event=btn_event to=$colId"
        
        pglet "set body:form:fullName value='John Smith'"
        
        while($true) {
            $e = pglet_event
            if ($e.Target -eq 'body:form:submit' -and $e.Name -eq 'click') {
                greet
                return
            }
        }
    }

    # entry point
    main

} finally {
    $pipe.Close()
    $eventPipe.Close()
}