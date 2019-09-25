# Clustering
This project implements [k-means](https://en.wikipedia.org/wiki/K-means_clustering) clustering algorithm.

# Description
The code programmatically optimizes the number of clusters (i.e. k) and at the end of the process, it stores the detected clusters to the disk.

# Build k-means
`go build`

`mkdir outputs`

# Input data
You can test the code with [San Francisco Crimes Data](https://data.sfgov.org/Public-Safety/Police-Department-Incident-Reports-Historical-2003/tmnf-yvry) file loctaed in "inputs" folder (crimes.csv.gz). 

If you want to test the code with San Francisco crimes data, you should unzip the file first:

`gunzip inputs/crimes.csv.gz`

When you want to run clustering with your locations data, you should store your locations file in "inputs" folder.  Also, your input should have a CSV format with two columns: "Lat,Lon".

Next, just pass your filename as the only parameter to the clustering program in the command line as shown below:

# Run k-means
`./clustering crimes.csv`
