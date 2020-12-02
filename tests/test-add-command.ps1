$ErrorActionPreference = "Stop"

$pipeName = $null
$pipe = $null
$pipeReader = $null
$pipeWriter = $null
$eventPipe = $null
$eventPipeReader = $null

function pglet_event {
    if ($IsLinux -or $IsMacOS) {
        $eventPipeReader = New-Object System.IO.StreamReader("$pipeName.events")
    }
    $line = $eventPipeReader.ReadLine()
    if ($IsLinux -or $IsMacOS) {
        $eventPipeReader.Close()
    }
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

    $waitResult = $true
    if ($command -match "(?<commandName>[^\s]+)\s(.*)") {
        $commandName = $Matches["commandName"]
        #Write-Host "COMMAND: $commandName"
        if ($commandName.toLower().endsWith("f")) {
            $waitResult = $false
        }
    }

    if ($IsLinux -or $IsMacOS) {
        $pipeWriter = New-Object System.IO.StreamWriter($pipeName)
    }

    # send command
    $pipeWriter.WriteLine($command)
    #$pipeWriter.Flush()

    #Write-Host "After write"

    if ($IsLinux -or $IsMacOS) {
        $pipeWriter.Close()
    }

    #Write-Host "After writer close"

    if ($waitResult) {
        # parse results
        $ERROR_RESULT = "error"

        try {
            if ($IsLinux -or $IsMacOS) {
                $pipeReader = New-Object System.IO.StreamReader($pipeName)
            }
            
            $result = $pipeReader.ReadLine()

            #Write-Host "After read"
        
            #Write-Host "Result: $result"
        
            if ($result.StartsWith("$ERROR_RESULT ")) {
                throw $result.Substring($ERROR_RESULT.Length + 1)
            } elseif ($result -match "(?<lines_count>[\d]+)\s(?<result>.*)") {
                $lines_count = [int]$Matches["lines_count"]
                $result = $Matches["result"]
        
                # read the rest of multi-line result
                for($i = 0; $i -lt $lines_count; $i++) {
                    $line = $pipeReader.ReadLine()
                    $result = "$result`n$line"
                }
            } else {
                throw "Invalid result: $result"
            }
            return $result
        } finally {
            if ($IsLinux -or $IsMacOS) {
                $pipeReader.Close()
            }
        }
    }
}

try {

    $res = (pglet page page1)

    if ($res -match "(?<pipeName>[^\s]+)\s(?<url>[^\s]+)") {
        $pipeName = $Matches["pipeName"]
        $pageUrl = $Matches["url"]
    } else {
        throw "Invalid event data: $res"
    }

    Write-Host "Page URL: $pageUrl"

    if (-not $IsLinux -and -not $IsMacOS) {
        # use Windows named pipes
        $pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
        $pipe.Connect(5000)
        $pipeReader = new-object System.IO.StreamReader($pipe)
        $pipeWriter = new-object System.IO.StreamWriter($pipe)
        $pipeWriter.AutoFlush = $true
        
        $eventPipe = new-object System.IO.Pipes.NamedPipeClientStream("$pipeName.events")
        $eventPipe.Connect(5000)
        $eventPipeReader = new-object System.IO.StreamReader($eventPipe)
    }
    
    pglet_send "clean page"
    Start-Sleep -s 2

    pglet_send "set page title='Hello, world!' gap=10 horizontalAlign=start"
    pglet_send "add text value='Your name' size=large"
    #pglet_send "add button id=submit text=Submit primary=yes event=btn_event"
    #$rowId = pglet_send "add row"
    # pglet_send "add
    #     button id=a1
    #     row
    #         button id=b1"

    pglet_send "add to=page at=0
        stack width=600px horizontalAlign=stretch
          textbox id=fullName value='someone' label=Name placeholder='Your name, please' description='That\'s your name'
          textbox id=bio label='Bio' description='A few words about yourself' value='Line1\nLine2' multiline=true
          dropdown id=color label='Your favorite color' placeholder='Select color'
            option key=red text=Red
            option key=green text=Green
            option key=blue text=Blue
          checkbox id=agree label='I agree to the terms of services'"

    pglet_send "add stack at=0 id=buttons horizontal=true
            button id=submit text=Submit primary=yes event=btn_event
            button id=cancel event=btn_event2"
    
    Start-Sleep -s 2

    pglet_send "add button id=b1 to=buttons"

    pglet_send "set fullName value='John Smith'"
    Start-Sleep -s 2
    pglet_send "set
        bio value='Some bio...'
        color value='green'"

    pglet_send "add progress id=prog label='Doing something...' width=400px"
    
    while($true) {
        pglet_event

        $fullName = pglet_send "get fullName value"
        Write-Host "Full name: $fullName"

        $bio = pglet_send "get bio value"
        Write-Host "Bio: $bio"

        for ($i = 0; $i -lt 101; $i++) {
            pglet_send "setf prog value=$($i) label='Step $i...'" | out-null
            Start-Sleep -ms 10
        }
    }
} catch {
    Write-Host "$_"
} finally {
    $pipe.Close()
    $eventPipe.Close()
}