# QCS SDK for python
# @Author: mmq88

import requests
import pyqcs_type

class QCSAdmin:
    def __init__(self, 
                 host: str,
                 port: int,
                 api_path: str = "/api/v1", 
                 tls: bool = False, 
                 access_token: str = "", 
                 runtime_code: str = ""
                 ) -> None:
        self.access_prefix =f'{host}:{port}{api_path}'
        self.access_token = access_token
        self.runtime_code = runtime_code

        if tls:
            self.access_prefix = "https://" + self.access_prefix
        else:
            self.access_prefix = "http://" + self.access_prefix

    def create_sn(self, sn: str, reason: str = "none") -> pyqcs_type.QCSCreateSNResponse:
        '''
        Add a serial number created by admin to QCS.
        sn: serial number.
        reason: reason for creating this serial number.
        '''

        url = self.access_prefix + "/sn/create"
        headers = {"X-Access-Token": self.access_token, "X-Runtime-Code": self.runtime_code}
        body = {"serial_number": sn, "reason": reason}
        res = requests.post(url, headers=headers, json=body)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()
            return pyqcs_type.QCSCreateResponse(data["msg"], data["serial_number"])

    def generate_sn(self, count: int, reason: str = "none") -> pyqcs_type.QCSGenerateSNResponse:
        '''
        Generate serial number(s) randomly.
        count: number of serial numbers to generate.
        reason: reason for generating these serial numbers.
        '''

        url = self.access_prefix + "/sn/generate"
        headers = {"X-Access-Token": self.access_token, "X-Runtime-Code": self.runtime_code}
        body = {"count": count, "reason": reason}
        res = requests.post(url, headers=headers, json=body)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()
            return pyqcs_type.QCSGenerateSNResponse(data["msg"], data["serial_numbers"])

    def get_all_record(self) -> pyqcs_type.QCSAllRecordResponse:
        '''
        Get all records in QCS.
        '''

        url = self.access_prefix + "/sn/get-all"
        headers = {"X-Access-Token": self.access_token, "X-Runtime-Code": self.runtime_code}
        res = requests.get(url, headers=headers)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()
            records = []
            for record in data["data"]:
                records.append(pyqcs_type.QCSRecord(record["sn"], record["key"], record["note"]))
            return pyqcs_type.QCSAllRecordResponse(records)
        
    def get_available_sn(self) -> pyqcs_type.QCSAvailableSNResponse:
        '''
        Get all available serial numbers in QCS.
        '''

        url = self.access_prefix + "/sn/get-available"
        headers = {"X-Access-Token": self.access_token, "X-Runtime-Code": self.runtime_code}
        res = requests.get(url, headers=headers)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()
            return pyqcs_type.QCSAvailableSNResponse(data["data"])
    
    def update_sn_note(self, target_sn: str, note: str) -> pyqcs_type.QCSUpdateSNNoteResponse:
        '''
        Update note of a serial number.
        target_sn: serial number to update.
        note: new note.
        '''

        url = self.access_prefix + "/sn/update-note"
        headers = {"X-Access-Token": self.access_token, "X-Runtime-Code": self.runtime_code}
        body = {"serial_number": target_sn, "note": note}
        res = requests.post(url, headers=headers, json=body)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()
            return pyqcs_type.QCSUpdateSNNoteResponse(data["msg"], data["note"])


class QCSClient:
    def __init__(self, 
                 host: str,
                 port: int,
                 api_path: str = "/api/v1", 
                 tls: bool = False, 
                 access_token: str = ""
                 ) -> None:
        self.access_prefix =f'{host}:{port}{api_path}'
        self.access_token = access_token

        if tls:
            self.access_prefix = "https://" + self.access_prefix
        else:
            self.access_prefix = "http://" + self.access_prefix

    def apply_cert(self, sn: str, board_producer: str, board_name: str, mac_address: str) \
        -> pyqcs_type.QCSApplyCertResponse:
        '''
        Use a serial number and device information to apply for a certificate.
        sn: serial number.
        board_producer: board producer.
        board_name: board name.
        mac_address: physical ethernet mac address.
        '''

        url = self.access_prefix + "/apply/cert"
        headers = {"X-Access-Token": self.access_token}
        body = {
            "serial_number": sn,
            "board_producer": board_producer,
            "board_name": board_name,
            "mac_address": mac_address
        }
        
        res = requests.post(url, headers=headers, json=body)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()   
            return pyqcs_type.QCSApplyCertResponse(data["key"], data["signature"])
        
    def apply_temp_permit(self, board_producer: str, board_name: str, mac_address: str) \
        -> pyqcs_type.QCSApplyTempPermitResponse:
        '''
        Use device information to apply for a temporary permit(with time limit certificate).
        board_producer: board producer.
        board_name: board name.
        mac_address: physical ethernet mac address.
        '''

        url = self.access_prefix + "/apply/temp-permit"
        headers = {"X-Access-Token": self.access_token}
        body = {
            "board_producer": board_producer,
            "board_name": board_name,
            "mac_address": mac_address
        }

        res = requests.post(url, headers=headers, json=body)

        if res.status_code != 200:
            raise Exception("QCS::Error:" + res.json()["error"])
        else:
            data = res.json()   
            return pyqcs_type.QCSApplyTempPermitResponse(data["remaining_time"], data["status"])