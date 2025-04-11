# HTTP Connection Manager (HCM) Abort

HTTP 요청을 파싱하고 처리하며, **라우팅(Route)**을 담당합니다.

http_filters를 사용하여 더 세부적인 처리를 설정할 수 있어요.
```
http_filters:
  - name: envoy.filters.http.router
```

음 HCM은 tls.Conn으로부터 받은 연결을 HTTP.request나 아니면 http.responsewriter를 이용해서 응답을 읽고 써야해.

이러한 작업을 수행하기 위해서 ServiceEngine을 이용해. ServiceEngine은 Method GET, POST를 지원해야해.
GET을 수행한 경우에는 WebRoot 접근
POST을 수행한 경우에는 WebRoot에 파일 업로드

그리고 Scipt suport가 되어야 해            