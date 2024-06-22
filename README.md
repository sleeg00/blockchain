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
### 노드별 블록 스토리지 효율성 비교 & 블록 전송 성능 비교
<img width="475" alt="스크린샷 2023-11-16 오후 4 17 28" src="https://github.com/sleeg00/blockchain/assets/96710732/421ba419-23e8-46e4-93f1-b89a8a7d62b2">
<img width="475" alt="스크린샷 2023-11-16 오후 4 17 44" src="https://github.com/sleeg00/blockchain/assets/96710732/0a365604-8944-4bf6-820d-4425ac48e50d">

### CPU 사용량 비교
<img width="615" alt="스크린샷 2023-11-16 오후 4 17 53" src="https://github.com/sleeg00/blockchain/assets/96710732/7ed9e52f-85d0-4a08-97fa-d4a953c84a8d">

### 블록 수에 따른 블록 복구 시간
<img width="815" alt="스크린샷 2024-06-22 오후 4 24 40" src="https://github.com/sleeg00/blockchain/assets/96710732/4c1e7475-bae6-45ca-886f-1505624f94a8">

### 블록 인코딩 및 디코딩을 위한 시퀀스 다이어그램
<img width="615" alt="스크린샷 2024-06-22 오후 4 25 32" src="https://github.com/sleeg00/blockchain/assets/96710732/b82281c0-4b86-4d28-81f0-ab47021b9760">
