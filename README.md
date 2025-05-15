CSVをgroup by countするツール

### 対応する形式
- ヘッダーありCSVの列名指定
- ヘッダーなしCSVの列番号指定(ヘッダーあってもスキップされないだけで集計はできるが)

### オプション
- file=: CSVファイルのパス(必須)
- column=: 列名 or 列番号(必須)
- header: 列名指定の場合必須
- delimiter=: 区切り文字1字(任意。デフォルトは`,`)
- buffer=: バッファサイズKB(任意。デフォルトは`4KB`)

### コマンド例
区切り文字指定、ヘッダー行あり、列名指定
```bash
$ go run main.go -file=sample_data/sample_500.csv -column=姓 -delimiter="," --header
```
結果
```text
2025/05/12 01:20:40 Column 姓 found at index 1
Group by 姓 counts

Davis: 54
Rodriguez: 57
Jones: 37
Wilson: 40
Miller: 46
Brown: 46
Williams: 58
Garcia: 56
Johnson: 43
Smith: 63

```

区切り文字指定、ヘッダー行なし、列番号指定
```bash
$ go run main.go -file=sample_data/sample_500.csv -column=4
```
結果
```text
Column index provided: 4
Group by 4 counts

性別: 1
男: 247
女: 253
```

other example
```bash
go run main.go -file=sample_data/sample.csv -column=住所 -delimiter="," --header --buffer=$((1024*10))
```