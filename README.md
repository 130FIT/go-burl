
# Burl 

**Burl** เป็นเครื่องมือสำหรับการทดสอบ API ที่รองรับทั้ง JSON และ XML โดยสามารถทำการทดสอบได้ทั้ง Unit Test และ Integration Test

## คุณสมบัติ

- **รองรับ JSON และ XML**: ทดสอบ API ที่ตอบสนองด้วย JSON หรือ XML
- **Unit Testing**: ทดสอบฟังก์ชันหรือโมดูลเฉพาะ
- **Integration Testing**: ทดสอบการทำงานร่วมกันของหลายๆ ส่วนในระบบ
- **รองรับการกำหนดค่าโดยใช้ไฟล์ JSON**: ใช้ไฟล์ JSON สำหรับการตั้งค่า API และการทดสอบ

## การใช้งาน
**การทดสอบ Unit**
1. สร้างไฟล์ทดสอบ เช่น 'test1.json'
```json
{
    "url": "https://pokeapi.co/api/v2/pokemon/pikachu",
    "method": "GET",
    "headers": {
        "Content-Type": "application/json"
    },
    "base_request":{
        "A":"A"
    },
    "type": "json",
    "captures":[
        {
            "pass_path": "cases[0].request.name",
            "capture_path": "abilities[0].ability.url"
        }
    ],
    "cases":[
        {
            "id": 1,
            "name": "get pikachu info should return 200",
            "description": "this case must success",
            "tags": ["success cases"],
            "status": 200,
            "request":{
                "A":"B"
            },
            "assert_response": {
                "name": "pikachu",
                "forms[0].name": "pikachu"
            }
        },
        {
            "id": 2,
            "name": "get pikachu info should return 400",
            "description": "this case must fail",
            "tags":["fail cases"],
            "status": 400,
            "request":{
                "A":"C"
            }
        }
    ]
}
```
2. ใช้คำสั่ง 

```console
burl test1.json
```

## Options
| **Option** |            **Description**            |
|:----------:|:-------------------------------------:|
| -flow      | to test multiple files in a sequence. |
| -runner    | Enable runner file to run test cases  |
| -help      | Show help                             |
| -version   | Show version                          |
