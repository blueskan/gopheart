import os
import sys

hostname = "batikansenemoglu.com"
response = os.system("ping -c 1 " + hostname)

if response == 0:
  sys.exit(0)
else:
  sys.exit(1)