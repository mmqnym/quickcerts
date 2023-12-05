export namespace QCSType {
  export type QCSAdminConfig = {
    host: string;
    port: number;
    apiPath: string;
    tls: boolean;
    accessToken: string;
    runtimeCode: string;
  };

  export type QCSCreateResponse = {
    msg: string;
    serial_number: string;
  };

  export type QCSGenerateSNResponse = {
    msg: string;
    serial_numbers: Array<string>;
  };

  export type QCSRecord = {
    sn: string;
    key: string;
    note: string;
  };

  export type QCSGetAllRecordResponse = {
    data: Array<QCSRecord>;
  };

  export type QCSAvailableSNResponse = {
    data: Array<string>;
  };

  export type QCSUpdateSNNoteResponse = {
    msg: string;
    note: string;
  };

  export type QCSClientConfig = {
    host: string;
    port: number;
    apiPath: string;
    tls: boolean;
    accessToken: string;
  };

  export type QCSApplyCertResponse = {
    key: string;
    signature: string;
  };

  export type QCSApplyTempPermitResponse = {
    remainingTime: number;
    status: string;
  };
}
