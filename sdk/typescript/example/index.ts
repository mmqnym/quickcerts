// This is an example for verifying the signature by the given public key.
// Used arguments:
// hash method: SHA3-512
// PSS salt length: rsa.PSSSaltLengthEqualsHash
// Message: 95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830
// Signature: upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g==

import * as crypto from "crypto";

const publicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzTuY9ePxSX533aa54/aY
Qobqzz0/alc40C31/fYgYXLQVeMJ4vXBHKFhWOaf+ZBf2bQBLx2aIa2ODZcH4ZNF
UIbSZu9jmWN6kcSCw5IMPuDW2YF0b0MlxCemPgCPdIioBa/qsgmy4/s6LpZ2JtUG
7+KBOJIBxuzt8k2XtfRK7k8HBL5v3pQI6IqgooN6cq/M9IOWges1RwLTsMcUbISm
pSOGIC57XmreGiOQik3IlWLYaDbo5nOhzhGtnz6FlAOscW3guYuMBiPjYnTERXNz
1rwx1dHM+t+K2/7poB477RoBEHeLYkEF2JkxVZAXdAg+5PKkMj+Cd/U867t83mDG
OQIDAQAB
-----END PUBLIC KEY-----`;

const message =
  "95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830";
const signatureBase64 =
  "upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g==";

const messageBuffer = Buffer.from(message, "utf8");
const signature = Buffer.from(signatureBase64, "base64");
const verifier = crypto.createVerify("sha3-512");
verifier.update(messageBuffer);

const isVerified = verifier.verify(
  {
    key: publicKeyPem,
    padding: crypto.constants.RSA_PKCS1_PSS_PADDING,
    saltLength: crypto.constants.RSA_PSS_SALTLEN_DIGEST,
  },
  signature
);

console.log(isVerified ? "PASS" : "FAIL");
