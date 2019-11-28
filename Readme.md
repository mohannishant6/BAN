Scalable computing assignment 4

this project contains 3 parts:

1. sink
2. sensors
3. PDA/Edge device



# sink

**sink is a simulator program for sink, it is written in Golang**

* directory: sink

* how to build:

  * install Golang

  * ```shell
    cd sink && go build -o sink main.go
    ```

* how to run

  ```shell
  ./sink -h
  Usage of ./sink:
    -local string
      	local addr (default "0.0.0.0:9090")
    -pda string
      	personal digital assistant address (default "127.0.0.1:9091")
  ```

  



# sensors

**sensor is a simulator program for sensor, it is written in Golang.**

**you have to start sink first, then start all the sensor and give them the right sink address**

sensor receive json files as dataset, you can use following files as dataset for sensors

```shell
sensorBloodAlcohol.json
sensorBloodPressure.json
sensorBodyOxygen.json
sensorBreathingRate.json
sensorInsulin.json
sensorPacemaker.json
sensorTemprature.json
```



* directory: sensors

* how to build: 

  * install Golang

  * ```shell
    cd sensors && go build -o sensor cmd/main.go 
    ```

* how to run:

  ```shell
  ./sensor -h 
  Usage of ./censor:
    -addr string
      	listen ip address (default "127.0.0.1")
    -dataset string
      	dataset file (default "data.txt")
    -duration int
      	working duration(seconds) (default 5)
    -id string
      	sensor id (default "1")
    -interval int
      	interval between working(seconds) (default 10)
    -port string
      	listen port (default "1234")
    -sink string
      	sink address (default "127.0.0.1:9090")
    -x int
      	coordinate x (default 1)
    -y int
      	coordinate y (default 1)
    -z int
      	coordinate z (default 1)
  
  ```
  
 #For Windows Systems
 
 There are two approaches -
 
 #A. Directly on the system using cmd
 1. Intsalling Go in the system.
 2. Starting the sink and builiding it via command prompt using GO libraries.
 3. After that starting sensors so as to connect to the sync.
 
 Note: There might be IP address issues in the above approach so as to avoid that second approach was used.
 
 #B. Using Virtualbox
 
 1. Install virtualbox.
 2. Download a Ubuntu iso image.
 3. Create a VM using iso image.
 4. Check the IP address(WiFi) of the main machine(i.e Windows) and provide it to the Edge node.
 5. Install tmux in the VM using
  * ```shell
    sudo apt install tmux 
    ```
 
 6. Install git in the VM
  * ```shell
    sudo apt install git 
    ```
    
 7. Cloning this repository using git clone (---)
 8. Go to the clone directory using cd
 9. Running the shell script using 
  * ```shell
    sh start.sh
    ``` 
 
 10. This script will autmatically turn on sink node as well as the sensors.
 

  

  
