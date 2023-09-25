import os
import sys
import uselist
                
args = sys.argv[1:]

def remUseList():
    lines = uselist.getlist_array()
    if lines != None:
        for line in lines:
            line = line.removeprefix('./')
            if line == '':
                line = '.'
            os.system(f'sh clean.sh {line}')

for arg in args:
    if arg == '-u':
        remUseList()
    else:
        os.system(f'sh clean.sh {arg}')