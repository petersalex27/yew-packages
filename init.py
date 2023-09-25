import os
import sys
import uselist
                
args = sys.argv[1:]
prefix = 'github.com/petersalex27/yew-packages'
if args[0] == '-p':
    prefix = args[1]
    args = args[2:]

def modUseList():
    lines = uselist.getlist_array()
    if lines != None:
        for line in lines:
            mod(line)

def mod(where: str):
    cwd = os.getcwd()
    where = where.removeprefix('./')
    if where == '.' or where == '':
        where = ''
    else:
        where = f'/{where}'

    os.chdir(f'.{where}')

    if os.path.isfile('go.mod'):
        os.system('rm -f go.mod')

    print(prefix)
    os.system(f'go mod init {prefix}{where}')
    os.chdir(cwd)

for arg in args:
    if arg == '-u':
        modUseList()
    else:
        mod(arg)