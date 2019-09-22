# Clustering
This project implements [k-means](https://en.wikipedia.org/wiki/K-means_clustering) clustering algorithm 

# Description
This implementation programmatically optimizes for the number of clusters (k) and at the end of clustering process stores the clusters to disk.

# Build
go build

# Data
You can test the code with San Francisco crimes locations data located in "inputs" folder (i.e. crimes.csv). Note that you need to copy your location file in CSV format into "inputs" folder and pass your filename as a parameter to the clustering program as shown below.  The CSV file has "Lat,Lon" coding.

# Run k-means
./clustering crimes.csv