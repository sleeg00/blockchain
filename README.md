# 오픈소스 블록체인 환경에서 리드 솔로몬 부호화 된 블록의 복구 성능 평가

## 🖥 프로젝트 소개️
- blockchain_go 오픈소스 환경에 gRPC, Reed-Solomon 부호화를 적용해 블록의 확장성을 넓히고 그로인한 오버헤드가 얼마나 되는지 성능 평가를 진행한다.
- 한국정보처리학회 ACK 2023(추계학술발표대회) 발표 논문에 사용되었습니다.
  - 부경대 총장상 수상
## ⏱️ 개발 기간
- 23-07 ~ 23.09
## ⚙️ 개발 환경
- Go 1.18.9 
- Database: Bolt DB
- IDE: VScode
- Open Source Link: https://github.com/Jeiwan/blockchain_go.git
## 🔦 기존 오픈소스에서 추가한 것
- TCP 통신 -> gRPC로 번경
- Reed-Solomon 부호화 적용
- Transaction 검증 방식 변경
- ... 
