from sklearn.datasets import make_blobs
import csv

N = 5000
D = 160

def main():
    data,_= make_blobs(n_samples=N, n_features=D, centers=10, random_state=42)


    filename = f"data-N:{N}-D:{D}.csv"
    with open (filename, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerows(data)

main()