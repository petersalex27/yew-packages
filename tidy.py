import os
import sys
import uselist
                
args = sys.argv[1:]

def tidyUseList():
    lines = uselist.getlist_array()
    if lines != None:
        for line in lines:
            tidy(line)

def tidy(where: str):
    cwd = os.getcwd()
    where = where.removeprefix('./')
    os.chdir(f'./{where}')
    os.system(f'go mod tidy')
    os.chdir(cwd)

for arg in args:
    if arg == '-u':
        tidyUseList()
    else:
        tidy(arg)