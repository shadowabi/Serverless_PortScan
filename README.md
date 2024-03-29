# Serverless_PortScan
利用云函数实现端口扫描

## 配置
当首次运行Serverless_PortScan时，会检测config.json文件是否存在，不存在则会自动创建

config.json的填写内容应该如下：  
```
{
    "ServerUrl":"http://",
    "PortList":"21,22,23,25,80,135,139,389,443,445,873,888,1433,1521,2049,2181,2375,2379,3306,3389,3690,5432,5900,5985,5986,6379,6443,7001,8000,8061,8080,8089,8161,8500,8443,8649,8888,9080,10250,10255,11211,13389,16379,27017,27019,23791,30000,50070,63791"
}
```
ServerUrl为你的云函数地址  


## 云函数配置
将以下内容配置到云函数,并建议将云函数执行时间和api网关后端超时设置为15秒或以上：

```python
import socket
from concurrent.futures import ThreadPoolExecutor
from collections import defaultdict
import threading

# 使用defaultdict存储每个线程找到的开放端口
results_per_thread = defaultdict(list)
ip = ""
ports = ""

def Scan(port):
    try:
        conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        conn.settimeout(1)
        res = conn.connect_ex((str(ip), int(port)))
        if res == 0:
            # 将结果添加到当前线程局部变量中
            results_per_thread[threading.current_thread()].append(port)
    except Exception as err:
        print(err)
    finally:
        conn.close()

def main_handler(event, context):
    global ip, ports
    results_per_thread.clear()
    ip = event["queryString"]["ip"]
    ports = event["queryString"]["port"].split(",")

    with ThreadPoolExecutor(max_workers=20) as executor:
        executor.map(Scan, ports)

    # 合并所有线程的结果
    open_ports = []
    for thread_results in results_per_thread.values():
        open_ports.extend(thread_results)

    # 将端口列表转换为逗号分隔的字符串
    a = ",".join(str(port) for port in open_ports)
    return a
```

触发管理配置如下：  
![image](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/899e0445-dd7c-4c2b-9bdd-26c248fa0eb6)


## 安装
下载release中的对应版本

使用Makefile进行编译


## 用法
```
Usage:  

  Serverless_PortScan [flags]  


Flags:  

  -f, --file string       从文件中读取目标地址 (Input filename)  
  -h, --help              help for Serverless_PortScan  
      --logLevel string   设置日志等级 (Set log level) [trace|debug|info|warn|error|fatal|panic] (default "info")  
  -o, --output string     输入结果文件输出的位置 (Enter the location of the scan result output) (default "./result.txt")  
  -p, --port string       输入需要被扫描的端口，逗号分割 (Enter the port to be scanned, separated by commas (,))  
  -t, --timeout int       输入每个 http 请求的超时时间 (Enter the timeout period for every http request) (default 5)  
  -u, --url string        输入目标地址 (Input [ip|domain|url]) 
```


## 功能列表

1. 利用云函数特性扫描端口，防止封ip
2. 本地多线程+云函数多线程发包，提高扫描速度
3. 自动去重
4. 文件输出时为：ip+端口号形式，方便利用其他工具如指纹识别工具进行扫描
5. 采用go语言编写，提升性能


## 旧版本

python版本在python分支  

旧的go版本在go-old
