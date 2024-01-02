import dis
import marshal

with open('__pycache__/hvmc.cpython-311.pyc', 'rb') as f:
    f.seek(16)
    dis.dis(marshal.load(f))