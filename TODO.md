## 개선점들

1. 스택 트레이스 출력시, 어디 파일이 잘못되었는지, 몇 번째 줄인지, 몇 번째 칼럼인지에 대한 정보를 출력하지 않음
2. ASC코드만 지원한다.
3. 알파벳과 _ 만을 변수명에 허용한다.
4. 십진수로 이루어진 정수만을 취급한다.
5. &&, || 연산자의 부재
   > makeTwoCharToken 메서드의 생성?
6. 후위 연산자의 부재
   > ++, -- 가 반드시 필요할까?