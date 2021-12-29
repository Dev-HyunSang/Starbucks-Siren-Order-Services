# Siren Order Services(System)
**스타벅스의 사이렌 오더 서비스(기능)들을 보고 클론 코딩을 하는 프로젝트입니다 ☕️**

## ToDo
### User Account `/api/member`
- [X] DataBase 설정
    - [X] 구조체를 이용한 테이블 생성함
- [X] POST `/register`: 기본적인 회원가입
    - [X] UUID를 통한 회원 식별을 할 수 있도록 함.
    - [X] Password 암호화를 함.
    - [X] 회원의 생일을 입력 받을 수 있도록 함.
    - [X] 가입되어 있는 메일에 대한 중복 오류 처리
    - [X] 가입되어 있는 연락처에 대한 중복 오류 처리
    - [X] 사용되고 있는 닉네임에 대한 중복 오류 처리
- [X] POST `/login`: 기본적인 로그인
    - [X] 회원 정보 맞을 시 
        -  [X] `JWT` Token 발행
    - [X] 회원 정보가 안 맞을 시
        - [X] 오류 출력
    - [X] 회원정보 GORM를 통해서 불러오기
- [X] POST `/logout`: 로그아웃
    -  생성 되어 있는 쿠키 무력화
- [ ] `/edit`: 회원의 기본적인 정보를 수정할 수 있는 기능
    - [ ] 이메일 변경
    - [ ] 비밀번호 변경
    - [ ] 닉네임 변경
    - [ ] 전화번호 변경
- [ ] `/delete`: 회원 탈퇴


### Order `/order/{id}`
- [ ] 푸드마다 고유ID를 통한 푸드 및 음료를 식별할 수 있도록 개발
- [ ] 주문 설정
    - [ ] 사이즈
    - [ ] 컵

### Payment `/payment`
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

### Push `/push`
- [ ] 주문 시
    - 주문 확인 후 주문 몇 번째인지 확인 할 수 있도록 개발
        - Ex) 1번째 메뉴로 준비중입니다. (A-14)
    - 음료 및 푸드 픽업 대기 중
        - Ex) 메뉴가 모두 준비되었습니다. (A-14)  
        픽업대에서 메뉴를 픽업해주세요!  
        매장 방문시 마스크를 꼭 착용해주세요.
        - 준비 되었지만 5분 안에 픽업이 안 된 경우 2분마다 알람 재송신
            - Ex) 픽업대에서 기다리고 있어요.
- [ ] 광고 시
    - 상세 내용 추후 추가 예정

## Functions
### POST `/api/member/register`
#### Request
```json
{
    "name": "HyunSang Park",
    "nickname": "이것은 닉네임",
    "birthday": "2004-06-25",
    "email": "parkhyunsang@kakao.com",
    "password": "hello!@#"
}
```

#### Response
```json
{
    "uuid": "4eceee7e-2c74-490c-b094-9f5c959a654a",
    "name": "HyunSang Park",
    "nickname": "이것은 닉네임",
    "birthday": "2004-06-25",
    "email": "parkhyunsang@kakao.com",
    "password": "$2a$14$FtVJ.pqkp7BukLqnQRWehudF1l5CjUgVVfRIK8GQH9ust61Jpe0sO",
    "created_at": "2021-12-26T16:13:08.154225+09:00",
    "updated_at": "2021-12-26T16:13:08.169+09:00"
}
```
**중복되는 메일이 있는 경우:**
```json
{
    "message": "중복되는 메일이 있습니다."
}
```
**중복되는 닉네임이 있는 경우:**
```json
{
    "message": "중복되는 닉네임이 있습니다, 다시 확인 해 주세요."
}
```
**중복되는 전화번호가 있는 경우:**
```json
{
    "message": "중복되는 전화번호가 있습니다, 다시 확인 해 주세요."
}
```

### POST `/api/member/login`
#### Reqeust
```json
{
    "email": "parkhyunsang@kakao.com",
    "password": "hello!@#"
}
```

#### Response
```json
{
    "message": "success"
}
```

| Name| Value|Domain| Path |Expires|
|:----------:|------------|:------------:|:------------:|:------------:|
|jwt|eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA1OTIyNzIsImlzcyI6IjA1M2E0OTNmLWVhYjItNDExNC04ZmJkLWQwOTVkYWQxYmIyZCJ9.wdwZIbIGzlutmDpog1FYV0RX7aH-y2QyQDaGkJXopBA|localhost|/|Sun, 26 Dec 2021 09:04:32 GMT|

**메일 혹은 비밀번호가 틀린 경우:**
```json
{
    "message": "메일 혹은 비밀번호가 올바르지 않습니다. 다시 확인 해 주세요."
}
```

## 참고한 문서 및 글
- [GO 언어로 JWT 인증서버 만들기](https://covenant.tistory.com/203)
- [<go> fcm push 발송 (android, ios)](https://www.byfuls.com/programming/read?id=25)
- [WATCHA 푸시 서버 개선하기](https://medium.com/watcha/watcha-%ED%91%B8%EC%8B%9C-%EC%84%9C%EB%B2%84-%EA%B0%9C%EC%84%A0%ED%95%98%EA%B8%B0-56070b73c287)