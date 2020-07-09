# diep.io-clone
![GitHub License](https://img.shields.io/github/license/Diep-clone/diep-clone_golang?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/diep-clone/diep-clone_golang?style=flat-square)
[![HitCount](http://hits.dwyl.com/Diep-clone/diep-clone_golang.svg)](http://hits.dwyl.com/Diep-clone/diep-clone_golang)

This repo is Diep.io's Clone Version.

The goal of this project is to copy the web game diep.io as same as possible.
Currently, this project is stopped practically due to the vacancy of lead programmer.
If you want to be the maintainer of this project although this might be poor because it is my first github project, please contact:

ybs1164@naver.com

만약 이가 아니더라도 프로젝트에 도움을 주고 싶다면 언제든지 연락해주시면 감사하겠습니다.

## 📚 Project Stack
### Front end (JS)
- Webpack
- babel

### Back end (Golang)
- gorilla/websocket
- sirupsen/logrus

Before introducing this project,
these are the nicknames of the "researchers" who gave their attention to this project and have researched various needed information about diep.io.
Without them, this project would not have progressed.
I would like to express my gratitude to the following researchers:

[K] YouTube.
피로⚛Firo_SF.
항해사.
케빈.

Plus, I would like to thank 레인우드, who made the foundation for this project.

현재 이 프로젝트에 구현되지 못한 점들은 다음과 같습니다. (추가 될 시 수정하겠습니다.)

1. 반동 시스템.
- 이 시스템이 아직 제대로 연구되지 못했습니다. 사실상 이 프로젝트가 어색한 가장 큰 이유입니다.
2. 도형.
- 네모를 제외한 다른 도형들과 (물론 이건 구현은 가능하겠지만) 도형들의 움직임, 도형들이 스폰되는 위치 (플레이어 근방에 스폰되지 않아야 함), 그리고 도형들의 스폰율이 구현되지 못했습니다.
3. UI.
- 스코어보드를 비롯한 모든 UI 가 구현되지 않았습니다. 또한 오브젝트의 이름과 점수, 체력 바의 위치 같은 것들이 제대로 연구되지 못했습니다.
4. 보스.
5. 게임모드.
6. 서버.
- 원래 클라이언트는 여러 서버에 유동적으로 접속할 수 있는 형태로 제작되어야 하는데, 기술적인 문제로 이게 구현되지 못했습니다.
7. 전직.
- 서버 쪽에서 탱크가 다른 탱크로 바뀌는 시스템을 구현하지 못했습니다.


이 외에 말씀드릴 점은 다른 유저 분들의 자작탱들이 프로젝트에 일부 구현되어 있습니다.
헤로쿠 서버에서 플레이 가능하실 것 같네요.

---
Copyright © 2019-2020
