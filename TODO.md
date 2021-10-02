## 개선점들

1. ASCII만 지원
2. 십진수만을 취급한다.
3. multi-line 입력을 지원하지 않는다.
4. return 하나만 해서 void값 만들기
5. assignment가 array나, hash일 경우에 index operation이 제대로 동작 안한다.
6. BlockStatement가 문법적으로 완전하지 않아도 에러가 안뜸
7. 함수 정의 시, 이름을 추가할 수 있도록 
9. type(nil) == type(void) 입력시 NULL 값 출력하도록 수정하기. 