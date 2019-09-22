# clustering
Clustering algorithm using [k-means](https://en.wikipedia.org/wiki/K-means_clustering)

# description
This implementation optimize the number of clusters and write the cluster assignments and clusters centroids to disk at the end of process.

# build
go build

# data
You can test the code with San Francisco crimes locations data located in "inputs" folder (i.e. crimes.csv). Note that you need to copy your location file in CSV format into "inputs" folder and pass your filename as a parameter to the clustering program as shown below.  The CSV file has "Lat,Lon" coding.

# run k-means
./clustering crimes.csv