echo "# gocoin" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/Guiwoo/gocoin.git
git push -u origin main

# push an existing repository from the command line

git remote add origin https://github.com/Guiwoo/gocoin.git
git branch -M main
git push -u origin main

# Writen by Guiwoo

|---------------------------------------------------------|
|1. fmt.Sprintf("%b,%x","etetete") => formating the string|
|---------------------------------------------------------|

|-----------------------------------------------------|
|2. One funcing wokrs one thing, Receiver Funcion |
|*Receve Func ? ex) Do not want to copy |
|type blockchain struct {blocks []block} |
|func (b *blockchain) getLastHash() string {return""} |
|------------------------------------------------------|

|------------------------------------------------------------|
|3. Sigleton_pattern|
| Do not share instance directly , make a function to control|
| 최초 한번만 메모리를 할당해서 메모리 낭비를 방지하고,|
| blockchain에 다른 패키지에서 접근하고 공유하는 걸 쉽게 하고,|
| 전역으로 선언해서 해당 인스턴스가 절대적으로 한개만 존재함을 보증 하기 때문에|
|-------------------------------------------------------------|

|4. Pacaga Sync|
|병렬적으로 실행가능한 go, 동기적 처리해야하는부분,한번만실행|
