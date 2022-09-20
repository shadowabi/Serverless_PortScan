#!/usr/bin/env python
# -*- coding:utf-8 -*-

# import grequests
import requests
from requests.adapters import HTTPAdapter
import socket
from re import *
from sys import argv
import os
import argparse
from time import sleep
from config import *
import readline
from time import sleep
from queue import Queue
from threading import Thread
import traceback

rs = [] #存放结果
ap = argparse.ArgumentParser()
group = ap.add_mutually_exclusive_group()
group.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "127.0.0.1")
group.add_argument("-f", "--file", help = "Input FILENAME", metavar = "1.txt")
ap.add_argument("-p", "--port", help = "Input PORTS", metavar= "80,443", default = default_port)

q = Queue()
s = requests.Session()
s.mount('http://', HTTPAdapter(max_retries=3))
s.mount('https://', HTTPAdapter(max_retries=3))

def check():
    if(q.empty() == False):
        ip = q.get()
        flag = 0
        if(search(r"(http|https)\:\/\/", ip)): # 当输入URL时提取出域名
            ip = sub(r"(http|https)\:\/\/", "", ip)
            if (search(r"(\/|\\).*", ip)):
                ip = sub(r"(\/|\\).*", "", ip)
                
        if match(r"^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$", ip): #匹配IP
            flag = 1
        if match(r"^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$", ip): #匹配域名
            flag = 1
        if flag == 0:
            ip = ""

        scan(ip)

def scan(ip):
    grs = []
    try:
        if (ip == ""):
            raise Exception("[-]ERROR DOMAIN OR IP!")
        target2 = serverless + "?ip=" + ip + "&port=" + port_list
        response = s.get(target2, timeout = 5, verify = False)
        print("The result of IP={}".format(ip))
        if response.text and response.text != "\"\"":
            for i in response.text.split(","):
                rs.append(ip + ":" + i.replace('\"',''))
                print('[+]{}/TCP OPEN.'.format(i.replace('\"',''))) #读取扫描结果，回显扫描成功的端口信息
    except Exception as err:
        # traceback.print_exc()
        pass
    finally:
        sleep(0.1)
        grs.clear()

if __name__ == '__main__':
    try:
        args = ap.parse_args()
        port_list = args.port
        target = args.url or args.file
        if args.file:
            for i in open(target):
                q.put(i.strip())
        else:
            q.put(target)

        for i in range(20):
            t = Thread(target = check)
            # sleep(0.1)
            t.start()
            t.join()
        
        print("SCAN END!")
        with open("result.txt","w+",encoding='utf8') as f:
            for i in rs:
                f.write(i + "\n")
        print("已保存到result.txt文件中，按回车键退出程序！")
        input() 
        os._exit(0)
    except Exception as err:
        traceback.print_exc()
