그러면 FilterChain 구조에서 TLS Socket Filter는 어떻게 구현되어야 할까?

우선 TLS Socket을 구현하기에 앞서서 TLSConfig를 사용해야 하는데,  요구사항은 다음과 같아.
0. Filter는 다음 Filter 인터페이스를 구현해야 한다. TLSSocket은 Filter이다.
   type Filter interface {
   Init(config map[string]interface{}) error
   Handle(conn net.Conn) error
   SetNext(next Filter)
   }
1. ListenerFilter에서 net.Conn 객체를 전달받아, tls.Server를 통해서 통신할 수 있어야 한다.
2. TLSSocket은 궁극적으로 다음 필터에게 tls.Conn을 반환해야 한다.
3. tls.Server는 tls.Config와 net.Conn을 인자로 받는다.
4. TLSConfig는 KeyProvider와 TLS 통신을 위한 인자 값들을 가진다.
5. KeyProvider는 인증서와 키가 존재하는 위치를 string 형태로 가지고 있다.
6. KeyProvider는 os로부터 키 파일과 인증서 파일을 불러와, x509.Certificate와 rsa.PrivateKey를 반환할 수 있어야 한다.
7. 키 파일을 불러오고 나서, 별도의 x509.Certificate와 rsa.PrivateKey를 반환하는 함수를 구현한다.
8. TLSConfig는 궁극적으로 KeyProvider로부터 키와 인증서 파일을 반환 받아야 한다.
9. 궁극적으로 TLSSocket은 YAML을 통해서 구성될 수 있으나, TLSConfig가 이 부분을 해결한다. TLSSocket 객체는 처리에 관한 구현만 가진다. 

