import { QCSType } from "./qcs_type";

export class QCSAdmin {
  accessPrefix: string;
  accessToken: string;
  runtimeCode: string;

  constructor(qcsAdminConfig: QCSType.QCSAdminConfig) {
    const c = qcsAdminConfig;
    this.accessPrefix = `${c.host}:${c.port}${c.apiPath}`;
    this.accessToken = c.accessToken ? c.accessToken : "";
    this.runtimeCode = c.runtimeCode ? c.runtimeCode : "";

    if (c.tls) {
      this.accessPrefix = "https://" + this.accessPrefix;
    } else {
      this.accessPrefix = "http://" + this.accessPrefix;
    }
  }

  /**
   * Add a serial number created by admin to QCS.
   *
   * @param sn: serial number.
   * @param reason: reason for creating this serial number.
   *
   * @return: Message and the sn were created if success, error if failed.
   */
  async createSN(
    sn: string,
    reason: string = "none"
  ): Promise<QCSType.QCSCreateResponse> {
    const url = this.accessPrefix + "/sn/create";
    console.log(url);
    const headers = {
      "X-Access-Token": this.accessToken,
      "X-Runtime-Code": this.runtimeCode,
    };

    const body = {
      serial_number: sn,
      reason: reason,
    };

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as { msg: string; serial_number: string };
      const result = {
        msg: data["msg"],
        serial_number: data["serial_number"],
      } as QCSType.QCSCreateResponse;

      return result;
    }
  }

  /**
   * Generate serial number(s) randomly.
   *
   * @param count: number of serial numbers to generate.
   * @param reason: reason for creating this serial number.
   *
   * @return: Message and the sn(s) were created if success, error if failed.
   */
  async generateSN(
    count: number,
    reason: string = "none"
  ): Promise<QCSType.QCSGenerateSNResponse> {
    const url = this.accessPrefix + "/sn/generate";
    const headers = {
      "X-Access-Token": this.accessToken,
      "X-Runtime-Code": this.runtimeCode,
    };

    const body = {
      count: count,
      reason: reason,
    };

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        msg: string;
        serial_numbers: Array<string>;
      };
      const result = {
        msg: data["msg"],
        serial_numbers: data["serial_numbers"],
      } as QCSType.QCSGenerateSNResponse;

      return result;
    }
  }

  /**
   * Get all records in QCS.
   *
   * @return: All records in QCS.
   */
  async getAllRecords(): Promise<QCSType.QCSGetAllRecordResponse> {
    const url = this.accessPrefix + "/sn/get-all";
    const headers = {
      "X-Access-Token": this.accessToken,
      "X-Runtime-Code": this.runtimeCode,
    };

    const res = await fetch(url, {
      method: "GET",
      headers: headers,
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        data: Array<QCSType.QCSRecord>;
      };

      const result = {
        data: data["data"],
      } as QCSType.QCSGetAllRecordResponse;

      return result;
    }
  }

  /**
   * Get available serial numbers.
   *
   * @return: Available serial numbers.
   */
  async getAvailableSN(): Promise<QCSType.QCSAvailableSNResponse> {
    const url = this.accessPrefix + "/sn/get-available";
    const headers = {
      "X-Access-Token": this.accessToken,
      "X-Runtime-Code": this.runtimeCode,
    };

    const res = await fetch(url, {
      method: "GET",
      headers: headers,
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        data: Array<string>;
      };

      const result = {
        data: data["data"],
      } as QCSType.QCSAvailableSNResponse;

      return result;
    }
  }

  /**
   * Update note of a serial number.
   *
   * @param sn: The serial number to update.
   * @param note: The note to update.
   * @returns: Message and the note were updated if success, error if failed.
   */
  async updateSNNote(
    sn: string,
    note: string
  ): Promise<QCSType.QCSUpdateSNNoteResponse> {
    const url = this.accessPrefix + "/sn/update";
    const headers = {
      "X-Access-Token": this.accessToken,
      "X-Runtime-Code": this.runtimeCode,
    };

    const body = {
      serial_number: sn,
      note: note,
    };

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        msg: string;
        note: string;
      };

      const result = {
        msg: data["msg"],
        note: data["note"],
      } as QCSType.QCSUpdateSNNoteResponse;

      return result;
    }
  }
}

export class QCSClient {
  accessPrefix: string;
  accessToken: string;

  constructor(qcsClientConfig: QCSType.QCSClientConfig) {
    const c = qcsClientConfig;
    this.accessPrefix = `${c.host}:${c.port}${c.apiPath}`;
    this.accessToken = c.accessToken ? c.accessToken : "";

    if (c.tls) {
      this.accessPrefix = "https://" + this.accessPrefix;
    } else {
      this.accessPrefix = "http://" + this.accessPrefix;
    }
  }

  /**
   * Apply a certificate for a device.
   *
   * @param sn: The serial number of the device.
   * @param board_producer: The producer of the board.
   * @param board_name: The name of the board.
   * @param mac_address: Physical ethernet mac address.
   *
   * @returns: The key and signature of the certificate.
   */
  async applyCert(
    sn: string,
    board_producer: string,
    board_name: string,
    mac_address: string
  ): Promise<QCSType.QCSApplyCertResponse> {
    const url = this.accessPrefix + "/apply/cert";
    const headers = {
      "X-Access-Token": this.accessToken,
    };

    const body = {
      serial_number: sn,
      board_producer: board_producer,
      board_name: board_name,
      mac_address: mac_address,
    };

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        key: string;
        signature: string;
      };

      const result = {
        key: data["key"],
        signature: data["signature"],
      } as QCSType.QCSApplyCertResponse;

      return result;
    }
  }

  /**
   * Apply a temporary permit(With time limit certificate) for a device.
   *
   * @param board_producer: The producer of the board.
   * @param board_name: The name of the board.
   * @param mac_address: Physical ethernet mac address.
   *
   * @returns: Remaining time and status of the temporary permit.
   */
  async applyTempPermit(
    board_producer: string,
    board_name: string,
    mac_address: string
  ): Promise<QCSType.QCSApplyTempPermitResponse> {
    const url = this.accessPrefix + "/apply/temp-permit";
    const headers = {
      "X-Access-Token": this.accessToken,
    };

    const body = {
      board_producer: board_producer,
      board_name: board_name,
      mac_address: mac_address,
    };

    const res = await fetch(url, {
      method: "POST",
      headers: headers,
      body: JSON.stringify(body),
    });

    if (res.status != 200) {
      const errorObj = (await res.json()) as { error: string };
      throw new Error("QCS::Error:" + errorObj["error"]);
    } else {
      const data = (await res.json()) as {
        remaining_time: number;
        status: string;
      };

      const result = {
        remainingTime: data["remaining_time"],
        status: data["status"],
      } as QCSType.QCSApplyTempPermitResponse;

      return result;
    }
  }
}
