## Benchmarking Results

1. With
`
5ms diff ; 0 retries after 0 seconds 
`
<br> 
```
Benchmark finished: 4m40.185756708s
Total Connections: 50000
Successful Connections: 48043
Failed Connections: 1957
Total Messages Sent: 50000
Successful Messages Sent: 48043
Failed Messages Sent: 0
Message Throughput: 171.47 messages/second
Data Throughput: 857.27 bytes/second
```
2. With
`
5ms diff ; 3 retries after 3 seconds 
`
<br> 
```
Benchmark finished: 4m43.286261959s
Total Connections: 50000
Successful Connections: 47944
Failed Connections: 2056
Total Messages Sent: 50000
Successful Messages Sent: 47944
Failed Messages Sent: 0
Message Throughput: 169.24 messages/second
Data Throughput: 845.51 bytes/second
```
3. With
`
10ms diff ; 3 retries after 3 seconds 
`
<br> 
```
Benchmark finished: 9m12.487814042s
Total Connections: 50000
Successful Connections: 50000
Failed Connections: 0
Total Messages Sent: 50000
Successful Messages Sent: 50000
Failed Messages Sent: 0
Message Throughput: 90.50 messages/second
Data Throughput: 452.50 bytes/second
```
3. With
`
8ms diff ; 3 retries after 3 seconds 
`
<br> 
```
Benchmark finished: 7m29.506731625s
Total Connections: 50000
Successful Connections: 50000
Failed Connections: 0
Total Messages Sent: 50000
Successful Messages Sent: 50000
Failed Messages Sent: 0
Message Throughput: 111.23 messages/second
Data Throughput: 556.11 bytes/second
```

