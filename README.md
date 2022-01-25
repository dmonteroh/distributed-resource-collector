# Distributed Resource Collector & Heartbeat
Simple webserver that collects relevant system data

# v0.1
 - Performs required (known) data collection tasks
 - Outputs data as JSON via a GET handler

Data collection takes a fraction of a second to complete (Linux environment)

2022-01-DD HH:26:17.305794417 +0200 EET

2022-01-DD HH:26:17.547917812 +0200 EET

Web server and handler add very little overhead to the operation
![image](https://user-images.githubusercontent.com/16642619/151067744-43a6913b-775a-4c7e-91fc-db7e30474bda.png)

# Planned
 - Post data automatically to Central Service (Heartbeat)
 - Perform latency checking tasks
    - Receive (and maybe Store) parameters from Central Service

# Possible expansions
 - Add authentication mechanism (CA & Central Service) to securize endpoints
