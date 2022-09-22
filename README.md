# Serverless_PortScan  
利用云函数实现端口扫描

## 配置
打开config.py，配置云函数地址

## 云函数配置  
将以下内容配置到云函数,并建议将云函数执行时间和api网关后端超时设置为15秒或以上：

```python
from socket import *
from threading import Thread
from queue import Queue
from time import sleep

q = Queue()
Tport = ""

def Scan(ip):
    global Tport
    if(q.empty() == False):
        port = q.get()
        try:
            conn = socket(AF_INET,SOCK_STREAM)
            conn.settimeout(0.2)
            # conn.setsockopt(SOL_SOCKET, SO_REUSEPORT, 1)
            res = conn.connect_ex((str(ip),int(port)))
            if res == 0:
                Tport =  Tport + "," + port
        except Exception as err:
            print(err)
        finally:
            conn.close()
            q.task_done()

def main_handler(event, context):
    global Tport
    ip = event["queryString"]["ip"]
    ports = event["queryString"]["port"].split(",")

    for i in ports:
        q.put(i)
    for i in range(len(ports)):
        t = Thread(target = Scan, args = [ip])
        # sleep(0.1)
        t.start()
        t.join()
    a = Tport[1:]
    Tport = ""
    return a
```




## 安装
pip install -r requirements.txt



## 用法  
python Serverless_PortScan.py [-h] [-u 127.0.0.1 | -f 1.txt] [-p 80,443]
