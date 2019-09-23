# Clustering
This project implements [k-means](https://en.wikipedia.org/wiki/K-means_clustering) clustering algorithm 

# Description
This implementation programmatically optimizes for the number of clusters (k) and at the end of clustering process stores the clusters to disk.

# Build
`go build`

`mkdir outputs`

# Data
You can test the code with [San Francisco crimes data](https://data.sfgov.org/Public-Safety/Police-Department-Incident-Reports-Historical-2003/tmnf-yvry) in "inputs" folder (i.e. crimes.csv.gz). 

Note that if you want to test with your own location data, you need to copy your location CSV format file into "inputs" folder first.  Next, pass your filename as a parameter to the clustering program as shown below.  Your CSV file should have "Lat,Lon" format.  

If you want to test with San Francisco crimes data, you should unzip the crimes file first:

`gunzip inputs/crimes.csv.gz`

# Run k-means
`./clustering crimes.csv`