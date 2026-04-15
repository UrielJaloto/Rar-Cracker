# RAR5 Binary Architecture (Reverse Engineering)

### General archive structure
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
    * `0x0001`: Extra area is present in the end of header 
    * `0x0002`: Data area is present in the end of header 
    * `0x0004`: Blocks with unknown type and this flag must be skipped when updating an archive 
    * `0x0008`: Data area is continuing from previous volume 
    * `0x0010`: Data area is continuing in next volume 
    * `0x0020`: Block depends on preceding file block 
    * `0x0040`: Preserve a child block if host block is modified 
  * **Extra Area Size:** VarInt, optional field (present if flag bit `0x0001` is set)
  * **Data Area Size:** VarInt, optional field (present if flag bit `0x0002` is set)
  * **[Fields specific for current block type]**
  * **Extra Area:** Área opcional para campos adicionais e metadados
    * **Size:** VarInt (Tamanho da area, iniciando a partir do *Type*)
    * **Type:** VarInt (Identificador exclusivo do tipo de registro)
    * **Data:** Variável (Conteúdo específico do registro. Pode ser omitido se o registro necessitar apenas de Size e Type)
  * **Data Area:** Área opcional para armazenamento do payload/arquivos comprimidos (não contada no Header CRC e Header Size)

---
### Encryption Meta Data 
  * **Encryption Version:** VarInt, usually 0 for AES-256
  * **Encryption Flags:** VarInt, bit 1 signals Password Check presence
  * **KDF Count:** 1 byte, base-2 log of PBKDF2 iterations
  * **Salt:** 16 bytes, random cryptographic data to protect the password
  * **IV (present if in extra area):** 16 bytes, random cryptographic data to protect the data
  * **Password Check:** 12 bytes, optional quick validation hash

---
### Extra Area Types
* **Main Archive Header:**
  * `0x01`: **Locator:** Offsets para blocos de serviço
  * `0x02`: **Metadata:** Nome original do arquivo e timestamp de criação

* **File/Service Header:**
  * `0x01`: **File Encryption:** Estrutura de criptografia por arquivo. Contém IV.
  * `0x02`: **File Hash:** Armazena hashes criptográficos mais complexos que o CRC32 do header
  * `0x03`: **File Time:** Timestamps de criação, modificação e acesso em alta precisão
  * `0x04`: **File Version:** Número da versão do arquivo
  * `0x05`: **Redirection:** Informações de symlinks, hard links e cópias de arquivos
  * `0x06`: **Unix Owner:** Nomes e IDs de usuário e grupo para sistemas Unix
  * `0x07`: **Service Data:** Parâmetros adicionais exclusivos para Service Headers