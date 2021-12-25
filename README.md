# Siren Order Services(System)
**스타벅스의 사이렌 오더 서비스(기능)들을 보고 클론 코딩을 하는 프로젝트입니다 ☕️**

## Getting Started
### FireBase Settings
#### Add SDK
```shell
# Install as a module dependency
$ go get firebase.google.com/go/v4
# Install to $GOPATH
$ go get firebase.google.com/go
```
#### Reset SDK
```shell
$ export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json"
```
**서비스 계정의 비공개 키 파일을 생성하려면 다음 안내를 따르세요.**
1. Firebase Console에서 설정 > 서비스 계정을 엽니다.
2. 새 비공개 키 생성을 클릭한 다음 키 생성을 클릭하여 확인합니다.
3. 키가 들어 있는 JSON 파일을 안전하게 저장합니다.

```shell
$ brew install --cask google-cloud-sdk
$ gcloud auth application-default login
```
- [Could not load the default credentials? (Node.js Google Compute Engine tutorial)](https://stackoverflow.com/questions/42043611/could-not-load-the-default-credentials-node-js-google-compute-engine-tutorial)

#### Create GCP Storage Bucket
- [Storage Bucket Locations](https://cloud.google.com/storage/docs/locations)

## ToDo
### User Account `/member`
- [ ] DataBase 설정
    - [ ] 구조체를 이용한 테이블 생성함
- [ ] Redis를 이용한 JWT 인증 서버
    - [ ] JWT 메타데이터 정의 
    - [ ] JWT 메타데이터 저장
- [ ] `/register`: 기본적인 회원가입
    - [ ] UUID를 통한 회원 식별을 할 수 있도록 함.
    - [ ] Password 암호화를 함.
    - [ ] 회원의 생일을 입력 받을 수 있도록 함.
- [ ] `/login`: 기본적인 로그인
    - [ ] 회원 정보 맞을 시 
        - [ ] `JWT` Token 발행
    - [ ] 회원 정보가 안 맞을 시
        - [ ] 오류 출력
    - [ ] 회원정보 GORM를 통해서 불러오기
- [ ] `/edit`: 회원의 기본적인 정보를 수정할 수 있는 기능
    - [ ] 이메일 변경
    - [ ] 닉네임 변경
    - [ ] 전화번호 변경
- [ ] `/delete`: 회원 탈퇴

### Auth `/auth`
- 개발하면서 기능에 대한 설명 추가 예정.

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
- [] 주문 시
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
### POST `/member/register`
#### Request
```json
{
    "email": "parkhyusang@kakao.com",
    "password": "test!@#",
    "phone_number": "+15555550100",
    "name": "HyunSang Park"
}
```
#### Response
```json
{
    "data": {
        "displayName": "HyunSang Park",
        "email": "parkhyusang@kakao.com",
        "phoneNumber": "+15555550100",
        "providerId": "firebase",
        "rawId": "qF93zq9M2XOc9YcB4q0jO5rXiKh1",
        "CustomClaims": null,
        "Disabled": false,
        "EmailVerified": false,
        "ProviderUserInfo": [
            {
                "phoneNumber": "+15555550100",
                "providerId": "phone",
                "rawId": "+15555550100"
            },
            {
                "displayName": "HyunSang Park",
                "email": "parkhyusang@kakao.com",
                "providerId": "password",
                "rawId": "parkhyusang@kakao.com"
            }
        ],
        "TokensValidAfterMillis": 1640421875000,
        "UserMetadata": {
            "CreationTimestamp": 1640421875059,
            "LastLogInTimestamp": 0,
            "LastRefreshTimestamp": 0
        },
        "TenantID": ""
    },
    "message": "Successfully Created User!",
    "status": 200,
    "time": "2021-12-25T17:44:35.611168+09:00"
}
```

## 참고한 문서 및 글
- [GO 언어로 JWT 인증서버 만들기](https://covenant.tistory.com/203)
- [<go> fcm push 발송 (android, ios)](https://www.byfuls.com/programming/read?id=25)
- [WATCHA 푸시 서버 개선하기](https://medium.com/watcha/watcha-%ED%91%B8%EC%8B%9C-%EC%84%9C%EB%B2%84-%EA%B0%9C%EC%84%A0%ED%95%98%EA%B8%B0-56070b73c287)