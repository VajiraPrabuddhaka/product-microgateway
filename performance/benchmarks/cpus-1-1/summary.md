# WSO2 Choreo Connect Performance Test Results

During each release, we execute various automated performance test scenarios and publish the results.

| Test Scenarios | Description |
| --- | --- |
| Choreo-Connect with throtteling enabled |  |

Our test client is [Apache JMeter](https://jmeter.apache.org/index.html). We test each scenario for a fixed duration of
time. We split the test results into warmup and measurement parts and use the measurement part to compute the
performance metrics.

Test scenarios use a [Netty](https://netty.io/) based back-end service which echoes back any request
posted to it after a specified period of time.

We run the performance tests under different numbers of concurrent users, message sizes (payloads) and back-end service
delays.

The main performance metrics:

1. **Throughput**: The number of requests that the WSO2 API Microgateway processes during a specific time interval (e.g. per second).
2. **Response Time**: The end-to-end latency for an operation of invoking an API. The complete distribution of response times was recorded.

In addition to the above metrics, we measure the load average and several memory-related metrics.

The following are the test parameters.

| Test Parameter | Description | Values |
| --- | --- | --- |
| Scenario Name | The name of the test scenario. | Refer to the above table. |
| Heap Size | The amount of memory allocated to the `enforcer` | 512M |
| Concurrent Users | The number of users accessing the application at the same time. | 10, 50, 100, 200 |
| Message Size (Bytes) | The request payload size in Bytes. | 50, 1024, 10240 |
| Back-end Delay (ms) | The delay added by the back-end service. | 0 |

The duration of each test is **900 seconds**. The warm-up period is **300 seconds**.
The measurement results are collected after the warm-up period.

A [EKS](https://aws.amazon.com/eks) cluster with 3 nodes, each having [**c5.xlarge** Amazon EC2 instance](https://aws.amazon.com/ec2/instance-types/) was used to setup Choreo-Connect on k8s.
To limit the cpu utilization, the container configs' resource utilization set as follows.
```aidl
resources:
    requests:
      cpu: "756m"
    limits:
      cpu: "1000m"
```

The jmeter is configured such that the maximum waiting time for receiving a response to be 20 seconds.

The following figures shows how the Throughput changes for different number of concurrent users with different backend delays
![picture](plots/tps_0ms.png)

The following figures shows how the Average Response Time changes for different number of concurrent users with different backend delays.
![picture](plots/response_time_0ms.png)

Let’s look at the 90th, 95th, and 99th Response Time percentiles.
This is useful to measure the percentage of requests that exceeded the response time value for a given percentile.
A percentile can also tell the percentage of requests completed below the particular response time value.
![picture](plots/p99_0ms.png)
![picture](plots/p95_0ms.png)
![picture](plots/p90_0ms.png)

The GC Throughput was calculated for each test to check whether GC operations are not impacting the performance of the enforcer.
The GC Throughput is the time percentage of the application, which was not busy with GC operations.
![picture](plots/gc_tps_0ms.png)

The following are the measurements collected from each performance test conducted for a given combination of
test parameters.

| Measurement | Description |
| --- | --- |
| Error % | Percentage of requests with errors |
| Average Response Time (ms) | The average response time of a set of results |
| Standard Deviation of Response Time (ms) | The “Standard Deviation” of the response time. |
| 99th Percentile of Response Time (ms) | 99% of the requests took no more than this time. The remaining samples took at least as long as this |
| Throughput (Requests/sec) | The throughput measured in requests per second. |
| Average Memory Footprint After Full GC (M) | The average memory consumed by the application after a full garbage collection event. |

The following is the summary of performance test results collected for the measurement period.

|Scenario Name        |Heap Size|Concurrent Users|Message Size (Bytes)|Back-end Service Delay (ms)|Error %|Throughput (Requests/sec)|Average Response Time (ms)|Standard Deviation of Response Time (ms)|90th Percentile of Response Time (ms)|95th Percentile of Response Time (ms)|99th Percentile of Response Time (ms)|WSO2 API Microgateway GC Throughput (%)|Average WSO2 API Microgateway Memory Footprint After Full GC (M)|
|---------------------|---------|----------------|--------------------|---------------------------|-------|-------------------------|--------------------------|----------------------------------------|-------------------------------------|-------------------------------------|-------------------------------------|---------------------------------------|----------------------------------------------------------------|
|Choreo-connect on k8s|512      |10              |50                  |0                          |0      |2908.15                  |3.41                      |3.94                                    |3.0                                  |4.0                                  |26.0                                 |98.51                                  |116.0                                                           |
|Choreo-connect on k8s|512      |10              |1024                |0                          |0      |2865.83                  |3.46                      |3.95                                    |3.0                                  |4.0                                  |26.0                                 |98.4                                   |123.857                                                         |
|Choreo-connect on k8s|512      |200             |50                  |0                          |0      |3261.22                  |61.26                     |33.15                                   |93.0                                 |96.0                                 |103.0                                |98.29                                  |122.429                                                         |
|Choreo-connect on k8s|512      |100             |50                  |0                          |0      |3393.49                  |29.42                     |27.88                                   |74.0                                 |76.0                                 |80.0                                 |98.33                                  |121.0                                                           |
|Choreo-connect on k8s|512      |200             |10240               |0                          |0      |2542.22                  |78.57                     |32.11                                   |102.0                                |106.0                                |168.0                                |98.58                                  |123.333                                                         |
|Choreo-connect on k8s|512      |50              |1024                |0                          |0      |3288.9                   |15.16                     |19.79                                   |61.0                                 |63.0                                 |67.0                                 |98.31                                  |122.375                                                         |
|Choreo-connect on k8s|512      |50              |10240               |0                          |0      |2573.33                  |19.37                     |23.09                                   |66.0                                 |69.0                                 |72.0                                 |98.66                                  |118.167                                                         |
|Choreo-connect on k8s|512      |100             |10240               |0                          |0      |2588.02                  |38.57                     |31.08                                   |80.0                                 |83.0                                 |87.0                                 |98.63                                  |121.833                                                         |
|Choreo-connect on k8s|512      |50              |50                  |0                          |0      |3325.05                  |15.0                      |19.5                                    |60.0                                 |63.0                                 |66.0                                 |98.38                                  |120.286                                                         |
|Choreo-connect on k8s|512      |200             |1024                |0                          |0      |3311.87                  |60.32                     |33.06                                   |93.0                                 |96.0                                 |102.0                                |98.23                                  |123.125                                                         |
|Choreo-connect on k8s|512      |10              |10240               |0                          |0      |2361.36                  |4.2                       |5.38                                    |4.0                                  |5.0                                  |32.0                                 |98.69                                  |120.0                                                           |
|Choreo-connect on k8s|512      |100             |1024                |0                          |0      |3365.21                  |29.66                     |28.0                                    |74.0                                 |77.0                                 |80.0                                 |98.27                                  |119.875                                                         |
