import json
import mysql.connector


conf = ""
field_str = ""
value_str = ""


def init_conf(conf_path):
    fd = open(conf_path)
    try:
        conf_json = fd.read()
    finally:
        fd.close()
    global conf
    conf = json.loads(conf_json)


def init_field_str():
    global field_str
    field_dic = conf["database"]["field"]
    for key in field_dic.keys():
        field_str = field_str + key + ","
    data_list = conf["database"]["data"]
    if len(data_list) > 0:
        for key in data_list[0].keys():
            field_str = field_str + key + ","
    field_str = field_str[0:-1]


def init_value_str():
    global value_str
    field_val = ""
    field_dic = conf["database"]["field"]
    for val in field_dic.values():
        field_val = field_val + "'" + val + "'" + ","
    data_list = conf["database"]["data"]
    for data in data_list:
        if field_val != "":
            value_str = value_str + "(" + field_val
        for val in data.values():
            value_str = value_str + "'" + val + "'" + ","
        value_str = value_str[0:-1]
        value_str = value_str + "),"
    value_str = value_str[0:-1]


def get_insert_sql():
    db_name = conf["database"]["db"]
    table_name = conf["database"]["table"]
    table = db_name + "." + table_name
    global field_str
    global value_str
    insert_sql = "insert into {} ({}) values {}".format(table, field_str, value_str)
    return insert_sql


def run_sql():
    db_config = {}
    db_config['host'] = conf["server"]["ip"]
    db_config['port'] = conf["server"]["port"]
    db_config['user'] = conf["server"]["user"]
    db_config['password'] = conf["server"]["password"]

    try:
        db = mysql.connector.connect(**db_config)
        insert_sql = get_insert_sql()
        print("insert sql:", insert_sql)
        db.cursor().execute(insert_sql)
        db.commit()
    except Exception as e:
        print("insert error:", e)
    finally:
        db.close()


if __name__ == "__main__":
    init_conf("conf.json")
    init_field_str()
    init_value_str()
    run_sql()
