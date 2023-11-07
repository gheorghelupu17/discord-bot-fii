import csv
import random
import string
import json

import pymongo

choices = set()


def id_generator(size=6, chars=string.ascii_uppercase + string.digits):
    while True:
        choice = ''.join(random.choice(chars) for _ in range(size))
        if choice not in choices:
            choices.add(choice)
            return choice


con = pymongo.MongoClien('')
db = con.get_database('bot')
col = db.get_collection('stud_ids_to_roles')
rolesCol = db.get_collection('roles')
rolesCol.delete_many({})

with open('./data/licenta.csv', 'r') as f:
    for row in csv.reader(f, delimiter=';'):
        nickname = ' '.join(['-'.join(map(lambda x: x.capitalize(), row[3].split('-'))), row[4].capitalize()])
        roles = ['verified', 'licenta', f'year{row[0]}', f'{row[1].lower()}_series', ''.join(row[:3])]
        _id = row[5]
        doc = col.find_one({'_id': _id})
        if doc is None or row[0] == '1':
            rolesCol.insert_one({'_id': _id, 'roles': roles, 'nickname': nickname, 'discord_id': ''})
        else:
            rolesCol.insert_one({'_id': _id, 'roles': roles, 'nickname': nickname, 'discord_id': doc['discord_id']})
    print('licenta done')

with open('./data/master.csv', 'r') as f:
    for row in csv.reader(f, delimiter=';'):
        nickname = ' '.join(['-'.join(map(lambda x: x.capitalize(), row[2].split('-'))), row[3].capitalize()])
        roles = ['verified', 'master', f'{row[0]}M{row[1]}', f'M{row[1]}']
        _id = row[4]
        doc = col.find_one({'_id': _id})
        if doc is None:
            doc = col.find_one({'nickname': nickname})
        if doc is None:
            rolesCol.insert_one({'_id': _id, 'roles': roles, 'nickname': nickname, 'discord_id': ''})
        else:
            if 'year3' in doc['roles'] or 'master' in doc['roles']:
                rolesCol.insert_one({'_id': _id, 'roles': roles, 'nickname': nickname, 'discord_id': doc['discord_id']})
    print('master done')

with open('./data/optionale.csv', 'r') as f:
    for row in csv.reader(f, delimiter=';'):
        _id = row[5]
        role = row[0] + ''.join(row[6:])
        doc = rolesCol.find_one({'_id': _id})
        if doc is None:
            raise f'user with id: {_id} not found'
        rolesCol.update_one({'_id': _id}, {'$push': {'roles': role}})
    print('optionale done')

with open('./data/teacher.json') as json_file:
     data = json.load(json_file)
     for row in data:
           print(row)
           doc = rolesCol.find_one({'_id': row['_id']})
           if doc is None:
               rolesCol.insert_one(row)
           else:
               rolesCol.update_one({'_id': row['_id']}, {'$push': {'roles': row['roles']}})   
     print('teacher done')