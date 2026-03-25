# RAR5 Binary Architecture (Reverse Engineering)

* **Signature** (8 bytes: `52 61 72 21 1A 07 01 00`)
* **Variable Integers** (7 bits for data, 8th bit for continuation)
* **Base Block Header** (Standard header for all blocks)
  * **CRC32:** 4 bytes for header integrity
  * **Header Size:** VarInt for jumping/skipping blocks
  * **Header Type:** VarInt identifier
    * `0x01`: Main Archive Header
    * `0x02`: File Header
    * `0x03`: Service Header
    * `0x04`: Encryption Header
    * `0x05`: End of Archive Header
  * **Header Flags:** VarInt for general rules
    * **Extra Area Size:** VarInt, optional field (present if flag bit `0x0001` is set)
    * **Data Area Size:** VarInt, optional field (present if flag bit `0x0002` is set)
* **Encryption Block Header** (Type `0x04`)
  * **Encryption Version:** VarInt, usually 0 for AES-256
  * **Encryption Flags:** VarInt, bit 1 signals Password Check presence
  * **KDF Count:** 1 byte, base-2 log of PBKDF2 iterations
  * **Salt:** 16 bytes of random cryptographic data
  * **Password Check:** 8 bytes, optional quick validation hash