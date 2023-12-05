import { QCSAdmin, QCSClient } from "../qcs";
import { QCSType } from "../qcs_type";

const qcsAdminConfig = {
  host: "127.0.0.1",
  port: 33333,
  apiPath: "/api/v1",
  tls: false,
  accessToken:
    "0b09b6dc41f61813346ba76322d19e07a0b71ba939a1bf90211dfff40f552ed0",
  runtimeCode: "",
} as QCSType.QCSAdminConfig;

const qcsAdmin = new QCSAdmin(qcsAdminConfig);

const qcsClientConfig = {
  host: "127.0.0.1",
  port: 33333,
  apiPath: "/api/v1",
  tls: false,
  accessToken: "QcsTestToken********************************",
} as QCSType.QCSClientConfig;

const qcsClient = new QCSClient(qcsClientConfig);

(async () => {
  // Create a serial number.
  try {
    let res = await qcsAdmin.createSN("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX");
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  //   // Get a serial number.
  try {
    let res = await qcsAdmin.generateSN(2);
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  // Get all records in QCS.
  try {
    let res = await qcsAdmin.getAllRecords();
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  // Get available serial numbers.
  try {
    let res = await qcsAdmin.getAvailableSN();
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  // Update note of a serial number.
  try {
    let res = await qcsAdmin.updateSNNote(
      "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
      "test"
    );
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  // Apply a certificate.
  try {
    let res = await qcsClient.applyCert(
      "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
      "ASUSTeK Computer Inc.",
      "ROG STRIX Z790-A GAMING WIFI",
      "XXXXXXXXXXXX"
    );
    console.log(res);
  } catch (e) {
    console.log(e);
  }

  // Apply a temporary permit(with time limit certificate).
  try {
    let res = await qcsClient.applyTempPermit(
      "ASUSTeK Computer Inc.",
      "ROG STRIX Z790-A GAMING WIFI",
      "XXXXXXXXXXXX"
    );
    console.log(res);
  } catch (e) {
    console.log(e);
  }
})();
