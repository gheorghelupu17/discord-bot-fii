
import json
import csv

licenta = []
seria = {}
optionals = []
notInOptional = []
with open('./data/licenta.csv', 'r') as f:
    for row in csv.reader(f, delimiter=';'):
        seria[row[5]]=[row[1], row[2]]

with open('./data/optionale_nou.csv', 'r') as f:
    for row in csv.reader(f, delimiter=';'):
        tmpOptional = seria.get(row[5], None)
        if tmpOptional is not None:
            #3;E;4;VASILESCU;STEFAN;310910401ESL201088;TPM;4
            optionals.append([row[0],tmpOptional[0], tmpOptional[1],row[3],row[4],row[5],row[1],row[2]] )
        else:
            print('Optional not found for: ', row[5],'\n')
            notInOptional.append([row[5]])
        
with open('./data/optionale.csv', 'w') as f:
    writer = csv.writer(f, delimiter=';')
    writer.writerows(optionals)
with open('./data/notInOptional.csv', 'w') as f:
    writer = csv.writer(f, delimiter=';')
    writer.writerows(notInOptional)