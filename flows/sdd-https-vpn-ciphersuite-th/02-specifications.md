# รายละเอียดทางเทคนิค: sdd-https-vpn-ciphersuite-th

> เวอร์ชัน: 1.1  
> สถานะ: ร่าง (DRAFT)  
> อัปเดตล่าสุด: 2026-05-11  
> ข้อกำหนด: [01-requirements.md](01-requirements.md)

## ภาพรวม (Overview)

สถาปัตยกรรมความปลอดภัยนี้ใช้แนวทาง "Hybrid Cryptography" เพื่อรวมความเชื่อมั่นของอัลกอริทึมคลาสสิก (ECC) เข้ากับความทนทานต่อควอนตัมของ PQC โดยแบ่งตามระดับความเสี่ยงและความต้องการด้านประสิทธิภาพ

## สถาปัตยกรรม (Architecture)

### 1. การจัดการกุญแจ (Key Management - FIPS 203)

เราใช้ **ML-KEM** เป็นกลไกหลักในการแลกเปลี่ยนกุญแจ (Key Encapsulation):

| โปรไฟล์ (Profile) | อัลกอริทึมไฮบริด | การใช้งาน | เหตุผลทางเทคนิค |
|-------------------|------------------|-----------|-----------------|
| **Balanced (Default)** | X25519 + ML-KEM-768 | ทราฟฟิก VPN ทั่วไป, Session Setup | จุดสมดุลระหว่าง Security Strength (Level 3) และ Performance |
| **High-Assurance** | P-384 + ML-KEM-1024 | Admin Channels, Enrollment, Root Setup | ความปลอดภัยสูงสุด (Level 5) สำหรับโครงสร้างพื้นฐานวิกฤต |

### 2. ระบบการลงนามดิจิทัล (Digital Signatures - FIPS 204 & 205)

แบ่งตาม Use Cases เพื่อจัดการปัญหาเรื่องขนาด (Size Overhead):

- **ML-DSA-65 (FIPS 204):** ใช้ในโหมดไฮบริด (เช่น ECDSA+ML-DSA) สำหรับ:
    - End-entity Certificates
    - Control API messages
    - Release manifests / Config bundles
- **SLH-DSA (FIPS 205):** ใช้สำหรับโครงสร้างพื้นฐานที่มีความคงทนสูง (High Persistence):
    - Offline Root CAs
    - Firmware Signing (Fallback)
    - Disaster Recovery Keys
    - *หมายเหตุ: SLH-DSA มีขนาดลายเซ็นใหญ่กว่า แต่ให้ความเชื่อมั่นสูงมากเนื่องจากเป็น Hash-based*

### 3. ระบบสำรองและอนาคต (Backup & Roadmap)

- **HQC (Hamming Quasi-Cyclic):** เตรียมพร้อมเป็น Backup KEM ในกรณีที่ Lattice-based (ML-KEM) ถูกโจมตีทางคณิตศาสตร์
- **Falcon:** ติดตามความคืบหน้าของ NIST เพื่อรองรับในฐานะ Signature scheme ที่มีขนาดกะทัดรัด (Compact Signature)

## รายละเอียดส่วนประกอบ (Component Details)

### Hybrid TLS Extension

ต้องมีการเพิ่มการรองรับ Named Groups ใหม่ใน TLS 1.3:
- `X25519_MLKEM768` (0xEE42 - สมมติ)
- `P384_MLKEM1024` (0xEE43 - สมมติ)

### ช่องว่างที่เหลืออยู่ (Remaining Gaps)

1. **Standardization (IANA Codepoints):** รหัสสำหรับ Hybrid Groups (เช่น X25519 + ML-KEM-768) ยังไม่ถูกกำหนดเป็นมาตรฐานสากลที่แน่นอนใน IANA ปัจจุบัน เราอาจต้องใช้ช่วง "Private Use" (0xEE00-0xEEEE) ไปก่อนในช่วงแรก
2. **MTU & Fragmentation:** เนื่องจากกุญแจและลายเซ็น PQC มีขนาดใหญ่กว่า ECC มาก (เช่น SLH-DSA อาจมีขนาดถึง 30KB+) ซึ่งเกินขนาดมาตรฐานของ MTU (1500 bytes) ของอินเทอร์เน็ตทั่วไป หากส่งผ่าน UDP (ในกรณี DTLS) อาจเกิดปัญหา Packet Loss จากการทำ IP Fragmentation
3. **Hardware Acceleration:** ปัจจุบัน CPU ส่วนใหญ่ยังไม่มีชุดคำสั่งพิเศษ (เช่น AES-NI) สำหรับ Lattice-based Cryptography ทำให้การประมวลผลใช้ CPU Cycles สูงกว่าอัลกอริทึมคลาสสิก
4. **Library Certification:** แม้อัลกอริทึมจะถูกรับรองโดย NIST (FIPS) แต่ไลบรารีในภาษา Go ที่ผ่านการรับรองความถูกต้อง (Formal Verification) ยังอยู่ในช่วงเริ่มต้นพัฒนา
5. **Certificate Infrastructure:** การสร้างใบรับรองแบบ Hybrid (ที่มีทั้ง ECC และ ML-DSA signature) ยังไม่รองรับโดย Public CAs ส่วนใหญ่ ทำให้เราต้องบริหารจัดการ Private PKI ของเราเองทั้งหมดในเฟสแรก

## แผนผังการเชื่อมต่อ (Integration Diagram)

```
[Client] 
   |
   |-- (ClientHello: Support Hybrid Groups) --> [Server]
   |
   |-- (ServerHello: Selected X25519+MLKEM768) --|
   |
   |-- (Encrypted Extensions: Hybrid Certs) ----|
```

## กลยุทธ์การทดสอบ (Testing Strategy)

- **Interoperability:** ทดสอบการสื่อสารระหว่างเวอร์ชันที่เปิด PQC และเวอร์ชัน Classical เท่านั้น
- **Performance Benchmarking:** วัดความหน่วง (Latency) ของ ML-KEM-768 เทียบกับ ML-KEM-1024 ในสภาพเครือข่ายที่แตกต่างกัน
- **Recovery:** ทดสอบการสลับไปใช้ HQC เมื่อจำลองสถานการณ์ความล้มเหลวของ ML-KEM

---

## การอนุมัติ (Approval)

- [ ] ตรวจสอบโดย: [ชื่อ]
- [ ] อนุมัติเมื่อ: [วันที่]
- [ ] หมายเหตุ: [ข้อเสนอแนะเพิ่มเติม]
