"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.QCSClient = exports.QCSAdmin = void 0;
class QCSAdmin {
    constructor(qcsAdminConfig) {
        const c = qcsAdminConfig;
        this.accessPrefix = `${c.host}:${c.port}${c.apiPath}`;
        this.accessToken = c.accessToken ? c.accessToken : "";
        this.runtimeCode = c.runtimeCode ? c.runtimeCode : "";
        if (c.tls) {
            this.accessPrefix = "https://" + this.accessPrefix;
        }
        else {
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
    createSN(sn, reason = "none") {
        return __awaiter(this, void 0, void 0, function* () {
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
            const res = yield fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body),
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    msg: data["msg"],
                    serial_number: data["serial_number"],
                };
                return result;
            }
        });
    }
    /**
     * Generate serial number(s) randomly.
     *
     * @param count: number of serial numbers to generate.
     * @param reason: reason for creating this serial number.
     *
     * @return: Message and the sn(s) were created if success, error if failed.
     */
    generateSN(count, reason = "none") {
        return __awaiter(this, void 0, void 0, function* () {
            const url = this.accessPrefix + "/sn/generate";
            const headers = {
                "X-Access-Token": this.accessToken,
                "X-Runtime-Code": this.runtimeCode,
            };
            const body = {
                count: count,
                reason: reason,
            };
            const res = yield fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body),
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    msg: data["msg"],
                    serial_numbers: data["serial_numbers"],
                };
                return result;
            }
        });
    }
    /**
     * Get all records in QCS.
     *
     * @return: All records in QCS.
     */
    getAllRecords() {
        return __awaiter(this, void 0, void 0, function* () {
            const url = this.accessPrefix + "/sn/get-all";
            const headers = {
                "X-Access-Token": this.accessToken,
                "X-Runtime-Code": this.runtimeCode,
            };
            const res = yield fetch(url, {
                method: "GET",
                headers: headers,
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    data: data["data"],
                };
                return result;
            }
        });
    }
    /**
     * Get available serial numbers.
     *
     * @return: Available serial numbers.
     */
    getAvailableSN() {
        return __awaiter(this, void 0, void 0, function* () {
            const url = this.accessPrefix + "/sn/get-available";
            const headers = {
                "X-Access-Token": this.accessToken,
                "X-Runtime-Code": this.runtimeCode,
            };
            const res = yield fetch(url, {
                method: "GET",
                headers: headers,
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    data: data["data"],
                };
                return result;
            }
        });
    }
    /**
     * Update note of a serial number.
     *
     * @param sn: The serial number to update.
     * @param note: The note to update.
     * @returns: Message and the note were updated if success, error if failed.
     */
    updateSNNote(sn, note) {
        return __awaiter(this, void 0, void 0, function* () {
            const url = this.accessPrefix + "/sn/update";
            const headers = {
                "X-Access-Token": this.accessToken,
                "X-Runtime-Code": this.runtimeCode,
            };
            const body = {
                serial_number: sn,
                note: note,
            };
            const res = yield fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body),
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    msg: data["msg"],
                    note: data["note"],
                };
                return result;
            }
        });
    }
}
exports.QCSAdmin = QCSAdmin;
class QCSClient {
    constructor(qcsClientConfig) {
        const c = qcsClientConfig;
        this.accessPrefix = `${c.host}:${c.port}${c.apiPath}`;
        this.accessToken = c.accessToken ? c.accessToken : "";
        if (c.tls) {
            this.accessPrefix = "https://" + this.accessPrefix;
        }
        else {
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
    applyCert(sn, board_producer, board_name, mac_address) {
        return __awaiter(this, void 0, void 0, function* () {
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
            const res = yield fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body),
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    key: data["key"],
                    signature: data["signature"],
                };
                return result;
            }
        });
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
    applyTempPermit(board_producer, board_name, mac_address) {
        return __awaiter(this, void 0, void 0, function* () {
            const url = this.accessPrefix + "/apply/temp-permit";
            const headers = {
                "X-Access-Token": this.accessToken,
            };
            const body = {
                board_producer: board_producer,
                board_name: board_name,
                mac_address: mac_address,
            };
            const res = yield fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify(body),
            });
            if (res.status != 200) {
                const errorObj = (yield res.json());
                throw new Error("QCS::Error:" + errorObj["error"]);
            }
            else {
                const data = (yield res.json());
                const result = {
                    remaining_time: data["remaining_time"],
                    status: data["status"],
                };
                return result;
            }
        });
    }
}
exports.QCSClient = QCSClient;
