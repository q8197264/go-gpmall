import os
import sys
import subprocess
import shutil

"""
    功能：同步proto
        1. 拷贝python的proto到go的对应目录下
        2. 生成python的源码 - import . 
        3. 生成go的源码
"""
class cd:
    def __init__(self, newPath):
        # with init
        # 把path中包含的"~"和"~user"转换成用户目录
        self.newPath = os.path.expanduser(newPath)

    def __enter__(self):
        # with enter
        self.savedPath = os.getcwd()
        os.chdir(self.newPath)

    def __exit__(self, etype, value, traceback):
        # with exit
        os.chdir(self.savedPath)


def replace_file(file_name):
    # 修改 ?_pb2_grpc.py 文件中 import 为 from . import

    new_file_name= f"{file_name}-back"
    modify_times = 0 #统计修改次数
    f1 = open(file_name, "r", encoding="utf-8") #以r模式打开旧文件
    f2 = open(new_file_name, "w", encoding="utf-8") #以w模式重写或创建新文件
    for lines in f1: # 循环读内容
        if lines.startswith("import") and not lines.startswith("import grpc"):
            lines = f"from . {lines}"
            modify_times += 1
        # 把 import 改为 from . import 
        f2.write(lines)

    print("文件修改次数: ", modify_times)

    f1.close()
    f2.close() # 关闭文件 f2(同时打开多个文件, 先打开先关闭，后打开的后关闭)
    os.replace(new_file_name, file_name) #修改（替换）文件名


def proto_file_list(path):
    # 递归目录下所有文件
    flist = []
    lsdir = os.listdir(path)
    dirs = [ i for i in lsdir if os.path.isdir(os.path.join(path, i)) ]
    if dirs:
        for i in dirs:
            # 递归多级目录
            proto_file_list(os.path.join(path, i))
    
    # 获取当前目录下所有 proto 文件
    files = [ i for i in lsdir if os.path.isfile(os.path.join(path, i)) ]
    for file in files:
        if file.endswith(".proto"):
            flist.append(file)
    
    return flist


def copy_form_py_to_go(src_dir, dist_dir):
    # 拷贝proto源文件到目标目录

    proto_files = proto_file_list(src_dir)
    for proto_file in proto_files:
        try:
            # 从 src目录 拷贝到 dist目录
            shutil.copy(f"{src_dir}/{proto_file}", dist_dir)
        except IOError as e:
            print("Unable to copy file. %s" % e)
        except:
            print("Unexpected error:", sys.exc_info())


def generated_pyfile_list(path):
    # 获取目录下所有文件(递归所有目录)
    flist = []
    lsdir = os.listdir(path)

    # 递归目录
    dirs = [ i for i in lsdir if os.path.isdir(os.path.join(path, i)) ]
    if dirs:
        for i in dirs:
            proto_file_list(os.path.join(path, i))

    # 获取当前目录所有文件
    files = [ i for i in lsdir if os.path.isfile(os.path.join(path, i)) ]
    for file in files:
        if file.endswith(".py"):
            flist.append(file)

    return flist


class ProtoGenerator:
    def __init__(self, python_dir, go_dir):
        self.python_dir = python_dir
        self.go_dir = go_dir

    def generate(self):
        with cd(self.python_dir):
            files = proto_file_list(self.python_dir)
            # subprocess.call("workon daoserver", shell=True) #切换虚拟环境
            for file in files:
                command = f"python3 -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. {file}"
                subprocess.call(command, shell=True)

            # 改 ?_pb2_grpc.py中 import 为 from . import
            py_files = generated_pyfile_list(self.python_dir)
            for file_name in py_files:
                replace_file(file_name)
            
            # 展示生成的文件
            print(py_files)

        with cd(self.go_dir):
            files = proto_file_list(self.go_dir)
            for file in files:
                command = f"protoc -I . {file} --go_out=plugins=grpc:."
                subprocess.call(command, shell=True)


if "__main__"==__name__:
    #
    srv = "goods"
    if srv == "goods":
        python_dir = "/Users/www/learn-note/gpmall/daoServer/order/proto"
        go_dir = "/Users/www/learn-note/gpmall/webServer/order/proto"
    elif srv == "user":
        python_dir = "/Users/www/learn-note/gpmall/daoServer/order/proto"
        go_dir = "/Users/www/learn-note/gpmall/webServer/order/proto"

    # 将python项目拷贝到go项目下
    copy_form_py_to_go(python_dir, go_dir)

    # 编译proto源码
    pr = ProtoGenerator(python_dir, go_dir)
    pr.generate()