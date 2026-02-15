Client Server

---

1. TCP Handshake

---

SYN ---------------------------------------------->
<---------------------------------------------- SYN-ACK
ACK ---------------------------------------------->

(TCP connection established)

2. TLS Handshake Begins

---

ClientHello

- Supported TLS versions
- Cipher suites
- Client Random
- Client Ephemeral Public Key (A = g^a)
  --------------------------------------------------------->

                                                   ServerHello
                                                     - Chosen cipher
                                                     - Server Random
                                                     - Server Ephemeral Public Key (B = g^b)
                                                   <---------------------------------------------------------

                                                   Certificate
                                                     - Server public key
                                                     - Signed by CA
                                                   <---------------------------------------------------------

                                                   CertificateVerify
                                                     - Signature over handshake
                                                   <---------------------------------------------------------

                                                   Finished
                                                     - Handshake integrity check
                                                   <---------------------------------------------------------

3. Client Verifies

---

✔ Certificate is valid
✔ Domain matches
✔ Signed by trusted CA
✔ Signature is correct

4. Shared Secret Creation (Diffie-Hellman)

---

Client computes:
SharedSecret = B^a

Server computes:
SharedSecret = A^b

Both now have:
Same Shared Secret = g^(ab)

(Private keys never transmitted)

5. Key Derivation

---

SharedSecret
↓
HKDF (Key Derivation Function)
↓
Session Keys:

- Encryption Key
- Integrity Key
- IV

6. Secure Communication Starts

---

Encrypted HTTP Request
---------------------------------------------->

                                   Encrypted HTTP Response

<----------------------------------------------

Using:

- Symmetric encryption (AES-GCM / ChaCha20)
- Authenticated encryption (prevents tampering)
