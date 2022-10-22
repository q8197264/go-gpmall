def getTree(menu_list,pid):

    tree = []
    for item in menu_list:
        if item["pid"] == pid:
            item["child"] = getTree(menu_list, item["id"]) # 找儿子的过程
            tree.append(item)
    return tree


def getTree(data):
    map = {}
    for item in data:
        map.update({item.id:{
                "name":item.name.strip(),
                "id":item.id,
                "parent_id": item.parent_id.id,
                "level":item.level,
                "is_tab":item.is_tab
            }})

    m1 = {}
    for k,item in map.items():
        item["child"] = []
        m1[item["id"]] = item
    
    tree = []
    for item in m1.values():
        if m1.get(item["parent_id"]) and item["parent_id"]!= item["id"]:# 找儿子
            m1[item["parent_id"]]["child"].append(item)
        else: # 找出所有的顶级
            tree.append(item)
    print(tree)