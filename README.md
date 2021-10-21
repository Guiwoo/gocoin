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

-1. fmt.Sprintf("%b,%x","etetete") => formating the string
-2. One funcing wokrs one thing, Receiver Funcion
*Receve Func ? ex) Do not want to copy
type blockchain struct {blocks []block}
func (b *blockchain) getLastHash() string {return""}
