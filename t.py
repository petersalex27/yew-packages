import sys
import time

if len(sys.argv) < 3:
  print("len(argv) < 3")
  exit(1)

res = int(sys.argv[1])
res2 = time.gmtime(res)
hex_ = (sys.argv[2])[:12]
print('@v0.0.0-'+time.strftime('%Y%m%d%H%M%S', res2)+'-'+hex_)