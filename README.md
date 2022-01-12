# System monitor
Golang API for accessing system stats, linux-only

```shell
$ curl -s http://127.0.0.1:8080/cpu/info | jq
[
  {
    "Model": "Intel(R) Core(TM) i7-10875H CPU @ 2.30GHz",
    "MHz": 2303.759,
    "Cores": 8,
    "CoreID": 0,
    "PhysicalID": 0
  },
  {
    "Model": "Intel(R) Core(TM) i7-10875H CPU @ 2.30GHz",
    "MHz": 2303.759,
    "Cores": 8,
    "CoreID": 0,
    "PhysicalID": 0
  },
  {
    "Model": "Intel(R) Core(TM) i7-10875H CPU @ 2.30GHz",
    "MHz": 2303.759,
    "Cores": 8,
    "CoreID": 1,
    "PhysicalID": 0
  },
...
```

```shell
$ curl -s http://127.0.0.1:8080/cpu/utilization | jq
{
  "cpuUtilization": 0.001875
}
```

```shell
$ curl -s http://127.0.0.1:8080/memory/info | jq
{
  "TotalKB": 26065704,
  "FreeKB": 23878080
}
```

```shell
$ curl -s http://127.0.0.1:8080/memory/utilization | jq
{
  "memoryUtilization": 3.2072795731893525e-05
}
```