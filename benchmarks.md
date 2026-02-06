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
