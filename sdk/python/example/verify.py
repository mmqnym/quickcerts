# This is an example for verifying the signature by the given public key.
# Used arguments:
# hash method: SHA3-512
# PSS salt length: rsa.PSSSaltLengthEqualsHash
# Message: 95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830
# Signature: upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g==

from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.serialization import load_pem_public_key
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives.asymmetric import utils
import base64


public_key_pem = b"""-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzTuY9ePxSX533aa54/aY
Qobqzz0/alc40C31/fYgYXLQVeMJ4vXBHKFhWOaf+ZBf2bQBLx2aIa2ODZcH4ZNF
UIbSZu9jmWN6kcSCw5IMPuDW2YF0b0MlxCemPgCPdIioBa/qsgmy4/s6LpZ2JtUG
7+KBOJIBxuzt8k2XtfRK7k8HBL5v3pQI6IqgooN6cq/M9IOWges1RwLTsMcUbISm
pSOGIC57XmreGiOQik3IlWLYaDbo5nOhzhGtnz6FlAOscW3guYuMBiPjYnTERXNz
1rwx1dHM+t+K2/7poB477RoBEHeLYkEF2JkxVZAXdAg+5PKkMj+Cd/U867t83mDG
OQIDAQAB
-----END PUBLIC KEY-----
"""

public_key = load_pem_public_key(public_key_pem)

message = "95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830"
message_bytes = message.encode()
signature_base64 = "upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g=="

digest = hashes.Hash(hashes.SHA3_512())
digest.update(message_bytes)
message_hash = digest.finalize()

signature = base64.b64decode(signature_base64)
verification_result = ""

try:
    public_key.verify(
        signature,
        message_hash,
        padding.PSS(
            mgf=padding.MGF1(hashes.SHA3_512()),
            salt_length=len(message_hash)
        ),
        utils.Prehashed(hashes.SHA3_512())
    )
    verification_result = "PASS"
except InvalidSignature:
    verification_result = "FAIL"

print(verification_result)