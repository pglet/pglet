$pipeName=$args[0]
$pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
$pipe.Connect()
$pipeWriter = new-object System.IO.StreamWriter($pipe)
$pipeWriter.AutoFlush = $true
$pipeWriter.WriteLine("quit")
$pipeWriter.Flush()