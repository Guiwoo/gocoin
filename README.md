## Explore

## REST API

-[Strcut Field Tag](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)

- adapter pattern

## CLI

## Persistance

## Mining(POW)

- It's hard to find the answer for computer but it will be able to be verified easily

  - Need to reWatching

## Transaction

- How transaction work ?
  - Tx [Tx In,Tx Out]
  - Tx in ? means money that you have before tx
  - Tx out means the money everyone have the your money history

## Wallet

- 존나 재밌다 지갑

- be able to check own transaction input && output
- verified apporve the transactions

1. Verifing and Sign without blockchain
2. Persist Wallet on DB ,Restore
3. apply on transactions

## P2P

- [Dive In Channel](https://www.velotio.com/engineering-blog/understanding-golang-channels#:~:text=So%2C%20what%20are%20the%20channels,put%20or%20read%20the%20data.)
  - what haeppend make chan ?
    - Unbuffered Channel (사이즈를 정하지 않은경우)
    - Buffered Channel (사이즈를 정한경우)
      - 총 용량이 가득 찾을때 채널은 보내고, 받는 고루틴 그것을 받아들일수 있다.
      - 채널을 받는것을 한번에 총용량 만큼 받아서 실행하는것이 가능ㅎ다ㅏ.
  - 채널 로 전송할것이 끝나면 close() 닫아주고
  - 채널을 값으로 받는 상위 함수에서 채널이 있는지 없는지 확인해주는 함수로 함수를 종료할지 이어갈지 정해준다
- websocket is main !

- EverySingleNode will run go routine
