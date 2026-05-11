# ข้อกำหนด: sdd-https-vpn-ciphersuite-th (ชุดรหัสลับสำหรับประเทศไทย)

> เวอร์ชัน: 1.0  
> สถานะ: ร่าง (DRAFT)  
> อัปเดตล่าสุด: 2026-05-11

## ปัญหาที่ต้องการแก้ไข (Problem Statement)

เพื่อเสริมสร้างความปลอดภัยของ HTTPS VPN ให้รองรับมาตรฐานการเข้ารหัสลับสมัยใหม่ของประเทศไทย และเตรียมพร้อมสำหรับการเปลี่ยนแปลงสู่ยุค Quantum-Resistant Cryptography (QRC) ตามคำแนะนำของ NIST (FIPS 203, 204, 205) เพื่อให้มั่นใจว่าการสื่อสารจะยังคงปลอดภัยแม้ในยุคคอมพิวเตอร์ควอนตัม

## เรื่องราวของผู้ใช้ (User Stories)

### เรื่องราวหลัก (Primary)

**ในฐานะ** ผู้ดูแลระบบความปลอดภัย (Security Administrator)  
**ฉันต้องการ** เพิ่มชุดรหัสลับ (Ciphersuites) ที่รองรับอัลกอริทึม Post-Quantum Cryptography (PQC) เข้าไปใน HTTPS VPN  
**เพื่อให้** ข้อมูลการสื่อสารได้รับการปกป้องจากภัยคุกคามในอนาคตและเป็นไปตามมาตรฐานสากลและระดับประเทศ

## เกณฑ์การยอมรับ (Acceptance Criteria)

### สิ่งที่ต้องมี (Must Have)

1. **Given** ระบบ HTTPS VPN  
   **When** มีการตั้งค่าการเชื่อมต่อ  
   **Then** ต้องสามารถเลือกใช้งาน ML-KEM (FIPS 203) สำหรับการแลกเปลี่ยนกุญแจ (Key Encapsulation Mechanism) ได้

2. **Given** การลงนามดิจิทัล (Digital Signature)  
   **When** มีการตรวจสอบความถูกต้อง  
   **Then** ต้องรองรับ ML-DSA (FIPS 204) และ SLH-DSA (FIPS 205)

3. **Given** ความเสี่ยงจากช่องโหว่ของอัลกอริทึมหลัก  
   **When** ML-KEM มีปัญหา  
   **Then** ต้องมี HQC (Hamming Quasi-Cyclic) เป็นระบบสำรอง (Backup KEM)

### สิ่งที่ควรมี (Should Have)

- การจัดลำดับความสำคัญของชุดรหัสลับ (Cipher Suite Prioritization) สำหรับการใช้งานในประเทศไทย
- ตัวอย่างการกำหนดค่า (Configuration Examples) ในภาษาไทย

### สิ่งที่ยังไม่ทำในรอบนี้ (Won't Have)

- การรวมระบบเข้ากับ Hardware Security Module (HSM) ที่ยังไม่รองรับ PQC
- การรับรองมาตรฐานจากสำนักงานพัฒนาธุรกรรมทางอิเล็กทรอนิกส์ (ETDA) อย่างเป็นทางการในเฟสแรก

## ข้อจำกัด (Constraints)

- **Technical**: ต้องทำงานร่วมกับโครงสร้างพื้นฐานเดิมของ HTTPS VPN ได้
- **Performance**: อัลกอริทึม PQC อาจมีขนาดกุญแจที่ใหญ่ขึ้นและใช้เวลาประมวลผลนานขึ้น ต้องมีการทดสอบประสิทธิภาพ
- **Platform**: รองรับบนระบบปฏิบัติการที่กำหนด (Darwin, Linux, etc.)
- **Dependencies**: ขึ้นอยู่กับไลบรารีคริปโตกราฟีที่รองรับมาตรฐาน FIPS ใหม่

## คำถามที่ยังไม่มีคำตอบ (Open Questions)

- [ ] ไลบรารี Go มาตรฐานหรือไลบรารีภายนอกใดที่เสถียรที่สุดสำหรับ ML-KEM/ML-DSA ในปัจจุบัน?
- [ ] ความเข้ากันได้กับไคลเอนต์ VPN รุ่นเก่าเมื่อเปิดใช้งาน PQC?

## เอกสารอ้างอิง (References)

- FIPS 203: Module-Lattice-Based Key-Encapsulation Mechanism Standard
- FIPS 204: Module-Lattice-Based Digital Signature Standard
- FIPS 205: Stateless Hash-Based Digital Signature Standard
- NIST PQC Project: HQC (Hamming Quasi-Cyclic)

---

## การอนุมัติ (Approval)

- [ ] ตรวจสอบโดย: [ชื่อ]
- [ ] อนุมัติเมื่อ: [วันที่]
- [ ] หมายเหตุ: [ข้อเสนอแนะเพิ่มเติม]
