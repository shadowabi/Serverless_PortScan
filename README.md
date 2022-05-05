# Serverless_PortScan  
利用云函数实现端口扫描



# 云函数配置  
参考这篇文章：https://www.anquanke.com/post/id/261551

或者将以下内容配置到云函数：

```python
from socket import *

def main_handler(event, context):
    IP = event["queryString"]["ip"]
    port = event["queryString"]["port"]
    try:
        conn = socket(AF_INET,SOCK_STREAM)
        conn.setsockopt(SOL_SOCKET, SO_REUSEPORT, 1)
        res = conn.connect_ex((str(IP),int(port)))
        if res == 0:
            return port
    except Exception as err:
        print(err)
    finally:
        conn.close()
    return None
```




# 安装
pip install -r requirements.txt



# 用法  
python Serverless_PortScan.py -u 127.0.0.1 -p 80,443
