$pipeName=$args[0]
$pipe = $null
$pipeReader = $null
$pipeWriter = $null

try {
    $pipe = new-object System.IO.Pipes.NamedPipeClientStream("pglet_pipe_$pipeName")
    $pipe.Connect()

    $pipeReader = new-object System.IO.StreamReader($pipe)
    $pipeWriter = new-object System.IO.StreamWriter($pipe)
    $pipeWriter.AutoFlush = $true

    for ($i = 0; $i -le 10; $i++) {
        $pipeWriter.WriteLine("cmd $i`nline 2")
        $pipeWriter.Flush()

        $result = $pipeReader.ReadLine()
        Write-Host "received: $result"

        #Start-Sleep -s 1
    }

    $pipeWriter.WriteLine("Boom!`nZoom!")
    $pipeWriter.Flush()

    Start-Sleep -s 2

    $result = $pipeReader.ReadLine()
    Write-Host $result

    $pipeWriter.WriteLine("Hello!`nWorld!")
    $pipeWriter.Flush()

    Write-Host "After sending command"

    $result = $pipeReader.ReadLine()
    Write-Host $result

    Start-Sleep -s 2

    Write-Host "End!"

} finally {
    $pipeReader.Dispose()
    $pipe.Dispose()
}


