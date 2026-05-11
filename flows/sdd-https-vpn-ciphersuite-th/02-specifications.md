# รายละเอียดทางเทคนิค: sdd-https-vpn-ciphersuite-th

> เวอร์ชัน: 1.0  
> สถานะ: ร่าง (DRAFT)  
> อัปเดตล่าสุด: 2026-05-11  
> ข้อกำหนด: [01-requirements.md](01-requirements.md)

## ภาพรวม (Overview)

การพัฒนานี้จะเพิ่มแพ็กเกจ `crypto/th` เพื่อเป็นผู้ให้บริการคริปโตกราฟี (Cryptographic Provider) สำหรับประเทศไทย โดยเน้นการใช้งานอัลกอริทึม Post-Quantum Cryptography (PQC) ตามมาตรฐาน NIST FIPS 203, 204 และ 205 พร้อมระบบสำรอง HQC เพื่อความมั่นคงปลอดภัยสูงสุด

## ระบบที่ได้รับผลกระทบ (Affected Systems)

| ระบบ (System) | ผลกระทบ (Impact) | หมายเหตุ (Notes) |
|---------------|------------------|------------------|
| `crypto/th` | สร้างใหม่ (Create) | แพ็กเกจหลักสำหรับผู้ให้บริการไทย |
| `crypto/th/tls` | สร้างใหม่ (Create) | นิยาม Cipher Suites และค่าคงที่สำหรับ TLS |
| `crypto/provider.go` | แก้ไข (Modify) | ลงทะเบียนโปรไวเดอร์ใหม่ |

## สถาปัตยกรรม (Architecture)

### ไดอะแกรมส่วนประกอบ (Component Diagram)

```
[Application] -> [crypto.Provider (Interface)]
                        ^
                        |
                [crypto/th.Provider]
                /        |         \
        [ML-KEM]    [ML-DSA]    [HQC (Backup)]
        (FIPS 203)  (FIPS 204)  (NIST R4)
```

### การไหลของข้อมูล (Data Flow)

1. เมื่อเริ่มต้นระบบ `crypto/th` จะลงทะเบียนตัวเองกับ `crypto.Registry`
2. เมื่อมีการสร้าง TLS Config ระบบจะเรียก `Provider("th").ConfigureTLS()`
3. การตกลงกุญแจ (Handshake) จะใช้โหมดไฮบริด (Hybrid Mode) ระหว่าง X25519 และ ML-KEM-768

## อินเทอร์เฟซ (Interfaces)

### อินเทอร์เฟซใหม่ (New Interfaces)

ในแพ็กเกจ `crypto/th`:
```go
type Provider struct{}
func (p *Provider) Name() string
func (p *Provider) ConfigureTLS(cfg *tls.Config) error
func (p *Provider) SupportedCipherSuites() []uint16
```

## รูปแบบข้อมูล (Data Models)

### ค่าคงที่ใหม่ (New Constants)

ใน `crypto/th/tls/cipher_suites.go`:
```go
const (
    // รหัสสำหรับชุดรหัสลับไทย (สมมติเพื่อการทดสอบ)
    TLS_TH_PQC_WITH_AES_256_GCM_SHA384 uint16 = 0xEA01 
)
```

## รายละเอียดพฤติกรรม (Behavior Specifications)

### กรณีปกติ (Happy Path)

1. ผู้ใช้เลือกโปรไวเดอร์ "th"
2. ระบบตั้งค่า `CurvePreferences` ให้มี `X25519MLKEM768` เป็นอันดับแรก
3. ระบบใช้ `ML-DSA` สำหรับการลงนามในใบรับรอง (Certificates)

### กรณีขอบเขต (Edge Cases)

| กรณี (Case) | สาเหตุ (Trigger) | พฤติกรรมที่คาดหวัง (Expected Behavior) |
|-------------|-----------------|----------------------------------------|
| ML-KEM ล้มเหลว | พบช่องโหว่หรือข้อผิดพลาดในการประมวลผล | สลับไปใช้ HQC (Hamming Quasi-Cyclic) เป็น Backup KEM |
| ไคลเอนต์ไม่รองรับ PQC | ไคลเอนต์รุ่นเก่าเชื่อมต่อ | Fallback กลับไปใช้ X25519 (Classical) เพื่อความเข้ากันได้ |

## การจัดการข้อผิดพลาด (Error Handling)

- หากอัลกอริทึม PQC ไม่สามารถใช้งานได้ ระบบต้องแจ้งเตือนใน Log และใช้โหมด Classical เท่านั้น (Fail-safe to Classical)

## กลยุทธ์การทดสอบ (Testing Strategy)

### Unit Tests
- [ ] `crypto/th/provider_test.go`: ทดสอบการกำหนดค่า TLS
- [ ] ทดสอบการสลับใช้งานระหว่าง ML-KEM และ HQC

### Integration Tests
- [ ] ทดสอบการเชื่อมต่อ VPN จริงโดยใช้โปรไวเดอร์ "th"

---

## การอนุมัติ (Approval)

- [ ] ตรวจสอบโดย: [ชื่อ]
- [ ] อนุมัติเมื่อ: [วันที่]
- [ ] หมายเหตุ: [ข้อเสนอแนะเพิ่มเติม]
