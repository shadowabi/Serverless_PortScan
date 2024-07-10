# Serverless_PortScan
利用云函数实现端口扫描

## 安装

下载 release 中的二进制文件使用

或使用 Makefile 进行编译二进制文件后使用

## 配置
当首次运行 Serverless_PortScan 时，会检测 config.yaml 文件是否存在，不存在则会自动创建

config.yaml 中的 ServerUrl 需要在云函数配置好后才能配置，请看下面：

## 云函数配置
控制台搜索云函数，进入云函数控制台后，新建函数
![img](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/a7eed0d2-7ccc-4eb3-9acf-a4166bf57ab7)
点击从头开始，选择事件函数，运行环境选择python3.x（都可以）
![img_1](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/555cfd63-5a5a-40a1-9d15-4ee5d791c6ff)
在函数代码中中复制粘贴下面的代码：

```python
import socket
from concurrent.futures import ThreadPoolExecutor
from collections import defaultdict
import threading

# 使用defaultdict存储每个线程找到的开放端口
results_per_thread = defaultdict(list)
ip = ""
lock = threading.Lock()

def Scan(port):
    try:
        conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        conn.settimeout(3) # 3秒的话，大概1k端口47秒，2秒则是1k端口33秒，延时越大越不容易漏但时间越长
        res = conn.connect_ex((str(ip), int(port)))
        if res == 0:
            # 将结果添加到当前线程局部变量中
            lock.acquire()
            results_per_thread[threading.current_thread()].append(port)
            lock.release()
    except Exception as err:
        print(err)
    finally:
        conn.close()

def main_handler(event, context):
    global ip
    
    try:
        ip = event["queryString"]["ip"]
        port_str = event["queryString"]["port"]
        path = event["path"]
        if path != "/portscan" :
            return
    except:
        return event

    ports = []
    for port_range in port_str.split(","):
        if "-" in port_range:
            start, end = map(int, port_range.split("-"))
            ports.extend(range(start, end+1))
        else:
            ports.append(int(port_range))

    with ThreadPoolExecutor(max_workers=30) as executor:
        executor.map(Scan, ports)

    # 合并所有线程的结果
    open_ports = []
    for thread_results in results_per_thread.values():
        open_ports.extend(thread_results)

    # 将端口列表转换为逗号分隔的字符串
    a = ",".join(str(port) for port in open_ports)
    
    results_per_thread.clear()
    return a
```
![img_2](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/2b17267a-484b-4642-a58a-fa87eabe3a9e)
往下拉，点击高级配置，设置执行超时时间为 900 秒，也可以自定义一个时间。内存选择 384MB 或以下即可(如果扫描端口过多会爆内存)
![img_3](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/862006a0-cbe5-4d8b-a1e3-200bdf164160)
基础配置如上，最后勾选“我已阅读并同意 《腾讯云云函数网络服务协议》”，点击完成。

创建完毕后，会弹出下面界面：
![img_5](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/f1774aa3-586d-43f4-862e-28d3aa6203ee)
点留在本页，然后点函数URL，创建URL，启用公网访问
![img_6](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/6d084286-226f-45af-9156-faf2dee413f7)
![img_7](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/e64147a9-23cc-40c9-bc81-8694768b086e)
点击提交后，会生成一个地址
![img_8](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/26f1917f-066f-46e3-8d90-219af185a127)
例如这里：
```
https://1304875526-enqewmqsrf-hk.scf.tencentcs.com
```
在 URL 地址最后拼接上 /portscan，配置到 config.yaml 中的 ServerUrl，即
```
ServerUrl: https://1304875526-enqewmqsrf-hk.scf.tencentcs.com/portscan
```


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
  -t, --timeout int       输入每个 http 请求的超时时间 (Enter the timeout period for every http request) (default 900)  
  -u, --url string        输入目标地址 (Input [ip|domain|url]) 
```


## 功能列表

1. 利用云函数特性扫描端口，防止封 ip
2. 本地多线程 + 云函数多线程发包，提高扫描速度
3. 自动去重
4. 文件输出时为：ip + 端口号形式，方便利用其他工具如指纹识别工具进行扫描
5. 采用 go 语言编写，提升性能


## 旧版本

python 版本在 python 分支  

旧的 go 版本在 go-old
