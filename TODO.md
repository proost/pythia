## 개선점들

1. 스택 트레이스 출력시, 어디 파일이 잘못되었는지, 몇 번째 줄인지, 몇 번째 칼럼인지에 대한 정보를 출력하지 않음
2. ASCII만 지원
3. 십진수만을 취급한다.
4. Hash의 key로 유효한 것은 문자열, 정수, boolean뿐이다. -> null, float 추가 
5. non-null type 추가
6. multi-line 입력을 지원하지 않는다.
7. bit operation이 없다.
8. ==, != 에서 객체 간의 비교를 어떻게 정의할 것인가?
9. 함수 argument에 let 명시하도록 바꾸기