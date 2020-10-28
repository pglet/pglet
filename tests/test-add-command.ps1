$ErrorActionPreference = "Stop"
$pipeName=$args[0]
$pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
$pipe.Connect(5000)
$pipeReader = new-object System.IO.StreamReader($pipe)
$pipeWriter = new-object System.IO.StreamWriter($pipe)
$pipeWriter.AutoFlush = $true

function PgletCommand {
    param (
        $line
    )
    $pipeWriter.WriteLine($line)
    $pipeWriter.Flush()
    
    $result = $pipeReader.ReadLine()
    Write-Host $result
}

PgletCommand "add row id=body"
PgletCommand "add col id=form to=body"
PgletCommand "add text value='Enter your name:' to=body:form"
PgletCommand "add textbox id=fullName to=body:form"
PgletCommand "add button id=submit text=Submit to=body:form"

$pipe.Close()
