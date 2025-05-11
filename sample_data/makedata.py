import pandas as pd
import random
import time
import json
import os

# Base32 Crockford alphabet for ULID
CROCKFORD = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

def encode_base32(value, length):
    encoded = ""
    for _ in range(length):
        encoded = CROCKFORD[value & 0x1F] + encoded
        value >>= 5
    return encoded

def generate_ulid():
    # Timestamp: 48 bits (milliseconds)
    ts = int(time.time() * 1000) & ((1 << 48) - 1)
    # Randomness: 80 bits
    rand = random.getrandbits(80)
    return encode_base32(ts, 10) + encode_base32(rand, 16)

# Sample data pools
last_names = ["Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia", "Rodriguez", "Wilson"]
first_names = ["James", "Mary", "Robert", "Patricia", "John", "Jennifer", "Michael", "Linda", "William", "Elizabeth"]
prefectures = ["北海道","青森県","岩手県","宮城県","秋田県","山形県","福島県","茨城県","栃木県","群馬県",
               "埼玉県","千葉県","東京都","神奈川県","新潟県","富山県","石川県","福井県","山梨県","長野県",
               "岐阜県","静岡県","愛知県","三重県","滋賀県","京都府","大阪府","兵庫県","奈良県","和歌山県",
               "鳥取県","島根県","岡山県","広島県","山口県","徳島県","香川県","愛媛県","高知県","福岡県",
               "佐賀県","長崎県","熊本県","大分県","宮崎県","鹿児島県","沖縄県"]
hobbies = ["reading", "traveling", "cooking", "sports", "music", "gaming", "art"]
styles = ["casual", "formal", "sporty", "elegant", "vintage"]
music = ["rock", "pop", "jazz", "classical", "hiphop", "electronic"]
books = ["1984", "Pride and Prejudice", "To Kill a Mockingbird", "The Great Gatsby", "Moby Dick"]

def create_dataset(n):
    rows = []
    for _ in range(n):
        uid = generate_ulid()
        last = random.choice(last_names)
        first = random.choice(first_names)
        age = random.randint(0, 60)
        gender = random.choice(["男", "女"])
        address = random.choice(prefectures)
        email = f"{first.lower()}.{last.lower()}{random.randint(1,1000)}@example.com"
        
        # Random subset of "other" keys
        other = {}
        if random.random() < 0.7:
            other["hobby"] = random.choice(hobbies)
        if random.random() < 0.5:
            other["style"] = random.choice(styles)
        if random.random() < 0.6:
            other["favorite_music"] = random.choice(music)
        if random.random() < 0.4:
            other["favorite_book"] = random.choice(books)
        other_json = json.dumps(other, ensure_ascii=False)
        
        rows.append([uid, last, first, age, gender, address, email, other_json])
    df = pd.DataFrame(rows, columns=["ID", "姓", "名", "年齢", "性別", "住所", "メールアドレス", "その他"])
    return df

# Create datasets
df_small = create_dataset(500)
df_large = create_dataset(20000)

# Save CSVs
small_path = './sample_500.csv'
large_path = './sample_20000.csv'
df_small.to_csv(small_path, index=False)
df_large.to_csv(large_path, index=False)

# Display a preview
import ace_tools as tools; tools.display_dataframe_to_user(name="Sample CSV (500 rows) Preview", dataframe=df_small.head())



