from datetime import datetime

# Data pertemuan (id, tanggal)
pertemuan_data = [
    (803, "2025-02-22"),
    (804, "2025-03-01"),
    (805, "2025-03-08"),
    (806, "2025-03-15"),
    (807, "2025-03-22"),
    (808, "2025-03-29"),
    (809, "2025-04-05"),
    (810, "2025-04-12"),
    (811, "2025-04-19"),
    (812, "2025-04-26"),
    (813, "2025-05-03"),
    (814, "2025-05-10"),
    (815, "2025-05-17"),
    (816, "2025-05-24"),
    (817, "2025-02-22"),
    (818, "2025-03-01"),
    (819, "2025-03-08"),
    (820, "2025-03-15"),
    (821, "2025-03-22"),
    (822, "2025-03-29"),
    (823, "2025-04-05"),
    (824, "2025-04-12"),
    (825, "2025-04-19"),
    (826, "2025-04-26"),
    (827, "2025-05-03"),
    (828, "2025-05-10"),
    (829, "2025-05-17"),
    (830, "2025-05-24"),
]

# Mahasiswa
npm = "210711333"
status = "hadir"

# Batas tanggal sekarang
now = datetime(2025, 4, 7)

# Buat query INSERT untuk presensi yang tanggalnya sudah lewat
for id_pertemuan, tanggal_str in pertemuan_data:
    tanggal = datetime.strptime(tanggal_str, "%Y-%m-%d")
    if tanggal < now:
        waktu = tanggal.replace(hour=8, minute=0, second=0)
        print(f"INSERT INTO presensi (npm, waktu_presensi, status, id_pertemuan) "
              f"VALUES ('{npm}', '{waktu.strftime('%Y-%m-%d %H:%M:%S')}', '{status}', {id_pertemuan});")
