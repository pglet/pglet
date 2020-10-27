$ErrorActionPreference = "Stop"
$pipeName=$args[0]
$pipe = new-object System.IO.Pipes.NamedPipeClientStream($pipeName)
$pipe.Connect(5000)
$pipeWriter = new-object System.IO.StreamWriter($pipe)
$pipeWriter.AutoFlush = $true
$pipeWriter.WriteLine("add text value='hello, world'")
$pipeWriter.Flush()