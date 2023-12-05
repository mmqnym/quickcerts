# QCS SDK for python
# @Author: mmq88

from typing import List

class QCSCreateSNResponse:
    def __init__(self, msg: str, serial_number: str) -> None:
        self.msg = msg
        self.serial_number = serial_number

class QCSGenerateSNResponse:
    def __init__(self, msg: str, serial_numbers: str) -> None:
        self.msg = msg
        self.serial_numbers = serial_numbers

class QCSRecord:
    def __init__(self, sn: str, key: str, note: str) -> None:
        self.sn = sn
        self.key = key
        self.note = note

class QCSAllRecordsResponse:
    def __init__(self, data: List[QCSRecord]) -> None:
        self.data = data

class QCSAvailableSNResponse:
    def __init__(self, data: List[str]) -> None:
        self.data = data

class QCSUpdateSNNoteResponse:
    def __init__(self, msg: str, note: str) -> None:
        self.msg = msg
        self.note = note

class QCSApplyCertResponse:
    def __init__(self, key: str, signature: str) -> None:
        self.key = key
        self.signature = signature

class QCSApplyTempPermitResponse:
    def __init__(self, remaining_time: str, status: str) -> None:
        self.remaining_time = remaining_time
        self.status = status