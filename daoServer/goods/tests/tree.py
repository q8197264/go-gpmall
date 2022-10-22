import json


def generate_tree(source, parent, cache=[]):
    tree = []
    for item in source:
        if item["id"] in cache:
            continue
        if item["parent"] == parent:
            cache.append(item["id"])
            item["child"] = generate_tree(source, item["id"], cache)
            tree.append(item)
    return tree 


if __name__ == '__main__':
    source = [
        {"id": 1, "name": '电器', "parent": 0},
        {"id": 2, "name": '水果', "parent": 0},
        {"id": 3, "name": '家用电器', "parent": 1},
        {"id": 4, "name": '电吹风', "parent": 2},
        {"id": 5, "name": '电风扇', "parent": 3},
        {"id": 6, "name": '台灯', "parent": 3},
        {"id": 7, "name": '商用电器', "parent": 1},
        {"id": 8, "name": '大型电热锅', "parent": 7},
    ]
    arr = [
        {"id": 1,"name": "水果","parent_id": 1},
        {"id": 2,"name": "李子","parent_id": 1,"level": 1},
        {"id": 3,"name": "小李子","parent_id": 2,"level": 2},
        {"id": 4,"name": "桃子","parent_id": 1,"level": 1},
        {"id": 5,"name": "苹果","parent_id": 1,"level": 1},
        {"id": 6,"name": "蔬菜","parent_id": 6},
        {"id": 7,"name": "波菜","parent_id": 6,"level": 1},
        {"id": 8,"name": "韭菜","parent_id": 6,"level": 1},
        {"id": 9,"name": "日用","parent_id": 9,"is_tab": True},
        {"id": 10,"name": "小桃子","parent_id": 2,"level": 2},
    ]
    rsp = generate_tree(arr, 1)
    print(rsp)
    # print(json.dumps(permission_tree, ensure_ascii=False)) 