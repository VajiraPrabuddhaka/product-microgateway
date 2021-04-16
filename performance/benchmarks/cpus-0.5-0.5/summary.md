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
      cpu: "320m"
    limits:
      cpu: "500m"
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

|Scenario Name        |Heap Size|Concurrent Users|Message Size (Bytes)|Back-end Service Delay (ms)|Error %|Throughput (Requests/sec)|Average Response Time (ms)|Standard Deviation of Response Time (ms)|90th Percentile of Response Time (ms)|95th Percentile of Response Time (ms)|99th Percentile of Response Time (ms)|enforcer GC Throughput (%)|Average enforcer Memory Footprint After Full GC (M)|
|---------------------|---------|----------------|--------------------|---------------------------|-------|-------------------------|--------------------------|----------------------------------------|-------------------------------------|-------------------------------------|-------------------------------------|---------------------------------------|----------------------------------------------------------------|
|Choreo-connect on k8s|512      |10              |50                  |0                          |0      |1518.12                  |6.56                      |14.74                                   |3.0                                  |59.0                                 |65.0                                 |98.83                                  |117.333                                                         |
|Choreo-connect on k8s|512      |10              |1024                |0                          |0      |1456.82                  |6.83                      |15.24                                   |4.0                                  |61.0                                 |66.0                                 |98.67                                  |122.0                                                           |
|Choreo-connect on k8s|512      |200             |50                  |0                          |0      |1565.2                   |127.69                    |45.96                                   |195.0                                |198.0                                |204.0                                |99.09                                  |126.667                                                         |
|Choreo-connect on k8s|512      |100             |50                  |0                          |0      |1631.93                  |61.22                     |40.56                                   |97.0                                 |98.0                                 |101.0                                |99.01                                  |120.75                                                          |
|Choreo-connect on k8s|512      |200             |10240               |0                          |0      |1249.42                  |159.96                    |54.43                                   |204.0                                |213.0                                |294.0                                |99.23                                  |121.0                                                           |
|Choreo-connect on k8s|512      |50              |1024                |0                          |0      |1610.22                  |31.0                      |36.39                                   |86.0                                 |88.0                                 |90.0                                 |99.0                                   |120.5                                                           |
|Choreo-connect on k8s|512      |50              |10240               |0                          |0      |1258.77                  |39.66                     |39.25                                   |90.0                                 |91.0                                 |93.0                                 |99.3                                   |117.0                                                           |
|Choreo-connect on k8s|512      |100             |10240               |0                          |0      |1272.74                  |78.49                     |39.28                                   |101.0                                |104.0                                |185.0                                |99.23                                  |118.0                                                           |
|Choreo-connect on k8s|512      |50              |50                  |0                          |0      |1638.7                   |30.46                     |36.17                                   |86.0                                 |87.0                                 |90.0                                 |99.1                                   |118.0                                                           |
|Choreo-connect on k8s|512      |200             |1024                |0                          |0      |1567.41                  |127.51                    |45.98                                   |195.0                                |198.0                                |204.0                                |99.04                                  |123.0                                                           |
|Choreo-connect on k8s|512      |10              |10240               |0                          |0      |1209.22                  |8.24                      |17.31                                   |4.0                                  |65.0                                 |69.0                                 |99.12                                  |120.333                                                         |
|Choreo-connect on k8s|512      |100             |1024                |0                          |0      |1617.76                  |61.76                     |40.4                                    |97.0                                 |98.0                                 |102.0                                |99.07                                  |119.333                                                         |

