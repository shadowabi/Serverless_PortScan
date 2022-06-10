#!/usr/bin/env python
# -*- coding:utf-8 -*-
#author:Sh4d0w_小白

import grequests
import socket
from re import *
from sys import argv
import os
import argparse
from time import sleep

serverless = "" #云函数地址
port_list = []
default_port = "21,22,23,25,80,135,139,443,445,888,1433,1521,3306,3389,5985,5986,6379,8080,27019" #默认端口
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
        flag = 1
    if match(r"^([a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,11}$", ip): #匹配域名
        flag = 1
    if flag == 0:
        ip = ""

    scan(ip)

def scan(ip):
    rs = []
    try:
        for i in port_list:
            if (ip == ""):
                raise Exception("[-]ERROR DOMAIN OR IP!")
            if not i.isdigit():
                raise Exception("[-]ERROR PORTS！")
            target2 = serverless + "?ip=" + ip + "&port=" + str(i)
            rs.append(grequests.get(target2, timeout = 3, verify = False)) #扫描
        print("The result of IP={}".format(ip))
        for j in grequests.map(rs):
            if j != None and j.text != "null" and j.text.find("errorCode") == -1:
                print('[+]{}/TCP OPEN.'.format(j.text.replace('\"',''))) #读取扫描结果，回显扫描成功的端口信息
    except Exception as err:
        print(err)
        pass
    finally:
        sleep(0.1)
        rs.clear()

if __name__ == '__main__':
    try:
        args = ap.parse_args()
        port_list = args.port.split(",")
        target = args.url or args.file
        if args.file:
            for i in open(target):
                check(i.strip())
                print("SCAN END!")
        else:
            check(target)
            print("SCAN END!")
        
        os._exit(0)
    except Exception as err:
        print(err)
