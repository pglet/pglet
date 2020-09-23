$pipeName=$args[0]
$pipe = $null
$pipeReader = $null

try {
    $pipe = new-object System.IO.Pipes.NamedPipeClientStream("pglet_pipe_$pipeName.events")
    $pipe.Connect()

    $pipeReader = new-object System.IO.StreamReader($pipe)
    #$pipeWriter = new-object System.IO.StreamWriter($pipe)
    #$pipeWriter.AutoFlush = $true

    #$pipeWriter.WriteLine("Boom!")

    for ($i = 0; $i -lt 5; $i++) {
        $line = $pipeReader.ReadLine()
        if ($line -eq $null) {
            break
        }
        Write-Host $line

        Start-Sleep -s 1
    }

} finally {
    $pipeReader.Dispose()
    $pipe.Dispose()
}


