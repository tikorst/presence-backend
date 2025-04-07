from datetime import datetime, timedelta
import random

# Konfigurasi
npm_list = ["210711333"]  # Bisa ditambah kalau perlu banyak NPM
status_list = ["hadir"]   # Bisa kembangkan ke "izin", "sakit", "alpa" kalau mau variasi
presensi_dummy = []

# Range id_pertemuan
for id_pertemuan in range(145, 803):  # 146 sampai 802
    for npm in npm_list:
        # Buat waktu antara jam 07.00 - 09.00 secara random
        jam_random = random.randint(7, 9)
        menit_random = random.randint(0, 59)
        waktu_presensi = datetime(2024, 3, 1, jam_random, menit_random)  # base tanggal dummy
        waktu_str = waktu_presensi.strftime('%Y-%m-%d %H:%M:%S')

        status = random.choice(status_list)

        presensi_dummy.append(
            f"INSERT INTO presensi (id_pertemuan, npm, waktu_presensi, status) VALUES ({id_pertemuan}, '{npm}', '{waktu_str}', '{status}');"
        )

# Simpan ke file
with open("presensi_dummy.sql", "w") as f:
    f.write("\n".join(presensi_dummy))
