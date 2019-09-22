package main

import (
    "fmt"
    "encoding/csv"
    "os"
    "io"
    "strconv"
    "math"
    "math/rand"
    "time"
)

type Point struct {
	Id int
    Lat, Lon float64
}

var __DEBUG__ bool = false

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a < b {
        return b
    }
    return a
}

func LoadLocations(path string) []Point {
	// load a csv file and parse locations
	var points []Point
    f, _ := os.Open(path)

    // create a new reader
    r := csv.NewReader(f)
    id := 0
    for {
        record, err := r.Read()
        // stop at EOF
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }

        lat, err := strconv.ParseFloat(record[0], 64)
        lon, err := strconv.ParseFloat(record[1], 64)
        p := Point{id, lat, lon}
        id += 1
        points = append(points, p)
    }
    return points
}

func ComputeCentroids(points []Point, clusters_map map[int]int) []Point {
	// take points and current clusters assignment and calculate new centroids
	var centroids []Point
    clusters_sum_p := make(map[int]Point)
    clusters_count := make(map[int]int)
    for p_id, c_id := range clusters_map { 
        if _, ok := clusters_count[c_id]; ok {
            clusters_count[c_id] += 1
            p := clusters_sum_p[c_id]
            p.Lat += points[p_id].Lat
            p.Lon += points[p_id].Lon
            clusters_sum_p[c_id] = p
        } else {
            clusters_count[c_id] = 1
            clusters_sum_p[c_id] = points[p_id]
        }
    }
    for c_id, _cluster_size := range clusters_count {
        p := clusters_sum_p[c_id]
        cluster_size := float64(_cluster_size)
        centroid := Point{-1, p.Lat/cluster_size, p.Lon/cluster_size}
        centroids = append(centroids, centroid)
    }
    return centroids
}

func UpdateClusters(points []Point, centroids []Point, c chan map[int]int) {
    clusters_map := make(map[int]int)
    for _, p := range points {
        p_id := p.Id
        min_dist := math.Inf(1)
        centroid_inx := 0
        for c_id, c := range centroids {
            dist := math.Sqrt(math.Pow((p.Lat - c.Lat), 2.0) + math.Pow((p.Lon - c.Lon), 2.0))
            // fmt.Println(dist)
            if dist < min_dist {
                min_dist = dist
                centroid_inx = c_id
            }
        }
        clusters_map[p_id] = centroid_inx
    }
    c <- clusters_map //send clusters_map to c
}

func ComputeCentroidsChanges(centroids_l []Point, centroids_r []Point) float64 {
    lat_delta := 0.0
    lon_delta := 0.0
    matched_centroids_r := make(map[int]float64)  // use this map to save matched points from centroids_r
    matched_centroids_l := make(map[int]int)      // use this map to save matched points from centroids_l  
    for inx_l, p_l := range centroids_l {
        // find closest centroid from right to this point from left centroids
        min_dist := 10000.0
        matched_centroid_inx := -1
        for inx_r, p_r := range centroids_r {
            lat_delta = math.Abs(p_l.Lat - p_r.Lat)
            lon_delta = math.Abs(p_l.Lon - p_r.Lon)
            dist := math.Sqrt(math.Pow(lat_delta, 2.0) + math.Pow(lon_delta, 2.0))
            if dist < min_dist {
                if _, ok := matched_centroids_r[inx_r]; ok {
                    //This right centroid has been already matched so skip this point
                    continue
                } else {
                    min_dist = dist
                    matched_centroid_inx = inx_r
                }
            }
        }
        if matched_centroid_inx >= 0 {
            matched_centroids_r[matched_centroid_inx] = min_dist
            matched_centroids_l[inx_l] = matched_centroid_inx
        }
    }
    delta := 0.0
    for inx_r, _ := range centroids_r {
        if dist, ok := matched_centroids_r[inx_r]; ok {
            delta += dist
        } else {
            // this centroid from right does not have a match - panalize
            delta += 1000.0
        }
    }
    for inx_l, _ := range centroids_l {
        if _, ok := matched_centroids_l[inx_l]; ok {
            continue
        } else {
            // this centroid from left does not have a match - panalize
            delta += 1000.0
        }
    }
    return delta
}


func ComputeSSE(points []Point, clusters_map map[int]int, centroids []Point) float64 {
    var sse float64 = 0.0
    for i, p := range points {
        cid := clusters_map[i]
        centroid := centroids[cid]
        dist := math.Pow((p.Lat - centroid.Lat), 2.0) + math.Pow((p.Lon - centroid.Lon), 2.0)
        sse += dist
    }
    return sse
}

