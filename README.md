# 오픈소스 블록체인 환경에서 리드 솔로몬 부호화 된 블록의 복구 성능 평가

## 🖥 프로젝트 소개️
- blockchain_go 오픈소스 환경에 gRPC, Reed-Solomon 부호화를 적용해 블록의 확장성을 넓히고 그로인한 오버헤드가 얼마나 되는지 성능 평가를 진행한다.
- 한국정보처리학회 ACK 2023(추계학술발표대회) 발표 논문에 사용되었습니다.
  - 부경대 총장상 수상
  - TKIPS 등재지 게재
## ⏱️ 개발 기간
- 23-07 ~ 24.01
## ⚙️ 개발 환경
- Go 1.18.9 
- Database: Bolt DB
- IDE: VScode
- Open Source Link: https://github.com/Jeiwan/blockchain_go.git
## 🔦 기존 오픈소스에서 추가한 것
- gRPC로 데이터 전송 기능 변경
- Reed-Solomon 부호화 적용
- Transaction 검증 방식 변경 (경량 스레드인 Go루틴 추가)
- CPU 오버헤드 성능 평가
- ... 

## 🖼️ 주요 그림
<img width="475" alt="스크린샷 2023-11-16 오후 4 17 28" src="https://github.com/sleeg00/blockchain/assets/96710732/46468a99-1c63-4be3-89ec-4367cd9a8bee">
<img width="320" alt="스크린샷 2023-11-16 오후 4 17 44" src="https://github.com/sleeg00/blockchain/assets/96710732/603bc79c-0d4a-4cda-96c3-fb7d6b047b7c">
<img width="310" alt="스크린샷 2023-11-16 오후 4 17 53" src="https://github.com/sleeg00/blockchain/assets/96710732/755cb5e4-f647-45e0-a839-5365f46ec5ea">
<img width="250" alt="스크린샷 2023-11-16 오후 4 18 01" src="https://github.com/sleeg00/blockchain/assets/96710732/034fcdc6-bfb5-4c8a-9c59-21fc0f2ec70b">
