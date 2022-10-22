import json


def toTree(arr, parent_id, cache=[], level=0):
    tree = []
    level += 1
    print(cache, level)
    for item in arr:
        if item["id"] in cache:
            continue
        # 46次运算
        if item["parent_id"]==parent_id or (item["parent_id"]==item["id"] and level==parent_id):
            # level==parent_id 把顶层分类控制在第一层循环
            cache.append(item["id"])
            item["child"] = toTree(arr, item["id"], cache, level)
            tree.append(item)
    return tree


# 递归 yield
def toNode(arr, parent_id, cache=[]):
    for item in arr:
        if item["id"] in cache:
            continue

        if item["parent_id"]==parent_id:
            cache.append(item["id"])
            for i in toNode(arr, item["id"], cache):
                yield i
            
            yield item


if __name__ == "__main__":
    arr = [
        {"id": 1,"name": "水果","parent_id": 1,"level": 0},
        {"id": 2,"name": "李子","parent_id": 1,"level": 1},
        {"id": 3,"name": "小李子","parent_id": 2,"level": 2},
        {"id": 4,"name": "桃子","parent_id": 1,"level": 1},
        {"id": 5,"name": "苹果","parent_id": 1,"level": 1},
        {"id": 6,"name": "蔬菜","parent_id": 6,"level": 0},
        {"id": 7,"name": "波菜","parent_id": 6,"level": 1},
        {"id": 8,"name": "韭菜","parent_id": 6,"level": 1},
        {"id": 9,"name": "日用","parent_id": 9,"level": 0,"is_tab": True},
        {"id": 10,"name": "小桃子","parent_id": 2,"level": 2},
    ]
    # rsp = toTree(arr, 1, [])
    # print(rsp)
    # rsp = toTree(arr, 1, [])
    # print(rsp)


    rsp = toNode(arr, 1, [])
    for item in rsp:
        print(item["id"])
        
