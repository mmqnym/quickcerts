import sys
sys.path.append('..')

from pyqcs import QCSAdmin, QCSClient


qcs_admin = QCSAdmin("localhost", 33333, access_token="0b09b6dc41f61813346ba76322d19e07a0b71ba939a1bf90211dfff40f552ed0")

# Create a serial number.
try:
    res = qcs_admin.create_sn("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX")
    print(res.msg, res.serial_number)
except Exception as e:
    print(e.args[0])

# Generate serial number(s).
try:
    res = qcs_admin.generate_sn(5)
    print(res.msg, res.serial_numbers)
except Exception as e:
    print(e.args[0])

# Get all records.
try:
    res = qcs_admin.get_all_records()
    for record in res.data:
        print(record.serial_number, record.key, record.note)
except Exception as e:
    print(e.args[0])

# Get all available serial numbers.
try:
    res = qcs_admin.get_available_sn()
    print(res.data)
except Exception as e:
    print(e.args[0])

# Update note of a serial number.
try:
    res = qcs_admin.update_sn_note(target_sn="XXXX-XXXX-XXXX-XXXX-XXXX-XXXX", note="test")
    print(res.msg, res.note)
except Exception as e:
    print(e.args[0])



qcs_client = QCSClient("localhost", 33333, access_token="QcsTestToken********************************")

# Apply certificate.
try:
    res = qcs_client.apply_cert(
        sn="XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
        board_producer="ASUSTeK Computer Inc.",
        board_name="ROG STRIX Z790-A GAMING WIFI",
        mac_address="XXXXXXXXXXXX"
    )
    print(res.key, res.signature)
except Exception as e:
    print(e.args[0])

# Apply temporary permit(with time limit certificate).
try:
    res = qcs_client.apply_temp_permit(
        board_producer="ASUSTeK Computer Inc.",
        board_name="ROG STRIX Z790-A GAMING WIFI",
        mac_address="XXXXXXXXXXXX"
    )
    print(res.remaining_time, res.status)
except Exception as e:
    print(e.args[0])