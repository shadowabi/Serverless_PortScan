#!/usr/bin/env python
# -*- coding:utf-8 -*-
#author:Sh4d0w_小白

import grequests
import socket
from re import *
from sys import argv
import os
import argparse

serverless = "" #云函数地址
port_list = []
default_port = "21,22,23,25,80,135,139,443,445,888,1433,1521,3306,3389,5985,5986,6379,8080,27019" #默认端口
ap = argparse.ArgumentParser()
ap.add_argument("-u", "--url", help = "Input IP/DOMAIN/URL", metavar = "127.0.0.1", required = True)
ap.add_argument("-p", "--port", help = "Input PORTS", metavar= "80,443", default = default_port)

def check(ip):
    if(search(r"(http|https)\:\/\/", ip)): # 当输入URL时提取出域名
        ip = sub(r"(http|https)\:\/\/", "", ip)
        if (search(r"(\/|\\).*", ip)):
            ip = sub(r"(\/|\\).*", "", ip)
    try:
        res = socket.getaddrinfo(ip, None) # 若输入域名则提取IP
        return(res[0][4][0])
    except Exception as err: # 非法输入直接返回空值
        return ""

def main():
    rs = []
    try:
        for i in port_list:
            if (ip == ""):
                raise Exception("[-]ERROR DOMAIN OR IP!")
            target = serverless + "?ip=" + ip + "&port=" + str(i)
            rs.append(grequests.get(target, timeout = 3, verify = False)) #扫描
        print("The result of IP={}".format(ip))
        for i in grequests.map(rs):
            if i != None and i.text != "null" and i.text.find("errorCode") == -1:
                print('[+]{}/TCP OPEN.'.format(i.text.replace('\"',''))) #读取扫描结果，回显扫描成功的端口信息
    except Exception as err:
        print(err)
        pass
    finally:
        print("SCAN END!")
        os._exit(0)

if __name__ == '__main__':
    try:
        args = ap.parse_args()
        ip = check(args.url)
        port_list = args.port.split(",")
        main()
    except Exception as err:
        print(err)
