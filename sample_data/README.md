```bash
python3 -m venv .venv
./.venv/bin/pip install pandas
./.venv/bin/python makedata.py -n 200000 -o sample.csv
```
 
```text
$ ./.venv/bin/python makedata.py -h
usage: makedata.py [-h] [-n ROWS] -o OUTPUT

サンプルCSVを生成します

options:
  -h, --help            show this help message and exit
  -n ROWS, --rows ROWS  生成する行数 (デフォルト: 500)
  -o OUTPUT, --output OUTPUT
                        出力ファイルパス (例: sample.csv)
```