# Siren Order Services(System)
**스타벅스의 사이렌 오더 서비스(기능)들을 보고 클론 코딩을 하는 프로젝트입니다 ☕️**

## ToDo
- [X] DataBase 설정
    - [X] 구조체를 이용한 테이블 생성함
- [X] `/register`: 기본적인 회원가입
    - [X] UUID를 통한 회원 식별을 할 수 있도록 함.
    - [X] Password 암호화를 함.
    - [X] 회원의 생일을 입력 받을 수 있도록 함.

## Functions
### `/register`
#### Request
```json
{
    "name": "HyunSang Park",
    "nickname": "박현상",
    "birthday": "2006-01-02T00:00:00Z",
    "email": "helloworld@helloworld.com",
    "password": "helloworld!"
}
```

#### Response
```json
{
    "exp": 1639209153,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMDkxNTMsInVzZXJfdXVpZCI6IjFiYzA0NTNlLTVhNTMtMTFlYy05NWE3LWFjZGU0ODAwMTEyMiJ9.zoRUInDcKQrHjhZ0twR6NpPzrLcRa1h4iowhl6hKBLI",
    "user": {
        "uuid": "1bc0453e-5a53-11ec-95a7-acde48001122",
        "name": "HyunSang Park",
        "nickname": "박현상",
        "birthday": "2004-06-25T00:00:00Z",
        "email": "parkhyunsang@kakao.com",
        "password": "$2a$10$bWq76y30poYnf399DsFac.1CWVXsr0cgT.17GQy1B.rt4LuYQ519y",
        "created_at": "2021-12-11T16:22:33.601364+09:00",
        "updated_at": "2021-12-11T16:22:33.608+09:00"
    },
    "uuid": "1bc0453e-5a53-11ec-95a7-acde48001122"
}
```
