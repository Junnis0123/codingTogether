## API 스펙 요약

- [Auth](#auth)
> 1. [POST] /auth/login
> 2. [GET] /auth/duplication/{userID}

- [User control](#user)
> 1. [POST] /users
> 2. [GET] /users
> 3. [GET] /users/{userID}

- [CodingTogether control](#codingTogether)
> 1. [GET] /codingTogethers
> 2. [POST] /codingTogether/
> 3. [GET] /codingTogether/me
> 4. [GET] /codingTogether/{codingTogetherIdx}  

## API 스펙 상세 설명
### Auth
1. 로그인
>- method : POST
>- endpoint : /auth/login
>- Description : ID 와 PW 를 이용하여 로그인
>- Example :
> > URL : [POST] /auth/login
> > Request Body :  
> > ```
> > { 'userID' : 'sool' ,
> >   'userPW' : base64('1234')
> > }
> > ```
> > Response Body : 
> > ```
> > {
> >     "success": {true} or {false} //로그인 여부
> >     "message": "로그인 성공" or "로그인 실패",
> >     "errors": "" or "Not correct ID or PW",
> >     "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1OThkZGI2MzIyYWMxMDExZTA3MDJjYjAiLCJ1c2VybmFtZSI6InRlc3QxIiwibmFtZSI6InRlc3QxIiwiZW1haWwiOiIiLCJpYXQiOjE1MDQ3MzI2NzcsImV4cCI6MTUwNDgxOTA3N30.4eG2zGpSeY2XezKB4Djf6usy7DdygIybR1VKUBj-ScE"
> >     "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTgzNDI4NjgsInJlYWxfaXAiOiI6OjEiLCJ1c2VyX2lkeCI6IjMifQ.YQvmXfT1hLwGre-2MUaaH2XzLPJOICMxy63_grPoLGs"
> > }
> > ```

2. 아이디 중복 체크
>- method : GET
>- endpoint : /auth/duplication/{userID}
>- Description : ID를 이용하여 중복 체크
>- Example :
> > URL : [POST] /auth/duplication/sool
> > Response Body : 
> > HTTP Response : 200(중복 없음) or 409(중복 존재)
> > ```
> > {
> >     "success": {true} or {false} //true = 중복 X, false = 중복
> >     "message": "사용가능한 ID입니다." or "중복된 ID입니다.",
> >     "errors": "" or "Duplicate ID",
> >     "data": 
> > }
> > ```

### User
1. 회원가입
>- method : POST
>- endpoint : /users/
>- Description : ID, PW, Nickname으로 회원가입을 합니다.
>- Example :
> > URL : [POST] /users/
> > Request Body :  
> > ```
> > { 'userID' : 'junnis0123' ,
> >   'userPW' : base64('0123'),
> >   'userNickname' : '킹선진'
> > }
> > ```
> > Response Body :  
> > HTTP Response : 200(성공) or 500(실패)
> > ```
> > {
> >     "success": {true} or {false} //성공 여부
> >     "message": "회원 가입 완료" or "회원 가입 실패",
> >     "errors": "" or "Not correct ID",
> >     "data": ""
> > }
> > ```

2. 회원 id로 조회 - 테스트용임, 실제로는 불가능해야 함.
>- method : GET
>- endpoint : /users/test/{userID}
>- Description : 테스트용 특정 USER 조회 기능
>- Example :
> > URL : [GET] /users/test/{userID}  
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > /users/sool
> > ```
> > Response Body :
> > ```
> > {
> >     "success": {true} or {false} //성공 여부
> >     "message": "",
> >     "errors": "",
> >     "data": "기무서무지니"
> > }
> > ```

3. 전체 회원 조회
>- method : GET
>- endpoint : /users/
>- Description : 회원 전체를 조회하여 userID : userNickname 으로 반환합니다.
>- Example :
> > URL : [GET] /users/
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > /users/
> > ```
> > Response Body : 
> > ```
> > {
> >     "success": {true}
> >     "message": "전체 조회 완료",
> >     "errors": "",
> >     "data": "[{\"userID\":\"sool\",\"userNickname\":\"기무서무지니\"},{\"userID\":\"duck\",\"userNickname\":\"덕\"}]"
> > }
> > ```

4. 자기정보 조회
>- method : GET
>- endpoint : /users/me
>- Description : 자기 자신 닉네임 확인.
>- Example :
> > URL : [GET] /users/me  
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > ```
> > Response Body :  
> > HTTP Response : 200(성공) or 204(존재 없음) or 401(토큰 인증 문제)
> > ```
> > {
> >     "success": {true} or {false} //성공 여부
> >     "message": "조회 성공" or "조회 실패",
> >     "errors": "" or "Not correct user idx",
> >     "data": "기무서무지니"
> > }
> > ```

### CodingTogether
1. 전체 모임 조회
>- method : GET
>- endpoint : /codingTogethers/
>- Description : 모임 전체 목록을 반환합니다.
>- Example :
> > URL : [GET] /codingTogethers/
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > /codingTogethers/
> > ```
> > Response Body : 
> > HTTP Response : 200(성공) or 401(토큰 인증 문제)
> > ```
> > {
> >     "success": {true}
> >     "message": "조회 성공", or "조회 실패"
> >     "errors": "Lookup Failure",
> >     "data": "[{\"codingTogetherOrgnizerName\":\"기무서무지니\",\"codingTogetherCreateTime\":\"2020-07-27 17:05:50\",\"codingTogetherIdx\":1,\"codingTogetherImgURL\":{\"String\":\"\",\"Valid\":false},\"codingTogetherMemberCount\":4,\"codingTogetherName\":\"테스트1\"},{\"codingTogetherOrgnizerName\":\"덕\",\"codingTogetherCreateTime\":\"2020-07-27 17:06:55\",\"codingTogetherIdx\":2,\"codingTogetherImgURL\":{\"String\":\"\",\"Valid\":false},\"codingTogetherMemberCount\":3,\"codingTogetherName\":\"테스트2\"}]"
> > }
> > ```

2. 모임 만들기
>- method : POST
>- endpoint : /codingTogethers/
>- Description : 모임을 생성한다. //좀 더 보안이 필요해보이지만 일단 넘어간다.
>- Example :
> > URL : [POST] /codingTogethers/
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > { 'codingTogetherName' : '모각코모임테스트1' ,
> >   'codingTogetherContents' : '테스트입니다. 확인 !!!!!'
> > }
> > ```
> > Response Body : 
> > HTTP Response : 200(성공) or 401(토큰 인증 문제) or 500(알 수 없는 에러)
> > ```
> > {
> >     "success": {true} or {false}
> >     "message": "모각코 생성 성공", or "모각코 생성 실패"
> >     "errors": "CodingTogether Create Failure",
> >     "data": ""
> > ```

3. 자기 가입 모임 조회
>- method : GET
>- endpoint : /codingTogethers/me
>- Description : 본인이 가입한 모임을 조회한다.
>- Example :
> > URL : [POST] /codingTogethers/me
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > ```
> > Response Body : 
> > HTTP Response : 200(성공) or 401(토큰 인증 문제) or 500(알 수 없는 에러)
> > ```
> > {
> >     "success": {true} or {false}
> >     "message": "참가 모각코 조회 섣공", or "참가 모각코 조회 실패"
> >     "errors": "Lookup Failure",
> >     "data": "[{\"codingTogetherOrgnizerName\":\"기무서무지니\",\"codingTogetherCreateTime\":\"2020-07-27 17:05:50\",\"codingTogetherIdx\":1,\"codingTogetherImgURL\":{\"String\":\"\",\"Valid\":false},\"codingTogetherMemberCount\":4,\"codingTogetherName\":\"테스트1\"},{\"codingTogetherOrgnizerName\":\"덕\",\"codingTogetherCreateTime\":\"2020-07-27 17:06:55\",\"codingTogetherIdx\":2,\"codingTogetherImgURL\":{\"String\":\"\",\"Valid\":false},\"codingTogetherMemberCount\":3,\"codingTogetherName\":\"테스트2\"}]"
> > }
> > ```

4. 모임 상세 정보 조회
>- method : GET
>- endpoint : /codingTogethers/{codingTogetherIdx}
>- Description : 본인이 가입한 모임을 조회한다.
>- Example :
> > URL : [POST] /codingTogethers/1
> > Request Header : "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY"  
> > Request Body :  
> > ```
> > ```
> > Response Body : 
> > HTTP Response : 200(성공) or 401(토큰 인증 문제) or 500(알 수 없는 에러)
> > ```
> > {
> >     "success": {true} or {false}
> >     "message": "조회 섣공", or "조회 실패"
> >     "errors": "Lookup Failure",
> >     "data": ""{\"codingTogetherContents\":\"1번모임이다능\"}"
> > }
> > ```