#!/usr/bin python3

import os
import uselist

if os.path.isfile('go.work'):
    os.system('rm -f go.work')

res = uselist.getlist()
if None == res:
    print("failed :(")
    exit(1)

os.system(f'go work init {res}')
