import os

def getlist_array(dir='.'):
    if not os.path.isfile(f'{dir}/use.list'):
        print('no file named use.list')
        return None
    
    f = open(f'{dir}/use.list')
    lines = f.readlines()
    f.close()

    for i in range(len(lines)):
        lines[i] = lines[i].strip()
    
    return lines


def getlist(dir='.'):
    lines = getlist_array(dir)
    if lines == None:
        return None
    
    return " ".join(lines)