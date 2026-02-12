# Version 2

We benchmarked Version 2 using the following command:

```powershell
Measure-Command { python .\GenerateData.py }
```

The tests were run on Bastian's PC with early termination.
The results are shown below:

| k   | n    | Total Time (seconds) |
| --- | ---- | -------------------- |
| 10  | 200  | 3.12                 |
| 10  | 400  | 3.64                 |
| 20  | 200  | 3.53                 |
| 20  | 2000 | 25.29                |
| 20  | 1000 | 13.30                |
| 40  | 1000 | 23.36                |

We see that with the old parameters the execution time was mostly overhead, so we decided to do more test with bigger parameters

# Version 1

We benchmarked Version 1 using the following command:

```powershell
Measure-Command { python .\GenerateData.py }
```

The tests were run on Bastian's PC with 10 iterations.
The results are shown below:

| k   | n   | Total Time (seconds) |
| --- | --- | -------------------- |
| 10  | 200 | 22.12                |
| 10  | 400 | 40.29                |
| 20  | 200 | 68.97                |