func RunKMeans(points []Point, k int) (map[int]int, []Point, float64) {
    var centroids []Point
    rand.Seed(time.Now().UTC().UnixNano())
    rand_indexes := rand.Perm(len(points))[:k]
    for inx := 0; inx < k; inx++ {
        centroids = append(centroids, points[rand_indexes[inx]])
    }
    if __DEBUG__ {
        fmt.Println("Initial centroids:", centroids)
    }
    change := 1000.0
    threshold := 0.01
    clusters_map := make(map[int]int)
    c := make(chan map[int]int)
    parallelization := 4
    for change > threshold {
        // compute clusters mapping
        for i := 0; i < parallelization; i++ {
            go UpdateClusters(points[i*len(points)/parallelization:(i+1)*len(points)/parallelization], centroids, c)    
        }
        for i := 0; i < parallelization; i++ {
            c_map := <- c
            for k, v := range c_map {
                clusters_map[k] = v
            }
        }
        new_centroids := ComputeCentroids(points, clusters_map)
        if __DEBUG__ {
            fmt.Println("Clusters map:", clusters_map)
            fmt.Println("New centroids:", new_centroids)
        }
        change = ComputeCentroidsChanges(centroids, new_centroids)
        if __DEBUG__ {
            fmt.Println("Centroids change:", change)
        }
        centroids = new_centroids
    }
    // since we have left the loop we need to update cluster map for the last time before calculating SSE
    for i := 0; i < parallelization; i++ {
        go UpdateClusters(points[i*len(points)/parallelization:(i+1)*len(points)/parallelization], centroids, c)    
    }
    for i := 0; i < parallelization; i++ {
        c_map := <- c
        for k, v := range c_map {
            clusters_map[k] = v
        }
    }
    if __DEBUG__ {
        fmt.Println("Clusters map:", clusters_map)
    }
    sse := ComputeSSE(points, clusters_map, centroids)
    return clusters_map, centroids, sse
}

func main() {
    file_name := os.Args[1]
    fmt.Println("Loading points ...")
	points := LoadLocations("inputs/" + file_name)
    if len(points) <= 5 {
        fmt.Println("Not enough points to run K-Means!")
        os.Exit(1)
    }
    points_no := len(points)
    fmt.Println("Loaded",  points_no,"points")
    max_k := int(2.0*math.Pow(math.Log2(float64(points_no)),1))
    if __DEBUG__ {
        fmt.Println("max_k:", max_k)
    }
    fmt.Println("Searching for optimal clustering ...")
    var errors []float64
    _, _, cur_sse := RunKMeans(points, 1)
    fmt.Println("k:1 SSE:", cur_sse)
    errors = append(errors, cur_sse)
    optimal_k := 1
    for k := 2; k < max(points_no, max_k); k++ {
        _, _, next_sse := RunKMeans(points, k)
        fmt.Println("k:", k, "SSE:", next_sse)
        errors = append(errors, next_sse)
        slope := math.Abs(next_sse - cur_sse)
        if slope < 1.0 {
            optimal_k = k
            if __DEBUG__{
                fmt.Println("k:", k, "slope:", slope, "breaking!")
            }
            break
        }
        cur_sse = next_sse
    }
    fmt.Println("Optimal no of clusters (k):", optimal_k)
    clusters_map, centroids, _ := RunKMeans(points, optimal_k)
    // write cluster map to a text file
    fmt.Println("Writing clusters/centroids to disk ...")
    f_1, err_1 := os.Create("outputs/clusters.csv")
    if err_1 != nil {
        panic(err_1)
    }
    defer f_1.Close()
    for pid := 0; pid < len(points); pid++ {
        cid := clusters_map[pid]
        out := strconv.Itoa(pid) + "," + strconv.Itoa(cid) + "\n"
        _, err := f_1.WriteString(out)
        if err != nil {
            panic(err)
        }
    }
    f_2, err_2 := os.Create("outputs/centroids.csv")
    if err_2 != nil {
        panic(err_2)
    }
    defer f_2.Close()
    for inx := 0; inx < len(centroids); inx++ {
        centroid := centroids[inx]
        out := fmt.Sprintf("%f", centroid.Lon)  + "," + fmt.Sprintf("%f", centroid.Lat) + "\n"
        _, err := f_2.WriteString(out)
        if err != nil {
            panic(err)
        }
    }
}