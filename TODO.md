## 개선점들

1. 스택 트레이스 출력시, 어디 파일이 잘못되었는지, 몇 번째 줄인지, 몇 번째 칼럼인지에 대한 정보를 출력하지 않음
2. ASCII만 지원
3. 십진수만을 취급한다.
4. Block Statement는 별도의 Environment를 갖지 않는다.
5. Hash의 key로 유효한 것은 문자열, 정수, boolean뿐이다.
6. fn -> func으로 바꾸기
7. non-null type 추가
8. multi-line 입력을 지원하지 않는다.
9. bit operation이 없다.
10. ==, != 에서 객체 간의 비교를 어떻게 정의할 것인가?