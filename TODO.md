## 개선점들

1. 스택 트레이스 출력시, 어디 파일이 잘못되었는지, 몇 번째 줄인지, 몇 번째 칼럼인지에 대한 정보를 출력하지 않음
2. ASCII만 지원
3. 십진수만을 취급한다.
4. makeTwoCharToken 메서드를 생성한다.
   1. 논리 연산자 &&, || 연산자의 부재
   2. 논리 연산자 <=, >= 의 부재
5. Block Statement는 별도의 Environment를 갖지 않는다.
6. Array access시에 size를 넘어가는 값에 Null을 리턴한다.
7. Hash의 key로 유효한 것은 문자열, 정수, boolean뿐이다.
8. type() built-in 함수를 만들 것
9. non-null type 추가
10. multi-line 입력을 지원하지 않는다.