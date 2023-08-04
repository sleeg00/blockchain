version: '3'
services:
  app-1:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-1
    ports:
      - "3000:3000"
    depends_on:
      - app-2
  app-2:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-2
    ports:
      - "3001:3000"
    depends_on:
      - app-3
  app-3:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-3
    ports:
      - "3002:3000"
    depends_on:
      - app-4
  app-4:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-4
    ports:
      - "3003:3000"
    depends_on:
      - app-5
  app-5:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-5
    ports:
      - "3004:3000"
    depends_on:
      - app-6
  app-6:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-6
    ports:
      - "3005:3000"
    depends_on:
      - app-7
  app-7:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-7
    ports:
      - "3006:3000"
    depends_on:
      - app-8
  app-8:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-8
    ports:
      - "3007:3000"
    depends_on:
      - app-9
  app-9:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-9
    ports:
      - "3008:3000"
    depends_on:
      - app-10
  app-10:
    build: /Users/user/Documents/my/dev/go/src/sleeg00/blockchain/app-10
    ports:
      - "3009:3000"
