#!/usr/bin/env python
# -*- coding:utf-8 -*-

import grequests
# from requests import get
import socket
from re import *
from sys import argv
import os
import argparse
from time import sleep,time
from config import *
import readline
from queue import Queue
from threading import Thread
import traceback

rs = [] #存放结果
url = [] #存放正确url
grs = [] #存放扫描结果
ap = argparse.ArgumentParser()
group = ap.add_mutually_exclusive_group()
group.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "127.0.0.1")
group.add_argument("-f", "--file", help = "Input FILENAME", metavar = "1.txt")
ap.add_argument("-p", "--port", help = "Input PORTS", metavar= "80,443", default = default_port)


def check(ip):
        flag = 0
        if(search(r"(http|https)\:\/\/", ip)): # 当输入URL时提取出域名
            ip = sub(r"(http|https)\:\/\/", "", ip)
            if (search(r"(\/|\\).*", ip)):
                ip = sub(r"(\/|\\).*", "", ip)
                
        if match(r"^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$", ip): #匹配IP
            url.append(ip)
        elif match(r"^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$", ip): #匹配域名
            url.append(ip)

def scan():
    try:
        for i in url:
            target2 = serverless + "?ip=" + i + "&port=" + port_list
            grs.append(grequests.get(target2, timeout = 5, verify = False))
    except Exception as err:
        traceback.print_exc()
        pass


if __name__ == '__main__':
    try:
        args = ap.parse_args()
        port_list = args.port
        target = args.url or args.file
        if args.file:
            for i in open(target):
                check(i.strip())
        else:
            check(target.strip())

        if (url == []):
            raise Exception("[-]ERROR DOMAIN OR IP!")
        else:
            scan()

        a = grequests.map(grs)
        for i in range(len(url)):
            print("The result of IP={}".format(url[i]))
            if a[i] != None and a[i].text != "\"\"":
                for j in a[i].text.replace('\"','').split(","):
                    rs.append(url[i] + ":" + j.replace('\"',''))
                    print('[+]{}/TCP OPEN.'.format(j)) #读取扫描结果，回显扫描成功的端口信息
        
        print("SCAN END!")
        with open("result.txt","w+",encoding='utf8') as f:
            for i in rs:
                f.write(i + "\n")
        print("已保存到result.txt文件中，按回车键退出程序！")
        input() 
        os._exit(0)
    except Exception as err:
        traceback.print_exc()