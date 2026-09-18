import sys
import re

with open('.github/workflows/ci.yaml', 'r') as f:
    content = f.read()

content = content.replace('uses: actions/checkout@v7', 'uses: actions/checkout@v4')
content = content.replace('uses: actions/setup-go@v7', 'uses: actions/setup-go@v4')

# Keep login-action@v3 and upload-artifact@v4

with open('.github/workflows/ci.yaml', 'w') as f:
    f.write(content)
