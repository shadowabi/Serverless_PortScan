# Serverless_PortScan  
利用云函数实现端口扫描

## 配置
打开config.py，配置云函数地址

## 云函数配置  
将以下内容配置到云函数,并建议将云函数执行时间和api网关后端超时设置为15秒或以上：

```python
from socket import *
from concurrent.futures import ThreadPoolExecutor

Tport = ""
ip = ""
ports = ""

def Scan(port):
    global Tport
    try:
        conn = socket(AF_INET,SOCK_STREAM)
        conn.settimeout(1)
        res = conn.connect_ex((str(ip),int(port)))
        if res == 0:
            Tport =  Tport + "," + port
    except Exception as err:
        print(err)
    finally:
        conn.close()

def main_handler(event, context):
    global Tport,ip,ports

    ip = event["queryString"]["ip"]
    ports = event["queryString"]["port"].split(",")

    with ThreadPoolExecutor(max_workers = 20) as executor:
        executor.map(Scan, ports)

    a = Tport[1:]
    Tport = ""
    return a
```

触发管理配置如下：  
![image](https://github.com/shadowabi/Serverless_PortScan/assets/50265741/899e0445-dd7c-4c2b-9bdd-26c248fa0eb6)


## 安装
pip install -r requirements.txt



## 用法  
python Serverless_PortScan.py [-h] [-u 127.0.0.1 | -f 1.txt] [-p 80,443]



## 功能列表

1. 利用云函数特性扫描端口，防止封ip
2. 本地多线程、协程+云函数多线程发包，提高扫描速度
3. 自动去重
4. 文件输出时为：ip+端口号形式，方便利用其他工具如指纹识别工具进行扫描



## 将要更新

1. 采用go语言重构本程序
