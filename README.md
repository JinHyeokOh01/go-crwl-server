# go-crwl-server

경희대학교 컴퓨터공학과 공지사항 & 소프트웨어중심대학사업단 공지사항을 크롤링하여

게시글 번호, 제목, 등록일, 링크를 응답하는 서버입니다.

최초 크롤링 후 DB에 저장하고, 이후 다시 크롤링할 때 DB에 있는 내용과 비교하여

새롭게 크롤링된 내용만 응답으로 처리합니다.

이 과정에서 DB와의 동기화가 이루어집니다.

위 모든 과정은 자동으로 동작하며, 크롤링 주기는 1시간으로 설정했습니다. main.go 파일에서 수정 가능합니다.

## 실행 방법
```
git clone https://github.com/JinHyeokOh01/go-crwl-server.git
```

```
make up
```

또는

```
 docker compose up -d
```

|method|URL|기능|
|------|---|---|
|GET|localhost:5000/cse|컴퓨터공학과 공지사항 수동 크롤링|
|GET|localhost:5000/sw|소프트웨어중심대학사업단 공지사항 수동 크롤링|
|GET|localhost:5000/notices/cse|현재 DB에 저장된 컴퓨터공학과 크롤링 내용|
|GET|localhost:5000/notices/sw|현재 DB에 저장된 소프트웨어중심대학사업단 크롤링 내용|
|DELETE|localhost:5000/notices|DB에 저장된 모든 내용 삭제|
|DELETE|localhost:5000/notices/cse|DB에 저장된 컴퓨터공학과 내용 삭제|
|DELETE|localhost:5000/notices/sw|DB에 저장된 소프트웨어중심대학사업단 내용 삭제|

실행 종료 시

```
make down
```

또는

```
 docker compose down
```

