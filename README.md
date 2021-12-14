# Siren Order Services(System)
**스타벅스의 사이렌 오더 서비스(기능)들을 보고 클론 코딩을 하는 프로젝트입니다 ☕️**

## ToDo
### User Account `/member`
- [X] DataBase 설정
    - [X] 구조체를 이용한 테이블 생성함
- [ ] Redis를 이용한 JWT 인증 서버
    - [ ] JWT 메타데이터 정의 
    - [ ] JWT 메타데이터 저장
- [X] `/register`: 기본적인 회원가입
    - [X] UUID를 통한 회원 식별을 할 수 있도록 함.
    - [X] Password 암호화를 함.
    - [X] 회원의 생일을 입력 받을 수 있도록 함.
- [X] `/login`: 기본적인 로그인
    - [X] 회원 정보 맞을 시 
        - [X] `JWT` Token 발행
    - [X] 회원 정보가 안 맞을 시
        - [X] 오류 출력
    - [X] 회원정보 GORM를 통해서 불러오기
- [ ] `/edit`: 회원의 기본적인 정보를 수정할 수 있도록 함.
    - [ ] 닉네임 변경
    - [ ] 전화번호 변경
### Order `/order/{id}`
- [ ] 푸드마다 고유ID를 통한 푸드 및 음료를 식별할 수 있도록 함.

## Payment `/payment`
- [ ] 멤버쉽 카드 생성
    - [ ] 생성 유저 식별 및 기록
    - [ ] 멤버쉽 카드의 상태
        - [ ] 분실
            - [ ] 삭제 or 보존
        - [ ] 충전
            - [ ] 충전 방식 결정
        - [ ] 잔액 확인

### Store `/store`
- [ ] 주문 시 결정할 수 있도록 함.
    - [ ] 드라이브 스루 사용 가능 매장 확인
    - [ ] 매장 운영 시간 확인
        - [ ] 매장이 열지 않았거나 닫은 경우 오류 출력

## Functions
### `/register`
#### Request
```json
{
    "name": "HyunSang Park",
    "nickname": "박현상",
    "birthday": "2006-01-02T00:00:00Z",
    "email": "parkhyunsang@kakao.com",
    "password": "helloworld!"
}
```

#### Response
```json
{
    "exp": "2021-12-12T17:02:00.486146+09:00",
    "isOk": true,
    "status": 200
}
```

### `/login`
#### Request
```json
{
    "email": "parkhyunsang@kakao.com",
    "password": "helloworld!"
}
```

#### Response

```json
{
    "exp": 1639236951,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMzY5NTEsInVzZXJfdXVpZCI6IjI1NWM2ZjU4LTVhNTUtMTFlYy1iZWEzLWFjZGU0ODAwMTEyMiJ9.ytPJKpOpRZ98w093k3FDZ1wfTR8ybrrxhJ84tmp8R0Y",
    "user_nickname": "박현상",
    "user_uuid": "255c6f58-5a55-11ec-bea3-acde48001122"
}
```

## 참고한 문서 및 글
- [GO 언어로 JWT 인증서버 만들기](https://covenant.tistory.com/203)
- [<go> fcm push 발송 (android, ios)](https://www.byfuls.com/programming/read?id=25)
- [WATCHA 푸시 서버 개선하기](https://medium.com/watcha/watcha-%ED%91%B8%EC%8B%9C-%EC%84%9C%EB%B2%84-%EA%B0%9C%EC%84%A0%ED%95%98%EA%B8%B0-56070b73c287)