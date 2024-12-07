#@get_sub_districts@
POST https://uniqlo-customerapi.brainergy.io/reqft/get_subdistrict
{
    "district_id": "{{district_id}}"
}
HTTP 200
[Asserts]
header "Content-Type" contains  "application/json"
jsonpath  "$.status" == 200
jsonpath  "$.msg" == "SUCCESS"
jsonpath  "$.data" count > 0
[Captures]
subdistrict_id : jsonpath "$.data[0].subdistrict_id"